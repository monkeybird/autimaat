// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

/*
Package weather provides bindings for some API alls at
https://www.wunderground.com/weather/api/

This service requires the registration of a free account in order to get a
valid API key. The API key you receive should be assigned to the
`WeatherApiKey` field in the bot profile.
*/
package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/irc/proto"
	"monkeybird/mod"
	"monkeybird/tr"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

const cacheTimeout = time.Minute * 10

type module struct {
	lock                sync.RWMutex
	apiKeyFunc          func() string
	commands            *cmd.Set
	currentWeatherCache map[string]*currentWeatherResponse
	forecastCache       map[string]*forecastResponse
}

// New returns a new module.
func New() mod.Module {
	return &module{
		currentWeatherCache: make(map[string]*currentWeatherResponse),
		forecastCache:       make(map[string]*forecastResponse),
	}
}

// Load initializes the library and binds commands.
func (m *module) Load(pb irc.ProtocolBinder, prof irc.Profile) {
	m.lock.Lock()

	pb.Bind("PRIVMSG", m.onPrivMsg)

	m.commands = cmd.New(
		prof.CommandPrefix(),
		func(r *irc.Request) bool {
			return prof.IsWhitelisted(r.SenderMask)
		},
	)

	m.commands.Bind(tr.WeatherName, tr.WeatherDesc, false, m.cmdCurrentWeather).
		Add(tr.WeatherLocationName, tr.WeatherLocationName, true, cmd.RegAny)

	m.commands.Bind(tr.ForecastName, tr.ForecastDesc, false, m.cmdForecast).
		Add(tr.WeatherLocationName, tr.WeatherLocationName, true, cmd.RegAny)

	m.apiKeyFunc = prof.WeatherApiKey
	m.lock.Unlock()
}

// Unload cleans up any library resources and unbinds commands.
func (m *module) Unload(pb irc.ProtocolBinder, prof irc.Profile) {
	m.lock.Lock()
	m.commands.Clear()
	m.apiKeyFunc = nil
	pb.Unbind("PRIVMSG", m.onPrivMsg)
	m.lock.Unlock()
}

// Help displays help on custom commands.
func (m *module) Help(w irc.ResponseWriter, r *cmd.Request) {
	m.commands.HelpHandler(w, r)
}

func (m *module) onPrivMsg(w irc.ResponseWriter, r *irc.Request) {
	m.commands.Dispatch(w, r)
}

// sendLocations sends location suggestions to the request's sender.
func sendLocations(w irc.ResponseWriter, r *cmd.Request, locs []location) {
	set := make([]string, 0, len(locs))

	// Add location descriptors to the set, provided they are unique.
	for _, l := range locs {
		value := fmt.Sprintf("%s %s %s", l.City, l.Country, l.State)
		if !hasString(set, value) {
			set = append(set, value)
		}
	}

	sort.Strings(set)

	proto.PrivMsg(w, r.Target, tr.WeatherLocationsText,
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
func (m *module) fetch(serviceURL, query string, v interface{}) bool {
	// Fetch new response.
	url := fmt.Sprintf(
		serviceURL,
		m.apiKeyFunc(),
		tr.LanguageISO,
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
