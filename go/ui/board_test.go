package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
)

// The board is one list: story first, then offers, then jobs already running.
// Accepting an offer must not make the row jump somewhere else - it gains a
// marker and stays where the player found it.
func TestBoardKeepsRowsInPlace(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	g := game.New()
	g.NewGame("TEST")
	g.State = "docked"
	g.Screen = "contracts"
	g.Offers = []game.Contract{
		{Kind: game.ContractCourier, Good: 0, Qty: 2, System: 42, Galaxy: g.Cmd.Galaxy, Pay: 500, Due: g.Cmd.Day + 5},
	}
	S.Contract = 1 // the offer, after the story row
	before := len(boardRows())
	Key("Enter")
	rows := boardRows()
	if len(rows) != before {
		t.Fatalf("board changed size on accept: %d -> %d", before, len(rows))
	}
	if !g.Offers[0].Taken {
		t.Fatal("the offer was not marked taken")
	}
	if rows[S.Contract].kind != "offer" {
		t.Fatalf("cursor now points at %q", rows[S.Contract].kind)
	}
	if len(g.Cmd.Contracts) != 1 {
		t.Fatalf("%d contracts running", len(g.Cmd.Contracts))
	}
}

// The description describes the job in words, with the days left on it.
func TestDescriptionSaysWhatToDo(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	g := game.New()
	g.NewGame("TEST")
	c := game.Contract{Kind: game.ContractCourier, Good: 0, Qty: 3, System: 42, Galaxy: g.Cmd.Galaxy, Pay: 700, Due: g.Cmd.Day + 4}
	lines := contractWhat(c)
	if len(lines) < 4 {
		t.Fatalf("%d lines, want a real description", len(lines))
	}
	joined := ""
	for _, l := range lines {
		joined += l + " | "
	}
	for _, want := range []string{"ДОСТАВЬ", "4", "КУДА"} {
		found := false
		for _, l := range lines {
			if len(l) > 0 && (want == "4" && containsDigit(l, '4') || want != "4" && contains(l, want)) {
				found = true
			}
		}
		if !found {
			t.Errorf("description has no %q in: %s", want, joined)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()
}

func containsDigit(s string, d byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == d {
			return true
		}
	}
	return false
}
