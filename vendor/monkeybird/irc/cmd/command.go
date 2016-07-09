// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cmd

import (
	"monkeybird/irc"
	"regexp"
	"strings"
)

type Handler func(irc.ResponseWriter, *Request)

// Command defines a single command which can be called by IRC users.
type Command struct {
	Name        string  // Name by which the command is called.
	Description string  // Command description. Used in help listing.
	Handler     Handler // Command handler.
	Params      []Param // Command parameter list.
	Restricted  bool    // Command may only be run by authorized users.
}

// newCommand creates a new command.
func newCommand(name, description string, restricted bool, handler Handler) *Command {
	c := new(Command)
	c.Name = strings.ToLower(name)
	c.Description = description
	c.Restricted = restricted
	c.Handler = handler
	return c
}

// Add adds a new command parameter.
func (c *Command) Add(name, description string, required bool, pattern *regexp.Regexp) *Command {
	var p Param

	p.Name = strings.ToLower(name)
	p.Description = description
	p.Required = required

	if pattern == nil {
		p.Pattern = RegAny
	} else {
		p.Pattern = pattern
	}

	c.Params = append(c.Params, p)
	return c
}

// RequiredParamCount returns the amunt of required parameters for this command.
func (c *Command) RequiredParamCount() int {
	var count int

	for i := range c.Params {
		if c.Params[i].Required {
			count++
		}
	}

	return count
}
