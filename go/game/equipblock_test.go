package game

import "testing"

// Pressing buy twice on the same fitting must take the money once. The mining
// laser charged every press, and the rig looked broken because nothing said it
// needed the laser first.
func TestFittingsCannotBeBoughtTwice(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.World.Sys.Tech = 12
	g.Cmd.Cash = 100000

	for _, id := range []string{"mininglaser", "miningrig", "shieldboost1"} {
		if why := g.EquipBlocked(id); why != "" {
			t.Fatalf("%s blocked before buying: %s", id, why)
		}
		g.BuyEquipment(id)
		after := g.Cmd.Cash
		g.BuyEquipment(id) // the second press
		if g.Cmd.Cash != after {
			t.Fatalf("%s charged twice: %.0f -> %.0f", id, after, g.Cmd.Cash)
		}
		if why := g.EquipBlocked(id); why != "ALREADY FITTED" {
			t.Fatalf("%s after buying: %q", id, why)
		}
	}
}

// The rig needs the laser, and the shop has to be able to say so.
func TestRigNeedsLaserFirst(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.World.Sys.Tech = 12
	g.Cmd.Cash = 100000
	if why := g.EquipBlocked("miningrig"); why != "FIT A MINING LASER FIRST" {
		t.Fatalf("rig before laser: %q", why)
	}
	before := g.Cmd.Cash
	g.BuyEquipment("miningrig")
	if g.Cmd.Cash != before {
		t.Fatal("the rig took money without being fitted")
	}
	g.BuyEquipment("mininglaser")
	if why := g.EquipBlocked("miningrig"); why != "" {
		t.Fatalf("rig after laser: %q", why)
	}
}
