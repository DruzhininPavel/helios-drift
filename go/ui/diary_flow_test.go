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
	for i := 0; i < 40; i++ {
		g.Note("trade", "ALLOYS", 1, 100)
	}
	// the diary is the last tab: arrow over to it
	for i := 0; i < len(logFilters)-1; i++ {
		Key("ArrowRight")
	}
	if S.LogFilter != logDiaryTab() {
		t.Fatalf("filter %d after arrowing right, want the diary tab", S.LogFilter)
	}
	Draw(render.NewFB(320, 224), 1.0/60.0)
	if !containsText(lang.TranslateLine(Hint), "ENTER") {
		t.Fatalf("the diary tab has no hint: %q", Hint)
	}
	Key("Enter")
	if !S.DiaryOpen {
		t.Fatal("enter did not open a page")
	}
	Draw(render.NewFB(320, 224), 1.0/60.0)
	Key("Escape")
	if S.DiaryOpen {
		t.Fatal("escape did not close the page")
	}
	Key("Escape")
	if g.Screen != "launch" {
		t.Fatalf("escape left the log for %q", g.Screen)
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

// The tab row grew a sixth tab: it still has to fit inside the plate.
func TestLogTabsFit(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	x := 12
	for _, name := range logFilters {
		w := render.TextWidth(lang.TranslateLine(name), 1)
		x += w + 6
	}
	if x > 310 {
		t.Fatalf("the tab row ends at %d, the plate at 310", x)
	}
}
