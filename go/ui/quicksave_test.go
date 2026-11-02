package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
)

// F5 writes to the slot the commander came from, F7 reads it back.
func TestQuickSaveFollowsTheSlot(t *testing.T) {
	lang.Set("en")
	defer lang.Set("ru")
	dir := t.TempDir()
	game.SaveDirOverride = dir
	defer func() { game.SaveDirOverride = "" }()

	g := game.New()
	g.NewGame("QUICK")
	g.State = "docked"
	g.Screen = "launch"
	g.Cmd.Cash = 1234
	if g.SaveSlot(2) != true {
		t.Fatal("could not write slot 3")
	}
	// the commander now belongs to slot 2
	if g.Slot() != 2 {
		t.Fatalf("slot %d after saving to 2", g.Slot())
	}
	g.Cmd.Cash = 9999
	Key("F5")
	g.Cmd.Cash = 1
	Key("F7")
	if g.Cmd.Cash != 9999 {
		t.Fatalf("cash %.0f after quick load, want the quick save", g.Cmd.Cash)
	}
	if g.Slot() != 2 {
		t.Fatalf("quick load moved the commander to slot %d", g.Slot())
	}
}

// Docking remembers where the pilot would come back to, and dying offers it.
func TestRespawnAtTheLastStation(t *testing.T) {
	lang.Set("en")
	defer lang.Set("ru")
	dir := t.TempDir()
	game.SaveDirOverride = dir
	defer func() { game.SaveDirOverride = "" }()

	g := game.New()
	g.NewGame("DEATH")
	g.State = "docked"
	g.Screen = "launch"
	g.Cmd.Cash = 5000
	g.DockSalt = 0
	g.Docked()
	where := g.SnapshotWhere()
	if where == "" {
		t.Fatal("docking did not take a snapshot")
	}
	// die in flight
	g.Launch()
	g.Cmd.Cash = 10
	g.Player.Invulnerable = 0 // launching grants a few seconds of grace
	g.DamagePlayer(9999, true, "laser")
	if g.State != "dead" {
		t.Fatalf("state %q after a fatal hit", g.State)
	}
	Key("1")
	if g.State != "docked" {
		t.Fatalf("respawn left the game in %q", g.State)
	}
	if g.Cmd.Cash != 5000 {
		t.Fatalf("cash %.0f after respawn, want the state at docking", g.Cmd.Cash)
	}
	if g.Cmd.Deaths != 1 {
		t.Fatalf("deaths %d after respawn, want the death counted", g.Cmd.Deaths)
	}
}
