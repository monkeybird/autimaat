// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package test implements a few utility types and functions, which
// should make the writing of unit tests for commands a little easier.
package test

import (
	"bytes"
	"monkeybird/irc"
	"strings"
	"testing"
)

// Some predefined sample values for various applications.
const (
	Nickname        = "test"
	VersionMajor    = 1
	VersionMinor    = 2
	VersionRevision = 3
	ChannelName     = "#test"
	SenderName      = "bob"
	SenderMask      = "~bob@server.com"
)

// NewRequest returns a new irc request with some sample data.
func NewRequest() *irc.Request {
	return &irc.Request{
		SenderName: SenderName,
		SenderMask: SenderMask,
		Target:     ChannelName,
		Type:       "PRIVMSG",
	}
}

// MockWriter defines a fake network stream. It qualifies as a
// irc.ResponseWriter implementation and as such, it can be passed straight
// into protocol and command handlers.
//
// It buffers all output in memory for later inspection.
type MockWriter struct {
	buf bytes.Buffer
}

func (w *MockWriter) Close() error                { return nil }
func (w *MockWriter) Write(p []byte) (int, error) { return w.buf.Write(p) }

// Verify compares the contents of the writer's buffer with that of the
// given lines of text.
func (w *MockWriter) Verify(t *testing.T, want ...string) {
	have := w.lines()

	if len(have) != len(want) {
		t.Fatalf("result count mismatch; want: %d, have: %d",
			len(want), len(have))
	}

	for i, wantValue := range want {
		haveValue := have[i]
		if !strings.EqualFold(wantValue, haveValue) {
			t.Fatalf("result mismatch at %d;\nwant: %q\nhave: %q",
				i, wantValue, haveValue)
		}
	}
}

// lines returns the buffer contents as a list of separate lines.
// This omits empty lines.
func (w *MockWriter) lines() []string {
	lines := strings.Split(w.buf.String(), "\n")
	out := make([]string, 0, len(lines))

	for _, v := range lines {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			out = append(out, v)
		}
	}

	return out
}
