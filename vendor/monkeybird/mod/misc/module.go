// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package misc defines a bunch of random, silly things in command form.
package misc

import (
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

	m.commands.Bind(tr.BeerName, tr.BeerDesc, false, m.actionCommand(tr.BeerAnswers)).
		Add(tr.GiveUserName, tr.GiveUserDesc, false, cmd.RegAny)

	m.commands.Bind(tr.WineName, tr.WineDesc, false, m.actionCommand(tr.WineAnswers)).
		Add(tr.GiveUserName, tr.GiveUserDesc, false, cmd.RegAny)

	m.commands.Bind(tr.CoffeeName, tr.CoffeeDesc, false, m.actionCommand(tr.CoffeeAnswers)).
		Add(tr.GiveUserName, tr.GiveUserDesc, false, cmd.RegAny)

	m.commands.Bind(tr.TeaName, tr.TeaDesc, false, m.actionCommand(tr.TeaAnswers)).
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
		idx := m.rng.Intn(len(set))
		msg := text.Action(set[idx], r.SenderName)
		proto.PrivMsg(w, r.Target, msg)
	}
}
