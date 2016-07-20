// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package stats

import (
	"sort"
	"strings"
	"time"
)

// UserStats defines a single user and its statistics.
type UserStats struct {
	Hostmask  string    `json:",omitempty"` // User's hostmask -- primary means of identification by bot.
	Nicknames []string  `json:",omitempty"` // Last known nicknames used by this user.
	FirstSeen time.Time `json:",omitempty"` // Date/time at which user was first seen by the bot.
	LastSeen  time.Time `json:",omitempty"` // Date/time at which user was last seen by the bot.
}

// newUser returns a new userStats struct.
func newUser(hostmask string, nicknames ...string) *UserStats {
	sort.Strings(nicknames)

	return &UserStats{
		Hostmask:  strings.ToLower(hostmask),
		Nicknames: nicknames,
		FirstSeen: time.Now(),
		LastSeen:  time.Now(),
	}
}

// AddNickname adds the given nickname to the list of known
// nicknames, provided it doesn't yet exist.
func (us *UserStats) AddNickname(name string) {
	for _, v := range us.Nicknames {
		if strings.EqualFold(v, name) {
			return
		}
	}

	us.Nicknames = append(us.Nicknames, name)
	sort.Strings(us.Nicknames)
}

// UsersByHostmask defines aset of user, sortable by first-seen date.
type UsersByFirstSeen []*UserStats

func (v UsersByFirstSeen) Len() int           { return len(v) }
func (v UsersByFirstSeen) Less(i, j int) bool { return v[i].FirstSeen.Before(v[j].FirstSeen) }
func (v UsersByFirstSeen) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }

// UsersByHostmask defines aset of user, sortable by hostmask.
type UsersByHostmask []*UserStats

func (v UsersByHostmask) Len() int           { return len(v) }
func (v UsersByHostmask) Less(i, j int) bool { return v[i].Hostmask < v[j].Hostmask }
func (v UsersByHostmask) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }

// Find returns the userstats for the given name/hostmask.
// Returns nil if no match could be found.
func (v UsersByHostmask) Find(name string) *UserStats {
	idx := v.index(name)
	if idx != -1 {
		return v[idx]
	}

	for _, us := range v {
		for _, nick := range us.Nicknames {
			if strings.EqualFold(name, nick) {
				return us
			}
		}
	}

	return nil
}

// Get returns the userstats for the given hostmask. Implicitely
// creates a new entry, if the given user could not be found.
func (v *UsersByHostmask) Get(hostmask string) *UserStats {
	idx := v.index(hostmask)
	if idx != -1 {
		return (*v)[idx]
	}

	us := newUser(hostmask)
	*v = append(*v, us)
	sort.Sort(*v)
	return us
}

// index returns the index of the given user bu its hostmask.
// This expects the list to be sorted, as it performs a binary search.
func (v UsersByHostmask) index(hostmask string) int {
	var lo int
	hi := len(v) - 1

	hostmask = strings.ToLower(hostmask)

	for lo < hi {
		mid := lo + ((hi - lo) / 2)

		if v[mid].Hostmask < hostmask {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	if hi == lo && v[lo].Hostmask == hostmask {
		return lo
	}

	return -1
}
