package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
)

// Opening the journal and coming back must land on the station menu, not back
// inside the journal - that is what looked like a hung menu.
func TestJournalComesBackToTheMenu(t *testing.T) {
	lang.Set("en")
	defer lang.Set("ru")
	g := game.New()
	g.NewGame("JOUR")
	g.State = "docked"
	g.Screen = "launch"
	g.Cmd.Contracts = append(g.Cmd.Contracts, game.Contract{
		Kind: game.ContractCourier, Good: 0, Qty: 2, System: 3, Galaxy: g.Cmd.Galaxy,
		Pay: 400, Due: g.Cmd.Day + 5,
	})

	openJournal()
	if g.Screen != "journal" {
		t.Fatalf("journal did not open: %q", g.Screen)
	}
	Key("Escape")
	if g.Screen != "launch" {
		t.Fatalf("Escape left the game in %q", g.Screen)
	}

	// and quitting to the title from a station page and coming back lands on the
	// menu rather than the page
	g.Screen = "journal"
	Key("q")
	if g.State != "title" {
		t.Fatalf("Q left the state in %q", g.State)
	}
	ContinueFlight()
	if g.Screen != "launch" {
		t.Fatalf("came back into %q instead of the station menu", g.Screen)
	}
}

// The journal answers Enter as well, because a pad's A is Enter.
func TestJournalEnterLeaves(t *testing.T) {
	lang.Set("en")
	defer lang.Set("ru")
	g := game.New()
	g.NewGame("JOUR")
	g.State = "docked"
	g.Screen = "journal"
	Key("Enter")
	if g.Screen != "launch" {
		t.Fatalf("Enter left the game in %q", g.Screen)
	}
}
