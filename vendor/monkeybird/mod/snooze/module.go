// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package snooze defines the !snooze command. It allows users to
// set (and unset) timed alarms. These persist across bot restarts.
package snooze

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/irc/proto"
	"monkeybird/mod"
	"monkeybird/text"
	"monkeybird/tr"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// alarm defines a single scheduled alarm.
type alarm struct {
	SenderMask string
	Target     string
	Message    string
	When       time.Time
}

type module struct {
	m        sync.RWMutex
	root     string
	commands *cmd.Set
	quitOnce sync.Once
	quit     chan struct{}
	table    map[string]alarm
	writer   irc.ResponseWriter
}

// New returns a new module.
func New(rw irc.ResponseWriter) mod.Module {
	return &module{
		table:  make(map[string]alarm),
		quit:   make(chan struct{}),
		writer: rw,
	}
}

// Load loads module resources and binds commands.
func (m *module) Load(pb irc.ProtocolBinder, prof irc.Profile) {
	pb.Bind("PRIVMSG", m.onPrivMsg)

	m.root = prof.Root()
	m.commands = cmd.New(prof.CommandPrefix(), nil)
	//m.commands.Bind(tr.SnoozeName, tr.SnoozeDesc, false, m.cmdSnooze).
	//	Add(tr.SnoozeTimeName, tr.SnoozeTimeDesc, true, cmd.RegAny).
	//	Add(tr.SnoozeMessageName, tr.SnoozeMessageDesc, false, cmd.RegAny)

	//m.commands.Bind(tr.UnsnoozeName, tr.UnsnoozeDesc, false, m.cmdUnsnooze).
	//	Add(tr.UnsnoozeIDName, tr.UnsnoozeIDDesc, true, cmd.RegAny)

	m.load()
	go m.poll()
}

// Unload cleans up library resources and unbinds commands.
func (m *module) Unload(pb irc.ProtocolBinder, prof irc.Profile) {
	m.quitOnce.Do(func() {
		close(m.quit)

		m.save()
		m.writer = nil
		m.table = nil
		m.commands.Clear()
		pb.Unbind("PRIVMSG", m.onPrivMsg)
	})
}

func (m *module) Help(w irc.ResponseWriter, r *cmd.Request) {
	m.commands.HelpHandler(w, r)
}

func (m *module) onPrivMsg(w irc.ResponseWriter, r *irc.Request) {
	m.commands.Dispatch(w, r)
}

// cmdUnsnooze unregisters a scheduled alarm. Only the person with the
// sae hostmask as the one who initially set the alarm, can do this.
func (m *module) cmdUnsnooze(w irc.ResponseWriter, r *cmd.Request) {
	id := strings.ToLower(r.String(0))

	m.m.Lock()

	a, ok := m.table[id]
	if ok && strings.EqualFold(a.SenderMask, r.SenderMask) {
		delete(m.table, id)
		proto.PrivMsg(w, r.Target, tr.SnoozeAlarmUnset, r.SenderName)
		m.save()
	}

	m.m.Unlock()
}

// cmdSnooze lets the user schedule an alarm at a specific time, or after
// N minutes. When triggered, the alarm sends a user-defined message
// to the channel, while pinging the user.
//
// Once scheduled, this gives the user a cancelation code they can
// use with the unsnooze command. It unschedules the alarm.
func (m *module) cmdSnooze(w irc.ResponseWriter, r *cmd.Request) {
	id := m.createID()

	if !m.addSnooze(w, r, id) {
		m.deleteID(id)
	}
}

// addSnooze does what the docs on cmdSnooze describe. This is a separate
// method with the unique id as added parameter to make unit test code
// easier to write. This returns false if the alarm was not scheduled.
// This can happen when the tim value is invalid. If this is the case, the
// given id should either be removed from the table, or reused.
func (m *module) addSnooze(w irc.ResponseWriter, r *cmd.Request, id string) bool {
	when := parseTime(r.String(0))
	if when <= 0 {
		proto.PrivMsg(w, r.Target, tr.SnoozeInvalidTime, r.SenderName, r.String(0))
		return false
	}

	msg := r.Remainder(2)
	if len(msg) == 0 {
		msg = tr.SnoozeDefaultMessage
	} else {
		msg = tr.SnoozeMessagePrefix + msg
	}

	m.m.Lock()

	m.table[id] = alarm{
		Target:     r.Target,
		SenderMask: r.SenderMask,
		Message:    fmt.Sprintf(msg, r.SenderName),
		When:       time.Now().Add(when),
	}

	m.save()
	m.m.Unlock()

	proto.PrivMsg(w, r.Target, tr.SnoozeAlarmSet, r.SenderName, text.Bold(id))
	return true
}

