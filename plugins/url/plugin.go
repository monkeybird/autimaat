// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package url defines a plugin, which finds and extracts URLs from
// incoming chat data. It performs a HTTP lookup to the found URL and
// attempts to determine the page title of the link. This title is then
// returned to the channel from which the message came.
package url

import (
	"github.com/monkeybird/autimaat/irc"
	"github.com/monkeybird/autimaat/plugins"
	"github.com/monkeybird/autimaat/util"
)

func init() { plugins.Register(&plugin{}) }

type plugin struct {
	data struct {
		YoutubeApiKey string
	}
}

// Load initializes the module and loads any internal resources
// which may be required.
func (p *plugin) Load(prof irc.Profile) error {
	return util.ReadFile("url.cfg", &p.data, false)
}

// Unload cleans the module up and unloads any internal resources.
func (p *plugin) Unload(prof irc.Profile) error {
	p.data.YoutubeApiKey = ""
	return nil
}

// Dispatch sends the given, incoming IRC message to the plugin for
// processing as it sees fit.
func (p *plugin) Dispatch(w irc.ResponseWriter, r *irc.Request) {
	if !r.IsPrivMsg() {
		return
	}

	// Find all URLs in the message body.
	list := regUrl.FindAllString(r.Data, -1)
	if len(list) == 0 {
		return
	}

	// Fetch title data for each of them.
	for _, url := range list {
		go fetchTitle(w, r, url, p.data.YoutubeApiKey)
	}
}
