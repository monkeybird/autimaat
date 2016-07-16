// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package weather provides bindings for some APIs at
// http://openweathermap.org/api
//
// This service requires the registration of a free account in order to get
// a valid API key. The API key you receive should be placed in a file named
// `openweathermap.cfg` in the profile directory for your bot. Its contents
// are expected to be as follows:
//
//	{
//		"ApiKey": "XXXXXXXX"
//	}
//
// Where `XXXXXXXX` is the API key you received from the openweathermap site.
package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"monkeybird/irc"
	"monkeybird/irc/cmd"
	"monkeybird/mod"
	"monkeybird/tr"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const cacheTimeout = time.Minute * 10

type module struct {
	lock                sync.RWMutex
	apiKeyFunc          func() string
	commands            *cmd.Set
	currentWeatherCache map[string]*CurrentWeatherResponse
	forecastCache       map[string]*ForecastResponse
}

// New returns a new module.
func New() mod.Module {
	return &module{
		currentWeatherCache: make(map[string]*CurrentWeatherResponse),
		forecastCache:       make(map[string]*ForecastResponse),
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
		Add(tr.WeatherLocationName, tr.WeatherLocationDesc, true, cmd.RegAny)

	m.commands.Bind(tr.ForecastName, tr.ForecastDesc, false, m.cmdForecast).
		Add(tr.ForecastLocationName, tr.ForecastLocationDesc, true, cmd.RegAny)

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
	m.lock.RLock()
	m.commands.HelpHandler(w, r)
	m.lock.RUnlock()
}

func (m *module) onPrivMsg(w irc.ResponseWriter, r *irc.Request) {
	m.commands.Dispatch(w, r)
}

// fetch fetches the given URL contents and unmarshals them into the
// specified struct. This returns false if the fetch failed.
func (m *module) fetch(serviceURL, query string, v interface{}) bool {
	// Fetch new response.
	url := fmt.Sprintf(
		serviceURL,
		url.QueryEscape(query),
		m.apiKeyFunc(),
	)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("[openweather] fetch: http.Get:", err)
		return false
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println("[openweather] fetch: ioutil.ReadAll:", err)
		return false
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		log.Println("[openweather] fetch: json.Unmarshal:", err)
		return false
	}

	return true
}

// weatherName returns the string representation of the given weather type ID.
func weatherName(id int) string {
	switch id {
	case 200:
		return tr.OpenWeather200
	case 201:
		return tr.OpenWeather201
	case 202:
		return tr.OpenWeather202
	case 210:
		return tr.OpenWeather210
	case 211:
		return tr.OpenWeather211
	case 212:
		return tr.OpenWeather212
	case 221:
		return tr.OpenWeather221
	case 230:
		return tr.OpenWeather230
	case 231:
		return tr.OpenWeather231
	case 232:
		return tr.OpenWeather232
	case 300:
		return tr.OpenWeather300
	case 301:
		return tr.OpenWeather301
	case 302:
		return tr.OpenWeather302
	case 310:
		return tr.OpenWeather310
	case 311:
		return tr.OpenWeather311
	case 312:
		return tr.OpenWeather312
	case 313:
		return tr.OpenWeather313
	case 314:
		return tr.OpenWeather314
	case 321:
		return tr.OpenWeather321
	case 500:
		return tr.OpenWeather500
	case 501:
		return tr.OpenWeather501
	case 502:
		return tr.OpenWeather502
	case 503:
		return tr.OpenWeather503
	case 504:
		return tr.OpenWeather504
	case 511:
		return tr.OpenWeather511
	case 520:
		return tr.OpenWeather520
	case 521:
		return tr.OpenWeather521
	case 522:
		return tr.OpenWeather522
	case 531:
		return tr.OpenWeather531
	case 600:
		return tr.OpenWeather600
	case 601:
		return tr.OpenWeather601
	case 602:
		return tr.OpenWeather602
	case 611:
		return tr.OpenWeather611
	case 612:
		return tr.OpenWeather612
	case 615:
		return tr.OpenWeather615
	case 616:
		return tr.OpenWeather616
	case 620:
		return tr.OpenWeather620
	case 621:
		return tr.OpenWeather621
	case 622:
		return tr.OpenWeather622
	case 701:
		return tr.OpenWeather701
	case 711:
		return tr.OpenWeather711
	case 721:
		return tr.OpenWeather721
	case 731:
		return tr.OpenWeather731
	case 741:
		return tr.OpenWeather741
	case 751:
		return tr.OpenWeather751
	case 761:
		return tr.OpenWeather761
	case 762:
		return tr.OpenWeather762
	case 771:
		return tr.OpenWeather771
	case 781:
		return tr.OpenWeather781
	case 800:
		return tr.OpenWeather800
	case 801:
		return tr.OpenWeather801
	case 802:
		return tr.OpenWeather802
	case 803:
		return tr.OpenWeather803
	case 804:
		return tr.OpenWeather804
	case 900:
		return tr.OpenWeather900
	case 901:
		return tr.OpenWeather901
	case 902:
		return tr.OpenWeather902
	case 903:
		return tr.OpenWeather903
	case 904:
		return tr.OpenWeather904
	case 905:
		return tr.OpenWeather905
	case 906:
		return tr.OpenWeather906
	case 951:
		return tr.OpenWeather951
	case 952:
		return tr.OpenWeather952
	case 953:
		return tr.OpenWeather953
	case 954:
		return tr.OpenWeather954
	case 955:
		return tr.OpenWeather955
	case 956:
		return tr.OpenWeather956
	case 957:
		return tr.OpenWeather957
	case 958:
		return tr.OpenWeather958
	case 959:
		return tr.OpenWeather959
	case 960:
		return tr.OpenWeather960
	case 961:
		return tr.OpenWeather961
	case 962:
		return tr.OpenWeather962
	}

	return ""
}

// direction turns a 0-360 degree angle into a wind direction name.
// E.g.: north, north-east, south-south-west, etc.
//
// ref: https://cdn.windfinder.com/prod/images/help/wind_directions.9d696e7e.png
func direction(v float32) string {
	const a = 11.25 // = (360.0 / 16.0) / 2

	if v >= a*31 && v <= a*1 {
		return tr.OpenWeatherN
	}

	if v > a*1 && v < a*3 {
		return tr.OpenWeatherNNE
	}

	if v > a*3 && v < a*5 {
		return tr.OpenWeatherNE
	}

	if v > a*5 && v < a*7 {
		return tr.OpenWeatherENE
	}

	if v > a*7 && v < a*9 {
		return tr.OpenWeatherE
	}

	if v > a*9 && v < a*11 {
		return tr.OpenWeatherESE
	}

	if v > a*11 && v < a*13 {
		return tr.OpenWeatherSE
	}

	if v > a*13 && v < a*15 {
		return tr.OpenWeatherSSE
	}

	if v > a*15 && v < a*17 {
		return tr.OpenWeatherS
	}

	if v > a*17 && v < a*19 {
		return tr.OpenWeatherSSW
	}

	if v > a*19 && v < a*21 {
		return tr.OpenWeatherSW
	}

	if v > a*21 && v < a*23 {
		return tr.OpenWeatherWSW
	}

	if v > a*23 && v < a*25 {
		return tr.OpenWeatherW
	}

	if v > a*25 && v < a*27 {
		return tr.OpenWeatherWNW
	}

	if v > a*27 && v < a*29 {
		return tr.OpenWeatherNW
	}

	if v > a*29 && v < a*31 {
		return tr.OpenWeatherNNW
	}

	return tr.OpenWeatherUnknown
}
