// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package proto defines convenience functions for IRC protocol requests.
package proto

import (
	"fmt"
	"io"
	"strings"

	"github.com/monkeybird/autimaat/irc"
)

// ref: https://en.wikipedia.org/wiki/List_of_Internet_Relay_Chat_commands#User_commands

// Raw sends the given, raw message data.
//
// The message being sent is reformatted to match the IRC specification.
// Meaning that it can not exceed 512 bytes and must end with `\r\n`.
// Any data beyond 512 bytes is simply discarded.
func Raw(w io.Writer, msg string, argv ...interface{}) error {
	data := []byte(fmt.Sprintf(msg, argv...) + "\r\n")
	sz := len(data)

	if sz <= 2 {
		return nil
	}

	if sz >= 512 {
		data = data[:512]

		if data[510] != '\r' {
			data[510] = '\r'
		}

		if data[511] != '\n' {
			data[511] = '\n'
		}
	}

	_, err := w.Write(data)
	return err
}

// Admin instructs the server to return information about the administrator of
// the server specified by <server>. If omitted, the current server is
// assumed.
func Admin(w io.Writer, server ...string) error {
	if len(server) > 0 {
		return Raw(w, "ADMIN %s", server[0])
	}
	return Raw(w, "ADMIN")
}

// Away marks us as being away, provided there is an away message.
// If the away message is empty, the away status is removed.
func Away(w io.Writer, message ...string) error {
	if len(message) > 0 {
		return Raw(w, "AWAY %s", message[0])
	}
	return Raw(w, "AWAY")
}

// CNotice sends a channel NOTICE message to <nickname> on <channel> that
// bypasses flood protection limits. The target nickname must be in the same
// channel as the client issuing the command, and the client must be a
// channel operator.
//
// Normally an IRC server will limit the number of different targets a client
// can send messages to within a certain time frame to prevent spammers or
// bots from mass-messaging users on the network, however this command can be
// used by channel operators to bypass that limit in their channel. For example,
// it is often used by help operators that may be communicating with a large
// number of users in a help channel at one time.
//
// This command is not formally defined in an RFC, but is in use by some IRC
// networks. Support is indicated in a RPL_ISUPPORT reply (numeric 005) with
// the CNOTICE keyword.
func CNotice(w io.Writer, nickname, channel, message string) error {
	return Raw(w, "CNOTICE %s %s :%s", nickname, channel, message)
}

// CPrivMsg sends a private message to <nickname> on <channel> that bypasses
// flood protection limits. The target nickname must be in the same channel as
// the client issuing the command, and the client must be a channel operator.
//
// Normally an IRC server will limit the number of different targets a client
// can send messages to within a certain time frame to prevent spammers or bots
// from mass-messaging users on the network, however this command can be used
// by channel operators to bypass that limit in their channel. For example, it
// is often used by help operators that may be communicating with a large
// number of users in a help channel at one time.
//
// This command is not formally defined in an RFC, but is in use by some IRC
// networks. Support is indicated in a RPL_ISUPPORT reply (numeric 005) with
// the CPRIVMSG keyword.
func CPrivMsg(w io.Writer, nickname, channel, message string) error {
	return Raw(w, "CPRIVMSG %s %s :%s", nickname, channel, message)
}

// Connect instructs the server <remote server> (or the current server, if
// <remote server> is omitted) to connect to <target server> on port <port>.
//
// This command should only be available to IRC Operators.
func Connect(w io.Writer, targetServer string, port int, remoteServer ...string) error {
	if len(remoteServer) > 0 {
		return Raw(w, "CONNECT %s %d %s", targetServer, port, remoteServer[0])
	}
	return Raw(w, "CONNECT %s %d", targetServer, port)
}

// Die instructs the server to shut down and may only be issued by
// IRC server operators.
func Die(w io.Writer) error { return Raw(w, "DIE") }

// Info requests information about the target server, or the current server if
// <server> is omitted. Information returned includes the server's version,
// when it was compiled, the patch level, when it was started, and any other
// information which may be considered to be relevant.
func Info(w io.Writer, server ...string) error {
	if len(server) > 0 {
		return Raw(w, "INFO %s", server[0])
	}
	return Raw(w, "INFO")
}

