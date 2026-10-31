package game

import (
	"testing"

	"heliosdrift/galaxy"
)

// Classic is the 1984 rule set: it must not quietly grow modern fittings just
// because the catalogue did.
func TestClassicShopHasNoModernFittings(t *testing.T) {
	sys := &galaxy.System{Tech: 20}
	for _, it := range galaxy.AvailableEquipment(sys, true) {
		if it.AdvancedOnly {
			t.Errorf("classic shop sells %s", it.ID)
		}
	}
	modern := 0
	for _, it := range galaxy.AvailableEquipment(sys, false) {
		if it.AdvancedOnly {
			modern++
		}
	}
	if modern == 0 {
		t.Fatal("advanced shop lost the modern fittings, the filter is backwards")
	}
}

// Nor does it sell pardons: a 1984 commander cannot buy his record clean.
func TestClassicHasNoPardon(t *testing.T) {
	g := New()
	g.NewGameMode(ModeClassic, "TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.World.Sys.Government = 3
	g.World.Sys.Tech = 12
	g.Cmd.Legal = 200
	g.Cmd.Cash = 50000
	if g.PardonAvailable() {
		t.Fatal("classic offered a pardon")
	}
	if g.BuyPardon() {
		t.Fatal("classic sold a pardon")
	}
	if g.Cmd.Legal != 200 {
		t.Fatal("record changed in classic")
	}
}
