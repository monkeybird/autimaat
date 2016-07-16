// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package weather

import (
	"math"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/irc/proto"
	"monkeybird/text"
	"monkeybird/tr"
	"strings"
	"time"
)

// This is the url used to fetch a weather forecast report.
const forecastURL = "http://api.openweathermap.org/data/2.5/forecast?q=%s&units=metric&type=accurate&APPID=%s"

// cmdForecast fetches a 3-day weather forecast for a specific location.
func (m *module) cmdForecast(w irc.ResponseWriter, r *cmd.Request) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if len(m.apiKeyFunc()) == 0 {
		proto.PrivMsg(w, r.Target, tr.OpenWeatherNotAvailable)
		return
	}

	key := strings.ToLower(r.String(0))

	if fr, ok := m.forecastCache[key]; ok {
		// If the cached result is younger than the timeout, print its
		// contents for the user and exit. Otherwise, consider it stale,
		// delete it and re-fetch.
		if time.Since(fr.Timestamp) <= cacheTimeout {
			sendForecast(w, r, fr)
			return
		}

		delete(m.currentWeatherCache, key)
	}

	var resp ForecastResponse
	resp.Timestamp = time.Now()

	if !m.fetch(forecastURL, key, &resp) {
		return
	}

	sendForecast(w, r, &resp)
	m.forecastCache[key] = &resp
}

// sendCurrentWeather formats a response for the user who invoked the
// weather request and sends it back to them.
func sendForecast(w irc.ResponseWriter, r *cmd.Request, fr *ForecastResponse) {
	if len(fr.List) == 0 {
		proto.PrivMsg(w, r.Target,
			tr.OpenWeatherForecastNotAvailable,
			r.SenderName,
			text.Bold(fr.City.Name),
			fr.City.Country,
		)
		return
	}

	proto.PrivMsg(w, r.SenderName,
		tr.OpenWeatherForecastText1,
		r.SenderName,
		text.Bold(fr.City.Name),
		fr.City.Country,
	)

	if fc := findForecast(fr.List, time.Hour*24); fc != nil {
		formatForecast(w, r, fc)
	}

	if fc := findForecast(fr.List, time.Hour*48); fc != nil {
		formatForecast(w, r, fc)
	}

	if fc := findForecast(fr.List, time.Hour*72); fc != nil {
		formatForecast(w, r, fc)
	}
}

// formatForecast formats and sends the given forecast to the caller.
func formatForecast(w irc.ResponseWriter, r *cmd.Request, fc *Forecast) {
	stamp := time.Unix(fc.When, 0)
	delta := stamp.Sub(time.Now().UTC())

	var weather string
	if len(fc.Weather) > 0 {
		weather = weatherName(fc.Weather[0].ID)
	}

	proto.PrivMsg(w, r.SenderName,
		tr.OpenWeatherForecastText2,

		r.SenderName,
		int(math.Abs(delta.Hours())),
		int(math.Abs(float64(fc.Main.Temp))),

		weather,
		fc.Wind.Speed*3.600,
		direction(fc.Wind.Direction),
		int(fc.Clouds.Percentage),
	)

	// Wait a little while before we continue with any additional
	// forecast entries. Nobody likes flooding.
	<-time.After(500 * time.Millisecond)
}

// findForecast scans the given list for a forecast at- or near the
// specified duration. If found, it returns the forecast. Returns
// nil otherwise.
func findForecast(set []*Forecast, when time.Duration) *Forecast {
	for _, fc := range set {
		stamp := time.Unix(fc.When, 0)
		delta := stamp.Sub(time.Now().UTC())

		if delta.Hours() >= when.Hours() {
			return fc
		}
	}

	return nil
}

// ForecastResponse defines an API response.
type ForecastResponse struct {
	Timestamp time.Time
	City      struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"city"`
	List []*Forecast `json:"list"`
}

// Forecast defines a single weather forecast.
type Forecast struct {
	When    int64 `json:"dt"`
	Weather []struct {
		ID int `json:"id"`
	} `json:"weather"`
	Main struct {
		Temp float32 `json:"temp"`
	} `json:"main"`
	Clouds struct {
		Percentage float32 `json:"all"`
	} `json:"clouds"`
	Wind struct {
		Speed     float32 `json:"speed"`
		Direction float32 `json:"deg"`
	} `json:"wind"`
}
