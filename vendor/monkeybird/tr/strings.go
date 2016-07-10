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

const (
	// ref: https://godoc.org/time#Time.Format
	DateFormat     = "02 January, 2006"
	TimeFormat     = "15:04 MST"
	DateTimeFormat = DateFormat + " " + TimeFormat

	CommandsRestricted       = "[beschermd]"
	CommandsOptional         = "[optioneel]"
	CommandsIntro1           = "Hier volgt een lijst met alle ondersteunde commandos. Commandos gemarkeerd met \x02*\x02 zijn beschermd en kunnen uitsluitend door een beheerder uitgevoerd worden."
	CommandsIntro2           = "Tenzij anders is aangegeven, dienen alle parameters met meerdere woorden ingevoerd te worden tussen dubbele aanhalingstekens, Bijv.: \x02\"Dit is 1 parameter\"\x02."
	CommandsMissingParameter = "Ontbrekende parameters voor commando: %s"
	CommandsInvalidParameter = "Commando %s: ongeldige waarde voor parameter %q"
	CommandsAccessDenied     = "Helaas, pindakaas. Het commando %q mag uitsluitend door beheerders uitgevoerd worden."

	UrlDisplayText = "De link van %s toont: %s"
)

const (
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
)

const (
	HelpName        = "help"
	HelpDesc        = "Toon algemene informatie voor alle ondersteunde commandos, of gedetaileerde informatie voor een specifiek commando."
	HelpCommandName = "commando"
	HelpCommandDesc = "Naam van het commando in kwestie."
)

const (
	LogName      = "log"
	LogDesc      = "Schakel het loggen van inkomende data in of uit, of geef de huidige logstatus weer."
	LogValueName = "status"
	LogValueDesc = "Een waarde die aangeeft of logging aan- of uit-geschakeld moet worden: ja/nee, j/n, aan/uit, a/u, 1/0"
	LogEnabled   = "Logging is ingeschakeld."
	LogDisabled  = "Logging is uitgeschakeld."
)

const (
	ReloadName = "herstart"
	ReloadDesc = "Instrueer de bot zichzelf te herstarten."
)

const (
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
)

const (
	VersionName        = "versie"
	VersionDesc        = "Geef versie informatie van de bot weer."
	VersionDisplayText = "%s, ik ben %s, versie %s. Mijn laatste revisie was op %s, om %s."
)
const (
	StatsNotInChannel = "Dit commando werkt alleen indien aangeroepen vanuit een kanaal."
	StatsNoSuchUser   = "%s, ik vond geen informatie over gebruiker %s."

	FirstOnName        = "firston"
	FirstOnDesc        = "Geef aan wanneer een bepaalde gebruiker voor het eerst in het kanaal gezien werd."
	FirstOnUserName    = "gebruiker"
	FirstOnUserDesc    = "Naam of hostmask van de gebruiker in kwestie."
	FirstOnDisplayText = "%s, ik heb %s voor het eerst gezien op %s, om %s (±%s geleden)."

	LastOnName        = "laston"
	LastOnDesc        = "Geef aan wanneer een bepaalde gebruiker voor het laatst in het kanaal gezien werd."
	LastOnUserName    = "gebruiker"
	LastOnUserDesc    = "Naam of hostmask van de gebruiker in kwestie."
	LastOnDisplayText = "%s, ik heb %s voor het laatst gezien op %s, om %s (±%s geleden)."
)

const (
	SnoozeTimeFormat     = "15:04"
	SnoozeName           = "snooze"
	SnoozeDesc           = "Plan een alarm voor de opgegeven tijd."
	SnoozeTimeName       = "tijd"
	SnoozeTimeDesc       = "De tijd waarop het alarm af dient te gaan. Dit is een absoluute tijd zoals \"16:32\", of een aantal minuten vanaf nu. Bijv.: \"30\""
	SnoozeMessageName    = "bericht"
	SnoozeMessageDesc    = "Het bericht dat weergegeven dient te worden als het alarm af gaat."
	SnoozeInvalidTime    = "%s, %q is geen geldige tijd waarde."
	SnoozeDefaultMessage = "%s, toooot! Snooze tijd!"
	SnoozeMessagePrefix  = "%s, het is %s: "
	SnoozeAlarmSet       = "%s, het alarm is ingesteld. Je kunt het verwijderen met: !unsnooze %s"
	SnoozeAlarmUnset     = "%s, het alarm is verwijderd."

	UnsnoozeName   = "unsnooze"
	UnsnoozeDesc   = "Verwijder een bestaand snooze alarm. Je kunt alleen alarmen verwijderen die je zelf hebt gemaakt."
	UnsnoozeIDName = "id"
	UnsnoozeIDDesc = "De code voor het alarm dat verwijderd dient te worden."
)

