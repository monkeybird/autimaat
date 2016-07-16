// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package irc

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"sync"
)

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

// profile defines bot configuration data.
//
// The fields are embedded in a sub struct to differentiate them from the
// method names needed to qualify as a Profile interface. I would rather
// just make these field names lower case, but Go's JSON decoder will not
// work on non-exported fields. Thus breaking the Load/Save functionality.
type profile struct {
	m    sync.RWMutex
	root string
	data profileData
}

// profileData defines the parts of the profile which are saved to
// an external configuration file.
type profileData struct {
	Whitelist          []string
	Channels           []Channel
	Address            string
	TLSKey             string
	TLSCert            string
	CAPemData          string
	Nickname           string
	NickservPassword   string
	OperPassword       string
	ConnectionPassword string
	CommandPrefix      string
	WeatherApiKey      string
	YoutubeApiKey      string
	Logging            bool
}

// NewProfile creates a new profile for the given root directory.
func NewProfile(root string) Profile {
	return &profile{
		root: root,
		data: profileData{
			Logging:  false,
			Address:  "server.net:6667",
			Nickname: "bot_name",
			Channels: []Channel{
				{Name: "#test_channel"},
			},
			Whitelist: []string{
				"~user@server.com",
			},
			CommandPrefix: "!",
		},
	}
}

func (p *profile) WeatherApiKey() string {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.data.WeatherApiKey
}

func (p *profile) YoutubeApiKey() string {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.data.YoutubeApiKey
}

func (p *profile) Root() string {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.root
}

func (p *profile) ForkArgs() []string {
	p.m.RLock()
	defer p.m.RUnlock()
	return []string{p.root}
}

func (p *profile) Channels() []Channel {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.data.Channels
}

func (p *profile) Address() string {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.data.Address
}

func (p *profile) TLSKey() string {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.data.TLSKey
}

func (p *profile) TLSCert() string {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.data.TLSCert
}

func (p *profile) CAPemData() string {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.data.CAPemData
}

func (p *profile) Nickname() string {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.data.Nickname
}

func (p *profile) SetNickname(v string) {
	p.m.Lock()
	p.data.Nickname = v
	p.m.Unlock()
}

func (p *profile) NickservPassword() string {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.data.NickservPassword
}

func (p *profile) OperPassword() string {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.data.OperPassword
}

func (p *profile) ConnectionPassword() string {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.data.ConnectionPassword
}

func (p *profile) CommandPrefix() string {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.data.CommandPrefix
}

func (p *profile) Whitelist() []string {
	p.m.RLock()
	defer p.m.RUnlock()

	out := make([]string, len(p.data.Whitelist))
	copy(out, p.data.Whitelist)
	return out
}

func (p *profile) WhitelistAdd(mask string) {
	p.m.Lock()

	for _, str := range p.data.Whitelist {
		if strings.EqualFold(str, mask) {
			p.m.Unlock()
			return
		}
	}

	p.data.Whitelist = append(p.data.Whitelist, mask)
	p.m.Unlock()
	p.Save()
}

func (p *profile) WhitelistRemove(mask string) {
	p.m.Lock()

	for i, str := range p.data.Whitelist {
		if !strings.EqualFold(str, mask) {
			continue
		}

		copy(p.data.Whitelist[i:], p.data.Whitelist[i+1:])
		p.data.Whitelist = p.data.Whitelist[:len(p.data.Whitelist)-1]
		break
	}

	p.m.Unlock()
	p.Save()
}

func (p *profile) IsWhitelisted(mask string) bool {
	p.m.RLock()
	defer p.m.RUnlock()

	for _, str := range p.data.Whitelist {
		if strings.EqualFold(str, mask) {
			return true
		}
	}

	return false
}

func (p *profile) IsNick(name string) bool {
	p.m.RLock()
	defer p.m.RUnlock()
	return strings.EqualFold(p.data.Nickname, name)
}

func (p *profile) Logging() bool {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.data.Logging
}

func (p *profile) SetLogging(v bool) {
	p.m.Lock()
	p.data.Logging = v
	p.m.Unlock()
	p.Save()
}

func (p *profile) Save() error {
	p.m.RLock()
	defer p.m.RUnlock()

	data, err := json.MarshalIndent(p.data, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile("profile.cfg", data, 0600)
}

func (p *profile) Load() error {
	p.m.Lock()
	defer p.m.Unlock()

	data, err := ioutil.ReadFile("profile.cfg")
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &p.data)
}
