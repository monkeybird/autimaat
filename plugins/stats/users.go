// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package stats

import (
	"sort"
	"strings"
	"time"
)

// User defines the hostmask and all known nicknames for a single user.
// Additionally, it defines some timestamps.
type User struct {
	Hostmask  string
	Nicknames []string
	FirstSeen time.Time
	LastSeen  time.Time
}

// AddNickname adds the given nickname to the user's name list,
// provided it is not already known. It returns true if the nickname
// is new.
func (u *User) AddNickname(v string) bool {
	v = strings.ToLower(v)
	idx := stringIndex(u.Nicknames, v)
	if idx > -1 {
		return false
	}

	u.Nicknames = append(u.Nicknames, v)
	sort.Strings(u.Nicknames)
	return true
}

// UserList defines a set of user descriptors, sortable by hostmask.
type UserList []*User

func (cl UserList) Len() int           { return len(cl) }
func (cl UserList) Less(i, j int) bool { return cl[i].Hostmask < cl[j].Hostmask }
func (cl UserList) Swap(i, j int)      { cl[i], cl[j] = cl[j], cl[i] }

// Get returns a user entry for the given hostmask. If it doesn't exist yet,
// a new entry is created and added to the list implicitely. It returns true
// if the user is new.
//
// This implicitely updates the LastSeen timestamp for the user.
func (cl *UserList) Get(mask string) (*User, bool) {
	idx := userIndex(*cl, mask)
	if idx > -1 {
		(*cl)[idx].LastSeen = time.Now()
		return (*cl)[idx], false
	}

	usr := &User{
		Hostmask:  strings.ToLower(mask),
		FirstSeen: time.Now(),
		LastSeen:  time.Now(),
	}

	*cl = append(*cl, usr)
	sort.Sort(*cl)
	return usr, true
}

// Find finds the command for the given hostmask or nickname.
// Returns nil if it was not found.
func (cl UserList) Find(name string) *User {
	name = strings.ToLower(name)

	// Known hostmask?
	idx := userIndex(cl, name)
	if idx > -1 {
		return cl[idx]
	}

	// Known nickname then perhaps? This will return the first
	// instance of the nickname we find.
	for _, usr := range cl {
		idx := stringIndex(usr.Nicknames, name)
		if idx > -1 {
			return usr
		}
	}

	return nil
}

// userIndex returns the index of user hostmask v in set.
// Returns -1 if it was not found. The list is expected to be sorted.
func userIndex(set []*User, v string) int {
	var lo int
	hi := len(set) - 1

	for lo < hi {
		mid := lo + ((hi - lo) / 2)

		if set[mid].Hostmask < v {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	if hi == lo && set[lo].Hostmask == v {
		return lo
	}

	return -1
}

// stringIndex returns the index of string v in set.
// Returns -1 if it was not found. The list is expected to be sorted.
func stringIndex(set []string, v string) int {
	var lo int
	hi := len(set) - 1

	for lo < hi {
		mid := lo + ((hi - lo) / 2)

		if set[mid] < v {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	if hi == lo && set[lo] == v {
		return lo
	}

	return -1
}
