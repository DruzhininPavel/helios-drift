package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// The diary has to be reachable from the log screen with D, and the hint has to
// advertise it on every branch of that screen - including the one with a full
// log, which used to overwrite the hint at the end.
func TestDiaryIsReachableFromTheLog(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	g := game.New()
	g.NewGame("D")
	g.State = "docked"
	g.Screen = "log"
	Draw(render.NewFB(320, 224), 1.0/60.0)
	if !containsText(lang.TranslateLine(Hint), "ДНЕВНИК") {
		t.Fatalf("hint %q does not mention the diary", lang.TranslateLine(Hint))
	}
	// a full log takes the other branch: still has to mention it
	for i := 0; i < 40; i++ {
		g.Note("trade", "ALLOYS", 1, 100)
	}
	Draw(render.NewFB(320, 224), 1.0/60.0)
	if !containsText(lang.TranslateLine(Hint), "ДНЕВНИК") {
		t.Fatalf("hint %q on a full log does not mention the diary", lang.TranslateLine(Hint))
	}
	Key("d")
	if g.Screen != "diary" {
		t.Fatalf("D opened %q", g.Screen)
	}
	Draw(render.NewFB(320, 224), 1.0/60.0)
	Key("Escape")
	if g.Screen != "log" {
		t.Fatalf("escape left the diary for %q", g.Screen)
	}
}

func containsText(hay, needle string) bool {
	if needle == "" {
		return true
	}
	for i := 0; i+len(needle) <= len(hay); i++ {
		if hay[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
