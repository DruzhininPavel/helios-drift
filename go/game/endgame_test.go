package game

import (
	"testing"

	"heliosdrift/galaxy"
)

// A promotion is announced, written down, and at DANGEROUS it hands over a gun.
func TestRankUpIsAnnouncedAndArmed(t *testing.T) {
	g := New()
	g.NewGame("RANK")
	g.Cmd.Kills = 7 // one short of MOSTLY HARMLESS
	before := galaxy.RankIndexFor(7)
	g.Cmd.Kills = 8
	g.RankUp(before)
	if len(g.Messages) == 0 {
		t.Fatal("a promotion said nothing")
	}
	if !logHas(g, "RANK MOSTLY HARMLESS") {
		t.Fatal("the promotion is not in the log")
	}
	// climbing to DANGEROUS fits the military laser
	g.Cmd.Kills = 600
	g.Cmd.Laser = 0
	from := galaxy.RankIndexFor(200) // DEADLY, one short of DANGEROUS
	g.Cmd.Kills = 512
	g.RankUp(from)
	if g.Cmd.Laser != 2 {
		t.Fatalf("laser %d after DANGEROUS, want the military one", g.Cmd.Laser)
	}
}

// The endings are a standing, not only an epitaph.
func TestEndingTitlesChangeTheMarket(t *testing.T) {
	g := New()
	g.NewGame("TITLE")
	g.State = "docked"
	g.EnterSystem(3, "arrival", 0, false)
	g.Market[0].Price = 100
	g.Market[0].Sell = 100
	g.Market[0].Qty = 10
	g.Cmd.Cargo[0] = 10

	plain := g.Cmd.Cash
	g.Sell(0, 1)
	plainEarn := g.Cmd.Cash - plain

	g.Cmd.Cash = 0
	g.Cmd.Cargo[0] = 10
	g.Cmd.Ending = "combine"
	g.Sell(0, 1)
	ownerEarn := g.Cmd.Cash
	if ownerEarn <= plainEarn {
		t.Fatalf("an owner earned %.0f, a free pilot %.0f", ownerEarn, plainEarn)
	}
	if g.TitleName() != "OWNER OF THE SECTOR" {
		t.Fatalf("title %q", g.TitleName())
	}

	g.Cmd.Ending = "council"
	g.Cmd.Cash = 10000
	g.Market[0].Qty = 10
	g.Buy(0, 1)
	if spent := 10000 - g.Cmd.Cash; spent >= 100 {
		t.Fatalf("a council member paid %.0f for a 100 crate", spent)
	}

	g.Cmd.Ending = "chaos"
	if g.dangerNow(1.0) >= dangerOf(t, g) {
		t.Fatal("chaos did not thin the pirates out")
	}
}

func dangerOf(t *testing.T, g *Game) float64 {
	t.Helper()
	saved := g.Cmd.Ending
	g.Cmd.Ending = ""
	v := g.dangerNow(1.0)
	g.Cmd.Ending = saved
	return v
}

func logHas(g *Game, what string) bool {
	for _, e := range g.Cmd.Log {
		if e.A == what || e.Kind == what {
			return true
		}
	}
	return false
}
