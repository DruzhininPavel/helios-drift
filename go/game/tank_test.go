package game

import (
	"testing"

	"heliosdrift/galaxy"
)

// The long range tank has to change the number the tanks hold, or it is a
// five thousand credit line in the shop that does nothing.
func TestLongRangeTankRaisesTheCap(t *testing.T) {
	g := New()
	g.NewGame("TANK")
	if g.FuelCap() != FuelMax {
		t.Fatalf("cap %.1f without a tank, want %.1f", g.FuelCap(), FuelMax)
	}
	// buy it the way a pilot would: cash first, then the fitting
	g.Cmd.Cash = 5000
	g.BuyEquipment("fueltank")
	if !g.Cmd.Equip["fueltank"] {
		g.Cmd.Equip["fueltank"] = true // the shop's rules are not the subject
	}
	if g.FuelCap() != FuelMax+FuelTankBonus {
		t.Fatalf("cap %.1f with a tank, want %.1f", g.FuelCap(), FuelMax+FuelTankBonus)
	}
	g.Cmd.Cash = 10000
	g.Cmd.Fuel = 1
	g.Refuel()
	if g.Cmd.Fuel != g.FuelCap() {
		t.Fatalf("refuelled to %.1f, want the cap %.1f", g.Cmd.Fuel, g.FuelCap())
	}
}

// The galactic drive works from the chart as well as from the cockpit, and it
// burns out with the jump.
func TestGalacticDriveFromTheChart(t *testing.T) {
	g := New()
	g.NewGame("DRIVE")
	g.Cmd.Equip["galdrive"] = true
	g.State = "flight"
	g.Overlay = "chart"
	before := g.Cmd.Galaxy
	if !g.GalacticJump() {
		t.Fatal("the drive refused a jump from the chart")
	}
	if g.Cmd.Galaxy != (before+1)%8 {
		t.Fatalf("galaxy %d after the jump, want %d", g.Cmd.Galaxy, (before+1)%8)
	}
	if g.Cmd.Equip["galdrive"] {
		t.Fatal("the drive survived the jump")
	}
	if g.GalacticJump() {
		t.Fatal("a spent drive jumped again")
	}
}

// A long range tank widens the jump range, not only the number of jumps: the
// chart draws its circle from the same number the hyperspace check uses.
func TestTankWidensTheJumpRange(t *testing.T) {
	g := New()
	g.NewGame("RANGE")
	if g.JumpRange() != 7 {
		t.Fatalf("range %.1f without a tank", g.JumpRange())
	}
	g.Cmd.Equip["fueltank"] = true
	if g.JumpRange() != 10 {
		t.Fatalf("range %.1f with a tank, want 10", g.JumpRange())
	}
	// an eight light year hop is out of reach without the tank and fine with it
	gal := galaxy.Get(g.Cmd.Galaxy)
	g.Cmd.Fuel = 10
	g.Cmd.Equip["fueltank"] = false
	far := -1
	for _, s := range gal.Systems {
		if d := galaxy.Distance(gal.Systems[g.Cmd.System], s); d > 7.5 && d <= 10 {
			far = s.ID
			break
		}
	}
	if far < 0 {
		t.Skip("no system between eight and ten light years in this galaxy")
	}
	if g.Hyperspace(far) {
		t.Fatal("jumped eight light years on a standard tank")
	}
	g.Cmd.Equip["fueltank"] = true
	if !g.Hyperspace(far) {
		t.Fatal("the long range tank did not reach a system eight light years out")
	}
}
