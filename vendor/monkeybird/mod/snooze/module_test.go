// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package snooze

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
	rw test.MockWriter
	tm *module
)

func TestMain(m *testing.M) {
	var bindings irc.BindingList
	profile := irc.NewProfile(".")

	tm = New(&rw).(*module)

	tm.Load(&bindings, profile)
	ret := m.Run()
	tm.Unload(&bindings, profile)

	// This module writes snooze data to a file when it is unloaded.
	// Get rid of the file.
	os.Remove("snooze.dat")
	os.Exit(ret)
}

func TestValidUnsnooze(t *testing.T) {
	req := test.NewRequest()
	req.SenderMask = test.SenderMask

	id := tm.createID()
	tm.addSnooze(&rw, &cmd.Request{
		Request: req,
		Params:  []cmd.Param{{Value: "12:13"}},
	}, id)

	rw.Reset()
	tm.cmdUnsnooze(&rw, &cmd.Request{
		Request: req,
		Params:  []cmd.Param{{Value: id}},
	})

	rw.Verify(t, fmt.Sprintf("PRIVMSG %s :%s", test.ChannelName,
		fmt.Sprintf(tr.SnoozeAlarmUnset, test.SenderName)))
}

func TestInvalidUnsnooze(t *testing.T) {
	req := test.NewRequest()
	id := tm.createID()

	tm.addSnooze(&rw, &cmd.Request{
		Request: req,
		Params:  []cmd.Param{{Value: "12:13"}},
	}, id)

	rw.Reset()
	req.SenderMask = "i@dont.exist"

	tm.cmdUnsnooze(&rw, &cmd.Request{
		Request: req,
		Params:  []cmd.Param{{Value: id}},
	})

	rw.Verify(t)
}

func TestSnooze(t *testing.T) {
	testSnoozeInvalidTime(t, "")
	testSnoozeInvalidTime(t, "abc")
	testSnoozeInvalidTime(t, "0")
	testSnoozeInvalidTime(t, "24:00")
	testSnoozeInvalidTime(t, "-11:20")
	testSnoozeInvalidTime(t, "-20")

	testSnoozeValidTime(t, "11:20")
	testSnoozeValidTime(t, "30")
	testSnoozeValidTime(t, "+20")
}

func testSnoozeValidTime(t *testing.T, timeValue string) {
	rw.Reset()

	id := tm.createID()

	ok := tm.addSnooze(&rw, &cmd.Request{
		Request: test.NewRequest(),
		Params:  []cmd.Param{{Value: timeValue}},
	}, id)

	if !ok {
		tm.deleteID(id)
		t.Fatalf("snooze for value %q failed. It should have succeeded.", timeValue)
	}

	rw.Verify(t, fmt.Sprintf("PRIVMSG %s :%s", test.ChannelName,
		fmt.Sprintf(tr.SnoozeAlarmSet, test.SenderName, text.Bold(id))))
}

func testSnoozeInvalidTime(t *testing.T, timeValue string) {
	rw.Reset()

	id := tm.createID()
	defer tm.deleteID(id)

	ok := tm.addSnooze(&rw, &cmd.Request{
		Request: test.NewRequest(),
		Params:  []cmd.Param{{Value: timeValue}},
	}, id)

	if ok {
		t.Fatalf("snooze for value %q succeeded. It should have failed.", timeValue)
	}

	rw.Verify(t, fmt.Sprintf("PRIVMSG %s :%s", test.ChannelName,
		fmt.Sprintf(tr.SnoozeInvalidTime, test.SenderName, timeValue)))
}
