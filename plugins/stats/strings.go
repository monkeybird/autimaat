// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package stats

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
	TextLastOnDisplay = "%s, ik heb %s het laatst gezien op %s, om %s."

	TextFirstOn        = "firston"
	TextFirstOnDisplay = "%s, ik heb %s voor het eerst gezien op %s, om %s."
)
