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

// This is the url used to fetch a weather forecast report.
const forecastURL = "https://api.wunderground.com/api/%s/forecast/lang:%s/q/%s.json"

// cmdForecast fetches a 3-day weather forecast for a specific location.
func (m *module) cmdForecast(w irc.ResponseWriter, r *cmd.Request) {
	proto.PrivMsg(w, r.Target, tr.WeatherNope, r.SenderName)
}

func (m *module) cmdForecast1(w irc.ResponseWriter, r *cmd.Request) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if len(m.apiKeyFunc()) == 0 {
		proto.PrivMsg(w, r.Target, tr.WeatherNotAvailable)
		return
	}

	loc := newLocation(r)
	key := strings.ToLower(loc.String())

	if fr, ok := m.forecastCache[key]; ok {
		// If the cached result is younger than the timeout, print its
		// contents for the user and exit. Otherwise, consider it stale,
		// delete it and re-fetch.
		if time.Since(fr.Timestamp) <= cacheTimeout {
			sendForecast(w, r, fr, loc)
			return
		}

		delete(m.currentWeatherCache, key)
	}

	var resp forecastResponse
	resp.Timestamp = time.Now()

	if !m.fetch(forecastURL, key, &resp) {
		return
	}

	// It is possible we received location suggestions, instead of weather
	// data. Present these suggestions to the user and exit. Do not cache
	// the response.
	if len(resp.Response.Results) > 0 {
		sendLocations(w, r, resp.Response.Results)
		return
	}

	sendForecast(w, r, &resp, loc)
	m.forecastCache[key] = &resp
}

// sendCurrentWeather formats a response for the user who invoked the
// weather request and sends it back to them.
func sendForecast(w irc.ResponseWriter, r *cmd.Request, fr *forecastResponse, loc *location) {
	location := text.Bold(loc.City)

	if len(loc.Country) > 0 {
		if len(loc.State) > 0 {
			location += fmt.Sprintf(" (%s, %s)", loc.Country, loc.State)
		} else {
			location += fmt.Sprintf(" (%s)", loc.Country)
		}
	}

	if len(fr.Forecast.TextForecast.ForecastDay) == 0 {
		proto.PrivMsg(w, r.SenderName, tr.WeatherForecastNotAvailable, location)
		return
	}

	proto.PrivMsg(w, r.SenderName, tr.WeatherForecastText, location)

	for _, v := range fr.Forecast.TextForecast.ForecastDay {
		proto.PrivMsg(w, r.SenderName, "%s: %s", text.Bold(v.Title), v.Text)
	}
}

// forecastResponse defines an API response.
type forecastResponse struct {
	Timestamp time.Time

	// This is filled if an ambiguous location name is provided to
	// the API. It will contain location suggestions for specific
	// places.
	Response struct {
		Results []location `json:"results"`
	} `json:"response"`

	// This defines actual forecast data for a specific location.
	// It will be empty if the Response.Results field is not.
	Forecast struct {
		TextForecast struct {
			ForecastDay []struct {
				Title string `json:"title"`
				Text  string `json:"fcttext_metric"`
			} `json:"forecastday"`
		} `json:"txt_forecast"`
	} `json:"forecast"`
}
