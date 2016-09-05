// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package stats retains a listing of user host names, mapped to
// nicknames they have ever been seen using, along with some other,
// rudimentary user statistics.
//
// This is intended to make it easier to pick out trolls, trying to
// present themselves as new users. While this is by no means fool-proof,
// it keeps the majority out.
package stats

import (
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/monkeybird/autimaat/irc"
	"github.com/monkeybird/autimaat/irc/cmd"
	"github.com/monkeybird/autimaat/irc/proto"
	"github.com/monkeybird/autimaat/plugins"
	"github.com/monkeybird/autimaat/util"
)

// SaveInterval determines the time interval after which we save stats data to disk.
const SaveInterval = time.Minute * 10

func init() { plugins.Register(&plugin{}) }

type plugin struct {
	m        sync.RWMutex
	cmd      *cmd.Set
	file     string
	users    UserList
	quitOnce sync.Once
	quit     chan struct{}
}

// Load initializes the module and loads any internal resources
// which may be required.
func (p *plugin) Load(prof irc.Profile) error {
	p.m.Lock()
	defer p.m.Unlock()

	p.quit = make(chan struct{})
	p.file = filepath.Join(prof.Root(), "stats.dat")
	p.cmd = cmd.New(prof.CommandPrefix(), nil)

	p.cmd.Bind(TextWhoisName, false, p.cmdWhois).
		Add(TextNick, true, cmd.RegAny)

	p.cmd.Bind(TextFirstOn, false, p.cmdFirstOn).
		Add(TextNick, true, cmd.RegAny)

	p.cmd.Bind(TextLastOn, false, p.cmdLastOn).
		Add(TextNick, true, cmd.RegAny)

	go p.periodicSave()
	return util.ReadFile(p.file, &p.users, true)
}

// Unload cleans the module up and unloads any internal resources.
func (p *plugin) Unload(prof irc.Profile) error {
	p.quitOnce.Do(func() {
		close(p.quit)
		p.saveFile()
	})
	return nil
}

// Dispatch sends the given, incoming IRC message to the plugin for
// processing as it sees fit.
func (p *plugin) Dispatch(w irc.ResponseWriter, r *irc.Request) {
	p.cmd.Dispatch(w, r)

	p.m.Lock()
	usr := p.users.Get(r.SenderMask)
	usr.AddNickname(r.SenderName)
	p.m.Unlock()
}

// periodicSave periodically saves the stats data to disk.
func (p *plugin) periodicSave() {
	for {
		select {
		case <-p.quit:
			return

		case <-time.After(SaveInterval):
			p.saveFile()
		}
	}
}

// saveFile saes the user data to disk.
func (p *plugin) saveFile() {
	p.m.RLock()
	err := util.WriteFile(p.file, p.users, true)
	p.m.RUnlock()

	if err != nil {
		log.Println("[stats] save:", err)
	}
}

// cmdWhois presents the caller with a list of usernames known for a specific
// user or hostmask.
func (p *plugin) cmdWhois(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	p.m.RLock()
	defer p.m.RUnlock()

	usr := p.users.Find(params.String(0))
	if usr == nil {
		proto.PrivMsg(w, r.Target, TextWhoisUnknownUser, r.SenderName,
			util.Bold(params.String(0)))
		return
	}

	proto.PrivMsg(w, r.Target,
		TextWhoisDisplay,
		r.SenderName,
		util.Bold(params.String(0)),
		usr.FirstSeen.Format(TextDateFormat),
		strings.Join(usr.Nicknames, ", "),
	)
}

// cmdFirstOn tells the caller when a specific user was first seen online.
func (p *plugin) cmdFirstOn(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	p.m.RLock()
	defer p.m.RUnlock()

	usr := p.users.Find(params.String(0))
	if usr == nil {
		proto.PrivMsg(w, r.Target, TextUnknownUser, r.SenderName,
			util.Bold(params.String(0)))
		return
	}

	proto.PrivMsg(w, r.Target,
		TextFirstOnDisplay,
		r.SenderName,
		util.Bold(params.String(0)),
		usr.FirstSeen.Format(TextDateFormat),
		usr.FirstSeen.Format(TextTimeFormat),
		FormatDuration(time.Since(usr.FirstSeen)),
	)
}

// cmdLastOn tells the caller when a specific user was last seen online.
func (p *plugin) cmdLastOn(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	p.m.RLock()
	defer p.m.RUnlock()

	usr := p.users.Find(params.String(0))
	if usr == nil {
		proto.PrivMsg(w, r.Target, TextUnknownUser, r.SenderName,
			util.Bold(params.String(0)))
		return
	}

	proto.PrivMsg(w, r.Target,
		TextLastOnDisplay,
		r.SenderName,
		util.Bold(params.String(0)),
		usr.LastSeen.Format(TextDateFormat),
		usr.LastSeen.Format(TextTimeFormat),
		FormatDuration(time.Since(usr.LastSeen)),
	)
}
