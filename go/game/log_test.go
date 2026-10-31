package game

import (
	"testing"

	"heliosdrift/render"
)

// Trading, fighting and jumping all have to leave a mark: the log is the only
// record a player has of what they did two hours ago.
func TestLogRecordsWhatHappened(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)

	before := len(g.Cmd.Log)
	g.Buy(0, 2)
	g.Sell(0, 1)
	if len(g.Cmd.Log) < before+2 {
		t.Fatalf("trade left no entries: %d -> %d", before, len(g.Cmd.Log))
	}
	if e := g.Cmd.Log[len(g.Cmd.Log)-2]; e.Kind != LogBuy || e.N != 2 {
		t.Fatalf("buy entry is %+v", e)
	}

	e := g.baseShip("viper", render.V3(0, 0, -400), "pirate")
	g.DestroyEntity(e, "laser")
	last := g.Cmd.Log[len(g.Cmd.Log)-1]
	if last.Kind != LogKill {
		t.Fatalf("kill not logged, last entry %+v", last)
	}
	if last.Sys == "" {
		t.Fatal("entry has no system, the log cannot say where it happened")
	}
}

// The log is bounded: a long run must not grow the save without end.
func TestLogStaysBounded(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	for i := 0; i < LogLimit+120; i++ {
		g.Note(LogJump, "SOMEWHERE", 0, 0)
	}
	if len(g.Cmd.Log) != LogLimit {
		t.Fatalf("log holds %d entries, want %d", len(g.Cmd.Log), LogLimit)
	}
}

// And it survives a save and a load, otherwise it is a diary you lose.
func TestLogSurvivesSave(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.Buy(0, 3)
	n := len(g.Cmd.Log)
	if !g.Save() {
		t.Fatal("save failed")
	}
	loaded := New()
	if !loaded.Load() {
		t.Skip("nothing to load in this build")
	}
	if len(loaded.Cmd.Log) < n {
		t.Fatalf("log lost in the save: %d -> %d", n, len(loaded.Cmd.Log))
	}
}

// A system is worth one line the first time, not one per jump, and a busy
// system is one pirate line rather than four.
func TestLogVisitOnceAndEncountersThrottled(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "arrive", 9000, false)
	visits := 0
	for _, e := range g.Cmd.Log {
		if e.Kind == LogVisit {
			visits++
		}
	}
	if visits != 1 {
		t.Fatalf("first arrival wrote %d visit lines, want 1", visits)
	}
	// a new system is another visit; coming back is a jump
	g.EnterSystem(4, "arrive", 9000, false)
	g.EnterSystem(3, "arrive", 9000, false)
	visits, jumps := 0, 0
	for _, e := range g.Cmd.Log {
		switch e.Kind {
		case LogVisit:
			visits++
		case LogJump:
			jumps++
		}
	}
	if visits != 2 || jumps != 1 {
		t.Fatalf("visit lines %d (want 2), jump lines %d (want 1)", visits, jumps)
	}

	// the arrival may have spawned pirates of its own; start the wave count
	// from a clean throttle
	g.notedPirate = false
	before := len(g.Cmd.Log)
	g.SpawnPirates(2)
	g.SpawnPirates(2)
	g.SpawnPirates(2)
	after := len(g.Cmd.Log)
	if after-before != 1 {
		t.Fatalf("three pirate waves wrote %d lines, want 1", after-before)
	}
	if g.Cmd.Log[after-1].Kind != LogPirate {
		t.Fatalf("last entry is %q", g.Cmd.Log[after-1].Kind)
	}
}
