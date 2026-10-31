package game

import (
	"math"
	"testing"
)

// hullGame starts a commander flying the given hull.
func hullGame(t *testing.T, ship string) *Game {
	t.Helper()
	g := New()
	g.NewGameMode(ModeAdvanced, "BOT")
	g.Cmd.Ship = ship
	g.Launch()
	return g
}

// TestHullsTurnDifferently is the point of the per-hull rates: the same stick
// input must turn a Fer-de-Lance much further than a Python.
func TestHullsTurnDifferently(t *testing.T) {
	turn := map[string]float64{}
	for _, id := range []string{"python", "cobra3", "ferdelance"} {
		g := hullGame(t, id)
		g.Player.Speed = 0
		start := g.Player.Mat.Col(2).Mul(-1)
		for i := 0; i < 30; i++ { // half a second of full pitch
			g.In = InputState{Pitch: 1}
			g.Update(1.0 / 60)
		}
		end := g.Player.Mat.Col(2).Mul(-1)
		angle := math.Acos(clampDot(start.Dot(end)))
		turn[id] = angle
		t.Logf("%-11s turned %.3f rad in 0.5 s", id, angle)
	}
	if turn["ferdelance"] <= turn["cobra3"]*1.1 {
		t.Errorf("fer-de-lance turns %.3f, cobra %.3f: no agility difference", turn["ferdelance"], turn["cobra3"])
	}
	if turn["python"] >= turn["cobra3"]*0.9 {
		t.Errorf("python turns %.3f, cobra %.3f: the hauler is not sluggish", turn["python"], turn["cobra3"])
	}
}

// TestHullsAccelerateDifferently checks the per-hull acceleration: in one
// second a ship reaches min(vmax, accel).
func TestHullsAccelerateDifferently(t *testing.T) {
	for _, id := range []string{"adder", "python", "ferdelance"} {
		h := ShipFor(id)
		g := hullGame(t, id)
		g.Player.Speed = 0
		for i := 0; i < 60; i++ {
			g.In = InputState{Accel: 1}
			g.Update(1.0 / 60)
		}
		want := math.Min(h.Speed, h.Accel)
		if math.Abs(g.Player.Speed-want) > 1 {
			t.Errorf("%s: speed %.1f after a second, want %.1f", id, g.Player.Speed, want)
		}
	}
}

// TestHullsHaveTheirOwnShields makes sure the shield bar is not the same for a
// hauler and a fighter.
func TestHullsHaveTheirOwnShields(t *testing.T) {
	seen := map[float64]bool{}
	for i := range ShipYard {
		h := ShipYard[i]
		g := hullGame(t, h.ID)
		if g.Player.ShieldF != h.Shield || g.Player.ShieldA != h.Shield {
			t.Errorf("%s: shields %.0f/%.0f, hull says %.0f",
				h.Name, g.Player.ShieldF, g.Player.ShieldA, h.Shield)
		}
		if g.Player.Energy != h.Energy {
			t.Errorf("%s: energy %.0f, hull says %.0f", h.Name, g.Player.Energy, h.Energy)
		}
		seen[h.Shield] = true
	}
	if len(seen) < 3 {
		t.Errorf("only %d distinct shield values across the yard", len(seen))
	}
}

// TestLaserLimitedByHull checks that a hull that cannot carry a military laser
// refuses it, and that a hull good enough keeps it.
func TestLaserLimitedByHull(t *testing.T) {
	g := New()
	g.NewGameMode(ModeAdvanced, "BOT")
	g.Cmd.Ship = "python"
	g.Cmd.Cash = 1e6
	g.World.Sys.Tech = 15 // a yard that sells everything: the hull decides
	g.BuyEquipment("lasermilitary")
	if g.Cmd.Laser == 2 {
		t.Fatal("python mounted a military laser")
	}

	g2 := New()
	g2.NewGameMode(ModeAdvanced, "BOT")
	g2.Cmd.Ship = "ferdelance"
	g2.Cmd.Cash = 1e6
	g2.World.Sys.Tech = 15
	g2.BuyEquipment("lasermilitary")
	if g2.Cmd.Laser != 2 {
		t.Fatalf("fer-de-lance could not mount a military laser (laser %d)", g2.Cmd.Laser)
	}
}

// TestSlowerHullDowngradesLaser covers the trade-in: a big gun in a small hull
// comes out of the ship it is moved into.
func TestSlowerHullDowngradesLaser(t *testing.T) {
	g := New()
	g.NewGameMode(ModeAdvanced, "BOT")
	g.Cmd.Ship = "cobra3"
	g.Cmd.Laser = 2
	g.Cmd.Cash = 1e6
	if !g.BuyShip("python") {
		t.Fatal("could not trade for a python")
	}
	if g.Cmd.Laser > ShipFor("python").MaxLaser {
		t.Fatalf("python kept a class %d laser", g.Cmd.Laser)
	}
	if g.Player.ShieldF != ShipFor("python").Shield {
		t.Fatalf("shields not reset to the new hull: %.0f", g.Player.ShieldF)
	}
}

func clampDot(v float64) float64 {
	if v > 1 {
		return 1
	}
	if v < -1 {
		return -1
	}
	return v
}
