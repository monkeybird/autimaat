// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package mod defines the interface for command modules. There are
// essentially plugins which provide specific functionality for the bot.
package mod

import (
	"monkeybird/irc"
	"monkeybird/irc/cmd"
)

// Module defines the interface for a single module.
type Module interface {
	// Load initializes the module and loads any internal resources
	// which may be required.
	Load(irc.ProtocolBinder, irc.Profile)

	// Unload cleans the module up and unloads any internal resources.
	Unload(irc.ProtocolBinder, irc.Profile)

	// Help is called whenever the module is to provide help on any
	// custom commands it may implement.
	Help(irc.ResponseWriter, *cmd.Request)
}
