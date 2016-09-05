// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package action

const (
	TextGiveUserName = "wie"

	TextSmokeName    = "peuk"
	TextBeerName     = "bier"
	TextWineName     = "wijn"
	TextPortName     = "port"
	TextWhiskeyName  = "whiskey"
	TextCoffeeName   = "koffie"
	TextTeaName      = "thee"
	TextLemonadeName = "fris"
	TextHugName      = "knuffel"
	TextPetName      = "aai"
)

// The strings below should be written as if part of an action.
// E.g.: "/me <something something...>"

var (
	TextHugAnswers = []string{
		"geeft %s een stevige knuffel.",
		"geeft %s een lieve knuffel.",
		"knuffelt %s eens flink.",
	}

	TextPetAnswers = []string{
		"geeft %s een aai over de bol.",
		"aait %s zachtjes over het hoofd.",
	}

	TextSmokeAnswers = []string{
		"geeft %s de asbak en aansteker.",
		"steekt een cubaan op en overhandigt deze aan %s.",
		"geeft %s een nieuw pakje sigaretten.",
		"geeft %s een nieuw pakje shag.",
	}

	TextBeerAnswers = []string{
		"geeft %s een lekker koud biertje.",
		"schenkt %s een trappistje naar keuze.",
		"sluit een nieuw fust aan en tapt een vers glas voor %s. Proost!",
	}

	TextWineAnswers = []string{
		"geeft %s een glaasje rode wijn.",
		"geeft %s een glaasje witte wijn.",
		"geeft %s de wijnkaart.",
	}

	TextPortAnswers = []string{
		"geeft %s een glaasje ruby Port.",
		"schenkt %s een glaasje Port naar keuze.",
	}

	TextWhiskeyAnswers = []string{
		"duikt de drankkast in om %s de oudste fles whiskey te halen die hij kan vinden.",
		"blaast het stof van de Whiskeyfles en schenkt %s een glas.",
	}

	TextCoffeeAnswers = []string{
		"schenkt %s een kopje koffie.",
		"geeft %s een stevige bak leut.",
		"zet de ketel op en rent naar de voorraadkast voor een nieuw pak koffie; Even geduld nog, %s!",
	}

	TextTeaAnswers = []string{
		"zet een kopje thee voor %s.",
		"loopt richting de keuken om voor %s de fluitketel op het vuur te zetten.",
		"brengt %s de doos met alle theesmaakjes; kies maar iets uit.",
		"geeft %s een kopje thee ;)",
	}

	TextLemonadeAnswers = []string{
		"schenkt %s een groot glas koude limonade.",
		"geeft %s de fles cola aan.",
		"geeft %s een scheut ijskoud aardbeiensap.",
		"geeft %s een scheut ijskoud bosvruchtensap.",
		"geeft %s een scheut ijskoud citroensap.",
		"geeft %s een scheut ijskoud appelsap.",
		"geeft %s een scheut ijskoud sinaasappelsap.",
	}
)
