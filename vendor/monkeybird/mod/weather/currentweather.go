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

// This is the url used to fetch a current weather report.
const currentWeatherURL = "http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&type=accurate&APPID=%s"

// cmdCurrentWeather fetches a current weather report for a specific location.
func (m *module) cmdCurrentWeather(w irc.ResponseWriter, r *cmd.Request) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if len(m.settings.ApiKey) == 0 {
		proto.PrivMsg(w, r.Target, tr.OpenWeatherNotAvailable)
		return
	}

	key := strings.ToLower(r.String(0))

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
	var resp CurrentWeatherResponse
	resp.Timestamp = time.Now()

	if !m.fetch(currentWeatherURL, key, &resp) {
		return
	}

	sendCurrentWeather(w, r, &resp)
	m.currentWeatherCache[key] = &resp
}

// sendCurrentWeather formats a response for the user who invoked the
// weather request and sends it back to them.
func sendCurrentWeather(w irc.ResponseWriter, r *cmd.Request, cwr *CurrentWeatherResponse) {
	var weather string
	if len(cwr.Weather) > 0 {
		weather = weatherName(cwr.Weather[0].ID)
	}

	proto.PrivMsg(w, r.Target,
		tr.OpenWeatherCurrentWeatherText,

		r.SenderName,
		text.Bold(cwr.Name),
		cwr.Sys.Country,
		int(math.Abs(float64(cwr.Main.Temp))),

		weather,
		int(cwr.Main.Pressure),
		int(cwr.Main.Humidity),

		cwr.Wind.Speed*3.600,
		direction(cwr.Wind.Direction),
		int(cwr.Clouds.Percentage),
	)
}

// CurrentWeatherResponse defines an API response.
type CurrentWeatherResponse struct {
	Timestamp time.Time
	Name      string `json:"name"`
	Weather   []struct {
		ID int `json:"id"`
	} `json:"weather"`
	Main struct {
		Temp     float32 `json:"temp"`
		Pressure float32 `json:"pressure"`
		Humidity float32 `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed     float32 `json:"speed"`
		Direction float32 `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		Percentage float32 `json:"all"`
	} `json:"clouds"`
	Sys struct {
		Country string `json:"country"`
	} `json:"sys"`
}
