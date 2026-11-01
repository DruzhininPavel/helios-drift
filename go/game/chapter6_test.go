package game

import "testing"

// Chapter six opens with the navy asking once. Three answers, three outcomes:
// papers end it quietly, surrender ends it badly, fighting spawns the vipers.
func TestChapterSixAnswers(t *testing.T) {
	open := func(t *testing.T) *Game {
		g := New()
		g.NewGame("CH6")
		g.Cmd.Campaign.Stage = 5
		g.Cmd.Campaign.Active = true
		g.EnterSystem(7, "launch", 0, false)
		g.Cmd.Campaign.Target = g.Cmd.System
		g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
		g.campaignArrive()
		if g.Screen != "talk" || g.Talk != "c6_brief" {
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

	t.Run("papers", func(t *testing.T) {
		g := open(t)
		g.Cmd.Papers = 2
		g.Cmd.Cash = 1000
		pick(g, "SHOW THE PAPERS")
		if g.Talk != "c6_papers" {
			t.Fatalf("after the papers: %q", g.Talk)
		}
		g.TalkChoose(0)
		if g.Talk != "c6_archive" {
			t.Fatalf("at the archive: %q", g.Talk)
		}
		// the form costs three hundred, the cabinet pays six thousand, and the
		// file costs one paper: 1000 - 300 + 6000
		pick(g, "PAY 300")
		if g.Cmd.Cash != 6700 {
			t.Fatalf("cash %.0f after the archive", g.Cmd.Cash)
		}
		if !g.Cmd.Flags["GOT_FORM"] {
			t.Fatal("the form was not recorded")
		}
		if g.Cmd.Papers != 1 {
			t.Fatalf("%d papers after filing", g.Cmd.Papers)
		}
		if g.Cmd.Campaign.Stage != 6 {
			t.Fatalf("stage %d, want six", g.Cmd.Campaign.Stage)
		}
	})

	t.Run("papers are required", func(t *testing.T) {
		g := open(t)
		g.Cmd.Papers = 0
		n := DialogueNode(g.Talk)
		for _, idx := range g.VisibleOptions(n) {
			if n.Options[idx].Text == "SHOW THE PAPERS" {
				t.Fatal("the papers option is offered with no papers")
			}
		}
	})

	t.Run("surrender", func(t *testing.T) {
		g := open(t)
		g.Cmd.Legal = 0
		pick(g, "SURRENDER")
		if g.Cmd.Legal < 100 {
			t.Fatalf("legal %d after cooperating, want a record", g.Cmd.Legal)
		}
		if g.Cmd.Campaign.Stage != 6 {
			t.Fatalf("stage %d, want six", g.Cmd.Campaign.Stage)
		}
	})

	t.Run("fight", func(t *testing.T) {
		g := open(t)
		pick(g, "FIGHT THEM")
		found := 0
		for _, e := range g.World.Entities {
			if e.CampaignTag == "NAVY" {
				found++
			}
		}
		if found != 2 {
			t.Fatalf("%d vipers spawned, want two", found)
		}
	})

	t.Run("spared convoy changes the briefing", func(t *testing.T) {
		g := open(t)
		g.Cmd.Flags["CONVOY_SPARED"] = true
		lines := g.NodeLines(DialogueNode("c6_brief"))
		if len(lines) == 0 || lines[0] != "THE CONVOY CONFIRMED YOU DID NOT ATTACK THEM." {
			t.Fatalf("briefing starts with %q", lines[0])
		}
	})
}
