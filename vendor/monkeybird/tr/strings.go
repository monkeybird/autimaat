// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package tr defines a list of all string literals used in the application.
// Specifically, those involved in messages which are sent to an IRC channel
// or user.
//
// They are kept here, in a centralized location to make changes and spelling
// corrections easier.
package tr

import "strings"

// ParseBool treats the given string as a boolean and returns its value.
// This is part of the tr package, because the way things like on/off are
// represented, differs per language.
func ParseBool(v string) bool {
	switch strings.ToLower(v) {
	case "1", "ja", "j", "aan", "a":
		return true
	default:
		return false
	}
}

// List of all known string literals.
const (
	// ref: https://godoc.org/time#Time.Format
	DateFormat     = "02-01-2006"
	TimeFormat     = "15:04:05 MST"
	DateTimeFormat = DateFormat + " " + TimeFormat

	CommandsRestricted       = "[beschermd]"
	CommandsOptional         = "[optioneel]"
	CommandsIntro1           = "Hier volgt een lijst met alle ondersteunde commandos. Commandos gemarkeerd met \x02*\x02 zijn beschermd en kunnen uitsluitend door een beheerder uitgevoerd worden."
	CommandsIntro2           = "Tenzij anders is aangegeven, dienen alle parameters met meerdere woorden ingevoerd te worden tussen dubbele aanhalingstekens, Bijv.: \x02\"Dit is 1 parameter\"\x02."
	CommandsMissingParameter = "Ontbrekende parameters voor commando: %s"
	CommandsInvalidParameter = "Commando %s: ongeldige waarde voor parameter %q"
	CommandsAccessDenied     = "Helaas, pindakaas. Het commando %q mag uitsluitend door beheerders uitgevoerd worden."

	UrlDisplayText = "De link van %s toont: %s"

	JoinName         = "join"
	JoinDesc         = "Instrueer de bot om het opgegeven kanaal te betreden."
	JoinChannelName  = "kanaal"
	JoinChannelDesc  = "Naam van het te betreden kanaal."
	JoinKeyName      = "sleutel"
	JoinKeyDesc      = "Toegangssleutel voor het kanaal."
	JoinPasswordName = "wachtwoord"
	JoinPasswordDesc = "Chanserv wachtwoord voor de bot in dit kanaal."

	PartName        = "part"
	PartDesc        = "Instrueer de bot het opgegeven kanaal te verlaten."
	PartChannelName = "kanaal"
	PartChannelDesc = "Naam van het te verlaten kanaal."

	HelpName        = "help"
	HelpDesc        = "Toon algemene informatie voor alle ondersteunde commandos, of gedetaileerde informatie voor een specifiek commando."
	HelpCommandName = "commando"
	HelpCommandDesc = "Naam van het commando in kwestie."

	LogName      = "log"
	LogDesc      = "Schakel het loggen van inkomende data in of uit, of geef de huidige logstatus weer."
	LogValueName = "status"
	LogValueDesc = "Een waarde die aangeeft of logging aan- of uit-geschakeld moet worden: ja/nee, j/n, aan/uit, a/u, 1/0"
	LogEnabled   = "Logging is ingeschakeld."
	LogDisabled  = "Logging is uitgeschakeld."

	ReloadName = "herstart"
	ReloadDesc = "Instrueer de bot zichzelf te herstarten."

	AuthListName        = "bazen"
	AuthListDesc        = "Geef een lijst weer met alle bekende bot beheerders."
	AuthListDisplayText = "De beheerders zijn: %s"

	AuthorizeName        = "baas"
	AuthorizeDesc        = "Geef een bepaalde gebruiker beheerderstoegang tot de bot."
	AuthorizeMaskName    = "hostmask"
	AuthorizeMaskDesc    = "Hostmask van de gebruiker in kwestie."
	AuthorizeDisplayText = "Gebruiker %q is toegevoegd aan de beheerderslijst."

	DeauthorizeName        = "ontbaas"
	DeauthorizeDesc        = "Ontneem een bepaalde gebruiker beheerderstoegang tot de bot."
	DeauthorizeMaskName    = "hostmask"
	DeauthorizeMaskDesc    = "Hostmask van de gebruiker in kwestie."
	DeauthorizeDisplayText = "Gebruiker %q is verwijderd van de beheerderslijst."

	VersionName        = "versie"
	VersionDesc        = "Geef versie informatie van de bot weer."
	VersionDisplayText = "%s, ik ben %s, versie %s. Mijn laatste revisie was op %s, om %s."

	EightballName         = "8ball"
	EightballDesc         = "Vraag De Magische 8ball een vraag en bereid je voor op ongezouten waarheid."
	EightballQuestionName = "vraag"
	EightballQuestionDesc = "De vraag die je De Magische 8ball wenst te stellen."
	Eightball1            = "%s, het is zeker."
	Eightball2            = "%s, het is absoluut zeker."
	Eightball3            = "%s, zonder twijfel."
	Eightball4            = "%s, ja, absoluut."
	Eightball5            = "%s, daar kun je van op aan."
	Eightball6            = "%s, zoals ik het zie, ja."
	Eightball7            = "%s, zeer waarschijnlijk."
	Eightball8            = "%s, vooruitzicht is goed."
	Eightball9            = "%s, ja."
	Eightball10           = "%s, tekenen wijzen op ja."
	Eightball11           = "%s, antwoord onduidelijk. Probeer later nog eens."
	Eightball12           = "%s, vraag later nog eens."
	Eightball13           = "%s, dat zeg ik je nu liever niet."
	Eightball14           = "%s, dat kan ik nu niet voorspellen."
	Eightball15           = "%s, concentreer je en vraag het nog eens."
	Eightball16           = "%s, reken er maar niet op."
	Eightball17           = "%s, mijn antwoord is nee."
	Eightball18           = "%s, mijn bronnen zeggen nee."
	Eightball19           = "%s, vooruitzicht is niet zo goed."
	Eightball20           = "%s, zeer twijfelachtig."

	StatsNotInChannel = "Dit commando werkt alleen indien aangeroepen vanuit een kanaal."
	StatsNoSuchUser   = "%s, ik vond geen informatie over gebruiker %s."

	FirstOnName        = "firston"
	FirstOnDesc        = "Geef aan wanneer een bepaalde gebruiker voor het eerst in het kanaal gezien werd."
	FirstOnUserName    = "gebruiker"
	FirstOnUserDesc    = "Naam of hostmask van de gebruiker in kwestie."
	FirstOnDisplayText = "%s, ik heb %s voor het eerst gezien op %s, op %s (±%s geleden)."

	LastOnName        = "laston"
	LastOnDesc        = "Geef aan wanneer een bepaalde gebruiker voor het laatst in het kanaal gezien werd."
	LastOnUserName    = "gebruiker"
	LastOnUserDesc    = "Naam of hostmask van de gebruiker in kwestie."
	LastOnDisplayText = "%s, ik heb %s voor het laatst gezien op %s, op %s (±%s geleden)."

	WeatherName         = "weer"
	WeatherDesc         = "Toon het huidige weer voor een specifieke lokatie."
	WeatherLocationName = "lokatie"
	WeatherLocationDesc = "Naam van de lokatie in kwestie. Dit is een dorp of stad, optioneel gevolgd door een land code. Bijv.: \"eindhoven\" of \"Amsterdam,NL\""

	ForecastName         = "weerfc"
	ForecastDesc         = "Toon een 3-daagse weersvoorspelling voor een specifieke lokatie."
	ForecastLocationName = "lokatie"
	ForecastLocationDesc = "Naam van de lokatie in kwestie. Dit is een dorp of stad, optioneel gevolgd door een land code. Bijv.: \"eindhoven\" of \"Amsterdam,NL\""

	OpenWeatherNotAvailable         = "Het weerbericht is momenteel niet beschikbaar."
	OpenWeatherCurrentWeatherText   = "%s, in %s (%s) is het %d℃, %s, luchtdruk: %d hPa, luchtvochtigheid: %d%%, wind: %.1f km/u uit %s richting, bewolking: %d%%."
	OpenWeatherForecastNotAvailable = "%s, er is momenteel geen weersvoorspelling voor %s (%s)."
	OpenWeatherForecastText1        = "%s, voorspelling voor de komende 3 dagen in %s (%s):"
	OpenWeatherForecastText2        = "%s, in \x02±%d uur\x02: %d℃, %s, wind: %.1f km/h uit %s richting, bewolking: %d%%."
	OpenWeatherN                    = "noordelijke"
	OpenWeatherNNE                  = "noord-noord-oostelijke"
	OpenWeatherNE                   = "noord-oostelijke"
	OpenWeatherENE                  = "oost-noord-oostelijke"
	OpenWeatherE                    = "oostelijke"
	OpenWeatherESE                  = "oost-zuid-oostelijke"
	OpenWeatherSE                   = "zuid-oostelijke"
	OpenWeatherSSE                  = "zuid-zuid-oostelijke"
	OpenWeatherS                    = "zuidelijke"
	OpenWeatherSSW                  = "zuid-zuid-westelijke"
	OpenWeatherSW                   = "zuid-westelijke"
	OpenWeatherWSW                  = "west-zuid-westelijke"
	OpenWeatherW                    = "westelijke"
	OpenWeatherWNW                  = "west-noord-westelijke"
	OpenWeatherNW                   = "noord-westelijke"
	OpenWeatherNNW                  = "noord-noord-westelijke"
	OpenWeatherUnknown              = "onbekende"
	OpenWeather200                  = "onweer met lichte regen"
	OpenWeather201                  = "onweer met regen"
	OpenWeather202                  = "onweer met zware regen"
	OpenWeather210                  = "lichte onweer"
	OpenWeather211                  = "onweer"
	OpenWeather212                  = "zware onweer"
	OpenWeather221                  = "ragged thunderstorm"
	OpenWeather230                  = "thunderstorm with light drizzle"
	OpenWeather231                  = "thunderstorm with drizzle "
	OpenWeather232                  = "thunderstorm with heavy drizzle"
	OpenWeather300                  = "light intensity drizzle"
	OpenWeather301                  = "drizzle"
	OpenWeather302                  = "heavy intensity drizzle"
	OpenWeather310                  = "light intensity drizzle rain"
	OpenWeather311                  = "drizzle rain"
	OpenWeather312                  = "heavy intensity drizzle rain"
	OpenWeather313                  = "shower rain and drizzle"
	OpenWeather314                  = "heavy shower rain and drizzle"
	OpenWeather321                  = "shower drizzle"
	OpenWeather500                  = "lichte regen"
	OpenWeather501                  = "matige regen"
	OpenWeather502                  = "hevige regen"
	OpenWeather503                  = "zeer hevige regen"
	OpenWeather504                  = "extreme regen"
	OpenWeather511                  = "vriezende regen"
	OpenWeather520                  = "light intensity shower rain"
	OpenWeather521                  = "shower rain"
	OpenWeather522                  = "heavy intensity shower rain"
	OpenWeather531                  = "ragged shower rain"
	OpenWeather600                  = "lichte sneeuw"
	OpenWeather601                  = "sneeuw"
	OpenWeather602                  = "zware sneeuw"
	OpenWeather611                  = "sleet"
	OpenWeather612                  = "shower sleet"
	OpenWeather615                  = "lichte regen en sneeuw"
	OpenWeather616                  = "regen en sneeuw"
	OpenWeather620                  = "light shower snow"
	OpenWeather621                  = "shower snow"
	OpenWeather622                  = "heavy shower snow"
	OpenWeather701                  = "mist"
	OpenWeather711                  = "rook"
	OpenWeather721                  = "mist"
	OpenWeather731                  = "zand, stoffige wolken"
	OpenWeather741                  = "mist"
	OpenWeather751                  = "zand"
	OpenWeather761                  = "stoffig"
	OpenWeather762                  = "vulkanische as"
	OpenWeather771                  = "squalls"
	OpenWeather781                  = "tornado"
	OpenWeather800                  = "heldere lucht"
	OpenWeather801                  = "enkele wolken"
	OpenWeather802                  = "scattered clouds"
	OpenWeather803                  = "broken clouds"
	OpenWeather804                  = "overcast clouds"
	OpenWeather900                  = "tornado"
	OpenWeather901                  = "tropische storm"
	OpenWeather902                  = "hurricane"
	OpenWeather903                  = "koud"
	OpenWeather904                  = "heet"
	OpenWeather905                  = "winderig"
	OpenWeather906                  = "hagel"
	OpenWeather951                  = "kalm"
	OpenWeather952                  = "light breeze"
	OpenWeather953                  = "gentle breeze"
	OpenWeather954                  = "moderate breeze"
	OpenWeather955                  = "fresh breeze"
	OpenWeather956                  = "strong breeze"
	OpenWeather957                  = "high wind, near gale"
	OpenWeather958                  = "gale"
	OpenWeather959                  = "severe gale"
	OpenWeather960                  = "storm"
	OpenWeather961                  = "hevige storm"
	OpenWeather962                  = "hurricane"
)