const (
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
	OpenWeather221                  = "extreem onweer"
	OpenWeather230                  = "onweer met lichte motregen"
	OpenWeather231                  = "onweer met motregen"
	OpenWeather232                  = "onweer met hevige motregen"
	OpenWeather300                  = "lichte motregen"
	OpenWeather301                  = "motregen"
	OpenWeather302                  = "hevige motregen"
	OpenWeather310                  = "lichte motregen"
	OpenWeather311                  = "motregen"
	OpenWeather312                  = "hevige motregen"
	OpenWeather313                  = "stortregen, afgewisseld met motregen"
	OpenWeather314                  = "hevige stortregen, afgewisseld met motregen"
	OpenWeather321                  = "stortregen"
	OpenWeather500                  = "lichte regen"
	OpenWeather501                  = "matige regen"
	OpenWeather502                  = "hevige regen"
	OpenWeather503                  = "zeer hevige regen"
	OpenWeather504                  = "extreme regen"
	OpenWeather511                  = "vriezende regen"
	OpenWeather520                  = "lichte stortregen"
	OpenWeather521                  = "stortregen"
	OpenWeather522                  = "hevige stortregen"
	OpenWeather531                  = "extreme stortregen"
	OpenWeather600                  = "lichte sneeuw"
	OpenWeather601                  = "sneeuw"
	OpenWeather602                  = "zware sneeuw"
	OpenWeather611                  = "ijzel"
	OpenWeather612                  = "ijzel"
	OpenWeather615                  = "lichte regen en sneeuw"
	OpenWeather616                  = "regen en sneeuw"
	OpenWeather620                  = "lichte sneeuw"
	OpenWeather621                  = "sneeuw"
	OpenWeather622                  = "hevige sneeuw"
	OpenWeather701                  = "mistig"
	OpenWeather711                  = "rook"
	OpenWeather721                  = "mistig"
	OpenWeather731                  = "zand, stoffige wolken"
	OpenWeather741                  = "mist"
	OpenWeather751                  = "zand"
	OpenWeather761                  = "stoffig"
	OpenWeather762                  = "vulkanische as"
	OpenWeather771                  = "rukwinden"
	OpenWeather781                  = "tornado"
	OpenWeather800                  = "heldere lucht"
	OpenWeather801                  = "enkele wolken"
	OpenWeather802                  = "verspreide bewolking"
	OpenWeather803                  = "gebroken bewolking"
	OpenWeather804                  = "bewolkt"
	OpenWeather900                  = "tornado"
	OpenWeather901                  = "tropische storm"
	OpenWeather902                  = "orkaan"
	OpenWeather903                  = "koud"
	OpenWeather904                  = "heet"
	OpenWeather905                  = "winderig"
	OpenWeather906                  = "hagel"
	OpenWeather951                  = "kalm"
	OpenWeather952                  = "lichte bries"
	OpenWeather953                  = "lichte bries"
	OpenWeather954                  = "matige bries"
	OpenWeather955                  = "fris briesje"
	OpenWeather956                  = "stevige bries"
	OpenWeather957                  = "stormachtig"
	OpenWeather958                  = "stormachtig"
	OpenWeather959                  = "hevige storm"
	OpenWeather960                  = "storm"
	OpenWeather961                  = "hevige storm"
	OpenWeather962                  = "orkaan"
)

const (
	EightballName         = "8ball"
	EightballDesc         = "Vraag De Magische 8ball een vraag en bereid je voor op ongezouten waarheid."
	EightballQuestionName = "vraag"
	EightballQuestionDesc = "De vraag die je De Magische 8ball wenst te stellen."
)

// EightBallAnswers defines the list of possible 8ball answers.
var EightBallAnswers = []string{
	"%s, het is zeker.",
	"%s, het is absoluut zeker.",
	"%s, zonder twijfel.",
	"%s, ja, absoluut.",
	"%s, daar kun je van op aan.",
	"%s, zoals ik het zie, ja.",
	"%s, zeer waarschijnlijk.",
	"%s, vooruitzicht is goed.",
	"%s, ja.",
	"%s, tekenen wijzen op ja.",
	"%s, antwoord onduidelijk. Probeer later nog eens.",
	"%s, vraag later nog eens.",
	"%s, dat zeg ik je nu liever niet.",
	"%s, dat kan ik nu niet voorspellen.",
	"%s, concentreer je en vraag het nog eens.",
	"%s, reken er maar niet op.",
	"%s, mijn antwoord is nee.",
	"%s, mijn bronnen zeggen nee.",
	"%s, vooruitzicht is niet zo goed.",
	"%s, zeer twijfelachtig.",
}

const (
	GiveUserName = "wie"
	GiveUserDesc = "Naam van de ontvanger."

	BeerName = "bier"
	BeerDesc = "Geef jezelf of iemand anders een biertje."

	WineName = "wijn"
	WineDesc = "Geef jezelf of iemand anders een wijntje."

	CoffeeName = "koffie"
	CoffeeDesc = "Geef jezelf of iemand anders een kopje koffie."

	TeaName = "thee"
	TeaDesc = "Geef jezelf of iemand anders een kopje thee."

	LemonadeName = "fris"
	LemonadeDesc = "Geef jezelf of iemand anders een glaasje fris."
)

// The strings below should be written as if part of an action.
// E.g.: "/me <something something...>"

var BeerAnswers = []string{
	"geeft %s een lekker koud biertje.",
	"geeft %s een lekkere blonde stoot.",
	"opent een fust en schuift het naar %s. Proost!",
}

var WineAnswers = []string{
	"geeft %s een glaasje rode wijn.",
	"geeft %s een glaasje witte wijn.",
	"geeft %s een glaasje Port.",
}

var CoffeeAnswers = []string{
	"schenkt %s een kopje verse koffie.",
	"geeft %s een stevige bak leut.",
}

var TeaAnswers = []string{
	"schenkt %s een vers kopje thee.",
	"biedt %s een glas hete thee aan.",
}

var LemonadeAnswers = []string{
	"schenkt %s een groot glas koude limonade.",
	"geeft %s de fles cola aan.",
	"geeft %s een scheut ijskoud aardbeiensap.",
	"geeft %s een scheut ijskoud bosvruchtensap.",
	"geeft %s een scheut ijskoud citroensap.",
	"geeft %s een scheut ijskoud appelsap.",
	"geeft %s een scheut ijskoud sinaasappelsap.",
}
