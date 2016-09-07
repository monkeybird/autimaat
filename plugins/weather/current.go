// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package weather

import (
	"fmt"
	"strings"
	"time"

	"github.com/monkeybird/autimaat/app/util"
	"github.com/monkeybird/autimaat/irc"
	"github.com/monkeybird/autimaat/irc/cmd"
	"github.com/monkeybird/autimaat/irc/proto"
)

const CurrentWeatherURL = "https://api.wunderground.com/api/%s/conditions/lang:%s/q/%s.json"

// cmdCurrentWeather yields current weather data for a given location.
func (p *plugin) cmdCurrentWeather(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	p.m.Lock()
	defer p.m.Unlock()

	if len(p.config.WundergroundApiKey) == 0 {
		proto.PrivMsg(w, r.Target, TextNoWeather)
		return
	}

	loc := newLocation(r)
	key := strings.ToLower(loc.String())

	if resp, ok := p.currentWeatherCache[key]; ok {
		// If the cached result is younger than the timeout, print its
		// contents for the user and exit. Otherwise, consider it stale,
		// delete it and re-fetch.
		if time.Since(resp.Timestamp) <= CacheTimeout {
			sendCurrentWeather(w, r, resp)
			return
		}

		delete(p.currentWeatherCache, key)
	}

	// Fetch new response.
	var resp currentWeatherResponse
	resp.Timestamp = time.Now()

	if !p.fetch(CurrentWeatherURL, key, &resp) {
		return
	}

	// It is possible we received location suggestions, instead of weather
	// data. Present these suggestions to the user and exit. Do not cache
	// the response.
	if len(resp.Response.Results) > 0 {
		sendLocations(w, r, resp.Response.Results)
		return
	}

	sendCurrentWeather(w, r, &resp)
	p.currentWeatherCache[key] = &resp
}

// sendCurrentWeather formats a response for the user who invoked the
// weather request and sends it back to them.
func sendCurrentWeather(w irc.ResponseWriter, r *irc.Request, cwr *currentWeatherResponse) {
	co := &cwr.CurrentObservation

	if strings.TrimSpace(co.DisplayLocation.City) == "" {
		proto.PrivMsg(w, r.Target, TextNoResult, r.SenderName)
		return
	}

	location := util.Bold(co.DisplayLocation.City)

	if len(co.DisplayLocation.Country) > 0 {
		if len(co.DisplayLocation.State) > 0 {
			location += fmt.Sprintf(" (%s, %s)",
				co.DisplayLocation.State, co.DisplayLocation.Country)
		} else {
			location += fmt.Sprintf(" (%s)", co.DisplayLocation.Country)
		}
	}

	proto.PrivMsg(w, r.Target, TextCurrentWeatherDisplay,
		r.SenderName,

		location,

		int(co.TempC),
		co.Weather,
		co.PressureMB,
		co.RelativeHumidity,
		co.WindKPH,
		co.WindDir,
	)
}

// currentWeatherResponse defines an API response.
type currentWeatherResponse struct {
	Timestamp time.Time

	// This is filled if an ambiguous location name is provided to
	// the API. It will contain location suggestions for specific
	// places.
	Response struct {
		Results []location `json:"results"`

		Error struct {
			Description string `json:"description"`
		} `json:"error"`
	} `json:"response"`

	// This is filled with actual weather data for a specific location.
	// It is only filled if the Response.Results field is empty.
	CurrentObservation struct {
		DisplayLocation  location `json:"display_location"`
		Weather          string   `json:"weather"`
		TempC            float32  `json:"temp_c"`
		RelativeHumidity string   `json:"relative_humidity"`
		WindDir          string   `json:"wind_dir"`
		WindKPH          float32  `json:"wind_kph"`
		PressureMB       string   `json:"pressure_mb"`
		FeelslikeC       string   `json:"feelslike_c"`
	} `json:"current_observation"`
}
