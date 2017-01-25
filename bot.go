// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/monkeybird/autimaat/app"
	"github.com/monkeybird/autimaat/app/logger"
	"github.com/monkeybird/autimaat/irc"
	"github.com/monkeybird/autimaat/irc/proto"
	"github.com/monkeybird/autimaat/plugins"

	_ "github.com/monkeybird/autimaat/plugins/action"
	_ "github.com/monkeybird/autimaat/plugins/admin"
	_ "github.com/monkeybird/autimaat/plugins/alarm"
	_ "github.com/monkeybird/autimaat/plugins/dictionary"
	_ "github.com/monkeybird/autimaat/plugins/stats"
	_ "github.com/monkeybird/autimaat/plugins/url"
	_ "github.com/monkeybird/autimaat/plugins/weather"
)

// connectionCount defines the number of connections passed into a forked
// process. Currently there is only 1 connection per bot implemented
// (N=1).
var connectionCount uint

// shuttingDown is true if and only if the bot is in the process of
// gracefully closing down
var shuttingDown bool = false

func init() {
	flag.UintVar(&connectionCount, "fork", 0, "Number of inherited file descriptors")
}

// Bot defines state for a single IRC bot.
type Bot struct {
	profile irc.Profile
	client  *Client
}

// Run creates a new connection to the server and begins processing
// incoming messages and OS signals. This call will not return for as long
// as the connection is active.
func Run(p irc.Profile) error {
	// Initialize the log and ensure it is properly stopped when we are done.
	logger.Init("logs")
	defer logger.Shutdown()

	log.Printf("[bot] Running %s version %d.%d.%s",
		app.Name, app.VersionMajor, app.VersionMinor, app.VersionRevision)
	defer log.Println("[bot] Shutting down")

	// Initialize plugins.
	plugins.Load(p)
	defer plugins.Unload(p)

	// Create te bot, open the connection and spin up the client's read loop
	// in a separate goroutine.
	var bot Bot
	bot.profile = p
	bot.client = NewClient(bot.payloadHandler)
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

		// err will always be non-nil here
		if e, ok := err.(*net.OpError); ok {
			if e.Err.Error() == "use of closed network connection" {
				// This can be the error value if the bot is in the
				// process of shutting down gracefully, the connection
				// is closed, and a pending read or write was
				// unblocked by that.  Just let the shutting down of
				// the bot continue and ignore the error.
				if shuttingDown {
					log.Printf("[bot] ignoring  '%+v'\n", e.Err)
					return
				}
			}
		}

		// Any other error is fatal, so a supervisor like systemd can
		// try to restart the bot.
		log.Fatal("[bot] exit 1: ", err)

	}()

	// Wait for external signals. Either to cleanly shut the bot down,
	// or to initiate the forking process.
	wait(b)
	shuttingDown = true
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

	files := inheritedFiles()

	// Are we a fork? Then we should inherit an existing connection.
	if len(files) > 0 {
		log.Println("[bot] Inherit connection to:", p.Address())

		err := b.client.OpenFd(files[0], config)
		if err != nil {
			return err
		}

		// We're done inheriting. Have the parent process break out of
		// its wait() call by sending SIGINT to it.
		syscall.Kill(os.Getppid(), syscall.SIGINT)
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
	return nil
}

// wait polls for OS signals to either kill or fork this process.
// The signals it waits for are: SIGINT, SIGTERM and SIGUSR1.
// The latter one being responsible for forking this process. The others
// are there so we may cleanly exit this process.
func wait(b *Bot) {
	signals := make(chan os.Signal, 1)
	signal.Notify(
		signals,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGUSR1,
	)

	// If the bot is run for the first time in a new session,
	// it should be forked at least once to play nice with systemd.
	// Forking is triggered by sending SIGUSR1 to the current process.
	if connectionCount == 0 {
		syscall.Kill(os.Getpid(), syscall.SIGUSR1)
	}

	log.Println("[bot] Waiting for signals...")
	for sig := range signals {
		log.Println("[bot] received signal:", sig)
		if sig != syscall.SIGUSR1 {
			return
		}

		log.Println("[bot] forking process...")
		err := doFork(b)
		if err != nil {
			log.Println("[bot]", err)
		}
	}
}

// doFork forks the current process into a child process and passes the
// given client connections along to be inherited.
//
// The forked process is called with the `-fork N` command line parameter.
// Where N is the number of file descriptors being passed along. This is
// used by the InheritedFiles() call to rebuild the files. Currently
// there is only one connection per bot implemented (N=1).
func doFork(b *Bot) error {

	// Build the command line arguments for our child process.
	// This includes any custom arguments defined in the profile.
	argv := b.profile.ForkArgs()
	args := append([]string{"-fork", "1"}, argv...)

	// Initialize the command runner.
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fd, _ := b.client.File()
	cmd.ExtraFiles = []*os.File{fd}

	// Fork the process.
	return cmd.Start()
}

// inheritedFiles returns a list of N file descriptors inherited from a
// previous session through the Fork call.
//
// This function assumes that flag.Parse() has been called at least once
// already. The `-fork` flag has been registered during initialization of
// this package.
func inheritedFiles() []*os.File {
	if connectionCount == 0 {
		return nil
	}

	out := make([]*os.File, connectionCount)

	for i := range out {
		out[i] = os.NewFile(3+uintptr(i), "conn"+strconv.Itoa(i))
	}

	return out
}
