// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

/*
Package sync provides a few IRC protocol calls in a synchronous form.
These can be used to send the server a query and wait for the response.

For instance, a WHOIS query:

	record, ok := sync.Whois(client, "bob")
	if ok {
		fmt.Println(record)
	}

*/
package sync

import (
	"monkeybird/irc"
	"time"
)

// ProtocolWriter defines the interface necessary in most of the
// functions in this package. It is a merger of irc.ResponseWriter and
// irc.ProtocolBinder.
type ProtocolWriter interface {
	irc.ResponseWriter
	irc.ProtocolBinder
}

// Timeout defines the time to wait for synchronous responses.
// Synchronous calls are canceled if the timeout passes before the
// full response has been received.
var Timeout = time.Second * 10

// filterName returns the given name minus any mode prefixes like @, +, etc
func filterName(name string) string {
	if len(name) < 2 {
		return name
	}

	for len(name) > 0 {
		switch name[0] {
		case ':', '@', '+', '~', '%':
			name = name[1:]
		default:
			return name
		}
	}

	return name
}
