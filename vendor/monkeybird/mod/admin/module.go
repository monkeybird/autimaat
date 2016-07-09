// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package admin defines basic administrative bot commands.
package admin

import (
	"fmt"
	"log"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/irc/proto"
	"monkeybird/mod"
	"monkeybird/proc"
	"monkeybird/text"
	"monkeybird/tr"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// version defines version info to be presented by the version command.
type version struct {
	Name     string
	Major    int
	Minor    int
	Revision int64
}

func (v *version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Revision)
}

type module struct {
	version      version
	commands     *cmd.Set
	authFunc     func(string)
	deauthFunc   func(string)
	authListFunc func() []string
	logging      uint32
}

// New returns a new admin module. The given application name and version
// numbers are used in the 'version' command.
func New(name string, major, minor int, rev string) mod.Module {
	var v version
	v.Name = name
	v.Major = major
	v.Minor = minor
	v.Revision, _ = strconv.ParseInt(rev, 10, 64)

	return &module{
		version: v,
		logging: 0,
	}
}

// Load loads module resources and binds commands.
func (m *module) Load(pb irc.ProtocolBinder, prof irc.Profile) {
	m.authFunc = prof.WhitelistAdd
	m.deauthFunc = prof.WhitelistRemove
	m.authListFunc = prof.Whitelist

	pb.Bind("*", m.onAny)
	pb.Bind("PRIVMSG", m.onPrivMsg)

	m.commands = cmd.New(
		prof.CommandPrefix(),
		func(r *irc.Request) bool {
			return prof.IsWhitelisted(r.SenderMask)
		},
	)

	m.commands.Bind(tr.JoinName, tr.JoinDesc, true, m.cmdJoin).
		Add(tr.JoinChannelName, tr.JoinChannelDesc, true, cmd.RegChannel).
		Add(tr.JoinKeyName, tr.JoinKeyDesc, false, cmd.RegAny).
		Add(tr.JoinPasswordName, tr.JoinPasswordDesc, false, cmd.RegAny)

	m.commands.Bind(tr.PartName, tr.PartDesc, true, m.cmdPart).
		Add(tr.PartChannelName, tr.PartChannelDesc, true, cmd.RegChannel)

	m.commands.Bind(tr.AuthListName, tr.AuthListDesc, true, m.cmdAuthList)

	m.commands.Bind(tr.AuthorizeName, tr.AuthorizeDesc, true, m.cmdAuthorize).
		Add(tr.AuthorizeMaskName, tr.AuthorizeMaskDesc, true, cmd.RegAny)

	m.commands.Bind(tr.DeauthorizeName, tr.DeauthorizeDesc, true, m.cmdDeauthorize).
		Add(tr.DeauthorizeMaskName, tr.DeauthorizeMaskDesc, true, cmd.RegAny)

	m.commands.Bind(tr.LogName, tr.LogDesc, true, m.cmdLog).
		Add(tr.LogValueName, tr.LogValueDesc, false, cmd.RegBool)

	m.commands.Bind(tr.ReloadName, tr.ReloadDesc, true, m.cmdReload)
	m.commands.Bind(tr.VersionName, tr.VersionDesc, false, m.cmdVersion)
}

// Unload cleans up library resources and unbinds commands.
func (m *module) Unload(pb irc.ProtocolBinder, prof irc.Profile) {
	m.commands.Clear()
	pb.Unbind("PRIVMSG", m.onPrivMsg)
	pb.Unbind("*", m.onAny)

	m.authFunc = nil
	m.deauthFunc = nil
	m.authListFunc = nil
}

func (m *module) Help(w irc.ResponseWriter, r *cmd.Request) {
	m.commands.HelpHandler(w, r)
}

// onPrivMsg ensures custom commands are executed.
func (m *module) onPrivMsg(w irc.ResponseWriter, r *irc.Request) {
	m.commands.Dispatch(w, r)
}

// onAny is called on /any/ incoming message type and optionally logs
// the incoming data. Whether to log or not can be toggled through a
// bot command.
func (m *module) onAny(w irc.ResponseWriter, r *irc.Request) {
	if atomic.LoadUint32(&m.logging) == 1 {
		log.Printf(
			"> Type: %s, SenderName: %s, SenderMask: %s, Target: %s, Data: %q",
			r.Type,
			r.SenderName,
			r.SenderMask,
			r.Target,
			r.Data,
		)
	}
}

// cmdLog changes and/or reports the current logging state.
func (m *module) cmdLog(w irc.ResponseWriter, r *cmd.Request) {
	if r.Len() > 0 {
		if r.Bool(0) {
			atomic.StoreUint32(&m.logging, 1)
		} else {
			atomic.StoreUint32(&m.logging, 0)
		}
	}

	if atomic.LoadUint32(&m.logging) == 1 {
		proto.PrivMsg(w, r.SenderName, tr.LogEnabled)
	} else {
		proto.PrivMsg(w, r.SenderName, tr.LogDisabled)
	}
}

// cmdJoin makes the bot join a new channel.
func (m *module) cmdJoin(w irc.ResponseWriter, r *cmd.Request) {
	var channel irc.Channel
	channel.Name = r.String(0)

	if r.Len() > 1 {
		channel.Key = r.String(1)
	}

	if r.Len() > 2 {
		channel.Password = r.String(2)
	}

	proto.Join(w, channel)
}

// cmdPart makes the bot leave a given channel.
func (m *module) cmdPart(w irc.ResponseWriter, r *cmd.Request) {
	proto.Part(w, irc.Channel{
		Name: r.String(0),
	})
}

// cmdReload forces the bot to fork itself.
func (m *module) cmdReload(w irc.ResponseWriter, r *cmd.Request) {
	proc.RequestFork()
}

// cmdAuthList lists all whitelisted users.
func (m *module) cmdAuthList(w irc.ResponseWriter, r *cmd.Request) {
	list := m.authListFunc()
	proto.PrivMsg(w, r.SenderName, tr.AuthListDisplayText,
		strings.Join(list, ", "))
}

// cmdAuthorize adds a new whitelisted user.
func (m *module) cmdAuthorize(w irc.ResponseWriter, r *cmd.Request) {
	m.authFunc(r.String(0))
	proto.PrivMsg(w, r.SenderName, tr.AuthorizeDisplayText, r.String(0))
}

// cmdDeauthorize removes a user from the whitelist.
func (m *module) cmdDeauthorize(w irc.ResponseWriter, r *cmd.Request) {
	m.deauthFunc(r.String(0))
	proto.PrivMsg(w, r.SenderName, tr.DeauthorizeDisplayText, r.String(0))
}

// cmdVersion prints version information.
func (m *module) cmdVersion(w irc.ResponseWriter, r *cmd.Request) {
	stamp := time.Unix(m.version.Revision, 0)

	proto.PrivMsg(
		w, r.Target,
		tr.VersionDisplayText,
		r.SenderName,
		text.Bold(m.version.Name),
		text.Bold(m.version.String()),
		stamp.Format(tr.DateFormat),
		stamp.Format(tr.TimeFormat),
	)
}
