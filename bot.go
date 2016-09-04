// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/monkeybird/autimaat/app"
	"github.com/monkeybird/autimaat/irc"
	"github.com/monkeybird/autimaat/irc/proto"
	"github.com/monkeybird/autimaat/plugins"
	"github.com/monkeybird/autimaat/proc"

	_ "github.com/monkeybird/autimaat/plugins/admin"
	_ "github.com/monkeybird/autimaat/plugins/dictionary"
	_ "github.com/monkeybird/autimaat/plugins/url"
)

// Bot defines state for a single IRC bot.
type Bot struct {
	profile irc.Profile
	client  *Client
}

// Run creates a new connection to the server and begins processing
// incoming messages and OS signals. This call will not return for as long
// as the connection is active.
func Run(p irc.Profile) error {
	log.Printf("[bot] Running %s version %d.%d.%s",
		app.Name, app.VersionMajor, app.VersionMinor, app.VersionRevision)
	defer log.Println("[bot] Shutting down")

	var bot Bot
	bot.profile = p
	bot.client = NewClient(bot.payloadHandler)

	// Initialize plugins.
	plugins.Load(p)
	defer plugins.Unload(p)

	// Open the connection and spin up the client's read loop in a
	// separate goroutine.
	return bot.run()
}

// run opens a new connection, or inherits an existing one and then begins
// the client's message poll routine.
func (b *Bot) run() error {
	// Initialize the connection.
	err := b.open()
	if err != nil {
		return err
	}

	// Spin up the connection's read loop.
	go func() {
		log.Println("[bot] Entering data loop...")

		err := b.client.Run()
		if err != nil {
			log.Println(err)
		}

		// Break out of the Wait() call below.
		proc.Kill()
	}()

	// Wait for external signals. Either to cleanly shut the bot down,
	// or to initiate the forking process.
	fd, _ := b.client.File()
	proc.Wait(b.profile.ForkArgs(), fd)

	return b.client.Close()
}

// payloadHandler handles incoming server messages.
func (b *Bot) payloadHandler(payload []byte) {
	var r irc.Request

	// Try to parse the payload into a request.
	if !parseRequest(&r, payload) {
		return
	}

	// If Target points to the bot's own name, then this message came from
	// a user as a PM. Change the Target to the sender's name, so any replies
	// we create, end up at the right destination. In any other case, the
	// target is set to the channel name from whence the message came.
	if b.profile.IsNick(r.Target) {
		r.Target = r.SenderName
	}

	// Run the appropriate handler for housekeeping.
	switch r.Type {
	case "ERROR":
		log.Println("[bot] Network error:", r.Data)
		return

	case "PING":
		proto.Pong(b.client, r.Data)
		return
	}

	// Notify plugins of message.
	plugins.Dispatch(b.client, &r)

	// Log request if applicable.
	//
	// Don't log PING requests. They only use up disk space and cause
	// noise in the log.
	if b.profile.Logging() {
		log.Println("[>]", r.String())
	}
}

// open either establishes a new connection or inherits an existing one
// from a parent process.
func (b *Bot) open() error {
	var config *tls.Config

	p := b.profile

	// Create TLS configuration, if applicable.
	if len(p.TLSCert()) > 0 && len(p.TLSKey()) > 0 {
		cert, err := tls.LoadX509KeyPair(p.TLSCert(), p.TLSKey())
		if err != nil {
			return err
		}

		config = &tls.Config{
			Certificates:             []tls.Certificate{cert},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       false,
		}

		// Should we replace the client's root CA pool?
		if len(p.CAPemData()) > 0 {
			config.RootCAs = x509.NewCertPool()

			data, err := ioutil.ReadFile(p.CAPemData())
			if err != nil {
				return err
			}

			if !config.RootCAs.AppendCertsFromPEM(data) {
				return fmt.Errorf("AppendCertsFromPEM: failed to add certificates in %s",
					p.CAPemData())
			}
		}
	}

	files := proc.InheritedFiles()

	// Are we a fork? Then we should inherit an existing connection.
	if len(files) > 0 {
		log.Println("[bot] Inherit connection to:", p.Address())

		err := b.client.OpenFd(files[0], config)
		if err != nil {
			return err
		}

		// We're done inheriting. Kill the parent process.
		proc.KillParent()
		return nil
	}

	log.Println("[bot] Opening new connection to:", p.Address())

	// Fresh session - create a new connection.
	err := b.client.Open(p.Address(), config)
	if err != nil {
		return err
	}

	// Perform initial handshake.
	proto.Pass(b.client, p.ConnectionPassword())
	proto.User(b.client, p.Nickname(), "8", p.Nickname())
	proto.Nick(b.client, p.Nickname(), p.NickservPassword())

	// We should fork ourselves at least once at this point.
	// This is done to provide behaviour like old fashioned
	// unix daemons. Systemd will be expecting this, as per
	// the service file for this program.
	//
	// TODO: find a less flaky way to do this.
	go time.AfterFunc(3*time.Second, proc.Fork)

	return nil
}
