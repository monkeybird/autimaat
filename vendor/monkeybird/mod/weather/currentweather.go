// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package weather

import (
	"fmt"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/irc/proto"
	"monkeybird/text"
	"monkeybird/tr"
	"strings"
	"time"
)

// This is the url used to fetch a current weather report.
const currentWeatherURL = "https://api.wunderground.com/api/%s/conditions/lang:%s/q/%s.json"

func (m *module) cmdCurrentWeather(w irc.ResponseWriter, r *cmd.Request) {
	proto.PrivMsg(w, r.Target, tr.WeatherNope, r.SenderName)
}

// cmdCurrentWeather fetches a current weather report for a specific location.
func (m *module) cmdCurrentWeather1(w irc.ResponseWriter, r *cmd.Request) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if len(m.apiKeyFunc()) == 0 {
		proto.PrivMsg(w, r.Target, tr.WeatherNotAvailable)
		return
	}

	loc := newLocation(r)
	key := strings.ToLower(loc.String())

	if resp, ok := m.currentWeatherCache[key]; ok {
		// If the cached result is younger than the timeout, print its
		// contents for the user and exit. Otherwise, consider it stale,
		// delete it and re-fetch.
		if time.Since(resp.Timestamp) <= cacheTimeout {
			sendCurrentWeather(w, r, resp)
			return
		}

		delete(m.currentWeatherCache, key)
	}

	// Fetch new response.
	var resp currentWeatherResponse
	resp.Timestamp = time.Now()

	if !m.fetch(currentWeatherURL, key, &resp) {
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
	m.currentWeatherCache[key] = &resp
}

// sendCurrentWeather formats a response for the user who invoked the
// weather request and sends it back to them.
func sendCurrentWeather(w irc.ResponseWriter, r *cmd.Request, cwr *currentWeatherResponse) {
	co := &cwr.CurrentObservation

	location := text.Bold(co.DisplayLocation.City)

	if len(co.DisplayLocation.Country) > 0 {
		if len(co.DisplayLocation.State) > 0 {
			location += fmt.Sprintf(" (%s, %s)",
				co.DisplayLocation.Country, co.DisplayLocation.State)
		} else {
			location += fmt.Sprintf(" (%s)", co.DisplayLocation.Country)
		}
	}

	proto.PrivMsg(w, r.Target, tr.WeatherCurrentWeatherText,
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
