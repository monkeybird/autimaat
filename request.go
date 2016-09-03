// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bytes"

	"github.com/monkeybird/autimaat/irc"
)

var (
	bNameSplitter = []byte{'!'}
	bSpace        = []byte{' '}
	bPING         = []byte("PING")
	bERROR        = []byte("ERROR")
	bQUIT         = []byte("QUIT")
)

// parseRequest reads the given message payload and parses it into the
// specified request structure. Returns false if the payload is not a valid
// protocol message.
func parseRequest(r *irc.Request, data []byte) bool {
	fields := bytes.Fields(data)
	if len(fields) == 0 {
		return false
	}

	// We may be dealing with utility messages like ERROR or PING.
	switch {
	case bytes.Index(data, bQUIT) > -1:
		return false

	case bytes.HasPrefix(data, bPING):
		r.Type = "PING"
		r.Data = string(fields[1][1:])
		r.SenderMask = ""
		r.SenderName = ""
		r.Target = ""
		return true

	case bytes.HasPrefix(data, bERROR):
		r.Type = "ERROR"
		r.Data = string(fields[1][1:])
		r.SenderMask = ""
		r.SenderName = ""
		r.Target = ""
		return true
	}

	// Strip leading ':' characters from all fields, except the actual
	// message contents.
	for i := 0; i < 4 && i < len(fields); i++ {
		if fields[i][0] == ':' {
			fields[i] = fields[i][1:]
		}
	}

	idx := bytes.Index(fields[0], bNameSplitter)
	if idx > -1 {
		r.SenderName = string(fields[0][:idx])
		r.SenderMask = string(fields[0][idx+1:])
	} else {
		r.SenderName = string(fields[0])
		r.SenderMask = r.SenderName
	}

	r.Type = string(fields[1])
	r.Target = string(fields[2])

	if len(fields) > 3 {
		r.Data = string(bytes.Join(fields[3:], bSpace))
	} else {
		r.Data = ""
	}

	return true
}
