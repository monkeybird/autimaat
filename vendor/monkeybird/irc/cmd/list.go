// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cmd

import (
	"strings"
)

// List defines a list of commands, sortable by name.
type List []*Command

func (cl List) Len() int           { return len(cl) }
func (cl List) Less(i, j int) bool { return cl[i].Name < cl[j].Name }
func (cl List) Swap(i, j int)      { cl[i], cl[j] = cl[j], cl[i] }

// Find finds the command for the given name.
// Returns nil if it was not found.
func (cl List) Find(name string) *Command {
	idx := cl.Index(name)
	if idx > -1 {
		return cl[idx]
	}
	return nil
}

// Index returns the index of the command for the given name.
// Returns -1 if it was not found.
func (cl List) Index(name string) int {
	var lo int
	hi := len(cl) - 1

	name = strings.ToLower(name)

	for lo < hi {
		mid := lo + ((hi - lo) / 2)

		if cl[mid].Name < name {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	if hi == lo && cl[lo].Name == name {
		return lo
	}

	return -1
}
