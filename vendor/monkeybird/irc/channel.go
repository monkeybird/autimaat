// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package irc

// Channel defines a single IRC channel.
type Channel struct {
	Name     string // Channel's name.
	Key      string // Authentication key for protected channel.
	Password string // Chanserv password.
}

// Returns true if the channel is local to the current server.
// This is the case when its name starts with '&'.
func (c *Channel) IsLocal() bool {
	return len(c.Name) > 0 && c.Name[0] == '&'
}
