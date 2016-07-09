// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package irc

import (
	"bytes"
	"strings"
)

var (
	bNameSplitter = []byte{'!'}
	bSpace        = []byte{' '}
	bPING         = []byte("PING")
	bERROR        = []byte("ERROR")
	bQUIT         = []byte("QUIT")
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

// parseRequest reads the given data and parses it into a new request.
// Returns nil if this is not a valid protocol message.
func parseRequest(data []byte) *Request {
	fields := bytes.Fields(data)
	if len(fields) == 0 {
		return nil
	}

	// We may be dealing with utility messages like ERROR or PING.
	switch {
	case bytes.Index(data, bQUIT) > -1:
		return nil

	case bytes.HasPrefix(data, bPING):
		return &Request{
			Type: "PING",
			Data: string(fields[1][1:]),
		}

	case bytes.HasPrefix(data, bERROR):
		return &Request{
			Type: "ERROR",
			Data: string(fields[1][1:]),
		}
	}

	// Strip leading ':' characters.
	for i := 0; i < 4 && i < len(fields); i++ {
		if fields[i][0] == ':' {
			fields[i] = fields[i][1:]
		}
	}

	var msg Request

	idx := bytes.Index(fields[0], bNameSplitter)
	if idx > -1 {
		msg.SenderName = string(fields[0][:idx])
		msg.SenderMask = string(fields[0][idx+1:])
	} else {
		msg.SenderName = string(fields[0])
		msg.SenderMask = msg.SenderName
	}

	msg.Type = string(fields[1])
	msg.Target = string(fields[2])

	if len(fields) > 3 {
		msg.Data = string(bytes.Join(fields[3:], bSpace))
	}

	return &msg
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

// Remainder returns the message payloud, but skips the first n words.
func (r *Request) Remainder(n int) string {
	words := strings.Fields(r.Data)
	if n < 0 || n >= len(words) {
		return ""
	}

	return strings.TrimSpace(strings.Join(words[n:], " "))
}
