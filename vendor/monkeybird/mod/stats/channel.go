// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package stats

import (
	"monkeybird/irc"
	"monkeybird/irc/proto/sync"
	"sort"
	"strings"
)

// ChannelStats defines stats for a channel.
type ChannelStats struct {
	Name  string          `json:",omitempty"` // Name of the target channel
	Users UsersByHostmask `json:",omitempty"` // List of current- and previously known users.
}

// newChannel returns a new hannelstats struct, filled with channel
// data, as best as possible,
func newChannel(w irc.ResponseWriter, name string) *ChannelStats {
	cs := &ChannelStats{
		Name: name,
	}

	if pw, ok := w.(sync.ProtocolWriter); ok {
		users, _ := sync.ChannelUsers(pw, name)
		cs.Users = make(UsersByHostmask, len(users))

		for i := range users {
			cs.Users[i] = newUser(users[i].Hostmask, users[i].Nickname)
		}

		sort.Sort(cs.Users)
	}

	return cs
}

// ChannelList defines a list of channels, sortable by name.
type ChannelList []*ChannelStats

func (v ChannelList) Len() int           { return len(v) }
func (v ChannelList) Less(i, j int) bool { return v[i].Name < v[j].Name }
func (v ChannelList) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }

// Get finds the stats for the given channel name. Implicitely
// creates a new entry, if the channel is not found. This requires access
// to the underlying network stream, hence the ResponseWriter parameter.
func (v *ChannelList) Get(w irc.ResponseWriter, name string) *ChannelStats {
	name = strings.ToLower(name)

	idx := v.index(name)
	if idx != -1 {
		return (*v)[idx]
	}

	cs := newChannel(w, name)
	*v = append(*v, cs)
	sort.Sort(*v)
	return cs
}

// index returns the index of the given channel by its name.
// This expects the list to be sorted, as it performs a binary search.
func (v ChannelList) index(name string) int {
	var lo int
	hi := len(v) - 1

	for lo < hi {
		mid := lo + ((hi - lo) / 2)

		if v[mid].Name < name {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	if hi == lo && v[lo].Name == name {
		return lo
	}

	return -1
}
