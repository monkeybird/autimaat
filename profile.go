// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"strings"
	"sync"

	"github.com/monkeybird/autimaat/irc"
	"github.com/monkeybird/autimaat/util"
)

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
	Channels           []irc.Channel
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
func NewProfile(root string) irc.Profile {
	return &profile{
		root: root,
		data: profileData{
			Logging:  false,
			Address:  "server.net:6667",
			Nickname: "bot_name",
			Channels: []irc.Channel{
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

func (p *profile) Channels() []irc.Channel {
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
	p.Save()
}

func (p *profile) NickservPassword() string {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.data.NickservPassword
}

func (p *profile) SetNickservPassword(v string) {
	p.m.Lock()
	p.data.NickservPassword = v
	p.m.Unlock()
	p.Save()
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
	err := util.WriteFile("profile.cfg", p.data, false)
	p.m.RUnlock()
	return err
}

func (p *profile) Load() error {
	p.m.Lock()
	err := util.ReadFile("profile.cfg", &p.data, false)
	p.m.Unlock()
	return err
}
