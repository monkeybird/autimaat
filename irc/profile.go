// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package irc defines some utility types and functions for an IRC bot.
package irc

// Profile defines bot configuration data.
type Profile interface {
	// Root defines the root directory with the bot's configuration data.
	Root() string

	// Channels yields all channels the bot should join on startup.
	Channels() []Channel

	// Address defines the host and port of the server/network to connect to.
	Address() string

	// TKSKey defines the TLS key file. Along with TLSCert, it defines the
	// components needed to establish an encrypted connection to the server.
	TLSKey() string

	// TLSCert defines the TLS certificate file. Along with TLSKey, it
	// defines the components needed to establish an encrypted connection
	// to the server.
	TLSCert() string

	// CAPemData defines one ore more, PEM encoded, server root certificates.
	// This is optional and is used to replace the client's existing root CA
	// pool. This can be useful if you are connecting to a server whos
	// certificate is not present in any system wide CA pools.
	CAPemData() string

	// Nickname yields the bot's nickname.
	Nickname() string

	// SetNickname sets the bot's nickname. This is generally only called
	// when the bot logs in and finds its name alredy in use. If the nick
	// can not be regained, this function is used to alter it to something
	// which is still available.
	SetNickname(string)

	// NickservPassword defines the bot's nickserv password. This will be
	// used to register the bot when it logs in. It is only relevant if the
	// bot has a registered nickname and nickserv exists on the server.
	NickservPassword() string

	// SetNickservPassword sets the bot's nickserv password. This is be
	// used to register the bot when it logs in. It is only relevant if the
	// bot has a registered nickname and nickserv exists on the server.
	SetNickservPassword(string)

	// OperPassword defines the bot's OPER password. If present, this will
	// register the bot as a server operator.
	OperPassword() string

	// Some connections may be secured and require a password to connect to.
	ConnectionPassword() string

	// CommandPrefix this is the prefix used for all bot commands. Whenever
	// the bot reads incoming PRIVMSG data, it looks for this prefix to
	// determine if a command call was issued or not.
	CommandPrefix() string

	// Save saves the profile to disk.
	Save() error

	// Load loads the profile from disk.
	Load() error

	// IsWhitelisted returns true if the given hostmask is in the whitelist.
	// This means the user to whom it belongs is allowed to execute restricted
	// commands. This performs a case-insensitive comparison.
	IsWhitelisted(string) bool

	// Whitelist returns a copy of the current whitelist.
	Whitelist() []string

	// WhitelistAdd adds the given hostmask to the whitelist,
	// provided it does not already exist.
	WhitelistAdd(string)

	// WhitelistRemove removes the given hostmask from the whitelist,
	// provided it exists.
	WhitelistRemove(string)

	// IsNick returns true if the given name equals the bot's nickname.
	// This is used in request handlers to quickly check if a request
	// is targeted specifically at this bot or not.
	IsNick(string) bool

	// ForkArgs returns a list of command line arguments which should be
	// passed to a forked child process.
	ForkArgs() []string

	// Logging returns true if incoming data logging is enabled.
	Logging() bool

	// Logging determines if logging of incoming data should be enabled or not.
	SetLogging(bool)

	// WeatherApiKey returns the API key for openweathermap.org.
	WeatherApiKey() string

	// YoutubeApiKey returns the API key for youtube.
	YoutubeApiKey() string
}
