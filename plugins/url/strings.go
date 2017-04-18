// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package url

const (
	TextDisplay         = "De link van %s toont: %s"
	TextYoutubeDuration = " (speelduur: %s)"
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
