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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// LogFormat defines the date layout for log file names.
	LogFormat = "20060102"

	// PurgeTimeout defines the timeout after which the bot should
	// check for stale log files.
	PurgeTimeout = time.Hour * 24

	// RefreshTimeout determines how often we should check if a new
	// log file should be opened.
	RefreshTimeout = time.Minute

	// LogExpiration defines how old a log file should be, before it
	// is considered stale.
	LogExpiration = time.Hour * 24 * 7 * 2
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
	getLogFunc   func() bool
	setLogFunc   func(bool)

	rootDir      string
	logPurgeQuit chan struct{}
	quitOnce     sync.Once
}

// New returns a new admin module. The given application name and version
// numbers are used in the 'version' command.
func New(name string, major, minor int, rev string) mod.Module {
	var v version
	v.Name = name
	v.Major = major
	v.Minor = minor
	v.Revision, _ = strconv.ParseInt(rev, 10, 64)
	return &module{version: v}
}

// Load loads module resources and binds commands.
func (m *module) Load(pb irc.ProtocolBinder, prof irc.Profile) {
	m.rootDir = prof.Root()
	m.logPurgeQuit = make(chan struct{})
	m.authFunc = prof.WhitelistAdd
	m.deauthFunc = prof.WhitelistRemove
	m.authListFunc = prof.Whitelist
	m.getLogFunc = prof.Logging
	m.setLogFunc = prof.SetLogging

	pb.Bind("*", m.onAny)
	pb.Bind("PRIVMSG", m.onPrivMsg)

	m.commands = cmd.New(
		prof.CommandPrefix(),
		func(r *irc.Request) bool {
			return prof.IsWhitelisted(r.SenderMask)
		},
	)

	// This command is invoked with the bot's nickname and directs the
	// user to the full !help command.
	m.commands.Bind(prof.Nickname(), tr.SimpleHelpText, false, m.cmdSimpleHelp)

	m.commands.Bind(tr.JoinName, tr.JoinDesc, true, m.cmdJoin).
		Add(tr.JoinChannelName, tr.JoinChannelDesc, true, cmd.RegChannel).
		Add(tr.JoinPasswordName, tr.JoinPasswordDesc, false, cmd.RegAny).
		Add(tr.JoinKeyName, tr.JoinKeyDesc, false, cmd.RegAny)

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

	go m.pollLogState()
}

// Unload cleans up library resources and unbinds commands.
func (m *module) Unload(pb irc.ProtocolBinder, prof irc.Profile) {
	m.quitOnce.Do(func() {
		close(m.logPurgeQuit)

		m.commands.Clear()
		pb.Unbind("PRIVMSG", m.onPrivMsg)
		pb.Unbind("*", m.onAny)

		m.authFunc = nil
		m.deauthFunc = nil
		m.authListFunc = nil
		m.getLogFunc = nil
		m.setLogFunc = nil
	})
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
	if m.getLogFunc() {
		if strings.EqualFold(r.Type, "PING") {
			return // Skip these. It's just noise.
		}

		fields := []string{
			r.SenderName,
			r.SenderMask,
			r.Type,
			r.Target,
			r.Data,
		}

		log.Printf("> %s", strings.Join(fields, ", "))
	}
}

// cmdSimpleHelp presents the user with a short message, pointing them
// to the full !help command.
func (m *module) cmdSimpleHelp(w irc.ResponseWriter, r *cmd.Request) {
	proto.PrivMsg(w, r.Target, tr.SimpleHelpText, r.SenderName)
}

// cmdLog changes and/or reports the current logging state.
func (m *module) cmdLog(w irc.ResponseWriter, r *cmd.Request) {
	if r.Len() > 0 {
		m.setLogFunc(r.Bool(0))
	}

	if m.getLogFunc() {
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
		channel.Password = r.String(1)
	}

	if r.Len() > 2 {
		channel.Key = r.String(2)
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
	proc.Fork()
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

// pollLogState periodically checks for stale logs and makes sure we
// open a new log file, once every day.
func (m *module) pollLogState() {
	for {
		select {
		case <-m.logPurgeQuit:
			return

		case <-time.After(RefreshTimeout):
			m.refreshLog()

		case <-time.After(PurgeTimeout):
			m.doLogPurge()
		}
	}
}

// doLogPurge checks the log file directory for files which are older
// than a predefined number of days. If found, the log file in question
// is deleted. This ensures we do not keep stale logs around unnecessarily.
func (m *module) doLogPurge() {
	log.Println("[admin] purging stale log files...")

	logDir := filepath.Join(m.rootDir, "logs")
	fd, err := os.Open(logDir)
	if err != nil {
		log.Println("[admin] purge log files:", err)
		return
	}

	files, err := fd.Readdir(-1)
	fd.Close()

	if err != nil {
		log.Println("[admin] purge log files:", err)
		return
	}

	for _, file := range files {
		if time.Since(file.ModTime()) < LogExpiration {
			continue
		}

		path := filepath.Join(logDir, file.Name())
		err = os.Remove(path)
		if err != nil {
			log.Printf("[admin] deleting log file %q: %v", file.Name(), err)
		}
	}
}

// refreshLog checks if we've started a new day. In which case,
// a new log file should be opened.
func (m *module) refreshLog() {
	err := InitLog(m.rootDir)
	if err != nil {
		log.Println("[admin]", err)
	}
}

// logFile has a handle to the currently used log output.
var logFile *os.File

// InitLog  initializes a new log file, if necessary. This happens once
// every day.
//
// This call is exported because it is called on program start in main()
func InitLog(root string) error {
	// Ensure the log file directory exists.
	dir := filepath.Join(root, "logs")
	err := os.Mkdir(dir, 0700)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Determine the name of the new log file.
	timeStamp := time.Now().Format(LogFormat)
	file := fmt.Sprintf("%s.txt", timeStamp)
	file = filepath.Join(dir, file)

	// Exit if we're already using this file.
	if logFile != nil && logFile.Name() == file {
		return nil
	}

	log.Println("[admin] opening new log file:", file)

	// Create the new logfile.
	fd, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	// Set the new log output.
	log.SetOutput(fd)

	// Close the old log file and assign the new one.
	if logFile != nil {
		logFile.Close()
	}

	logFile = fd
	return nil
}
