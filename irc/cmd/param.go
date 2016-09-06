// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cmd

import (
	"regexp"
	"strconv"
	"strings"
)

// ParamList defines a list of command parameters.
type ParamList []Param

func (p ParamList) Len() int            { return len(p) }
func (p ParamList) String(n int) string { return p[n].String() }
func (p ParamList) Int(n int) int64     { return p[n].Int() }
func (p ParamList) Uint(n int) uint64   { return p[n].Uint() }
func (p ParamList) Float(n int) float64 { return p[n].Float() }
func (p ParamList) Bool(n int) bool     { return p[n].Bool() }

// Join returns all parameter values, concatenated into a single string.
// Each entry is separated by a blank space.
func (p ParamList) Join() string {
	out := make([]string, len(p))

	for i := range p {
		out[i] = p[i].Value
	}

	return strings.Join(out, " ")
}

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
// True is represented by the values: "1", "t", "true", "y", "yes", "on"
// Any other value returns false.
func (p *Param) Bool() bool {
	switch strings.ToLower(p.Value) {
	case "1", "t", "true", "y", "yes", "on":
		return true
	default:
		return false
	}
}
