package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
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

// The summary row is always offered, and an empty log says so instead of looking
// like a broken screen.
func TestJournalRowAndEmptyLog(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	g := game.New()
	g.NewGame("JROW")
	g.State = "docked"
	g.Screen = "launch"
	found := false
	for _, it := range menuItems() {
		if it == "TASKS" {
			found = true
		}
	}
	if !found {
		t.Fatal("no summary row with nothing taken on")
	}
	// a fresh commander may already have a line or two (arrival in system)
	g.Cmd.Log = nil
	g.Cmd.Log = g.Cmd.Log[:0]
	g.Screen = "log"
	Draw(render.NewFB(320, 224), 1.0/60.0)
	if Hint == "" {
		t.Fatal("the empty log says nothing at all")
	}
}
