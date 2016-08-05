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
	LanguageISO = "NL"

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

	UrlDisplayText     = "De link van %s toont: %s"
	UrlYoutubeDuration = " (speelduur: %s)"
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

	SimpleHelpDesc = "Geeft een verwijzing naar !help"
	SimpleHelpText = "%s, voor een overzicht van alle commandos, type !help of kijk op http://jteeuwen.nl/autimaat.html"
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
	SnoozeDefaultMessage = "%s, het is %s: Snooze tijd!"
	SnoozeMessagePrefix  = "%s, het is %s: "
	SnoozeAlarmSet       = "%s, het alarm is ingesteld. Je kunt het verwijderen met: !unsnooze %s"
	SnoozeAlarmUnset     = "%s, het alarm is verwijderd."

	UnsnoozeName   = "unsnooze"
	UnsnoozeDesc   = "Verwijder een bestaand snooze alarm. Je kunt alleen alarmen verwijderen die je zelf hebt gemaakt."
	UnsnoozeIDName = "id"
	UnsnoozeIDDesc = "De code voor het alarm dat verwijderd dient te worden."
)

const (
	WeatherName = "weer"
	WeatherDesc = "Toon het huidige weer voor een specifieke lokatie. Verzorgd door: https://www.wunderground.com"

	ForecastName = "weerfc"
	ForecastDesc = "Toon een 3-daagse weersvoorspelling voor een specifieke lokatie. Verzorgd door: https://www.wunderground.com"

	WeatherNope = "%s, kijk maar uit het raam."

	WeatherLocationName               = "lokatie"
	WeatherLocationDesc               = "Naam van de lokatie in kwestie. Dit is een dorp of stad, optioneel gevolgd door een land code en/of staat/provincie code. Bijv.: \"Eindhoven\", \"Amsterdam NL\", \"London CA ON\""
	WeatherLocationsText              = "%s: er zijn meerdere lokaties met deze naam: %s"
	WeatherNotAvailable               = "Het weerbericht is momenteel niet beschikbaar."
	WeatherCurrentWeatherText         = "%s, in %s is het %d°C, %s, luchtdruk: %s hPa, luchtvochtigheid: %s, wind: %.1f km/u uit richting: %s."
	WeatherCurrentWeatherNotAvailable = "%s, er is momenteel geen weer beschikbaar voor deze lokatie."
	WeatherForecastNotAvailable       = "Er is momenteel geen weersvoorspelling voor %s."
	WeatherForecastText               = "Weersvoorspelling voor %s:"
)

const (
	EightballName         = "8ball"
	EightballDesc         = "Stel de Magische 8ball een vraag en bereid je voor op ongezouten waarheid."
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

	// GiveDesc has the name of a command passed to it as a format parameter.
	// Formulate the sentence so this makes sense. E.g.:
	//
	//    "Give yourself or someone else <X>"
	//
	GiveDesc = "Geef jezelf of iemand anders %s"

	SmokeName    = "peuk"
	BeerName     = "bier"
	WineName     = "wijn"
	CoffeeName   = "koffie"
	TeaName      = "thee"
	LemonadeName = "fris"
	HugName      = "knuffel"
	PetName      = "aai"
)

// The strings below should be written as if part of an action.
// E.g.: "/me <something something...>"

var HugAnswers = []string{
	"geeft %s een stevige knuffel.",
	"geeft %s een lieve knuffel.",
	"knuffelt %s eens flink.",
}

var PetAnswers = []string{
	"geeft %s een aai over de bol.",
	"aait %s zachtjes over het hoofd.",
}

var SmokeAnswers = []string{
	"geeft %s de asbak en aansteker.",
	"steekt een cubaan op en overhandigt deze aan %s.",
	"geeft %s een nieuw pakje cigaretten.",
	"geeft %s een nieuw pakje shag.",
}

var BeerAnswers = []string{
	"geeft %s een lekker koud biertje.",
	"schenkt %s een trappistje naar keuze.",
	"sluit een nieuw fust aan en tapt een vers glas voor %s. Proost!",
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

const (
	DefineName        = "watis"
	DefineDesc        = "Geef de definitie van de opgegeven term."
	DefineTermName    = "term"
	DefineTermDesc    = "De term waarvoor de definitie gewenst is."
	DefineNotFound    = "%s, de term %s is niet bekend."
	DefineDisplayText = "%s: [%d] %s"

	DefinitionsName    = "definities"
	DefinitionsDesc    = "Geeft een lijst met bekende zoek termen."
	DefinitionsDisplay = "Ik ken %s termen:"

	AddDefineName           = "definieer"
	AddDefineDesc           = "Voeg een nieuwe watis definitie toe."
	AddDefineTermName       = "term"
	AddDefineTermDesc       = "De term die toegevoegd dient te worden."
	AddDefineDefinitionName = "definitie"
	AddDefineDefinitionDesc = "De definitie van de term."
	AddDefineAllreadyUsed   = "Deze definitie voor de term %s bestaat al."
	AddDefineDisplayText    = "De term %s is toegevoegd."

	RemoveDefineName         = "ondefinieer"
	RemoveDefineDesc         = "Verwijder een bestaande watis definitie."
	RemoveDefineTermName     = "term"
	RemoveDefineTermDesc     = "De term die verwijdert dient te worden."
	RemoveDefineIndexName    = "index"
	RemoveDefineIndexDesc    = "Het nummer van de definitie die men wil verwijderen."
	RemoveDefineDisplayText1 = "De term %s is verwijderd."
	RemoveDefineDisplayText2 = "De sub-term %s, %d is verwijderd."
	RemoveDefineNotFound     = "De term %s is niet bekend."
	RemoveDefineInvalidIndex = "De opgeven index %s is niet geldig."
)
