// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package admin defines administrative bot commands.
package admin

import (
	"strconv"
	"strings"
	"time"

	"github.com/monkeybird/autimaat/app"
	"github.com/monkeybird/autimaat/irc"
	"github.com/monkeybird/autimaat/irc/cmd"
	"github.com/monkeybird/autimaat/irc/proto"
	"github.com/monkeybird/autimaat/plugins"
	"github.com/monkeybird/autimaat/proc"
)

// lastRestart defines the timestamp at which the bot was last restarted.
var lastRestart = time.Now()

func init() { plugins.Register(&plugin{}) }

type plugin struct {
	cmd *cmd.Set

	// This will store the bot's profile, but only as a subset of
	// the full interface. We only need access to these parts.
	profile interface {
		WhitelistAdd(string)
		WhitelistRemove(string)
		Whitelist() []string
		Logging() bool
		SetLogging(bool)
	}
}

// Load initializes the module and loads any internal resources
// which may be required.
func (p *plugin) Load(prof irc.Profile) error {
	p.profile = prof
	p.cmd = cmd.New(
		prof.CommandPrefix(),
		prof.IsWhitelisted,
	)

	// Two aliases for the same command. Can be invoked through
	// !help or !<bot nickname>
	p.cmd.Bind(TextHelpName, false, p.cmdHelp)
	p.cmd.Bind(prof.Nickname(), false, p.cmdHelp)

	p.cmd.Bind(TextNickName, true, p.cmdNick).
		Add(TextNickNickName, true, cmd.RegAny).
		Add(TextNickPassName, false, cmd.RegAny)

	p.cmd.Bind(TextJoinName, true, p.cmdJoin).
		Add(TextJoinChannelName, true, cmd.RegChannel).
		Add(TextJoinPasswordName, false, cmd.RegAny).
		Add(TextJoinKeyName, false, cmd.RegAny)

	p.cmd.Bind(TextPartName, true, p.cmdPart).
		Add(TextPartChannelName, true, cmd.RegChannel)

	p.cmd.Bind(TextAuthListName, true, p.cmdAuthList)

	p.cmd.Bind(TextAuthorizeName, true, p.cmdAuthorize).
		Add(TextAuthorizeMaskName, true, cmd.RegAny)

	p.cmd.Bind(TextDeauthorizeName, true, p.cmdDeauthorize).
		Add(TextDeauthorizeMaskName, true, cmd.RegAny)

	p.cmd.Bind(TextLogName, true, p.cmdLog).
		Add(TextLogValueName, false, cmd.RegBool)

	p.cmd.Bind(TextReloadName, true, p.cmdReload)
	p.cmd.Bind(TextVersionName, false, p.cmdVersion)

	return nil
}

// Unload cleans the module up and unloads any internal resources.
func (p *plugin) Unload(prof irc.Profile) error {
	p.profile = nil
	p.cmd = nil
	return nil
}

// Dispatch sends the given, incoming IRC message to the plugin for
// processing as it sees fit.
func (p *plugin) Dispatch(w irc.ResponseWriter, r *irc.Request) {
	p.cmd.Dispatch(w, r)
}

// cmdHelp presents the user with a short message, pointing them to
// a resource where the full bot help can be viewed.
func (p *plugin) cmdHelp(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	proto.PrivMsg(w, r.Target, TextHelpDisplay, r.SenderName)
}

// cmdNick allows the bot to change its name.
func (p *plugin) cmdNick(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	if params.Len() > 1 {
		proto.Nick(w, params.String(0), params.String(1))
	} else {
		proto.Nick(w, params.String(0))
	}
}

// cmdJoin makes the bot join a new channel.
func (p *plugin) cmdJoin(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	var channel irc.Channel
	channel.Name = params.String(0)

	if params.Len() > 1 {
		channel.Password = params.String(1)
	}

	if params.Len() > 2 {
		channel.Key = params.String(2)
	}

	proto.Join(w, channel)
}

// cmdPart makes the bot leave a given channel.
func (p *plugin) cmdPart(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	proto.Part(w, irc.Channel{
		Name: params.String(0),
	})
}

// cmdAuthList lists all whitelisted users.
func (p *plugin) cmdAuthList(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	list := p.profile.Whitelist()
	out := strings.Join(list, ", ")
	proto.PrivMsg(w, r.SenderName, TextAuthListDisplay, out)
}

// cmdAuthorize adds a new whitelisted user.
func (p *plugin) cmdAuthorize(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	p.profile.WhitelistAdd(params.String(0))
	proto.PrivMsg(w, r.SenderName, TextAuthorizeDisplay, params.String(0))
}

// cmdDeauthorize removes a user from the whitelist.
func (p *plugin) cmdDeauthorize(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	p.profile.WhitelistRemove(params.String(0))
	proto.PrivMsg(w, r.SenderName, TextDeauthorizeDisplay, params.String(0))
}

// cmdLog changes and/or reports the current logging state.
func (p *plugin) cmdLog(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	if params.Len() > 0 {
		p.profile.SetLogging(params.Bool(0))
	}

	if p.profile.Logging() {
		proto.PrivMsg(w, r.SenderName, TextLogEnabled)
	} else {
		proto.PrivMsg(w, r.SenderName, TextLogDisabled)
	}
}

// cmdReload forces the bot to fork itself.
func (p *plugin) cmdReload(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	proc.Fork()
}

// cmdVersion prints version information.
func (p *plugin) cmdVersion(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	rev, _ := strconv.ParseInt(app.VersionRevision, 10, 64)
	stamp := time.Unix(rev, 0)
	utime := time.Since(lastRestart).Hours()

	proto.PrivMsg(
		w, r.Target,
		TextVersionDisplay,
		r.SenderName,
		irc.Bold(app.Name),
		irc.Bold("%d.%d", app.VersionMajor, app.VersionMinor),
		stamp.Format(TextDateFormat),
		stamp.Format(TextTimeFormat),
		irc.Bold("%.3f", utime),
	)
}
