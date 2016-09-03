// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

/*
Package cmd allows the definition of command handlers to be called
by users from either a channel or a private message. These are messages
like the following:

   !join #test

Beginning with a predefined character (! in this case), followed by the
command name and an optional set of arguments. This package parses the
command data, verifies it refers to an existing command and ensures that
the parameter values have the correct formats.

For example:

	join := cmd.Bind("join", true, onJoin)
	join.Add("channel", true, cmd.RegChannel)
	join.Add("password", false, cmd.RegAny)

	...

	func onJoin(w irc.ResponseWriter, r *cmd.Request) {
		var c irc.Channel
		c.Name = r.String(0)

		if r.Len() > 1 {
			c.Password = r.String(1)
		}

		proto.Join(w, c)
	}

The name and description texts for the command and parameters, are there for
user documentation. You can bind a `!help` command to `cmd.HelpHandler`, which
will present the user either with an overview of all registered commands, or
get detailed help on a specific command.

The `cmd.RegXXX` values passed into the parameter definitions are predefined
regular expressions. You are free to pass in your own patterns. These are used
to ensure a parameter value given by a user, matches your expectations.
If this is not the case, the command handler is never executed and the user
is presented with an appropriate error response.

By the time the registered command handler is actually called, you may be
certain that the parameter value matches your definition.

The boolean value pass passed into each parameter definition determines if
that specific parameter is optional or not. This provides for rudimentary
varargs functionality.

In the example listed above, the user may call the `join` command in one
of two ways:

	!join #channel
	!join #channel somepassword

A parameter occupying multiple whitespace separated words, is to be supplied
in double quotes:

	!join #channel "some long password"

*/
package cmd
