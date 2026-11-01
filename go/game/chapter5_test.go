package game

import "testing"

// Chapter five: the convoy captain talks before anyone shoots. Breaking it pays,
// letting it pass costs nothing but earns the record, and a licence changes the
// price of everything afterwards.
func TestChapterFiveConviction(t *testing.T) {
	open := func(t *testing.T) *Game {
		g := New()
		g.NewGame("CH5")
		g.Cmd.Campaign.Stage = 4
		g.Cmd.Campaign.Active = true
		g.EnterSystem(9, "launch", 0, false)
		g.Cmd.Campaign.Target = g.Cmd.System
		g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
		g.campaignArrive()
		if g.Screen != "talk" || g.Talk != "c5_captain" {
			t.Fatalf("arrival opened %q/%q", g.Screen, g.Talk)
		}
		return g
	}
	pick := func(g *Game, text string) {
		n := DialogueNode(g.Talk)
		for vi, idx := range g.VisibleOptions(n) {
			if n.Options[idx].Text == text {
				g.TalkChoose(vi)
				return
			}
		}
		t.Fatalf("%q is not on %q", text, g.Talk)
	}

	t.Run("break", func(t *testing.T) {
		g := open(t)
		pick(g, "BREAK THE CONVOY")
		boas := 0
		for _, e := range g.World.Entities {
			if e.CampaignTag == "CONVOY" {
				boas++
			}
		}
		if boas != 3 {
			t.Fatalf("%d boas, want three", boas)
		}
		if !g.Cmd.Flags["CONVOY_ATTACKED"] {
			t.Fatal("the attack was not recorded")
		}
		// killing them pays the stage fee
		cash := g.Cmd.Cash
		g.campaignKill("CONVOY")
		if g.Cmd.Cash <= cash {
			t.Fatalf("cash %.0f -> %.0f after the convoy", cash, g.Cmd.Cash)
		}
		if g.Cmd.Campaign.Stage != 5 {
			t.Fatalf("stage %d, want five", g.Cmd.Campaign.Stage)
		}
	})

	t.Run("spare", func(t *testing.T) {
		g := open(t)
		cash := g.Cmd.Cash
		pick(g, "LET THEM PASS")
		if !g.Cmd.Flags["CONVOY_SPARED"] {
			t.Fatal("sparing was not recorded")
		}
		if g.Cmd.Cash != cash {
			t.Fatalf("cash changed to %.0f for sparing them", g.Cmd.Cash)
		}
		if g.Cmd.Campaign.Stage != 5 {
			t.Fatalf("stage %d after sparing", g.Cmd.Campaign.Stage)
		}
	})

	t.Run("licence", func(t *testing.T) {
		g := open(t)
		g.Cmd.Cash = 10000
		pick(g, "BUY A LICENCE - 4000")
		if g.Cmd.Cash != 6000 {
			t.Fatalf("cash %.0f after the licence", g.Cmd.Cash)
		}
		if !g.Cmd.Flags["BOUGHT_LICENSE"] {
			t.Fatal("the licence was not recorded")
		}
		// and the market remembers: five per cent off
		g.State = "docked"
		g.Market[0].Price = 100
		g.Market[0].Qty = 10
		before := g.Cmd.Cash
		g.Buy(0, 1)
		if spent := before - g.Cmd.Cash; spent != 95 {
			t.Fatalf("a licence holder paid %.0f for a 100 crate", spent)
		}
	})
}