// Invite invites <nickname> to <channel>. <channel> does not have to exist,
// but if it does, only members of the channel are allowed to invite other
// clients. If the channel mode i is set, only channel operators may invite
// other clients.
func Invite(w io.Writer, nickname, channel string) error {
	return Raw(w, "INVITE %s %s", nickname, channel)
}

// IsOn requests the server to see if the nicknames in the given list are
// currently on the network. The server returns only the nicknames which are on
// the network in a space-separated list. If none of the clients are on the
// network, it returns an empty list.
func IsOn(w io.Writer, nicknames ...string) error {
	return Raw(w, "ISON %s", strings.Join(nicknames, " "))
}

// Join joins the given channels.
func Join(w io.Writer, channels ...irc.Channel) (err error) {
	for _, ch := range channels {
		if err = Raw(w, "chanserv INVITE %s", ch.Name); err != nil {
			return
		}

		if len(ch.Key) > 0 {
			err = Raw(w, "JOIN %s %s", ch.Name, ch.Key)
		} else {
			err = Raw(w, "JOIN %s", ch.Name)
		}

		if err != nil {
			return
		}

		if len(ch.Password) > 0 {
			err = PrivMsg(w, "chanserv", "IDENTIFY %s %s", ch.Name, ch.Password)
			if err != nil {
				return
			}
		}
	}

	return
}

// Kick forcibly removes <client> from <channel>. This command may only be
// issued by channel operators. The optional reason tells the client why
// they were kicked.
func Kick(w io.Writer, channel, client string, reason ...string) error {
	if len(reason) > 0 {
		return Raw(w, "KICK %s %s :%s", channel, client, reason[0])
	}
	return Raw(w, "KICK %s %s", channel, client)
}

// Knock sends a NOTICE to an invitation-only <channel> with an optional
// <message>, requesting an invite.
//
// This command is not formally defined by an RFC, but is supported by most
// major IRC daemons. Support is indicated in a RPL_ISUPPORT reply (numeric 005)
// with the KNOCK keyword.
func Knock(w io.Writer, channel string, message ...string) error {
	if len(message) > 0 {
		return Raw(w, "KNOCK %s :%s", channel, message[0])
	}
	return Raw(w, "KNOCK %s", channel)
}

// List requests all channels on the server. If the comma-separated list
// <channels> is given, it will return the channel topics.
func List(w io.Writer, channels ...string) error {
	if len(channels) > 0 {
		return Raw(w, "LIST %s", strings.Join(channels, ","))
	}
	return Raw(w, "LIST")
}

// Mode changes the mode for the given user or channel.
// Optionally with the given argument.
func Mode(w io.Writer, target, mode string, argv ...string) error {
	if len(argv) > 0 {
		return Raw(w, "MODE %s %s %s", target, mode, argv[0])
	}
	return Raw(w, "MODE %s %s", target, mode)
}

// Names queries users in the given list of <channels>, If <channels> is
// omitted, all users are shown, grouped by channel name with all users who are
// not on a channel being shown as part of channel "*". If <server> is specified,
// the command is sent to <server> for evaluation.
//
// The response contains all nicknames in the channel, prefixed with the highest
// channel status prefix of that user, for example like this (with @ being the
// highest status prefix).
//
//     :irc.server.net 353 Phyre = #SomeChannel :@WiZ
//
// If a client wants to receive all the channel status prefixes of a user and
// not only their current highest one, the IRCv3 multi-prefix extension can
// be enabled (@ is the channel operator prefix, and + the lower voice status
// prefix):
//
//     :irc.server.net 353 Phyre = #SomeChannel :@+WiZ
//
func Names(w io.Writer, channels ...string) error {
	if len(channels) > 0 {
		return Raw(w, "NAMES %s", strings.Join(channels, ","))
	}
	return Raw(w, "NAMES")
}

// Nick allows a client to change their IRC nickname.
// The optional password is used to authenticate the user with nickserv.
func Nick(w io.Writer, nickname string, password ...string) error {
	err := Raw(w, "NICK %s", nickname)
	if err != nil {
		return err
	}

	if len(password) > 0 {
		return PrivMsg(w, "nickserv", "IDENTIFY %s", password[0])
	}

	return nil
}

