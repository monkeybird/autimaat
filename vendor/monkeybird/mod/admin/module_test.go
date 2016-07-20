// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package admin

import (
	"fmt"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/test"
	"monkeybird/text"
	"monkeybird/tr"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	bindings irc.BindingList
	profile  irc.Profile
	tm       *module
)

func TestMain(m *testing.M) {
	profile = irc.NewProfile(".")

	tm = New(test.Nickname, test.VersionMajor,
		test.VersionMinor, strconv.Itoa(test.VersionRevision)).(*module)

	tm.Load(&bindings, profile)
	ret := m.Run()
	tm.Unload(&bindings, profile)

	// The auth/deauth commands save the profile to disk.
	// Get rid of the file.
	os.Remove("profile.cfg")
	os.Exit(ret)
}

func TestVersion(t *testing.T) {
	var w test.MockWriter

	tm.cmdVersion(&w, &cmd.Request{
		Request: test.NewRequest(),
	})

	stamp := time.Unix(test.VersionRevision, 0)

	w.Verify(t, fmt.Sprintf("PRIVMSG %s :%s", test.ChannelName,
		fmt.Sprintf(tr.VersionDisplayText,
			test.SenderName,
			text.Bold(test.Nickname),
			text.Bold("%d.%d.%d", test.VersionMajor, test.VersionMinor, test.VersionRevision),
			stamp.Format(tr.DateFormat),
			stamp.Format(tr.TimeFormat))))
}

func TestAuthList(t *testing.T) {
	var w test.MockWriter

	tm.cmdAuthList(&w, &cmd.Request{
		Request: test.NewRequest(),
		Params: []cmd.Param{
			cmd.Param{Value: "1"},
		},
	})

	list := profile.Whitelist()
	w.Verify(t, fmt.Sprintf("PRIVMSG %s :%s", test.SenderName,
		fmt.Sprintf(tr.AuthListDisplayText, strings.Join(list, ", "))))
}

func TestAuthorize(t *testing.T) {
	var w test.MockWriter

	tm.cmdAuthorize(&w, &cmd.Request{
		Request: test.NewRequest(),
		Params: []cmd.Param{
			cmd.Param{Value: test.SenderMask},
		},
	})

	w.Verify(t, fmt.Sprintf("PRIVMSG %s :%s", test.SenderName,
		fmt.Sprintf(tr.AuthorizeDisplayText, test.SenderMask)))
}

func TestDeauthorize(t *testing.T) {
	var w test.MockWriter

	tm.cmdDeauthorize(&w, &cmd.Request{
		Request: test.NewRequest(),
		Params: []cmd.Param{
			cmd.Param{Value: test.SenderMask},
		},
	})

	w.Verify(t, fmt.Sprintf("PRIVMSG %s :%s", test.SenderName,
		fmt.Sprintf(tr.DeauthorizeDisplayText, test.SenderMask)))
}

func TestPart(t *testing.T) {
	const ChannelName = "#test"

	var w test.MockWriter
	tm.cmdPart(&w, &cmd.Request{
		Request: test.NewRequest(),
		Params: []cmd.Param{
			cmd.Param{Value: ChannelName},
		},
	})

	w.Verify(t, fmt.Sprintf("PART %s :", ChannelName))
}

func TestJoin(t *testing.T) {
	const ChannelName = "#test"
	const ChannelKey = "12345"
	const ChanservPassword = "somepassword"

	{
		var w test.MockWriter
		tm.cmdJoin(&w, &cmd.Request{
			Request: test.NewRequest(),
			Params: []cmd.Param{
				cmd.Param{Value: ChannelName},
			},
		})

		w.Verify(t,
			fmt.Sprintf("chanserv INVITE %s", ChannelName),
			fmt.Sprintf("JOIN %s", ChannelName))
	}

	{
		var w test.MockWriter
		tm.cmdJoin(&w, &cmd.Request{
			Request: test.NewRequest(),
			Params: []cmd.Param{
				cmd.Param{Value: ChannelName},
				cmd.Param{Value: ""},
				cmd.Param{Value: ChannelKey},
			},
		})

		w.Verify(t,
			fmt.Sprintf("chanserv INVITE %s", ChannelName),
			fmt.Sprintf("JOIN %s %s", ChannelName, ChannelKey),
		)
	}

	{
		var w test.MockWriter
		tm.cmdJoin(&w, &cmd.Request{
			Request: test.NewRequest(),
			Params: []cmd.Param{
				cmd.Param{Value: ChannelName},
				cmd.Param{Value: ChanservPassword},
				cmd.Param{Value: ChannelKey},
			},
		})

		w.Verify(t,
			fmt.Sprintf("chanserv INVITE %s", ChannelName),
			fmt.Sprintf("JOIN %s %s", ChannelName, ChannelKey),
			fmt.Sprintf("PRIVMSG chanserv :IDENTIFY %s %s",
				ChannelName, ChanservPassword),
		)
	}
}
