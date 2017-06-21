// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package url

const (
	TextDisplay         = "De link van %s toont: %s"
	TextYoutubeDuration = " (speelduur: %s)"

	// Not all user agents are created equal.
	//
	// Spotify will not return a song name in its <title>, if no user agent is
	// specified. But likewise when the Firefox UA for Ubuntu/Linux is provided.
	//
	// The UA listed here is found to work, so don't change it, unless you have
	// a replacement that passes the tests (TestTitle() in plugin_test.go).
	TextUserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:53.0) Gecko/20100101 Firefox/53.0"
)

// Ignore is a map of title strings to ignore. Only exact matches will
// be ignored.
var Ignore = map[string]bool{
	"Imgur: The most awesome images on the Internet": true,
	"Tweakers":                                     true,
	"Cookies op AD.nl | AD.nl":                     true,
	"Cookies op Trouw.nl":                          true,
	"Cookies op gelderlander.nl | gelderlander.nl": true,
	"Too Many Requests":                            true,
}