// Notice works similarly to PRIVMSG, except automatic replies must never be
// sent in reply to NOTICE messages.
func Notice(w io.Writer, target, f string, argv ...interface{}) error {
	return Raw(w, "NOTICE %s :%s", target, fmt.Sprintf(f, argv...))
}

// Oper authenticates a user as an IRC operator on a server/network.
func Oper(w io.Writer, nickname, password string) error {
	return Raw(w, "OPER %s %s", nickname, password)
}

// Part leaves the given channels.
func Part(w io.Writer, channels ...irc.Channel) (err error) {
	for _, ch := range channels {
		err = Raw(w, "PART %s :", ch.Name)
		if err != nil {
			return
		}
	}

	return
}

// Pass sets a connection password. This command must be sent before the
// NICK/USER registration combination. It is ignored if the given password
// is empty.
func Pass(w io.Writer, password string) error {
	if len(password) == 0 {
		return nil
	}
	return Raw(w, "PASS %s", password)
}

// Pong sends the given payload as a response to a PING message.
func Pong(w io.Writer, payload string) error {
	return Raw(w, "PONG %s", payload)
}

// PrivMsg sends the specified formatted message to the given target.
// The target may be a channel or nickname.
func PrivMsg(w io.Writer, target, f string, argv ...interface{}) error {
	return Raw(w, "PRIVMSG %s :%s", target, fmt.Sprintf(f, argv...))
}

// Quit disconnects from the server., optionally with the given message.
func Quit(w io.Writer, message ...string) error {
	if len(message) > 0 {
		return Raw(w, "QUIT %s", message[0])
	}
	return Raw(w, "QUIT")
}

// Recover attempts to re-authenticate our username, so we can
// regain the use of it. This is mostly useful after we received
// a NickInUse error and is only relevant if there is a nickserv.
func Recover(w io.Writer, nickname, password string) error {
	return Raw(w, "NS RECOVER %s %s", nickname, password)
}

// Rehash causes the server to re-read and re-process its configuration file(s).
// This command can only be sent by IRC Operators
func Rehash(w io.Writer, username string) error { return Raw(w, "REHASH") }

// Restart restarts a server. It may only be sent by IRC Operators.
func Restart(w io.Writer, username string) error { return Raw(w, "RESTART") }

// SQuit causes <server> to quit the network.
func SQuit(w io.Writer, server, message string) error {
	return Raw(w, "SQUIT %s %s", server, message)
}

// SetName allows a client to change the "real name" specified when registering
// a connection.
//
// This command is not formally defined by an RFC, but is in use by some IRC
// daemons. Support is indicated in a RPL_ISUPPORT reply (numeric 005) with the
// SETNAME keyword
func SetName(w io.Writer, name string) error {
	return Raw(w, "SETNAME %s", name)
}

// Silence adds or removes a host mask to a server-side ignore list that
// prevents matching users from sending the client messages. More than one mask
// may be specified. Each item prefixed with a "+" or "-" to designate whether
// it is being added or removed. Sending the command with no parameters returns
// the entries in the client's ignore list.
//
// This command is not formally defined in an RFC, but is supported by most
// major IRC daemons. Support is indicated in a RPL_ISUPPORT reply (numeric 005)
// with the SILENCE keyword and the maximum number of entries a client may have
// in its ignore list. For example:
//
//    :irc.server.net 005 WiZ WALLCHOPS WATCH=128 SILENCE=15 MODES=12 CHANTYPES=#
//
func Silence(w io.Writer, masks ...string) error {
	return Raw(w, "SILENCE %s", strings.Join(masks, " "))
}

// Summon gives users who are on the same host as <server> a message asking
// them to join IRC. If server is omitted, the current server is assumed.
// Channel is optional and will request them to join that specific channel.
func Summon(w io.Writer, user, server, channel string) error {
	if len(server) > 0 {
		if len(channel) > 0 {
			return Raw(w, "SUMMON %s %s %s", user, server, channel)
		}
		return Raw(w, "SUMMON %s %s", user, server)
	}
	return Raw(w, "SUMMON %s", user)
}

