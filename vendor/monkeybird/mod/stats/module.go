// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package stats maintains channel and user statistics.
package stats

import (
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/mod"
	"monkeybird/tr"
	"path/filepath"
	"sync"
	"time"
)

// saveInterval defines the interval after which the stats data should
// be written to disk.
const saveInterval = time.Minute * 5

type module struct {
	commands *cmd.Set
	file     string
	stats    Stats
	quitOnce sync.Once
	quit     chan struct{}
}

// New returns a new module.
func New() mod.Module { return &module{} }

// Load loads module resources and binds commands.
func (m *module) Load(pb irc.ProtocolBinder, prof irc.Profile) {
	pb.Bind("PRIVMSG", m.onPrivMsg)
	pb.Bind("JOIN", m.doUpdate)
	pb.Bind("PART", m.doUpdate)

	m.quit = make(chan struct{})
	m.file = filepath.Join(prof.Root(), "stats.dat")
	m.commands = cmd.New(
		prof.CommandPrefix(),
		func(r *irc.Request) bool {
			return prof.IsWhitelisted(r.SenderMask)
		},
	)

	m.commands.Bind(tr.FirstOnName, tr.FirstOnDesc, false, m.cmdFirstOn).
		Add(tr.FirstOnUserName, tr.FirstOnUserDesc, false, cmd.RegAny)

	m.commands.Bind(tr.LastOnName, tr.LastOnDesc, false, m.cmdLastOn).
		Add(tr.LastOnUserName, tr.LastOnUserDesc, false, cmd.RegAny)

	m.stats.Load(m.file)
	go m.periodicSave()
}

// Unload cleans up library resources and unbinds commands.
func (m *module) Unload(pb irc.ProtocolBinder, prof irc.Profile) {
	m.quitOnce.Do(func() {
		close(m.quit)

		m.stats.Save(m.file)
		m.commands.Clear()

		pb.Unbind("PRIVMSG", m.onPrivMsg)
		pb.Unbind("JOIN", m.doUpdate)
		pb.Unbind("PART", m.doUpdate)
	})
}

func (m *module) Help(w irc.ResponseWriter, r *cmd.Request) {
	m.commands.HelpHandler(w, r)
}

func (m *module) onPrivMsg(w irc.ResponseWriter, r *irc.Request) {
	m.commands.Dispatch(w, r)
	m.doUpdate(w, r)
}

// doUpdate is called on every PRIVMSG, PART and JOIN message.
// It updates the user/channel stats for the user sending the request.
func (m *module) doUpdate(w irc.ResponseWriter, r *irc.Request) {
	m.stats.Update(w, r)
}

// cmdFirstOn finds out when a specific user was first seen in
// the channel from whence this command was issued.
func (m *module) cmdFirstOn(w irc.ResponseWriter, r *cmd.Request) {
	m.stats.FirstOn(w, r)
}

// cmdLastOn finds out when a specific user was last seen in
// the channel from whence this command was issued.
func (m *module) cmdLastOn(w irc.ResponseWriter, r *cmd.Request) {
	m.stats.LastOn(w, r)
}

// periodicSave periodically saves the stats data to disk.
func (m *module) periodicSave() {
	for {
		select {
		case <-m.quit:
			return

		case <-time.After(saveInterval):
			m.stats.Save(m.file)
		}
	}
}
