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

// Arrows pick an answer and Enter takes it, so the dialogue works with no digits
// at all - which is what a gamepad has.
func TestArrowsChooseAndEnterConfirms(t *testing.T) {
	lang.Set("en")
	defer lang.Set("ru")
	g := game.New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.State = "flight"
	g.Cmd.Cash = 1000 // the two pay options are hidden: 1 is fight, 2 is jump
	g.TalkStart("c2_cmdA")

	Key("ArrowDown") // the second answer
	if S.TalkOpt != 1 {
		t.Fatalf("cursor at %d after one press", S.TalkOpt)
	}
	Key("ArrowUp") // back to the first
	Key("ArrowUp") // wraps to the last
	if S.TalkOpt != 1 {
		t.Fatalf("wrapping put the cursor at %d", S.TalkOpt)
	}
	Key("Enter")
	if g.Talk == "c2_cmdA" {
		t.Fatal("Enter did not take the answer")
	}
}

// Scrolling still belongs to the page keys, and does not move the cursor.
func TestPageKeysScrollTheTextOnly(t *testing.T) {
	lang.Set("en")
	defer lang.Set("ru")
	g := game.New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.State = "flight"
	g.TalkStart("c7_hive")
	before := S.TalkOpt
	Key("PageDown")
	if S.TalkOpt != before {
		t.Fatalf("scrolling moved the cursor to %d", S.TalkOpt)
	}
	if S.TalkScroll == 0 {
		t.Fatal("Page Down did not scroll a long monologue")
	}
}
