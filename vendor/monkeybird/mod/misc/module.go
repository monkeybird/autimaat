// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package misc defines a bunch of random, silly things
// in command form.
package misc

import (
	"math/rand"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/irc/proto"
	"monkeybird/mod"
	"monkeybird/tr"
	"strconv"
	"time"
)

// eightBallAnswers defines the list of possible 8ball answers.
var eightBallAnswers = []string{
	tr.Eightball1,
	tr.Eightball2,
	tr.Eightball3,
	tr.Eightball4,
	tr.Eightball5,
	tr.Eightball6,
	tr.Eightball7,
	tr.Eightball8,
	tr.Eightball9,
	tr.Eightball10,
	tr.Eightball11,
	tr.Eightball12,
	tr.Eightball13,
	tr.Eightball14,
	tr.Eightball15,
	tr.Eightball16,
	tr.Eightball17,
	tr.Eightball18,
	tr.Eightball19,
	tr.Eightball20,
}

type module struct {
	commands *cmd.Set
}

// New returns a new eightball module.
func New() mod.Module {
	return &module{}
}

// Load loads module resources and binds commands.
func (m *module) Load(pb irc.ProtocolBinder, prof irc.Profile) {
	pb.Bind("PRIVMSG", m.onPrivMsg)

	m.commands = cmd.New(prof.CommandPrefix(), nil)
	m.commands.Bind(tr.EightballName, tr.EightballDesc, false, m.cmd8Ball).
		Add(tr.EightballQuestionName, tr.EightballQuestionDesc, true, cmd.RegAny)

	m.commands.Bind(tr.SnoozeName, tr.SnoozeDesc, false, m.cmdSnooze).
		Add(tr.SnoozeTimeName, tr.SnoozeTimeDesc, true, cmd.RegAny).
		Add(tr.SnoozeMessageName, tr.SnoozeMessageDesc, false, cmd.RegAny)
}

// Unload cleans up library resources and unbinds commands.
func (m *module) Unload(pb irc.ProtocolBinder, prof irc.Profile) {
	m.commands.Clear()
	pb.Unbind("PRIVMSG", m.onPrivMsg)
}

func (m *module) Help(w irc.ResponseWriter, r *cmd.Request) {
	m.commands.HelpHandler(w, r)
}

func (m *module) onPrivMsg(w irc.ResponseWriter, r *irc.Request) {
	m.commands.Dispatch(w, r)
}

// cmdSnooze lets the user schedule an alarm at a specific time, or after
// N minutes. When triggered, the alarm simply sends a user-defined message
// to the channel.
func (m *module) cmdSnooze(w irc.ResponseWriter, r *cmd.Request) {
	when := parseTime(r.String(0))
	if when <= 0 {
		proto.PrivMsg(w, r.Target, tr.SnoozeInvalidTime, r.SenderName, r.String(0))
		return
	}

	msg := r.Remainder(2)
	if len(msg) == 0 {
		msg = tr.SnoozeDefaultMessage
	} else {
		msg = "%s, het is %s: " + msg
	}

	t := time.Now().Add(when)

	go func() {
		<-time.After(when)
		proto.PrivMsg(w, r.Target, msg, r.SenderName, t.Format(tr.TimeFormatShort))
	}()

	proto.PrivMsg(w, r.Target, tr.SnoozeAlarmSet, r.SenderName,
		t.Format(tr.DateFormat), t.Format(tr.TimeFormat))
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

// cmd8Ball asks the 8ball a question and presents the answer.
func (m *module) cmd8Ball(w irc.ResponseWriter, r *cmd.Request) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := rnd.Intn(len(eightBallAnswers))
	proto.PrivMsg(w, r.Target, eightBallAnswers[idx], r.SenderName)
}
