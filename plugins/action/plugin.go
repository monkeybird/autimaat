// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package action binds action commands. These are things like:
//
//    <steve> !beer
//    * bot hands steve a cold beer.
//
package action

import (
	"math/rand"
	"time"

	"github.com/monkeybird/autimaat/app/util"
	"github.com/monkeybird/autimaat/irc"
	"github.com/monkeybird/autimaat/irc/cmd"
	"github.com/monkeybird/autimaat/irc/proto"
	"github.com/monkeybird/autimaat/plugins"
)

func init() { plugins.Register(&plugin{}) }

type plugin struct {
	cmd *cmd.Set
	rng *rand.Rand
}

// Load initializes the module and loads any internal resources
// which may be required.
func (p *plugin) Load(prof irc.Profile) error {
	p.cmd = cmd.New(prof.CommandPrefix(), nil)
	p.rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	// action returns a command handler which presents a channel with
	// a random string from the given list.
	action := func(set []string) cmd.Handler {
		return func(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
			targ := r.SenderName

			if params.Len() > 0 {
				targ = params.String(0)
			}

			idx := p.rng.Intn(len(set))
			msg := util.Action(set[idx], targ)
			proto.PrivMsg(w, r.Target, msg)
		}
	}

	// Bind all known actions.
	for _, a := range TextActions {
		handler := action(a.Answers)

		for _, name := range a.Names {
			p.cmd.Bind(name, false, handler).
				Add(TextUserName, false, cmd.RegAny)
		}
	}

	return nil
}

// Unload cleans the module up and unloads any internal resources.
func (p *plugin) Unload(prof irc.Profile) error {
	return nil
}

// Dispatch sends the given, incoming IRC message to the plugin for
// processing as it sees fit.
func (p *plugin) Dispatch(w irc.ResponseWriter, r *irc.Request) {
	p.cmd.Dispatch(w, r)
}
