// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cmd

import (
	"monkeybird/tr"
	"regexp"
	"strconv"
)

// Param defines a parameter for a command.
type Param struct {
	Name        string         // Parameter name -- used in help listing.
	Description string         // Parameter description -- used in help listing.
	Value       string         // Parameter value.
	Pattern     *regexp.Regexp // Pattern defining the type of accepted value.
	Required    bool           // Parameter is required or not?
}

// validate returns true if the given value matches the param pattern.
func (p *Param) validate(v string) bool { return p.Pattern.MatchString(v) }

func (p *Param) String() string { return p.Value }

func (p *Param) Int() int64 {
	n, _ := strconv.ParseInt(p.Value, 0, 64)
	return n
}

func (p *Param) Uint() uint64 {
	n, _ := strconv.ParseUint(p.Value, 0, 64)
	return n
}

func (p *Param) Float() float64 {
	n, _ := strconv.ParseFloat(p.Value, 64)
	return n
}

// Bool returns the boolean value represented by the parameter.
// This expects any of: (0|1|t(rue)?|f(alse)?|y(es)?|no?|on|off)
// Any other value returns false.
func (p *Param) Bool() bool { return tr.ParseBool(p.Value) }
