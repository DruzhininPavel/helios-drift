package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
)

// Buying by amount: the prompt takes the keys, Enter buys that many, and the
// log records it as one line.
func TestAmountSelectorBuysWhatItSays(t *testing.T) {
	lang.Set("en")
	defer lang.Set("ru")
	g := game.New()
	g.NewGame("TEST")
	g.State = "docked"
	g.Screen = "market"
	g.Cmd.Cash = 10000
	g.Market[0].Qty = 50
	S.Market = 0
	S.MarketMode, S.MarketQty = "buy", 10

	Key("ArrowUp") // 11
	Key("ArrowUp") // 12
	Key("Enter")
	if g.Cmd.Cargo[0] != 12 {
		t.Fatalf("bought %d tonnes, want 12", g.Cmd.Cargo[0])
	}
	if S.MarketMode != "" {
		t.Fatal("the prompt stayed open after buying")
	}
}

// The market does not narrate success: one purchase is one log line, and the
// screen itself stays quiet about it.
func TestMarketIsQuietAboutSuccess(t *testing.T) {
	lang.Set("en")
	defer lang.Set("ru")
	g := game.New()
	g.NewGame("TEST")
	g.State = "docked"
	g.Screen = "market"
	g.Cmd.Cash = 10000
	g.Market[0].Qty = 50
	before := len(g.Cmd.Log)
	g.Buy(0, 4)
	g.Buy(0, 3)
	if got := len(g.Cmd.Log) - before; got != 1 {
		t.Fatalf("two buys wrote %d log lines, want 1", got)
	}
	// and a refusal still produces a message
	g.Messages = nil
	g.Market[0].Qty = 0
	g.Buy(0, 1)
	bad := false
	for _, m := range g.Messages {
		if m.Kind == "bad" {
			bad = true
		}
	}
	if !bad {
		t.Fatal("a refused purchase said nothing")
	}
}
