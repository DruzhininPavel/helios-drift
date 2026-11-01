package ui

import (
	"testing"
	"unicode/utf8"

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

// Escape closes the amount prompt and leaves the player in the market, so a
// purchase can be called off; it must not throw him out of the shop.
func TestEscapeClosesThePromptOnly(t *testing.T) {
	lang.Set("en")
	defer lang.Set("ru")
	g := game.New()
	g.NewGame("TEST")
	g.State = "docked"
	g.Screen = "market"
	S.Market = 0
	S.MarketMode, S.MarketQty = "buy", 5

	Key("Escape")
	if S.MarketMode != "" {
		t.Fatal("the prompt stayed open")
	}
	if g.Screen != "market" {
		t.Fatalf("Escape left the market: screen %q", g.Screen)
	}
	// and nothing was bought
	if game.CargoUsed(g.Cmd) != 0 {
		t.Fatal("cancelling a purchase bought something")
	}
	// a second Escape leaves the shop
	Key("Escape")
	if g.Screen != "launch" {
		t.Fatalf("second Escape: screen %q", g.Screen)
	}
}

// Russian must wrap by characters, not bytes: two bytes per letter wrapped
// every line at half the panel and bunched the text against the left edge.
func TestWrapCountsCharactersNotBytes(t *testing.T) {
	src := "ИХ КОРАБЛЬ УПАЛ, НО БОРТОВОЙ САМОПИСЕЦ УЦЕЛЕЛ. НАЙДИ ОБЛОМКИ И СОЖГИ."
	lines := wrapText(src, 30)
	for _, ln := range lines {
		if utf8.RuneCountInString(ln) > 30 {
			t.Fatalf("%q is %d characters", ln, utf8.RuneCountInString(ln))
		}
	}
	// and the lines actually use the width they were given
	if utf8.RuneCountInString(lines[0]) < 20 {
		t.Fatalf("first line is %q, only %d characters wide",
			lines[0], utf8.RuneCountInString(lines[0]))
	}
}
