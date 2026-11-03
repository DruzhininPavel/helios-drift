package ui

import (
	"testing"

	"heliosdrift/galaxy"
	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// The range ring must stay inside the map frame: with a long range tank and the
// cursor near an edge it used to be drawn across the panel and the screen.
func TestRangeRingStaysOnTheMap(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	g := game.New()
	g.NewGame("CLIP")
	g.State = "flight"
	g.EnterSystem(3, "launch", 0, false)
	g.Cmd.Equip["fueltank"] = true
	g.Overlay = "chart"

	// park the cursor in the corner of the map so the ring hangs over the edge
	gal := galaxy.Get(g.Cmd.Galaxy)
	corner, best := 0, -1
	for _, s := range gal.Systems {
		if s.X+s.Y > best {
			best = s.X + s.Y
			corner = s.ID
		}
	}
	S.ChartCursor = corner

	fb := render.NewFB(320, 224)
	Draw(fb, 1.0/60.0)

	ringCol := render.RGB(0x3fbcd8)
	outside := 0
	for y := 0; y < 224; y++ {
		for x := 0; x < 320; x++ {
			if fb.Buf[y*320+x] == ringCol && (x < 14 || x > 190 || y < 34 || y > 198) {
				outside++
			}
		}
	}
	if outside > 0 {
		t.Fatalf("%d ring pixels outside the map frame", outside)
	}
}
