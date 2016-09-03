// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package irc

import (
	"fmt"
	"strings"
)

// RequestFunc defines a handler for a request binding.
type RequestFunc func(ResponseWriter, *Request)

// Request defines a single incoming message from a server.
type Request struct {
	SenderName string // Nick name of sender.
	SenderMask string // Hostmask of sender.
	Type       string // Type of message: "001", "PRIVMSG", "PING", etc.
	Target     string // Receiver of reply.
	Data       string // Message content.
}

// FromChannel returns true if this request came from a channel context
// instead of a user or service.
func (r *Request) FromChannel() bool {
	if len(r.Target) == 0 {
		return false
	}

	c := r.Target[0]
	return c == '#' || c == '&' || c == '!' || c == '+'
}

// Fields returns the message payload, but skips the first n words.
// The result is returned as a slice of individual words.
func (r *Request) Fields(n int) []string {
	words := strings.Fields(r.Data)
	if n < 0 || n >= len(words) {
		return nil
	}
	return words[n:]
}

// String returns a string representation of the request data.
func (r *Request) String() string {
	return fmt.Sprintf("%s %s %s %s %s",
		r.SenderMask, r.SenderName, r.Type, r.Target, r.Data)
}

// IsPrivMsg returns true if the request comes from either a user or
// a channel, as a PRIVMSG. This has its own method, because it is a
// commonly used filter.
func (r *Request) IsPrivMsg() bool { return r.Type == "PRIVMSG" }
