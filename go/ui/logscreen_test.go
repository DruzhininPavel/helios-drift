package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// The log must draw itself: it used to fall through to the station menu, so the
// game said "log" while the screen stayed on the menu and looked frozen.
func TestLogScreenDrawsItself(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	g := game.New()
	g.NewGame("LOGDRAW")
	g.State = "docked"
	g.Screen = "launch"
	fb := render.NewFB(320, 224)
	Draw(fb, 1.0/60.0)
	menu := append([]uint32{}, fb.Buf...)

	g.Screen = "log"
	Draw(fb, 1.0/60.0)
	same := 0
	for i := range menu {
		if menu[i] == fb.Buf[i] {
			same++
		}
	}
	if same == len(menu) {
		t.Fatal("the log screen draws the station menu")
	}
}

// Every screen the game can be on has to have a case in the docked draw, or it
// silently renders the menu instead.
func TestEveryDockedScreenDraws(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	g := game.New()
	g.NewGame("SCREENS")
	g.State = "docked"
	fb := render.NewFB(320, 224)
	for _, screen := range []string{"market", "equip", "chart", "status", "contracts",
		"log", "shipyard", "slots", "closing", "ending"} {
		g.Screen = "launch"
		Draw(fb, 1.0/60.0)
		menu := append([]uint32{}, fb.Buf...)
		g.Screen = screen
		Draw(fb, 1.0/60.0)
		same := 0
		for i := range menu {
			if menu[i] == fb.Buf[i] {
				same++
			}
		}
		if same == len(menu) {
			t.Errorf("screen %q draws the station menu", screen)
		}
	}
}
