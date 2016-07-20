// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package dictionary

import (
	"fmt"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/test"
	"monkeybird/text"
	"monkeybird/tr"
	"os"
	"testing"
)

var (
	bindings irc.BindingList
	profile  irc.Profile
	tm       *module
)

func TestMain(m *testing.M) {
	profile = irc.NewProfile(".")

	tm = New().(*module)
	tm.Load(&bindings, profile)
	ret := m.Run()
	tm.Unload(&bindings, profile)

	// This module writes a dictionary file to disk; get rid of it.
	os.Remove("dictionary.dat")
	os.Exit(ret)
}

func TestAddDefine(t *testing.T) {
	var w test.MockWriter

	tm.cmdAddDefine(&w, &cmd.Request{
		Request: test.NewRequest(),
		Params: []cmd.Param{
			cmd.Param{Value: "a"},
			cmd.Param{Value: "b"},
		},
	})

	w.Verify(t, fmt.Sprintf("PRIVMSG %s :%s",
		test.SenderName,
		fmt.Sprintf(tr.AddDefineDisplayText, text.Bold("a")),
	))

	w.Reset()
	tm.cmdAddDefine(&w, &cmd.Request{
		Request: test.NewRequest(),
		Params: []cmd.Param{
			cmd.Param{Value: "a"},
			cmd.Param{Value: "b"},
		},
	})

	w.Verify(t, fmt.Sprintf("PRIVMSG %s :%s",
		test.SenderName,
		fmt.Sprintf(tr.AddDefineAllreadyUsed, text.Bold("a")),
	))
}

func TestRemoveDefine(t *testing.T) {
	var w test.MockWriter

	tm.cmdRemoveDefine(&w, &cmd.Request{
		Request: test.NewRequest(),
		Params: []cmd.Param{
			cmd.Param{Value: "a"},
		},
	})

	w.Verify(t, fmt.Sprintf("PRIVMSG %s :%s",
		test.SenderName,
		fmt.Sprintf(tr.RemoveDefineDisplayText1, text.Bold("a")),
	))

	w.Reset()
	tm.cmdRemoveDefine(&w, &cmd.Request{
		Request: test.NewRequest(),
		Params: []cmd.Param{
			cmd.Param{Value: "a"},
		},
	})

	w.Verify(t, fmt.Sprintf("PRIVMSG %s :%s",
		test.SenderName,
		fmt.Sprintf(tr.RemoveDefineNotFound, text.Bold("a")),
	))
}

func TestDefine(t *testing.T) {
	testBadDefine(t, "")
	testBadDefine(t, "abc")
}

func testBadDefine(t *testing.T, v string) {
	var w test.MockWriter

	tm.cmdDefine(&w, &cmd.Request{
		Request: test.NewRequest(),
		Params: []cmd.Param{
			cmd.Param{Value: v},
		},
	})

	w.Verify(t, fmt.Sprintf("PRIVMSG %s :%s",
		test.ChannelName,
		fmt.Sprintf(tr.DefineNotFound, test.SenderName, text.Bold(v)),
	))
}
