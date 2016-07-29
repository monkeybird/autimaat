// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cmd

import (
	"log"
	"monkeybird/irc"
	"monkeybird/irc/proto"
	"monkeybird/text"
	"monkeybird/tr"
	"sort"
	"strings"
	"sync"
	"time"
)

// CommandBinder defines a type which can bind/unbind commands.
type CommandBinder interface {
	Bind(string, string, bool, Handler) *Command
	Unbind(string)
	HelpHandler(irc.ResponseWriter, *Request)
	Clear()
}

// AuthFunc returns true if the given request comes from a
// whitelisted user. This function is used by the command dispatcher
// to ensure the user is allowed to execute a given, restricted command.
type AuthFunc func(*irc.Request) bool

// Set defines a set of bound commands.
type Set struct {
	m            sync.RWMutex
	authenticate AuthFunc
	data         List
	prefix       string
}

// New creates a new, empty set for the given prefix and auth handler.
// The auth handler is used to ensure a caller is allowed to run a
// restricted command.
func New(prefix string, authenticate AuthFunc) *Set {
	if authenticate == nil {
		authenticate = func(*irc.Request) bool { return false }
	}

	return &Set{
		prefix:       prefix,
		authenticate: authenticate,
	}
}

// Dispatch accepts the given message and issues command calls if applicable.
// Returns false if no command call was issued.
func (s *Set) Dispatch(w irc.ResponseWriter, r *irc.Request) bool {
	s.m.RLock()
	defer s.m.RUnlock()

	// Did we get a command request?
	if !strings.HasPrefix(r.Data, s.prefix) {
		return false
	}

	// Split message data into command name and individual arguments.
	name, args := split(r.Data[len(s.prefix):])
	if len(name) == 0 {
		return false
	}

	// Find the command instance.
	cmd := s.data.Find(name)
	if cmd == nil {
		return false
	}

	// Ensure the caller is authorized to run this command.
	if cmd.Restricted && !s.authenticate(r) {
		proto.PrivMsg(w, r.SenderName, tr.CommandsAccessDenied, cmd.Name)
		return false
	}

	// Ensure we have enough parameters.
	if cmd.RequiredParamCount() > len(args) {
		proto.PrivMsg(w, r.SenderName, tr.CommandsMissingParameter, cmd.Name)
		return false
	}

	// Process and validate each parameter value.
	params := make([]Param, 0, len(cmd.Params))

	for i := 0; i < len(args) && i < len(cmd.Params); i++ {
		if cmd.Params[i].validate(args[i]) {
			params = append(params, Param{Value: args[i]})
			continue
		}

		proto.PrivMsg(w, r.SenderName, tr.CommandsInvalidParameter,
			cmd.Name, cmd.Params[i].Name)
		return false
	}

	go func() {
		// Ensure command handlers don't bring the entire bot down when a panic occurs.
		defer func() {
			x := recover()
			if x == nil {
				return
			}

			log.Printf("Command error: %v", x)
			log.Printf("> %#v", r)
		}()

		cmd.Handler(w, &Request{Request: r, Params: params})
	}()

	return true
}

// HelpHandler is a builtin command handler which yields detailed help for a
// given command or a listing of all available commands.
func (s *Set) HelpHandler(w irc.ResponseWriter, r *Request) {
	s.m.RLock()
	defer s.m.RUnlock()

	if len(r.Params) == 0 {
		for _, cmd := range s.data {
			var status string
			if cmd.Restricted {
				status = "*"
			}

			proto.PrivMsg(w, r.SenderName,
				"%s%s%s: %s",
				s.prefix,
				text.Bold(cmd.Name),
				status,
				cmd.Description,
			)

			<-time.After(750 * time.Millisecond)
		}

		return
	}

	name := r.Params[0].String()
	cmd := s.data.Find(name)
	if cmd == nil {
		return
	}

	var status string
	if cmd.Restricted {
		status = " " + tr.CommandsRestricted
	}

	proto.PrivMsg(w, r.SenderName,
		"%s%s:%s %s",
		s.prefix,
		text.Bold(cmd.Name),
		status,
		cmd.Description,
	)

	for _, p := range cmd.Params {
		var required string

		if !p.Required {
			required = " " + tr.CommandsOptional
		}

		proto.PrivMsg(
			w,
			r.SenderName,
			" <%s>:%s %s",
			text.Bold(p.Name),
			required,
			p.Description,
		)

		<-time.After(750 * time.Millisecond)
	}
}

// Clear unbinds all commands.
func (s *Set) Clear() {
	s.m.Lock()
	s.data = nil
	s.m.Unlock()
}

// Bind binds the given command.
func (s *Set) Bind(name, description string, restricted bool, handler Handler) *Command {
	s.m.Lock()
	cmd := newCommand(name, description, restricted, handler)
	s.data = append(s.data, cmd)
	sort.Sort(s.data)
	s.m.Unlock()
	return cmd
}

// Unbind removes the given command.
func (s *Set) Unbind(name string) {
	s.m.Lock()

	idx := s.data.Index(name)
	if idx > -1 {
		copy(s.data[idx:], s.data[idx+1:])
		s.data[len(s.data)-1] = nil
		s.data = s.data[:len(s.data)-1]
	}

	s.m.Unlock()
}
