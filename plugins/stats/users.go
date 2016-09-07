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
// provided it is not already known.
func (u *User) AddNickname(v string) {
	v = strings.ToLower(v)
	idx := stringIndex(u.Nicknames, v)
	if idx > -1 {
		return
	}

	u.Nicknames = append(u.Nicknames, v)
	sort.Strings(u.Nicknames)
}

// UserList defines a set of user descriptors, sortable by hostmask.
type UserList []*User

func (cl UserList) Len() int           { return len(cl) }
func (cl UserList) Less(i, j int) bool { return cl[i].Hostmask < cl[j].Hostmask }
func (cl UserList) Swap(i, j int)      { cl[i], cl[j] = cl[j], cl[i] }

// Get returns a user entry for the given hostmask. If it doesn't exist yet,
// a new entry is created and added to the list implicitely.
//
// This implicitely updates the LastSeen timestamp for the user.
func (cl *UserList) Get(mask string) *User {
	idx := userIndex(*cl, mask)
	if idx > -1 {
		(*cl)[idx].LastSeen = time.Now()
		return (*cl)[idx]
	}

	usr := &User{
		Hostmask:  strings.ToLower(mask),
		FirstSeen: time.Now(),
		LastSeen:  time.Now(),
	}

	*cl = append(*cl, usr)
	sort.Sort(*cl)
	return usr
}

// Find finds the user which exactly matches the given hostmask,
// or all users which have a fuzzy match with the given nickname.
// It returns at most limit users.
func (cl UserList) Find(name string, limit int) []*User {
	name = strings.ToLower(name)

	// Known hostmask?
	idx := userIndex(cl, name)
	if idx > -1 {
		return []*User{cl[idx]}
	}

	// Known nickname then perhaps? This will return the first
	// instance of the nickname we find.
	out := make([]*User, 0, 4)

	// First find exact name matches.
	for _, usr := range cl {
		if userIndex(out, usr.Hostmask) > -1 {
			continue
		}

		if stringExactMatch(usr.Nicknames, name) {
			out = append(out, usr)

			if len(out) >= limit {
				break
			}
		}
	}

	// Then find any partial matches.
	for _, usr := range cl {
		if userIndex(out, usr.Hostmask) > -1 {
			continue
		}

		if stringPartialMatch(usr.Nicknames, name) {
			out = append(out, usr)

			if len(out) >= limit {
				break
			}
		}
	}

	return out
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

// stringPartialMatch returns true if any of the values in set are
// at least partially identical to v.
func stringPartialMatch(set []string, v string) bool {
	for _, s := range set {
		if strings.Index(s, v) > -1 {
			return true
		}
	}
	return false
}

// stringExactMatch returns true if any of the values in set are identical
// to v.
func stringExactMatch(set []string, v string) bool {
	for _, s := range set {
		if s == v {
			return true
		}
	}
	return false
}
