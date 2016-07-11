// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/irc/proto"
	"monkeybird/mod"
	"monkeybird/mod/admin"
	"monkeybird/mod/misc"
	"monkeybird/mod/snooze"
	"monkeybird/mod/stats"
	"monkeybird/mod/url"
	"monkeybird/mod/weather"
	"monkeybird/proc"
	"monkeybird/tr"
	"sync"
)

// Bot defines an implementation of a simple bot, connecting to one network.
// It combines functionality from various irc packages to create the bare
// minimum of something useful. It takes care of connection housekeeping.
//
// Once running, it waits for OS signals to arrive which can either fork
// the process or shut it down.
type Bot struct {
	profile  irc.Profile     // Bot profile.
	client   *irc.Client     // Underlying protocol bindings and network stream.
	bindings irc.BindingList // List of protocol bindings.
	commands *cmd.Set        // Custom commands.
	quit     sync.Once       // Ensures we close/quit only once.
	modules  []mod.Module    // List of connected modules.
}

// New creates a new bot instance for the given profile.
func New(profile irc.Profile) *Bot {
	var b Bot
	b.profile = profile
	b.client = irc.NewClient(profile, b.handleMessage)
	b.commands = cmd.New(profile.CommandPrefix(), nil)

	// Bind some protocol handlers.
	b.bindings.Bind("375", b.loginJoinChannels)
	b.bindings.Bind("422", b.loginJoinChannels)
	b.bindings.Bind("433", b.onNickInUse)
	b.bindings.Bind("ERROR", b.onError)
	b.bindings.Bind("PING", b.onPing)
	b.bindings.Bind("PRIVMSG", b.onPrivMsg)

	// Bind help command. This will print some leading help information and
	// then forward the call to sub packages, which may define their own
	// commands.
	b.commands.Bind(tr.HelpName, tr.HelpDesc, false, b.cmdHelp).
		Add(tr.HelpCommandName, tr.HelpCommandDesc, false, cmd.RegAny)

	// Create custom module instances.
	b.modules = []mod.Module{
		admin.New(
			AppName,
			AppVersionMajor,
			AppVersionMinor,
			AppVersionRevision,
		),
		weather.New(),
		url.New(),
		stats.New(),
		misc.New(),
		snooze.New(&b),
	}

	return &b
}

// Close closes the connection and cleans up resources.
func (b *Bot) Close() error {
	b.quit.Do(func() {
		b.client.Close()
		b.client = nil

		b.bindings.Clear()
		b.commands.Clear()

		for _, m := range b.modules {
			log.Printf("[bot] Unloading %T", m)
			m.Unload(&b.bindings, b.profile)
		}

		b.modules = nil
	})
	return nil
}

// Write helps implement irc.ResponseWriter
func (b *Bot) Write(p []byte) (int, error) {
	return b.client.Write(p)
}

// Bind helps implement irc.ProtocolBinder
func (b *Bot) Bind(mtype string, handler irc.RequestFunc) {
	b.bindings.Bind(mtype, handler)
}

// Unbind helps implement irc.ProtocolBinder
func (b *Bot) Unbind(mtype string, handler irc.RequestFunc) {
	b.bindings.Unbind(mtype, handler)
}

// Clear helps implement irc.ProtocolBinder
func (b *Bot) Clear() {
	b.bindings.Clear()
}

// Run opens the bot's client connection and begins the data loop.
//
// The connection is either a new one, or inherited from a previous session.
// This function will not return for as long as the bot is running.
func (b *Bot) Run() error {
	log.Printf("[bot] Running %s version %d.%d.%s",
		AppName, AppVersionMajor, AppVersionMinor, AppVersionRevision)
	defer log.Println("[bot] Shutting down")

	// Open connection.
	err := b.open()
	if err != nil {
		return err
	}

	// Load all modules.
	for _, m := range b.modules {
		log.Printf("[bot] Loading %T", m)
		m.Load(&b.bindings, b.profile)
	}

	// Spin up the client's data loop in a separate goroutine.
	// This is not strictly necessary, but allows us to monitor
	// the OS for kill signals, so we may exit cleanly.
	go func() {
		log.Println("[bot] Entering data loop...")

		err := b.client.Run()
		if err != nil {
			log.Println(err)
		}

		// Break out of the Wait() call below.
		proc.Kill()
	}()

	fd, _ := b.client.File()

	proc.Wait(b.profile.ForkArgs(), fd)
	return b.Close()
}

