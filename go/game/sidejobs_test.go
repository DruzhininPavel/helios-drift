package game

import (
	"testing"

	"heliosdrift/galaxy"
)

// A side job opens when its condition holds, walks its steps, pays, and never
// opens twice.
func TestClinicJobWalksAllSteps(t *testing.T) {
	g := New()
	g.NewGame("SIDE")
	g.Cmd.Flags["ENDING_SHOWN"] = true
	g.Cmd.Ending = "combine"
	g.noteEndingFlag("combine")
	g.State = "docked"

	// first dock at a high tech station opens it
	if !dockAtTech(g, 10) {
		t.Skip("no high tech system generated to test against")
	}
	if _, _, active := g.SideProgress("clinic"); !active {
		t.Fatal("the clinic job did not open")
	}
	if g.Screen != "talk" || g.Talk != "sj_clinic_a" {
		t.Fatalf("the first step opened %q/%q", g.Screen, g.Talk)
	}
	g.EndTalk()

	// the fight step spawns the two who collect the cut
	tag := "SIDE_CLINIC"
	found := 0
	for _, e := range g.World.Entities {
		if e.CampaignTag == tag {
			found++
		}
	}
	if found != 2 {
		t.Fatalf("%d marked ships for the fight, want two", found)
	}
	for i := 0; i < 2; i++ {
		g.SideOnKill(tag)
	}
	if _, _, active := g.SideProgress("clinic"); !active {
		t.Fatal("the job ended before the return trip")
	}

	// back to a high tech station: paid and closed
	cash := g.Cmd.Cash
	if !dockAtTech(g, 10) {
		t.Fatal("could not dock again")
	}
	if _, _, active := g.SideProgress("clinic"); active {
		t.Fatal("the clinic job did not finish")
	}
	if g.Cmd.Cash != cash+3000 {
		t.Fatalf("cash %.0f -> %.0f, want +3000", cash, g.Cmd.Cash)
	}
	if !g.Cmd.Flags["CLINIC_PAID"] {
		t.Fatal("the reward flag is missing")
	}
}

// The council job only exists after that ending.
func TestCouncilJobNeedsTheCouncilEnding(t *testing.T) {
	g := New()
	g.NewGame("SIDE2")
	g.Cmd.Flags["ENDING_SHOWN"] = true
	g.Cmd.Ending = "combine"
	g.noteEndingFlag("combine")
	g.State = "docked"
	if dockAtTech(g, 10) {
		if _, _, active := g.SideProgress("keys"); active {
			t.Fatal("the council job opened without the council ending")
		}
	}
}

// The courier job needs both the man and his list.
func TestListJobNeedsTheCourier(t *testing.T) {
	g := New()
	g.NewGame("SIDE3")
	g.Cmd.Flags["COURIER_ALIVE"] = true
	g.Cmd.Flags["HAS_LIST"] = true
	g.State = "docked"
	if dockAtTech(g, 12) {
		g.SideStartIfReady()
		if _, _, active := g.SideProgress("list"); !active {
			t.Fatal("the courier job did not open")
		}
	}
}

// dockAtTech puts the pilot in a system with at least this tech level and docks.
func dockAtTech(g *Game, tech int) bool {
	gal := galaxyAt(g)
	for _, sys := range gal.Systems {
		if sys.Tech >= tech {
			g.EnterSystem(sys.ID, "arrival", 0, false)
			if g.World.Sys.Tech < tech {
				continue
			}
			g.Docked()
			return true
		}
	}
	return false
}

func galaxyAt(g *Game) *galaxy.Galaxy { return galaxy.Get(g.Cmd.Galaxy) }

// The council job walks: three stations, a fight, then the first session.
func TestCouncilJobWalksAllSteps(t *testing.T) {
	g := New()
	g.NewGame("KEYS")
	g.Cmd.Flags["ENDING_SHOWN"] = true
	g.Cmd.Ending = "council"
	g.noteEndingFlag("council")
	g.State = "docked"
	if !dockAtTech(g, 8) {
		t.Skip("no suitable system")
	}
	if _, _, active := g.SideProgress("keys"); !active {
		t.Fatal("the council job did not open")
	}
	g.EndTalk()
	// three docks, any station will do
	for i := 0; i < 3; i++ {
		dockAny(g)
		g.EndTalk()
	}
	step, _, active := g.SideProgress("keys")
	if !active {
		t.Fatal("the council job closed before the fight")
	}
	if step != 1 {
		t.Fatalf("step %d after three stations, want the fight", step)
	}
	for i := 0; i < 2; i++ {
		g.SideOnKill("SIDE_COUNCIL")
	}
	if !dockAtTech(g, 8) {
		t.Fatal("could not dock for the session")
	}
	if !g.Cmd.Flags["COUNCIL_SESSION"] {
		t.Fatal("the first session never happened")
	}
}

// dockAny puts the pilot in the next system and docks.
func dockAny(g *Game) {
	gal := galaxyAt(g)
	next := (g.Cmd.System + 1) % len(gal.Systems)
	g.EnterSystem(next, "arrival", 0, false)
	g.Docked()
}
