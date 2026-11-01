package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// A world with no planet must draw: the station screens paint the sky behind them.
func TestFlightWithoutAPlanetDraws(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	fb := render.NewFB(320, 224)
	g := game.New()
	g.NewGame("NOPLANET")
	g.State = "flight"
	g.Screen = ""
	g.EnterSystem(3, "launch", 0, false)
	g.World.Planet = game.Planet{} // no planet at all
	g.World.PlanetFrames = nil
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("drawing a system with no planet panics: %v", r)
		}
	}()
	Draw(fb, 1.0/60.0)
}
