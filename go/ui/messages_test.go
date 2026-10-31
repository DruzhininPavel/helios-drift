package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/render"
)

// Refusing a contract for lack of hold space answered on a screen that never
// drew messages, so the board looked like it ignored the key. Every screen
// shows them now.
func TestRefusalsAreVisibleOnStationScreens(t *testing.T) {
	g := game.New()
	g.NewGame("TEST")
	g.State = "docked"
	g.Screen = "launch"
	g.Messages = nil

	quiet := render.NewFB(320, 224)
	Draw(quiet, 1.0/60)

	g.Message("NOT ENOUGH CARGO SPACE", "bad")
	loud := render.NewFB(320, 224)
	Draw(loud, 1.0/60)

	diff := 0
	for y := 188; y < 208; y++ {
		for x := 0; x < 320; x++ {
			if quiet.Buf[y*320+x] != loud.Buf[y*320+x] {
				diff++
			}
		}
	}
	if diff == 0 {
		t.Fatal("the message was not drawn on a station screen")
	}
}
