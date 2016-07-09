// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package sync

import (
	"monkeybird/irc"
	"monkeybird/irc/proto"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Channel defines a channel description, returned by the LIST query.
type Channel struct {
	Name        string
	Description string
	Users       int
}

// List lists all channels on the network.
//
// Returns false if the call times out before the full response has been received.
func List(w ProtocolWriter) ([]Channel, bool) {
	var set []Channel
	var quit sync.Once
	done := make(chan struct{})

	handler := func(w irc.ResponseWriter, r *irc.Request) {
		switch r.Type {
		case "322":
			fields := strings.Fields(r.Data)
			if len(fields) >= 3 {
				// Strip leading :
				fields[2] = fields[2][1:]

				users, _ := strconv.Atoi(fields[1])

				set = append(set, Channel{
					Name:        fields[0],
					Users:       users,
					Description: strings.Join(fields[2:], " "),
				})
			}

		case "323":
			quit.Do(func() {
				close(done)
			})
		}
	}

	// Bind temporary handlers.
	w.Bind("322", handler)
	w.Bind("323", handler)

	// Perform asychronous query.
	proto.List(w)

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
	w.Unbind("322", handler)
	w.Unbind("323", handler)
	return set, ok
}
