// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package sync

import (
	"monkeybird/irc"
	"monkeybird/irc/proto"
	"strings"
	"sync"
	"time"
)

// WhoisRecord defines a response for a WHOIS or WHOWAS query.
type WhoisRecord struct {
	Hostmask string // User's hostmask.
	Nickname string // User's nickname.

	// List of channels the user is on.
	// This may be empty if the user's mode is private.
	Channels []string
}

// Whois performs a WHOIS query for the given target and waits for the
// response. The given handler is called with said response.
//
// If a server name is specified, this is included in the query. It targets
// the query at this server specifically, instead of the entire network.
//
// Returns false if the call times out before the full response has been received.
func Whois(w ProtocolWriter, target string, server ...string) (WhoisRecord, bool) {
	var record WhoisRecord
	var quit sync.Once
	done := make(chan struct{})

	handler := func(w irc.ResponseWriter, r *irc.Request) {
		switch r.Type {
		case "311":
			fields := strings.Fields(r.Data)
			if len(fields) >= 3 {
				record.Nickname = fields[0]
				record.Hostmask = fields[1] + "@" + fields[2]
			}

		case "319":
			fields := strings.Fields(r.Data)
			if len(fields) >= 2 {
				fields = fields[1:]

				for i := range fields {
					fields[i] = filterName(fields[i])
				}

				record.Channels = fields
			}

		case "318":
			quit.Do(func() {
				close(done)
			})
		}
	}

	// Bind temporary handlers.
	w.Bind("311", handler)
	w.Bind("318", handler)
	w.Bind("319", handler)

	// Perform asychronous query.
	proto.Whois(w, target, server...)

	// Wait for the response to be collected or a timeout signal to occur.
	// Whichever comes first.
	var ok bool
	select {
	case <-done:
		ok = true
	case <-time.After(Timeout):
		ok = false
	}

	// Remove temporary handlers.
	w.Unbind("311", handler)
	w.Unbind("318", handler)
	w.Unbind("319", handler)
	return record, ok
}
