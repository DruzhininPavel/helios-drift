package game

import (
	"testing"

	"heliosdrift/galaxy"
)

// The catalogue itself has to be unambiguous: freight is resold, fittings are
// kept, ammunition is spent. A player cannot sell a mining laser back, so the
// shop must never sell a second one.
func TestEveryItemHasAKind(t *testing.T) {
	for _, it := range galaxy.EquipmentList {
		switch it.Kind {
		case galaxy.EquipUnique, galaxy.EquipConsumable, galaxy.EquipLaser:
		default:
			t.Errorf("%s has kind %q, the shop cannot reason about it", it.ID, it.Kind)
		}
		if it.Kind == galaxy.EquipConsumable && it.Needs != "" {
			t.Errorf("%s is consumable but needs something", it.ID)
		}
	}
}

// No fitting can be bought twice, whatever it is. This is the rule the whole
// distinction exists for, and it is checked over the catalogue rather than a
// hand-picked list.
func TestNoFittingCanBeBoughtTwice(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.World.Sys.Tech = 20
	g.Cmd.Cash = 1e9
	for _, it := range galaxy.EquipmentList {
		switch it.Kind {
		case galaxy.EquipConsumable:
			continue
		case galaxy.EquipLaser:
			// buying a better gun twice is a no-op: the level is already there
			continue
		}
		if why := g.EquipBlocked(it.ID); why != "" {
			continue // blocked by a prerequisite, covered elsewhere
		}
		g.BuyEquipment(it.ID)
		cash := g.Cmd.Cash
		g.BuyEquipment(it.ID)
		if g.Cmd.Cash != cash {
			t.Errorf("%s charged twice: %.0f -> %.0f", it.ID, cash, g.Cmd.Cash)
		}
	}
}

// Consumables are the opposite: they are supposed to be buyable again and
// again, and missiles stop at the rack limit rather than the wallet.
func TestConsumablesStillBuyable(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.World.Sys.Tech = 20
	g.Cmd.Cash = 10000
	for i := 0; i < MissilesMax+3; i++ {
		g.BuyEquipment("missile")
	}
	if g.Cmd.Missiles != MissilesMax {
		t.Fatalf("missiles %d, want %d", g.Cmd.Missiles, MissilesMax)
	}
}
