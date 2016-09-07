// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package stats

import (
	"fmt"
	"time"
)

const (
	// ref: https://godoc.org/time#Time.Format
	TextDateFormat  = "2 January, 2006"
	TextTimeFormat  = "15:04 MST"
	TextNick        = "gebruiker"
	TextUnknownUser = "%s, %s heb ik niet gezien."

	TextWhoisName        = "wieis"
	TextWhoisDisplay     = "%s, ik ken %s sinds %s als: %s"
	TextWhoisUnknownUser = "%s, ik herken %s niet."

	TextLastOn        = "laston"
	TextLastOnDisplay = "%s, ik heb %s (%s) het laatst gezien op %s, om %s (± %s geleden)."

	TextFirstOn        = "firston"
	TextFirstOnDisplay = "%s, ik heb %s (%s) voor het eerst gezien op %s, om %s (± %s geleden)."
)

// FormatDuration returns a custom, string representation of the
// given duration. It is in the correct language and slightly more
// readable than Go's default time.Duration.String() output.
func FormatDuration(d time.Duration) string {
	var dd time.Duration

	hh := d / time.Hour
	if hh >= 24 {
		dd = hh / 24
		hh = hh % 24
	}

	if dd > 0 {
		return fmt.Sprintf("%d dagen en %d uur", dd, hh)
	}

	return fmt.Sprintf("%d uur", hh)
}
