package game

import "testing"

// A new commander is told who he is before he is shown a market.
func TestIntroRunsOnce(t *testing.T) {
	g := New()
	g.NewGame("INTRO")
	g.StartIntro()
	if g.Screen != "talk" || g.Talk != "intro_who" {
		t.Fatalf("new game opened %q/%q", g.Screen, g.Talk)
	}
	// three screens, in order, then the station
	chooseText(t, g, "CONTINUE")
	if g.Talk != "intro_where" {
		t.Fatalf("second screen is %q", g.Talk)
	}
	chooseText(t, g, "CONTINUE")
	if g.Talk != "intro_what" {
		t.Fatalf("third screen is %q", g.Talk)
	}
	chooseText(t, g, "ENTER THE STATION")
	if g.Screen == "talk" {
		t.Fatal("the prologue did not close")
	}
	if !g.Cmd.Flags["INTRO_SEEN"] {
		t.Fatal("the prologue was not recorded as seen")
	}
	// and never again
	g.StartIntro()
	if g.Screen == "talk" {
		t.Fatal("the prologue opened twice")
	}
}

// Escape skips the rest without marking anything as read later than it was.
func TestIntroIsSkippable(t *testing.T) {
	g := New()
	g.NewGame("INTRO")
	g.StartIntro()
	g.EndTalk()
	if g.Screen == "talk" {
		t.Fatal("Escape did not close the prologue")
	}
}
