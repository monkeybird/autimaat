// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package weather

import (
	"fmt"
	"net/url"

	"github.com/monkeybird/autimaat/irc"
)

type location struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country_iso3166"`
}

// newLocation creates a new location from the given command request data.
func newLocation(r *irc.Request) *location {
	var l location

	fields := r.Fields(1)
	l.City = url.QueryEscape(fields[0])

	if len(fields) > 1 {
		l.Country = url.QueryEscape(fields[1])
	}

	if len(fields) > 2 {
		l.State = url.QueryEscape(fields[2])
	}

	return &l
}

func (l *location) String() string {
	if len(l.Country) == 0 {
		return l.City
	}

	if len(l.State) == 0 {
		return fmt.Sprintf("%s/%s", l.Country, l.City)
	}

	return fmt.Sprintf("%s/%s/%s", l.Country, l.State, l.City)
}
