// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cmd

import (
	"log"
	"runtime"
	"sort"
	"strings"

	"github.com/monkeybird/autimaat/irc"
	"github.com/monkeybird/autimaat/irc/proto"
)

// AuthFunc returns true if the given hostmask defines a whitelisted user.
// This function is used by the command dispatcher to ensure the user is
// allowed to execute a given, restricted command.
type AuthFunc func(string) bool

// Set defines a set of bound commands.
type Set struct {
	authenticate AuthFunc
	data         List
	prefix       string
}

// New creates a new, empty set for the given prefix and auth handler.
// The auth handler is used to ensure a caller is allowed to run a
// restricted command. This can be nil, which will outright deny access
// to all commands which have the restricted flag set.
func New(prefix string, authenticate AuthFunc) *Set {
	if authenticate == nil {
		authenticate = func(string) bool { return false }
	}

	return &Set{
		prefix:       prefix,
		authenticate: authenticate,
	}
}

// Dispatch accepts the given message and issues command calls if applicable.
// Returns false if no command call was issued.
func (s *Set) Dispatch(w irc.ResponseWriter, r *irc.Request) bool {
	// We are only interested in requests with the correct prefix.
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
	if cmd.Restricted && !s.authenticate(r.SenderMask) {
		proto.PrivMsg(w, r.SenderName, TextAccessDenied, cmd.Name)
		return false
	}

	// Ensure we have enough parameters.
	if cmd.RequiredParamCount() > len(args) {
		proto.PrivMsg(w, r.SenderName, TextMissingParameters, cmd.Name)
		return false
	}

	var params ParamList

	// Process and validate each parameter value.
	if len(cmd.Params) > 0 {
		params = make(ParamList, 0, len(cmd.Params))

		for i := 0; i < len(args) && i < len(cmd.Params); i++ {
			if cmd.Params[i].validate(args[i]) {
				params = append(params, Param{Value: args[i]})
				continue
			}

			proto.PrivMsg(w, r.SenderName, TextInvalidParameter,
				cmd.Name, cmd.Params[i].Name)
			return false
		}
	}

	go func() {
		// Ensure command handlers don't bring the entire bot down
		// when a panic occurs.
		defer func() {
			x := recover()
			if x != nil {
				// Go runtime errors should not be intercepted.
				if re, ok := x.(runtime.Error); ok {
					panic(re)
				}

				log.Printf("Command error: %v", x)
				log.Printf("> %#v", r)
			}
		}()

		cmd.Handler(w, r, params)
	}()

	return true
}

// Bind binds the given command.
func (s *Set) Bind(name string, restricted bool, handler Handler) *Command {
	cmd := newCommand(name, restricted, handler)
	s.data = append(s.data, cmd)
	sort.Sort(s.data)
	return cmd
}

// Unbind removes the given command.
func (s *Set) Unbind(name string) {
	idx := s.data.Index(name)
	if idx > -1 {
		copy(s.data[idx:], s.data[idx+1:])
		s.data[len(s.data)-1] = nil
		s.data = s.data[:len(s.data)-1]
	}
}

// split splits the given string into a command name and individual
// parameters. It ensures there are no empty entries from the parameter list.
func split(data string) (string, []string) {
	set := strings.Fields(data)

	for i := 0; i < len(set); i++ {
		if len(set[i]) == 0 {
			copy(set[i:], set[i+1:])
			set = set[:len(set)-1]
			i--
		}
	}

	if len(set) == 0 {
		return "", nil
	}

	return set[0], set[1:]
}