// deleteID removes the given id from the table, if it exists.
func (m *module) deleteID(id string) {
	m.m.Lock()
	delete(m.table, id)
	m.m.Unlock()
}

// createID returns a new, unique id for an alarm. This id can be used as
// a cancellation code. Note that this call will create a new table entry
// with the given id, so subsequent calls to createID() will not accidentally
// re-use the generated one before the caller can.
func (m *module) createID() string {
	m.m.Lock()
	defer m.m.Unlock()

	var key [5]byte

	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var generate = func() string {
		for i := 0; i < len(key); i++ {
			key[i] = alphabet[rng.Intn(len(alphabet))]
		}

		return string(key[:])
	}

	id := generate()
	for {
		if _, ok := m.table[id]; !ok {
			break
		}

		id = generate()
	}

	m.table[id] = alarm{}
	return id
}

// poll periodically checks the schedule table for expired alarms.
func (m *module) poll() {
	for {
		select {
		case <-m.quit:
			return

		case <-time.After(time.Minute):
			m.checkExpiredAlarms()
		}
	}
}

// checkExpiredAlarms checks for expired alarms.
// When found, it sends the appropriate notification.
func (m *module) checkExpiredAlarms() {
	m.m.Lock()
	defer m.m.Unlock()

	var deleted bool
	now := time.Now()

	for id, alarm := range m.table {
		if now.Before(alarm.When) {
			continue
		}

		proto.PrivMsg(m.writer, alarm.Target, alarm.Message)

		// Yes, this is safe.
		// ref: http://stackoverflow.com/a/23231539/357705
		delete(m.table, id)

		deleted = true
	}

	if deleted {
		m.save()
	}
}

// save writes scheduled alarm data to a file.
func (m *module) save() {
	file := filepath.Join(m.root, "snooze.dat")

	data, err := json.Marshal(m.table)
	if err != nil {
		log.Println("[snooze] json.Marshal:", err)
		return
	}

	err = ioutil.WriteFile(file, data, 0600)
	if err != nil {
		log.Println("[snooze] ioutil.WriteFile:", err)
		return
	}
}

// load loads scheduled alarm data from a file.
func (m *module) load() {
	file := filepath.Join(m.root, "snooze.dat")

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("[snooze] ioutil.ReadFile:", err)
		return
	}

	err = json.Unmarshal(data, &m.table)
	if err != nil {
		log.Println("[snooze] json.Unmarshal:", err)
		return
	}
}

// parseTime treats the given value as either an absolute time, or
// an offset in minutes. It returns the value which represents the
// duration between now and then.
func parseTime(v string) time.Duration {
	then, err := time.Parse(tr.SnoozeTimeFormat, v)

	if err == nil {
		// We expect the given time to include only the time.
		// We must set the date components manually.

		now := time.Now()
		then = time.Date(now.Year(), now.Month(), now.Day(),
			then.Hour(), then.Minute(), then.Second(), 0, now.Location())

		delta := then.Sub(now)

		// If delta is negative, we are probably dealing with a time which
		// is meant to mean 'tomorrow'. So add 24 hours to the clock and
		// recalculate the difference.
		if delta < 0 {
			then = then.Add(time.Hour * 24)
			delta = then.Sub(now)
		}

		return delta
	}

	// If not an absolute time, the value is expected to be an offset
	// from the current time, in minutes.
	num, err := strconv.ParseInt(v, 10, 32)
	if err == nil {
		// This can result in a negative duration, if someone specified
		// "-10" as the input. This is an error which is caught by the caller.
		return time.Duration(num) * time.Minute
	}

	return 0
}
