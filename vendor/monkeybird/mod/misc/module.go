// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package misc defines a bunch of random, silly things in command form.
package misc

import (
	"fmt"
	"math/rand"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/irc/proto"
	"monkeybird/mod"
	"monkeybird/text"
	"monkeybird/tr"
	"time"
)

type module struct {
	rng      *rand.Rand
	commands *cmd.Set
}

// New returns a new module.
func New() mod.Module {
	return &module{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Load loads module resources and binds commands.
func (m *module) Load(pb irc.ProtocolBinder, prof irc.Profile) {
	pb.Bind("PRIVMSG", m.onPrivMsg)

	m.commands = cmd.New(prof.CommandPrefix(), nil)
	m.commands.Bind(tr.EightballName, tr.EightballDesc, false, m.cmd8Ball).
		Add(tr.EightballQuestionName, tr.EightballQuestionDesc, true, cmd.RegAny)

	m.bindGiveAction(tr.SmokeName, tr.SmokeAnswers)
	m.bindGiveAction(tr.BeerName, tr.BeerAnswers)
	m.bindGiveAction(tr.WineName, tr.WineAnswers)
	m.bindGiveAction(tr.CoffeeName, tr.CoffeeAnswers)
	m.bindGiveAction(tr.TeaName, tr.TeaAnswers)
	m.bindGiveAction(tr.LemonadeName, tr.LemonadeAnswers)
	m.bindGiveAction(tr.HugName, tr.HugAnswers)
	m.bindGiveAction(tr.PetName, tr.PetAnswers)
}

// bindGiveAction provides a wrapper for the binding of a "X gives Y to Z"
// action commands. We bind quite a few of these, so the shortcut makes this
// process less verbose.
func (m *module) bindGiveAction(name string, set []string) {
	desc := fmt.Sprintf(tr.GiveDesc, name)
	m.commands.Bind(name, desc, false, m.actionCommand(set)).
		Add(tr.GiveUserName, tr.GiveUserDesc, false, cmd.RegAny)
}

// Unload cleans up library resources and unbinds commands.
func (m *module) Unload(pb irc.ProtocolBinder, prof irc.Profile) {
	m.commands.Clear()
	pb.Unbind("PRIVMSG", m.onPrivMsg)
}

func (m *module) Help(w irc.ResponseWriter, r *cmd.Request) {
	m.commands.HelpHandler(w, r)
}

func (m *module) onPrivMsg(w irc.ResponseWriter, r *irc.Request) {
	m.commands.Dispatch(w, r)
}

// cmd8Ball asks the 8ball a question and presents the answer.
func (m *module) cmd8Ball(w irc.ResponseWriter, r *cmd.Request) {
	idx := m.rng.Intn(len(tr.EightBallAnswers))
	proto.PrivMsg(w, r.Target, tr.EightBallAnswers[idx], r.SenderName)
}

// actionCommand returns a command handler which presents a channel with
// a random string from the given list.
func (m *module) actionCommand(set []string) func(w irc.ResponseWriter, r *cmd.Request) {
	return func(w irc.ResponseWriter, r *cmd.Request) {
		targ := r.SenderName

		if r.Len() > 0 {
			targ = r.String(0)
		}

		idx := m.rng.Intn(len(set))
		msg := text.Action(set[idx], targ)
		proto.PrivMsg(w, r.Target, msg)
	}
}
