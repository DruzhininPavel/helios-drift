package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// TestEveryHullHasItsOwnCockpit is the point of the table: five hulls, five
// different canopies, and none of them falling back to the default.
func TestEveryHullHasItsOwnCockpit(t *testing.T) {
	seen := map[string]string{}
	for _, sc := range game.ShipYard {
		st, ok := cockpitStyles[sc.ID]
		if !ok {
			t.Fatalf("%s has no cockpit", sc.ID)
		}
		key := string(rune(st.canopyY)) + string(rune(st.bezel)) + string(rune(st.footW)) + st.beam + st.glass
		if other, dup := seen[key]; dup {
			t.Errorf("%s and %s share a cockpit", sc.ID, other)
		}
		seen[key] = sc.ID
		if st.bezel >= game.ViewH {
			t.Errorf("%s: dashboard at %d covers the HUD", sc.ID, st.bezel)
		}
		if st.canopyY < 60 {
			t.Errorf("%s: canopy starts at %d, it would cut the sight", sc.ID, st.canopyY)
		}
	}
}

// TestCockpitCanBeSwitchedOff: the setting has to leave the view clean.
func TestCockpitCanBeSwitchedOff(t *testing.T) {
	fb := render.NewFB(320, 224)
	cockpit(fb, "cobra3")
	on := 0
	for _, p := range fb.Buf {
		if p != 0 {
			on++
		}
	}
	if on == 0 {
		t.Fatal("cockpit drew nothing")
	}
	was := lang.Cfg.NoCockpit
	lang.Cfg.NoCockpit = true
	defer func() { lang.Cfg.NoCockpit = was }()
	fb2 := render.NewFB(320, 224)
	cockpit(fb2, "cobra3")
	for _, p := range fb2.Buf {
		if p != 0 {
			t.Fatal("cockpit drew while switched off")
		}
	}
}
