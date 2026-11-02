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
		// the form costs three hundred
		pick(g, "PAY 300")
		if g.Cmd.Cash != 700 {
			t.Fatalf("cash %.0f after the form", g.Cmd.Cash)
		}
		if !g.Cmd.Flags["CLERK_1"] {
			t.Fatal("the form was not recorded")
		}
		if g.Talk != "c6_arch2" {
			t.Fatalf("second archive is %q", g.Talk)
		}
		// the second clerk wants a signature, a bribe or a forgery
		pick(g, "FORGE THE SIGNATURE")
		if !g.Cmd.Flags["FORGED"] {
			t.Fatal("the forgery was not recorded")
		}
		// the forgery has its own scene before the third archive
		if g.Talk != "c6_forged" {
			t.Fatalf("after the forgery: %q", g.Talk)
		}
		g.TalkChoose(0)
		if g.Talk != "c6_arch3" {
			t.Fatalf("third archive is %q", g.Talk)
		}
		pick(g, "TAKE THE COPY")
		if !g.Cmd.Flags["HAS_COPY"] {
			t.Fatal("the copy was not recorded")
		}
		// the cabinet pays, the file costs a paper
		if g.Cmd.Cash != 6700 {
			t.Fatalf("cash %.0f at the end", g.Cmd.Cash)
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

// The second archive opens three other ways: the receipt from chapter one, the
// courier as a witness, or two thousand credits.
func TestChapterSixArchivesOpenOtherWays(t *testing.T) {
	open := func(t *testing.T) *Game {
		g := New()
		g.NewGame("CH6B")
		g.Cmd.Campaign.Stage = 5
		g.Cmd.Campaign.Active = true
		g.EnterSystem(7, "launch", 0, false)
		g.Cmd.Campaign.Target = g.Cmd.System
		g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
		g.campaignArrive()
		g.Cmd.Papers = 2
		chooseText(t, g, "SHOW THE PAPERS")
		chooseText(t, g, "CONTINUE") // through the navy
		if g.Talk != "c6_archive" {
			t.Fatalf("not at the first archive: %q", g.Talk)
		}
		chooseText(t, g, "NO")
		if g.Talk != "c6_arch2" {
			t.Fatalf("not at the second archive: %q", g.Talk)
		}
		return g
	}
	has := func(g *Game, text string) bool {
		n := DialogueNode(g.Talk)
		for _, idx := range g.VisibleOptions(n) {
			if n.Options[idx].Text == text {
				return true
			}
		}
		return false
	}

	t.Run("receipt", func(t *testing.T) {
		g := open(t)
		g.Cmd.Flags["HAS_RECEIPT"] = true
		if !has(g, "SHOW THE RECEIPT") {
			t.Fatal("the receipt is not offered with the receipt in hand")
		}
	})
	t.Run("courier witness", func(t *testing.T) {
		g := open(t)
		g.Cmd.Flags["COURIER_ALIVE"] = true
		if !has(g, "BRING THE COURIER") {
			t.Fatal("the courier is not offered while he is alive")
		}
	})
	t.Run("bribe", func(t *testing.T) {
		g := open(t)
		g.Cmd.Cash = 5000
		if !has(g, "BRIBE HIM - 2000") {
			t.Fatal("the bribe is not offered with five thousand")
		}
	})
	t.Run("rig", func(t *testing.T) {
		g := New()
		g.NewGame("RIG")
		g.Cmd.Equip["miningrig"] = true
		g.Cmd.Campaign.Stage = 5
		g.Cmd.Campaign.Active = true
		g.EnterSystem(7, "launch", 0, false)
		g.Cmd.Campaign.Target = g.Cmd.System
		g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
		g.campaignArrive()
		g.Cmd.Papers = 1
		chooseText(t, g, "SHOW THE PAPERS")
		chooseText(t, g, "CONTINUE")
		if !has(g, "BRAKE THE CABINET") {
			t.Fatal("the rig does not open the cabinet")
		}
	})
	t.Run("auditor", func(t *testing.T) {
		g := open(t)
		chooseText(t, g, "FORGE THE SIGNATURE")
		chooseText(t, g, "CONTINUE")
		if g.Talk != "c6_arch3" {
			t.Fatalf("third archive is %q", g.Talk)
		}
		chooseText(t, g, "DEMAND AN AUDITOR")
		chooseText(t, g, "CONTINUE")
		if !g.Cmd.Flags["MET_AUDITOR"] {
			t.Fatal("the auditor was not recorded")
		}
		if g.Cmd.Campaign.Stage != 6 {
			t.Fatalf("stage %d after the auditor", g.Cmd.Campaign.Stage)
		}
	})
}

// chooseText presses the option with this label, whichever number it has.
func chooseText(t *testing.T, g *Game, text string) {
	t.Helper()
	n := DialogueNode(g.Talk)
	for vi, idx := range g.VisibleOptions(n) {
		if n.Options[idx].Text == text {
			g.TalkChoose(vi)
			return
		}
	}
	t.Fatalf("%q is not on %q", text, g.Talk)
}
