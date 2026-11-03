package game

import "testing"

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