// Time requests the local time on the current or given server.
func Time(w io.Writer, server ...string) error {
	if len(server) > 0 {
		return Raw(w, "TIME %s", server[0])
	}
	return Raw(w, "TIME")
}

// Topic allows the client to query or set the channel topic on <channel>.
// If channel mode +t is set, only a channel operator may set the topic.
func Topic(w io.Writer, channel string, topic ...string) error {
	if len(topic) > 0 {
		return Raw(w, "TOPIC %s %s", channel, topic[0])
	}
	return Raw(w, "TOPIC %s", channel)
}

// User is used at the beginning of a connection to specify the username,
// hostname, real name and initial user modes of the connecting client.
// <realname> may contain spaces.
//
//     E.g.: USER joe 8 * :joe smith
//
func User(w io.Writer, username, mode, realname string) error {
	return Raw(w, "USER %s %s * :%s", username, mode, realname)
}

// UserHost returns host information for the  specified nicknames.
func UserHost(w io.Writer, names ...string) error {
	return Raw(w, "USERHOST %s", strings.Join(names, " "))
}

// UserIP requests the direct IP address of the user with the specified nickname.
//
// This command is often used to obtain the IP of an abusive user to more
// effectively perform a ban. It is unclear what, if any, privileges are
// required to execute this command on a server.
//
// This command is not formally defined by an RFC, but is in use by some IRC
// daemons. Support is indicated in a RPL_ISUPPORT reply (numeric 005) with
// the USERIP keyword.
func UserIP(w io.Writer, nickname string) error {
	return Raw(w, "USERIP %s", nickname)
}

// Users requests a list of users and information about those users in a
// format similar to the UNIX commands who, rusers and finger. The command
// is optionally targeted at a specific server.
func Users(w io.Writer, server ...string) error {
	if len(server) > 0 {
		return Raw(w, "USERS %s", server[0])
	}
	return Raw(w, "USERS")
}

// Version requests version information for the current or given server.
func Version(w io.Writer, server ...string) error {
	if len(server) > 0 {
		return Raw(w, "VERSION %s", server[0])
	}
	return Raw(w, "VERSION")
}

// Wallops sends a formatted message to all operators connected to the server
// or all users with user mode 'w' set.
func Wallops(w io.Writer, f string, argv ...interface{}) error {
	return Raw(w, "WALLOPS %%s", fmt.Sprintf(f, argv...))
}

// Watch adds or removes a user to a client's server-side friends list.
// More than one nickname may be specified. Each item prefixed with a "+" or "-"
// to designate whether it is being added or removed. Sending the command
// with no parameters returns the entries in the client's friends list.
//
// This command is not formally defined in an RFC, but is supported by most
// major IRC daemons. Support is indicated in a RPL_ISUPPORT reply (numeric 005)
// with the WATCH keyword and the maximum number of entries a client may have in
// its friends list. For example:
//
//     :irc.server.net 005 WiZ WALLCHOPS WATCH=128 SILENCE=15 MODES=12 CHANTYPES=#
//
func Watch(w io.Writer, masks ...string) error {
	return Raw(w, "WATCH %s", strings.Join(masks, " "))
}

// Who requests a list of users who match <name>. If opOnly is truen, the
// server will only return information about IRC Operators.
func Who(w io.Writer, name string, opOnly bool) error {
	if opOnly {
		return Raw(w, "WHO %s o", name)
	}
	return Raw(w, "WHO %s", name)
}

// Whois requests information about the given nickname. If <server>
// is given, the command is forwarded to it for processing.
func Whois(w io.Writer, nickname string, server ...string) error {
	if len(server) > 0 {
		return Raw(w, "WHOIS %s %s", server[0], nickname)
	}
	return Raw(w, "WHOIS %s", nickname)
}

// Whowas requests information about a nickname that is no longer in use
// (due to client disconnection, or nickname changes). If <server> is given,
// the command is forwarded to it for processing.
func Whowas(w io.Writer, target string, server ...string) error {
	if len(server) > 0 {
		return Raw(w, "WHOWAS %s %s", server[0], target)
	}
	return Raw(w, "WHOWAS %s", target)
}