// open creates the bot's client connection. This is either a new connection
// or one inherited from a parent process.
func (b *Bot) open() error {
	var config *tls.Config

	// Create TLS configuration, if applicable.
	if len(b.profile.TLSCert()) > 0 && len(b.profile.TLSKey()) > 0 {
		cert, err := tls.LoadX509KeyPair(b.profile.TLSCert(), b.profile.TLSKey())
		if err != nil {
			return err
		}

		config = &tls.Config{
			Certificates:             []tls.Certificate{cert},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       false,
		}

		// Should we replace the client's root CA pool?
		if len(b.profile.CAPemData()) > 0 {
			config.RootCAs = x509.NewCertPool()

			data, err := ioutil.ReadFile(b.profile.CAPemData())
			if err != nil {
				return err
			}

			if !config.RootCAs.AppendCertsFromPEM(data) {
				return fmt.Errorf("AppendCertsFromPEM: failed to add certificates in %s",
					b.profile.CAPemData())
			}
		}
	}

	files := proc.InheritedFiles()

	// Are we a fork? Then we should inherit an existing connection.
	if len(files) > 0 {
		log.Println("[bot] Inherit connection to:", b.profile.Address())

		err := b.client.OpenFd(files[0], config)
		if err != nil {
			return err
		}

		// We're done inheriting. Kill the parent process.
		proc.KillParent()
		return nil
	}

	log.Println("[bot] Opening new connection to:", b.profile.Address())

	// Fresh session - create a new connection.
	err := b.client.Open(b.profile.Address(), config)
	if err != nil {
		return err
	}

	// Perform initial handshake.
	if len(b.profile.ConnectionPassword()) > 0 {
		proto.Pass(b.client, b.profile.ConnectionPassword())
	}

	proto.User(b.client, b.profile.Nickname(), "8", b.profile.Nickname())
	proto.Nick(b.client, b.profile.Nickname(), b.profile.NickservPassword())
	return nil
}

// handleMessage handles incoming server messages.
func (b *Bot) handleMessage(r *irc.Request) {
	// Check for handlers specific to this message type.
	for _, handler := range b.bindings.Find(r.Type) {
		go handler(b, r)
	}

	// Check for catch-all handlers.
	for _, handler := range b.bindings.Find("*") {
		go handler(b, r)
	}
}

// onPrivMsg is called on every incoming PRIVMSG and its purpose is to
// forward the message to the command parser/dispatcher.
func (b *Bot) onPrivMsg(w irc.ResponseWriter, r *irc.Request) {
	b.commands.Dispatch(w, r)
}

// onNickInUse signals that our nick is in use. If we can regain it, do so.
// Otherwise, change ours.
func (b *Bot) onNickInUse(w irc.ResponseWriter, r *irc.Request) {
	if len(b.profile.NickservPassword()) > 0 {
		log.Println("[bot] Nick in use: trying to recover")
		proto.Recover(w, b.profile.Nickname(), b.profile.NickservPassword())
		return
	}

	b.profile.SetNickname(b.profile.Nickname() + "_")

	log.Println("[bot] Nick in use: changing nick to:", b.profile.Nickname())
	proto.Nick(b.client, b.profile.Nickname())
}

// onPing response to PING requests in order to keep the connection alive.
func (b *Bot) onPing(w irc.ResponseWriter, r *irc.Request) {
	proto.Pong(w, r.Data)
}

// onError response to network error messages. It closes the connection
// after logging the error.
func (b *Bot) onError(w irc.ResponseWriter, r *irc.Request) {
	log.Println("[bot] Network error:", r.Data)
	b.Close()
}

// loginJoinChannels is called to complete the login sequence.
// It joins channels defined in the profile and is triggered when we receive
// either the STARTMOTD or NOMOTD messages.
func (b *Bot) loginJoinChannels(w irc.ResponseWriter, r *irc.Request) {
	proto.Join(w, b.profile.Channels()...)
}

// cmdHelp handles a call to the help command. It optionally prints a help
// header and then a listing of all supported commands, or help for a single,
// specific command.
func (b *Bot) cmdHelp(w irc.ResponseWriter, r *cmd.Request) {
	if r.Len() == 0 {
		proto.PrivMsg(w, r.SenderName, tr.CommandsIntro1)
		proto.PrivMsg(w, r.SenderName, tr.CommandsIntro2)
	}

	b.commands.HelpHandler(w, r)

	// Forward help call to all modules.
	for _, m := range b.modules {
		m.Help(w, r)
	}
}
