// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package alarm allows a user to schedule an alarm with a custom message.
// The alarm can be scheduled at an exact time or an offset from the
// current time. Once a scheduled alarm's time has come, the bot will notify
// the user who scheduled it. Alarms can be unscheduled by the user who
// scheduled it.
//
// Create a new alarm for 10 minutes from now:
//
//    <steve> !reminder 10 Make food.
//
// Create a new alarm for 18:15:
//
//    <steve> !reminder 18:15 Make food.
//
package alarm

import (
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/monkeybird/autimaat/app/util"
	"github.com/monkeybird/autimaat/irc"
	"github.com/monkeybird/autimaat/irc/cmd"
	"github.com/monkeybird/autimaat/irc/proto"
	"github.com/monkeybird/autimaat/plugins"
)

func init() { plugins.Register(&plugin{}) }

// alarm defines a single scheduled alarm.
type alarm struct {
	SenderMask string
	SenderName string
	Target     string
	Message    string
	When       time.Time
}

type plugin struct {
	m        sync.RWMutex
	file     string
	cmd      *cmd.Set
	table    map[string]alarm
	quitOnce sync.Once
	quit     chan struct{}
}

// Load initializes the module and loads any internal resources
// which may be required.
func (p *plugin) Load(prof irc.Profile) error {
	p.quit = make(chan struct{})
	p.table = make(map[string]alarm)
	p.file = filepath.Join(prof.Root(), "alarm.dat")

	p.cmd = cmd.New(prof.CommandPrefix(), nil)
	p.cmd.Bind(TextReminder, false, p.onReminder).
		Add(TextTimestamp, true, cmd.RegAny).
		Add(TextMessage, false, cmd.RegAny)
	p.cmd.Bind(TextClearReminder, false, p.onClearReminder).
		Add(TextID, true, cmd.RegAny)

	go p.pollReminders()
	return util.ReadFile(p.file, &p.table, true)
}

// Unload cleans the module up and unloads any internal resources.
func (p *plugin) Unload(prof irc.Profile) error {
	p.quitOnce.Do(func() {
		close(p.quit)
	})
	return nil
}

// Dispatch sends the given, incoming IRC message to the plugin for
// processing as it sees fit.
func (p *plugin) Dispatch(w irc.ResponseWriter, r *irc.Request) {
	p.cmd.Dispatch(w, r)
}

// onReminder lets a user schedule a new alarm.
func (p *plugin) onReminder(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	id := p.createID()

	if !p.addReminder(w, r, params, id) {
		p.deleteID(id)
	}
}

// onClearReminder lets a user remove an existing alarm.
func (p *plugin) onClearReminder(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	id := strings.ToLower(params.String(0))

	p.m.Lock()

	a, ok := p.table[id]
	if ok && strings.EqualFold(a.SenderMask, r.SenderMask) {
		delete(p.table, id)
		proto.PrivMsg(w, r.Target, TextAlarmUnset, r.SenderName)
		util.WriteFile(p.file, p.table, true)
	}

	p.m.Unlock()
}

// addReminder does what the docs on addReminder describe. This is a separate
// method with the unique id as added parameter to make unit test code
// easier to write. This returns false if the alarm was not scheduled.
// This can happen when the tim value is invalid. If this is the case, the
// given id should either be removed from the table, or reused.
func (p *plugin) addReminder(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList, id string) bool {
	when := parseTime(params.String(0))
	if when <= 0 {
		proto.PrivMsg(w, r.Target, TextInvalidTime, r.SenderName, params.String(0))
		return false
	}

	msg := strings.Join(r.Fields(2), " ")
	if len(msg) == 0 {
		msg = TextDefaultMessage
	} else {
		msg = TextMessagePrefix + msg
	}

	p.m.Lock()

	p.table[id] = alarm{
		Target:     r.Target,
		SenderMask: r.SenderMask,
		SenderName: r.SenderName,
		Message:    msg,
		When:       time.Now().Add(when),
	}

	util.WriteFile(p.file, p.table, true)
	p.m.Unlock()

	proto.PrivMsg(w, r.Target, TextAlarmSet, r.SenderName, util.Bold(id))
	return true
}

// pollReminders periodically checks if any of the defined reminders have expired.
func (p *plugin) pollReminders() {
	for {
		select {
		case <-p.quit:
			return

		case <-time.After(time.Minute):
			p.checkExpiredAlarms()
		}
	}
}

// deleteID removes the given id from the table, if it exists.
func (p *plugin) deleteID(id string) {
	p.m.Lock()
	delete(p.table, id)
	p.m.Unlock()
}

// createID returns a new, unique id for an alarm. This id can be used as
// a cancellation code. Note that this call will create a new table entry
// with the given id, so subsequent calls to createID() will not accidentally
// re-use the generated one before the caller can.
func (p *plugin) createID() string {
	p.m.Lock()
	defer p.m.Unlock()

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
		if _, ok := p.table[id]; !ok {
			break
		}

		id = generate()
	}

	p.table[id] = alarm{}
	return id
}

// checkExpiredAlarms checks for expired alarms.
// When found, it sends the appropriate notification.
func (p *plugin) checkExpiredAlarms() {
	p.m.Lock()
	defer p.m.Unlock()

	now := time.Now()

	c := irc.Connection
	if c == nil {
		return
	}

	for id, alarm := range p.table {
		if now.Before(alarm.When) {
			continue
		}

		proto.PrivMsg(c, alarm.Target, alarm.Message,
			alarm.SenderName, time.Now().Format(TextTimeFormat))

		delete(p.table, id)
		util.WriteFile(p.file, p.table, true)
	}
}

// parseTime treats the given value as either an absolute time, or
// an offset in minutes. It returns the value which represents the
// duration between now and then.
func parseTime(v string) time.Duration {
	then, err := time.Parse(TextTimeFormat, v)

	if err == nil {
		// We expect the given time to include only the time.
		// We must set the date components manually.

		now := time.Now()
		then = time.Date(now.Year(), now.Month(), now.Day(),
			then.Hour(), then.Minute(), 0, 0, now.Location())

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
	// in minutes from the current time.
	num, err := strconv.ParseInt(v, 10, 32)
	if err == nil {
		// This can result in a negative duration, if someone specified
		// "-10" as the input. This is an error which is caught by the caller.
		return time.Duration(num) * time.Minute
	}

	return 0
}
