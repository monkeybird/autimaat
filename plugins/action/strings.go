// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package action

const TextUserName = "wie"

// action defines a single action with a set of possible replies.
// One of which will be chosen at random, by the bot.
//
// The answers should be written as if part of an action.
// E.g.: "/me <something something...>"
type action struct {
	Names   []string // Name by which the action is invoked.
	Answers []string // Possible set of replies for this action.
}

// TextActions defines all known actions.
var TextActions = []action{
	{
		[]string{"peuk"},
		[]string{
			"geeft %s de asbak en aansteker.",
			"steekt een cubaan op en overhandigt deze aan %s.",
			"geeft %s een nieuw pakje sigaretten.",
			"geeft %s een nieuw pakje shag.",
		},
	},
	{
		[]string{"bier"},
		[]string{
			"geeft %s een lekker koud biertje.",
			"schenkt %s een trappistje naar keuze.",
			"sluit een nieuw fust aan en tapt een vers glas voor %s. Proost!",
		},
	},
	{
		[]string{"wijn"},
		[]string{
			"geeft %s een glaasje rode wijn.",
			"geeft %s een glaasje witte wijn.",
			"geeft %s de wijnkaart.",
		},
	},
	{
		[]string{"port"},
		[]string{
			"geeft %s een glaasje ruby Port.",
			"schenkt %s een glaasje Port naar keuze.",
		},
	},
	{
		[]string{"whiskey"},
		[]string{
			"duikt de drankkast in om %s de oudste fles whiskey te halen die hij kan vinden.",
			"blaast het stof van de Whiskeyfles en schenkt %s een glas.",
		},
	},
	{
		[]string{"koffie"},
		[]string{
			"schenkt %s een kopje koffie.",
			"geeft %s een stevige bak leut.",
			"zet de ketel op en rent naar de voorraadkast voor een nieuw pak koffie; Even geduld nog, %s!",
		},
	},
	{
		[]string{"thee"},
		[]string{
			"zet een kopje thee voor %s.",
			"loopt richting de keuken om voor %s de fluitketel op het vuur te zetten.",
			"brengt %s de doos met alle theesmaakjes; kies maar iets uit.",
			"geeft %s een kopje thee ;)",
		},
	},
	{
		[]string{"fris"},
		[]string{
			"schenkt %s een groot glas koude limonade.",
			"geeft %s de fles cola aan.",
			"geeft %s een scheut ijskoud aardbeiensap.",
			"geeft %s een scheut ijskoud bosvruchtensap.",
			"geeft %s een scheut ijskoud citroensap.",
			"geeft %s een scheut ijskoud appelsap.",
			"geeft %s een scheut ijskoud sinaasappelsap.",
		},
	},
	{
		[]string{"knuffel", "knuff", "knuf"},
		[]string{
			"geeft %s een stevige knuffel.",
			"geeft %s een lieve knuffel.",
			"knuffelt %s eens flink.",
		},
	},
	{
		[]string{"aai"},
		[]string{
			"geeft %s een aai over de bol.",
			"aait %s zachtjes over het hoofd.",
		},
	},
	{
		[]string{"muts"},
		[]string{
			"breidt een roze muts met gele stippen en een pluizige pompom en overhandigd deze aan %s.",
			"pakt een zelf gehaakte muts in met kadopapier en geeft het geheel, met een strikje aan %s.",
			"haalt voor %s een oude wintermuts uit de kast.",
		},
	},
}
