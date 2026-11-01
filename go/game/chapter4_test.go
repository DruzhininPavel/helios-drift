package game

import "testing"

// Chapter four: the courier broadcasts at the whole system and offers the list
// three ways - buy it, take it over his wreck, or read it and let him go.
func TestChapterFourBroadcast(t *testing.T) {
	open := func(t *testing.T) *Game {
		g := New()
		g.NewGame("CH4")
		g.Cmd.Campaign.Stage = 3
		g.Cmd.Campaign.Active = true
		g.EnterSystem(11, "launch", 0, false)
		g.Cmd.Campaign.Target = g.Cmd.System
		g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
		g.campaignArrive()
		if g.Screen != "talk" || g.Talk != "c4_broadcast" {
			t.Fatalf("arrival opened %q/%q", g.Screen, g.Talk)
		}
		g.TalkChoose(0) // through the show
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

	t.Run("buy", func(t *testing.T) {
		g := open(t)
		g.Cmd.Cash = 9000
		pick(g, "BUY THE LIST - 7000")
		if g.Cmd.Cash != 2000 {
			t.Fatalf("cash %.0f after buying", g.Cmd.Cash)
		}
		if !g.Cmd.Flags["HAS_LIST"] || !g.Cmd.Flags["COURIER_ALIVE"] {
			t.Fatalf("flags: %v", g.Cmd.Flags)
		}
		if g.Cmd.Legal == 0 {
			t.Fatal("buying from him on air left no record")
		}
		if g.Cmd.Campaign.Stage != 4 {
			t.Fatalf("stage %d, want four", g.Cmd.Campaign.Stage)
		}
	})

	t.Run("read and release", func(t *testing.T) {
		g := open(t)
		cash := g.Cmd.Cash
		pick(g, "READ IT AND LET HIM GO")
		if !g.Cmd.Flags["READ_LIST"] || !g.Cmd.Flags["COURIER_ALIVE"] {
			t.Fatalf("flags: %v", g.Cmd.Flags)
		}
		if g.Cmd.Cash != cash {
			t.Fatalf("cash changed to %.0f for reading", g.Cmd.Cash)
		}
		if g.Cmd.Campaign.Stage != 4 {
			t.Fatalf("stage %d after letting him go", g.Cmd.Campaign.Stage)
		}
	})

	t.Run("kill", func(t *testing.T) {
		g := open(t)
		pick(g, "TAKE THE LIST AND KILL")
		found := false
		for _, e := range g.World.Entities {
			if e.CampaignTag == "COURIER" {
				found = true
			}
		}
		if !found {
			t.Fatal("the courier is not in the sky after choosing to kill him")
		}
		if !g.Cmd.Flags["COURIER_DEAD"] {
			t.Fatal("the kill was not recorded")
		}
		cash := g.Cmd.Cash
		g.campaignKill("COURIER")
		if g.Cmd.Cash != cash+5200 {
			t.Fatalf("cash %.0f -> %.0f, want +5200", cash, g.Cmd.Cash)
		}
	})

	t.Run("buying needs the money", func(t *testing.T) {
		g := open(t)
		g.Cmd.Cash = 100
		n := DialogueNode(g.Talk)
		for _, idx := range g.VisibleOptions(n) {
			if n.Options[idx].Text == "BUY THE LIST - 7000" {
				t.Fatal("the list is on sale with a hundred credits")
			}
		}
	})
}
