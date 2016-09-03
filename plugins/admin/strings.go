// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package admin

const (
	// ref: https://godoc.org/time#Time.Format
	TextDateFormat = "02 January, 2006"
	TextTimeFormat = "15:04 MST"

	TextHelpName    = "help"
	TextHelpDisplay = "%s, voor een overzicht van de commandos die ik herken, kijk op http://jteeuwen.nl/autimaat.html"

	TextNickName     = "nick"
	TextNickNickName = "naam"
	TextNickPassName = "wachtwoord"

	TextJoinName         = "join"
	TextJoinChannelName  = "kanaal"
	TextJoinKeyName      = "sleutel"
	TextJoinPasswordName = "wachtwoord"

	TextPartName        = "part"
	TextPartChannelName = "kanaal"

	TextReloadName = "herstart"

	TextAuthListName    = "bazen"
	TextAuthListDisplay = "De beheerders zijn: %s"

	TextAuthorizeName     = "baas"
	TextAuthorizeMaskName = "hostmask"
	TextAuthorizeDisplay  = "Gebruiker %q is toegevoegd aan de beheerderslijst."

	TextDeauthorizeName     = "ontbaas"
	TextDeauthorizeMaskName = "hostmask"
	TextDeauthorizeDisplay  = "Gebruiker %q is verwijderd van de beheerderslijst."

	TextVersionName    = "versie"
	TextVersionDisplay = "%s, ik ben %s, versie %s. Mijn laatste revisie was op %s, om %s. De laatste herstart was %s uur geleden."

	TextLogName      = "log"
	TextLogValueName = "status"
	TextLogEnabled   = "Logging is ingeschakeld."
	TextLogDisabled  = "Logging is uitgeschakeld."
)
