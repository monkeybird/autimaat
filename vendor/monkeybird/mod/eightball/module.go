// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package eightball implements the magic 8ball.
package eightball

import (
	"math/rand"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/irc/proto"
	"monkeybird/mod"
	"monkeybird/tr"
	"time"
)

// answers defines the list of possible 8ball answers.
var answers = []string{
	tr.Eightball1,
	tr.Eightball2,
	tr.Eightball3,
	tr.Eightball4,
	tr.Eightball5,
	tr.Eightball6,
	tr.Eightball7,
	tr.Eightball8,
	tr.Eightball9,
	tr.Eightball10,
	tr.Eightball11,
	tr.Eightball12,
	tr.Eightball13,
	tr.Eightball14,
	tr.Eightball15,
	tr.Eightball16,
	tr.Eightball17,
	tr.Eightball18,
	tr.Eightball19,
	tr.Eightball20,
}

type module struct {
	commands *cmd.Set
}

// New returns a new eightball module.
func New() mod.Module {
	return &module{}
}

// Load loads module resources and binds commands.
func (m *module) Load(pb irc.ProtocolBinder, prof irc.Profile) {
	pb.Bind("PRIVMSG", m.onPrivMsg)

	m.commands = cmd.New(prof.CommandPrefix(), nil)
	m.commands.Bind(tr.EightballName, tr.EightballDesc, false, m.cmdAsk).
		Add(tr.EightballQuestionName, tr.EightballQuestionDesc, true, cmd.RegAny)
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

// cmdAsk asks the 8ball a question and presents the answer.
func (m *module) cmdAsk(w irc.ResponseWriter, r *cmd.Request) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := rnd.Intn(len(answers))
	proto.PrivMsg(w, r.Target, answers[idx], r.SenderName)
}
