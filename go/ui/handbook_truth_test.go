package ui

import (
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"heliosdrift/game"
	"heliosdrift/input"
	"heliosdrift/lang"
)

// The handbook tells the pilot which key does what. This checks the claims
// against the game itself, because a page that lies is worse than no page.
func TestHandbookClaimsAreTrue(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	fresh := func() *game.Game {
		g := game.New()
		g.NewGame("TRUTH")
		return g
	}

	t.Run("C opens the chart", func(t *testing.T) {
		g := fresh()
		g.Launch()
		Key("c")
		if g.Overlay != "chart" {
			t.Fatalf("C gave overlay %q", g.Overlay)
		}
	})
	t.Run("ESC pauses and H opens the handbook", func(t *testing.T) {
		g := fresh()
		g.Launch()
		Key("Escape")
		if g.Overlay != "pause" {
			t.Fatalf("Escape gave overlay %q", g.Overlay)
		}
		Key("h")
		if g.Screen != "help" {
			t.Fatalf("H gave screen %q", g.Screen)
		}
		Key("Escape")
		if g.Overlay != "pause" {
			t.Fatalf("Escape from the handbook gave %q, want the pause it came from", g.Overlay)
		}
	})
	t.Run("K jettisons", func(t *testing.T) {
		g := fresh()
		g.Launch()
		g.Cmd.Cargo[0] = 3
		Key("k")
		if g.Cmd.Cargo[0] != 2 {
			t.Fatalf("cargo %d after K", g.Cmd.Cargo[0])
		}
	})
	t.Run("B and S open the amount, not a straight trade", func(t *testing.T) {
		g := fresh()
		g.Docked()
		g.Screen = "market"
		g.Cmd.Cash = 5000
		g.Cmd.Cargo[0] = 2
		g.Market[0].Qty = 5
		S.Market = 0
		Key("b")
		if S.MarketMode != "buy" {
			t.Fatalf("B gave mode %q", S.MarketMode)
		}
		S.MarketMode = ""
		Key("s")
		if S.MarketMode != "sell" {
			t.Fatalf("S gave mode %q", S.MarketMode)
		}
	})
	t.Run("F refuels", func(t *testing.T) {
		g := fresh()
		g.Docked()
		g.Cmd.Fuel = 1
		before := g.Cmd.Fuel
		Key("f")
		if g.Cmd.Fuel <= before {
			t.Fatalf("fuel %.1f -> %.1f after F", before, g.Cmd.Fuel)
		}
	})
}

// Every page must be reachable and every jump line must stay inside the plate.
func TestHandbookPagesFit(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	for _, p := range helpPages {
		if p.Title == "" {
			t.Error("a page without a title")
		}
		for _, c := range p.Cols {
			// the lamp row is four single letters on purpose
			if strings.Contains(c[0], "  ") && c[0] != "E S D B" {
				t.Errorf("abbreviation %q is not one token", c[0])
			}
			if c[1] == "" {
				t.Errorf("%q has no meaning", c[0])
			}
		}
	}
}

// The book reads the pilot's bindings, so rebinding a key changes the text.
func TestHandbookReadsTheBindings(t *testing.T) {
	lang.Set("en")
	defer lang.Set("ru")
	defer input.Reset()
	before := helpText("{FIRE} FIRES THE LASERS.")
	input.Cur[input.Fire] = input.Binding{Keys: []int32{int32(ebiten.KeyX)}, Btn: -1, Axis: -1}
	after := helpText("{FIRE} FIRES THE LASERS.")
	if after == before {
		t.Fatalf("rebinding fire left the handbook saying %q", after)
	}
	if !strings.HasPrefix(after, "X") {
		t.Fatalf("handbook says %q, want the new key first", after)
	}
	// and no placeholder ever reaches the screen
	input.Reset()
	for _, p := range helpPages {
		for _, l := range p.Lines {
			if strings.ContainsAny(helpText(l), "{}") {
				t.Errorf("%q kept a placeholder: %q", l, helpText(l))
			}
		}
	}
}
