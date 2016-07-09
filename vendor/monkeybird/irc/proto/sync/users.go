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

// Names returns users in the given list of <channels>, If <channels> is
// omitted, all users are shown, grouped by channel name with all users who are
// not on a channel being shown as part of channel "*". If <server> is specified,
// the command is sent to <server> for evaluation.
//
// Returns false if the call times out before the full response has been received.
func Names(w ProtocolWriter, channels ...string) (map[string][]string, bool) {
	var quit sync.Once
	set := make(map[string][]string)
	done := make(chan struct{})

	handler := func(w irc.ResponseWriter, r *irc.Request) {
		switch r.Type {
		case "353": // ( '=' / '*' / '@' ) <channel> ' ' : [ '@' / '+' ] <nick> *( ' ' [ '@' / '+' ] <nick> )
			fields := strings.Fields(r.Data)
			if len(fields) < 3 {
				break
			}

			for i := range fields {
				fields[i] = filterName(fields[i])
			}

			key := fields[1]
			set[key] = append(set[key], fields[2:]...)

		case "366":
			quit.Do(func() {
				close(done)
			})
		}
	}

	// Bind temporary handlers.
	w.Bind("353", handler)
	w.Bind("366", handler)

	// Perform asychronous query.
	proto.Names(w, channels...)

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
	w.Unbind("353", handler)
	w.Unbind("366", handler)
	return set, ok
}

// ChannelUsers returns descriptors for all users in the given channel.
func ChannelUsers(w ProtocolWriter, channel string) ([]WhoisRecord, bool) {
	nameMap, ok := Names(w, channel)
	if !ok {
		return nil, false
	}

	out := make([]WhoisRecord, 0, len(nameMap))

	for _, names := range nameMap {
		for _, name := range names {
			record, ok := Whois(w, name)
			if ok {
				out = append(out, record)
			}
		}
	}

	return out, true
}
