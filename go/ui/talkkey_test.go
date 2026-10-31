package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
)

// The conversation in flight must take its keys: while it is open the cockpit
// is not flying, so digits and arrows belong to the dialogue.
func TestTalkKeysReachTheDialogueInFlight(t *testing.T) {
	lang.Set("en")
	defer lang.Set("ru")
	g := game.New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.State = "flight"
	g.Cmd.Cash = 1000 // no money: the pay options are hidden, 1 is the fight
	if !g.TalkStart("c2_cmdA") {
		t.Fatal("no such dialogue")
	}
	if g.Screen != "talk" {
		t.Fatalf("screen %q, want talk", g.Screen)
	}
	Key("2") // JUMP OUT
	if g.Talk == "c2_cmdA" {
		t.Fatal("the key did not reach the dialogue: still at the first node")
	}
}
