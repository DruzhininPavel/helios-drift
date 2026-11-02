package game

import (
	"strings"
	"testing"
)

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

// Sparing the convoy is not breaking it: the next briefing must not accuse the
// pilot of a massacre he did not commit, and the haulers must not stay behind.
func TestSparingTheConvoyClearsTheBlame(t *testing.T) {
	g := New()
	g.NewGame("SPARE")
	g.State = "flight"
	g.EnterSystem(3, "launch", 0, false)
	g.Cmd.Campaign.Stage = 4 // the convoy
	g.spawnCampaignGroup(QuestAt(4))
	if len(g.World.Entities) == 0 {
		t.Fatal("no convoy spawned")
	}
	if !g.TalkStart("c5_captain") {
		t.Fatal("the convoy captain does not talk")
	}
	// LET THEM PASS is the second visible answer
	vis := g.VisibleOptions(DialogueNode("c5_captain"))
	pick := -1
	for i, idx := range vis {
		if DialogueNode("c5_captain").Options[idx].Next == "c5_spared" {
			pick = i
		}
	}
	if pick < 0 {
		t.Fatal("no option spares the convoy")
	}
	if g.TalkChoose(pick) {
		g.EndTalk()
	}
	if !g.Cmd.Flags["CONVOY_SPARED"] {
		t.Fatal("the flag was not set")
	}
	for _, e := range g.World.Entities {
		if e.CampaignTag == "CONVOY" {
			t.Error("a spared hauler is still sitting in the system")
		}
	}
	next := QuestAt(g.Cmd.Campaign.Stage)
	if next.Talk != "c6_brief" {
		t.Fatalf("stage %d is %q", g.Cmd.Campaign.Stage, next.Talk)
	}
	if got := g.TitleFor(next); got == next.Title {
		t.Errorf("the job is still called %q after sparing the convoy", got)
	}
	blamed := false
	for _, ln := range g.BriefFor(next) {
		if strings.Contains(ln, "MASSACRE") || strings.Contains(ln, "FILTHY") {
			blamed = true
		}
	}
	if blamed {
		t.Errorf("the briefing still blames the pilot: %v", g.BriefFor(next))
	}
}
