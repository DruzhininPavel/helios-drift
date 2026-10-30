package bot

import (
	"testing"

	"heliosdrift/galaxy"
	"heliosdrift/game"
)

// TestNextHopMakesProgress checks the router: repeated hops must close the
// distance to a destination instead of stalling somewhere in the middle.
func TestNextHopMakesProgress(t *testing.T) {
	game.SetSeed(7)
	g := game.New()
	g.NewGameMode(game.ModeAdvanced, "BOT")
	g.Cmd.Fuel = 7
	pol := NewTrader(7)
	gal := galaxy.Get(0)
	target := 137

	at := g.Cmd.System
	dist := galaxy.Distance(gal.Systems[at], gal.Systems[target])
	t.Logf("start %d, target %d, %.2f ly", at, target, dist)
	for step := 0; step < 12; step++ {
		hop, ok := pol.nextHop(g, target)
		if !ok {
			t.Fatalf("no hop from %d with %.1f ly of fuel", at, g.Cmd.Fuel)
		}
		d := galaxy.Distance(gal.Systems[at], gal.Systems[hop])
		g.Cmd.System = hop
		at = hop
		dist = galaxy.Distance(gal.Systems[at], gal.Systems[target])
		t.Logf("  hop %d (%.2f ly) -> %.2f ly left", hop, d, dist)
		if dist < 0.01 {
			return
		}
	}
	t.Fatalf("still %.2f ly away after 12 hops", dist)
}
