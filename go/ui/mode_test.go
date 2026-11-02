package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
)

// A new commander gets the real game unless he says otherwise.
func TestAdvancedIsTheDefaultMode(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	g := game.New()
	g.NewGame("MODESEL")
	g.State = "title"
	g.Screen = "title"
	openMode()
	if S.Mode != 0 {
		t.Fatalf("cursor starts on row %d", S.Mode)
	}
	Key("Enter")
	if g.Classic() {
		t.Fatal("Enter on the default row started a classic commander")
	}
	// and classic is still one press away
	openMode()
	Key("ArrowDown")
	Key("Enter")
	if !game.G.Classic() {
		t.Fatal("the second row is not classic")
	}
}
