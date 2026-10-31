package game

import "testing"

// settleTo prepares a commander at stage two, in the target system, with the
// conversation already open.
func settleTo(t *testing.T, setup func(*Game)) *Game {
	t.Helper()
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	if setup != nil {
		setup(g)
	}
	g.Cmd.Campaign.Stage = 1
	g.Cmd.Campaign.Active = true
	g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
	g.Cmd.Campaign.Target = g.Cmd.System
	g.campaignArrive()
	return g
}

func chooseByText(g *Game, text string) bool {
	n := DialogueNode(g.Talk)
	if n == nil {
		return false
	}
	for vi, idx := range g.VisibleOptions(n) {
		if n.Options[idx].Text == text {
			g.TalkChoose(vi)
			return true
		}
	}
	return false
}

func countTagged(g *Game, tag string) int {
	n := 0
	for _, e := range g.World.Entities {
		if e.CampaignTag == tag {
			n++
		}
	}
	return n
}

// Paying for the permit ends the stage without a shot, and the receipt makes
// the price go up: the collectors get paid by the paper.
func TestChapterTwoPayToPass(t *testing.T) {
	g := settleTo(t, func(g *Game) { g.Cmd.Cash = 10000 })
	if g.Screen != "talk" || g.Talk != "c2_cmdA" {
		t.Fatalf("screen %q node %q, want the collectors", g.Screen, g.Talk)
	}
	if !chooseByText(g, "PAY 3000") {
		t.Fatal("no pay option with ten thousand in the bank")
	}
	if g.Cmd.Cash != 7000 {
		t.Fatalf("cash %.0f, want 7000", g.Cmd.Cash)
	}
	if !chooseByText(g, "LEAVE") {
		g.TalkChoose(0)
	}
	if g.Cmd.Campaign.Stage != 2 {
		t.Fatalf("stage %d, want 2", g.Cmd.Campaign.Stage)
	}
	if n := countTagged(g, "AMBUSH"); n != 0 {
		t.Fatalf("%d collectors spawned after paying", n)
	}
	if !g.Cmd.Flags["PAID_PERMIT"] {
		t.Fatal("no PAID_PERMIT flag")
	}
}

func TestChapterTwoReceiptRaisesPrice(t *testing.T) {
	g := settleTo(t, func(g *Game) {
		g.Cmd.Cash = 10000
		g.Cmd.Flags = map[string]bool{"HAS_RECEIPT": true}
	})
	if g.Talk != "c2_cmdAR" {
		t.Fatalf("node %q, want the custodian variant", g.Talk)
	}
	if !chooseByText(g, "PAY 4500") {
		t.Fatal("no 4500 option")
	}
	if g.Cmd.Cash != 5500 {
		t.Fatalf("cash %.0f, want 5500", g.Cmd.Cash)
	}
}

func TestChapterTwoRefusedBranch(t *testing.T) {
	g := settleTo(t, func(g *Game) {
		g.Cmd.Cash = 10000
		g.Cmd.Flags = map[string]bool{"REFUSED_SEAL": true}
	})
	if g.Talk != "c2_cmdB" {
		t.Fatalf("node %q, want the branch where the clinic was left to them", g.Talk)
	}
}

// Fighting spawns the three collectors and pays the bounty when they die.
func TestChapterTwoFight(t *testing.T) {
	g := settleTo(t, func(g *Game) { g.Cmd.Cash = 10000 })
	if !chooseByText(g, "FIGHT THEM") {
		t.Fatal("no fight option")
	}
	if n := countTagged(g, "AMBUSH"); n != 3 {
		t.Fatalf("%d contacts spawned, want 3", n)
	}
	before := g.Cmd.Cash
	for _, e := range append([]*Entity(nil), g.World.Entities...) {
		if e.CampaignTag == "AMBUSH" {
			g.DestroyEntity(e, "laser")
		}
	}
	if g.Cmd.Campaign.Stage != 2 {
		t.Fatalf("stage %d, want 2", g.Cmd.Campaign.Stage)
	}
	if got := g.Cmd.Cash - before; got < 2000 {
		t.Fatalf("bounty %.0f, want the stage to pay", got)
	}
}

// Running away costs nothing but the pay, and the guarantor pays instead.
func TestChapterTwoJumpOut(t *testing.T) {
	g := settleTo(t, func(g *Game) { g.Cmd.Cash = 1000 })
	if !chooseByText(g, "JUMP OUT") {
		t.Fatal("no jump option")
	}
	if g.Screen == "talk" {
		g.TalkChoose(0)
	}
	if g.Cmd.Campaign.Stage != 2 {
		t.Fatalf("stage %d, want 2: the run is failed, not stuck", g.Cmd.Campaign.Stage)
	}
	if !g.Cmd.Flags["DISP_OUT"] {
		t.Fatal("no DISP_OUT flag")
	}
	if g.Cmd.Cash > 1200 {
		t.Fatalf("cash %.0f: running away paid out", g.Cmd.Cash)
	}
}
