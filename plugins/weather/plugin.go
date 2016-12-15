// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package weather provides commands to do current weather lookups,
// as well as weather forecasts for specific locations.
package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/monkeybird/autimaat/app/util"
	"github.com/monkeybird/autimaat/irc"
	"github.com/monkeybird/autimaat/irc/cmd"
	"github.com/monkeybird/autimaat/irc/proto"
	"github.com/monkeybird/autimaat/plugins"
)

func init() { plugins.Register(&plugin{}) }

// CacheTimeout defines the time after which a cache entry is
// considered stale and it must be re-fetched.
const CacheTimeout = time.Minute * 10

// LookupTimeout defines the timeout after which a service request
// is considered failed.
const LookupTimeout = time.Second * 5

type plugin struct {
	m                   sync.Mutex
	cmd                 *cmd.Set
	currentWeatherCache map[string]*currentWeatherResponse
	forecastCache       map[string]*forecastResponse
	config              struct {
		WundergroundApiKey string
	}
}

// Load initializes the module and loads any internal resources
// which may be required.
func (p *plugin) Load(prof irc.Profile) error {
	p.currentWeatherCache = make(map[string]*currentWeatherResponse)
	p.forecastCache = make(map[string]*forecastResponse)

	p.cmd = cmd.New(prof.CommandPrefix(), nil)
	p.cmd.Bind(TextCurrentWeatherName, false, p.cmdCurrentWeather).
		Add(TextLocation, true, cmd.RegAny)
	p.cmd.Bind(TextForecastName, false, p.cmdForecast).
		Add(TextLocation, true, cmd.RegAny)

	file := filepath.Join(prof.Root(), "weather.cfg")
	return util.ReadFile(file, &p.config, false)
}

// Unload cleans the module up and unloads any internal resources.
func (p *plugin) Unload(prof irc.Profile) error {
	p.config.WundergroundApiKey = ""
	return nil
}

// Dispatch sends the given, incoming IRC message to the plugin for
// processing as it sees fit.
func (p *plugin) Dispatch(w irc.ResponseWriter, r *irc.Request) {
	if len(p.config.WundergroundApiKey) > 0 {
		p.cmd.Dispatch(w, r)
	}
}

// sendLocations sends location suggestions to the request's sender.
func sendLocations(w irc.ResponseWriter, r *irc.Request, locs []location) {
	set := make([]string, 0, len(locs))

	// Add location descriptors to the set, provided they are unique.
	for _, l := range locs {
		value := fmt.Sprintf("%s %s %s", l.City, l.Country, l.State)
		if !hasString(set, value) {
			set = append(set, value)
		}
	}

	sort.Strings(set)

	proto.PrivMsg(w, r.Target, TextLocationsText,
		r.SenderName, strings.Join(set, ", "))
}

// hasString returnstrue if p contains a case-insensitive version of v,
func hasString(p []string, v string) bool {
	for _, pv := range p {
		if strings.EqualFold(pv, v) {
			return true
		}
	}
	return false
}

// fetch fetches the given URL contents and unmarshals them into the
// specified struct. This returns false if the fetch failed.
func (p *plugin) fetch(serviceURL, query string, v interface{}) bool {
	// Fetch new response.
	url := fmt.Sprintf(
		serviceURL,
		p.config.WundergroundApiKey,
		TextLanguageISO,
		query,
	)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("[weather] fetch: http.Get:", err)
		return false
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println("[weather] fetch: ioutil.ReadAll:", err)
		return false
	}

	//log.Println(string(data))

	err = json.Unmarshal(data, v)
	if err != nil {
		log.Println("[weather] fetch: json.Unmarshal:", err)
		return false
	}

	return true
}
