// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package irc

import "fmt"

// Action returns the given, formatted message as a user action.
// E.g.: /me <something something>
func Action(f string, argv ...interface{}) string {
	return fmt.Sprintf("\x01ACTION %s\x01", fmt.Sprintf(f, argv...))
}

// Bold returns the given value as bold text.
func Bold(f string, argv ...interface{}) string {
	return fmt.Sprintf("\x02%s\x02", fmt.Sprintf(f, argv...))
}

// Italic returns the given value as italicized text.
func Italic(f string, argv ...interface{}) string {
	return fmt.Sprintf("\x1d%s\x1d", fmt.Sprintf(f, argv...))
}

// Underline returns the given value as underlined text.
func Underline(f string, argv ...interface{}) string {
	return fmt.Sprintf("\x1f%s\x1f", fmt.Sprintf(f, argv...))
}
