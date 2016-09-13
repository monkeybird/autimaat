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

const ForecastURL = "https://api.wunderground.com/api/%s/forecast/lang:%s/q/%s.json"

// cmdCurrentWeather yields weather forecast data for a given location.
func (p *plugin) cmdForecast(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	p.m.Lock()
	defer p.m.Unlock()

	if len(p.config.WundergroundApiKey) == 0 {
		proto.PrivMsg(w, r.Target, TextNoWeather)
		return
	}

	loc := newLocation(r)
	key := strings.ToLower(loc.String())

	if fr, ok := p.forecastCache[key]; ok {
		// If the cached result is younger than the timeout, print its
		// contents for the user and exit. Otherwise, consider it stale,
		// delete it and re-fetch.
		if time.Since(fr.Timestamp) <= CacheTimeout {
			sendForecast(w, r, fr, loc)
			return
		}

		delete(p.currentWeatherCache, key)
	}

	var resp forecastResponse
	resp.Timestamp = time.Now()

	if !p.fetch(ForecastURL, key, &resp) {
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
	p.forecastCache[key] = &resp
}

// sendCurrentWeather formats a response for the user who invoked the
// weather request and sends it back to them.
func sendForecast(w irc.ResponseWriter, r *irc.Request, fr *forecastResponse, loc *location) {
	location := util.Bold(loc.City)

	if len(loc.Country) > 0 {
		if len(loc.State) > 0 {
			location += fmt.Sprintf(" (%s, %s)", loc.Country, loc.State)
		} else {
			location += fmt.Sprintf(" (%s)", loc.Country)
		}
	}

	if len(fr.Forecast.TextForecast.ForecastDay) == 0 {
		proto.PrivMsg(w, r.SenderName, TextNoResult, r.SenderName)
		return
	}

	proto.PrivMsg(w, r.SenderName, TextForecastDisplay, location)

	for _, v := range fr.Forecast.TextForecast.ForecastDay {
		proto.PrivMsg(w, r.SenderName, "%s: %s", util.Bold(v.Title), v.Text)
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
