package bot

import (
	"testing"

	"heliosdrift/game"
)

// play runs one game from a seed with the trader policy.
func play(seed int64, frames int) (*game.Game, *Script) {
	game.SetSeed(seed)
	g := game.New()
	g.NewGameMode(game.ModeAdvanced, "BOT")
	return g, Run(g, frames, NewTrader(seed))
}

// TestSameSeedSameRun is the foundation of the whole harness: a world plus a
// policy plus a seed must always play out the same way.
func TestSameSeedSameRun(t *testing.T) {
	_, first := play(7, 900)
	_, second := play(7, 900)
	if first.Hash != second.Hash {
		t.Fatalf("same seed gave different runs: %016x vs %016x", first.Hash, second.Hash)
	}
}

// TestDifferentSeedsDiffer guards against the opposite mistake: a seed that is
// ignored would make every run identical and hide bugs.
func TestDifferentSeedsDiffer(t *testing.T) {
	_, a := play(7, 900)
	_, b := play(8, 900)
	if a.Hash == b.Hash {
		t.Fatal("two seeds produced the same run; the seed is not being used")
	}
}

// TestReplayReproducesRun replays a recording and demands the same world.
func TestReplayReproducesRun(t *testing.T) {
	g, sc := play(11, 1200)
	game.SetSeed(sc.Seed)
	rep := game.New()
	rep.NewGameMode(game.ModeAdvanced, "BOT")
	got, diverge := Replay(rep, sc)
	if diverge >= 0 {
		t.Fatalf("replay drifted at frame %d", diverge)
	}
	if got != sc.Hash {
		t.Fatalf("replay hash %016x, recording %016x", got, sc.Hash)
	}
	if g.State != sc.Report.Ended {
		t.Fatalf("state after run %q, report says %q", g.State, sc.Report.Ended)
	}
}
