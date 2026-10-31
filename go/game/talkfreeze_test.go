package game

import "testing"

// A conversation in flight stops the world: nothing moves, nothing shoots and
// nothing spawns until the answer is given. Before this, pirates kept firing
// while the player was reading, and the answer buttons looked dead.
func TestTalkFreezesTheWorld(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.Cmd.Campaign.Stage = 1
	g.Cmd.Campaign.Active = true
	g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
	g.Cmd.Campaign.Target = g.Cmd.System
	g.campaignArrive()
	if g.Screen != "talk" {
		t.Fatalf("no conversation on arrival: screen %q", g.Screen)
	}
	before := g.Time
	enemies := len(g.World.Entities)
	pos := g.Player.Pos
	for i := 0; i < 120; i++ {
		g.Update(1.0 / 60.0)
	}
	if g.Time != before {
		t.Fatalf("time advanced during the conversation: %.3f -> %.3f", before, g.Time)
	}
	if len(g.World.Entities) != enemies {
		t.Fatalf("the world changed while talking: %d -> %d entities", enemies, len(g.World.Entities))
	}
	if g.Player.Pos != pos {
		t.Fatal("the ship drifted during the conversation")
	}
	// and answering it puts the world back in motion
	g.TalkChoose(0)
	if g.Screen == "talk" {
		g.TalkChoose(0)
	}
	later := g.Time
	g.Update(1.0 / 60.0)
	if g.Time == later {
		t.Fatal("the world stayed frozen after the conversation")
	}
}
