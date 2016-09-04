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

	"github.com/monkeybird/autimaat/irc"
	"github.com/monkeybird/autimaat/irc/cmd"
	"github.com/monkeybird/autimaat/irc/proto"
	"github.com/monkeybird/autimaat/plugins"
	"github.com/monkeybird/autimaat/util"
)

func init() { plugins.Register(&plugin{}) }

type plugin struct {
	m     sync.RWMutex
	cmd   *cmd.Set
	file  string
	table map[string]UserList
}

// Load initializes the module and loads any internal resources
// which may be required.
func (p *plugin) Load(prof irc.Profile) error {
	p.m.Lock()
	defer p.m.Unlock()

	p.table = make(map[string]UserList)
	p.file = filepath.Join(prof.Root(), "stats.dat")
	p.cmd = cmd.New(prof.CommandPrefix(), nil)

	p.cmd.Bind(TextWhoisName, false, p.cmdWhois).
		Add(TextWhoisNick, true, cmd.RegAny)

	return util.ReadFile(p.file, &p.table, true)
}

// Unload cleans the module up and unloads any internal resources.
func (p *plugin) Unload(prof irc.Profile) error {
	return nil
}

// Dispatch sends the given, incoming IRC message to the plugin for
// processing as it sees fit.
func (p *plugin) Dispatch(w irc.ResponseWriter, r *irc.Request) {
	if r.IsPrivMsg() {
		p.cmd.Dispatch(w, r)
	}

	switch r.Type {
	case "JOIN", "PART", "QUIT", "NICK":
		p.update(r)
	}
}

// update is called whenever a user joins a channel the bot is in, or the
// user changed their nickname. It is used to update the user database.
func (p *plugin) update(r *irc.Request) {
	if !r.FromChannel() {
		return
	}

	p.m.Lock()
	defer p.m.Unlock()

	// Find the right user list and update the user information,
	// where applicable.
	users, newChannel := p.table[r.Target]
	usr, newUser := users.Get(r.SenderMask)
	newNickname := usr.AddNickname(r.SenderName)

	// If no new data has been added, just exit.
	if !newChannel && !newUser && !newNickname {
		return
	}

	// Otherwise, update the table and save it to disk.
	p.table[r.Target] = users

	err := util.WriteFile(p.file, p.table, true)
	if err != nil {
		log.Println("[whois]", err)
	}
}

// cmdWhois presents the caller with a list of usernames known for a specific
// user or hostmask.
func (p *plugin) cmdWhois(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	if !r.FromChannel() {
		proto.PrivMsg(w, r.SenderName, TextWhoisNotInChannel)
		return
	}

	p.m.RLock()
	defer p.m.RUnlock()

	users := p.table[r.Target]
	usr := users.Find(params.String(0))

	if usr == nil {
		proto.PrivMsg(w, r.Target, TextWhoisUnknownUser, r.SenderName,
			util.Bold(params.String(0)))
		return
	}

	proto.PrivMsg(w, r.Target,
		TextWhoisDisplay,
		r.SenderName,
		util.Bold(params.String(0)),
		usr.FirstSeen.Format(TextWhoisDateFormat),
		strings.Join(usr.Nicknames, ", "),
	)
}
