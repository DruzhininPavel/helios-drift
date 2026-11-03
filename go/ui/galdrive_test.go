package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// A pilot who bought the drive has to be told where the button is: the chart
// footer names it when it is fitted, and pressing it docked answers instead of
// doing nothing.
func TestGalacticDriveTellsYouHowToUseIt(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	g := game.New()
	g.NewGame("GAL")
	g.State = "flight"
	g.EnterSystem(3, "launch", 0, false)
	g.Cmd.Equip["galdrive"] = true
	g.Overlay = "chart"
	Draw(render.NewFB(320, 224), 1.0/60.0)
	if !containsText(lang.TranslateLine(Hint), "G ") {
		t.Fatalf("the chart says %q and does not name the galactic jump", Hint)
	}

	// docked: a message, not silence
	g.Overlay = ""
	g.State = "docked"
	g.Screen = "launch"
	if g.GalacticJump() {
		t.Fatal("jumped from a station")
	}
	found := false
	for _, m := range g.Messages {
		if containsText(m.Text, "ПОЛЁТЕ") || containsText(m.Text, "FLIGHT") {
			found = true
		}
	}
	if !found {
		t.Fatalf("no message explained the refusal: %v", g.Messages)
	}

	// in flight: the jump happens and the drive is spent
	g.State = "flight"
	before := g.Cmd.Galaxy
	if !g.GalacticJump() {
		t.Fatal("the drive refused a jump in flight")
	}
	if g.Cmd.Galaxy == before {
		t.Fatal("the galaxy did not change")
	}
	if g.Cmd.Equip["galdrive"] {
		t.Fatal("the drive survived")
	}
}
