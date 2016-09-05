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

	// bind provides a wrapper for the binding of a "X gives Y to Z" action
	// commands. We bind quite a few of these, so the shortcut makes this
	// process less verbose.
	bind := func(name string, set []string) {
		p.cmd.Bind(name, false, action(set)).
			Add(TextGiveUserName, false, cmd.RegAny)
	}

	// Bind a bunch of action commands.
	bind(TextSmokeName, TextSmokeAnswers)
	bind(TextBeerName, TextBeerAnswers)
	bind(TextWineName, TextWineAnswers)
	bind(TextPortName, TextPortAnswers)
	bind(TextWhiskeyName, TextWhiskeyAnswers)
	bind(TextCoffeeName, TextCoffeeAnswers)
	bind(TextTeaName, TextTeaAnswers)
	bind(TextLemonadeName, TextLemonadeAnswers)
	bind(TextHugName, TextHugAnswers)
	bind(TextPetName, TextPetAnswers)
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
