// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package url

const (
	TextDisplay         = "De link van %s toont: %s"
	TextYoutubeDuration = " (speelduur: %s)"
)

var Ignore = map[string]bool{
	"Imgur: The most awesome images on the Internet": true,
	"Cookies op AD.nl | AD.nl":                       true,
	"Cookies op Trouw.nl":                            true,
	"Cookies op gelderlander.nl | gelderlander.nl":   true,
}
