package game

import (
	"os"
	"testing"

	"heliosdrift/galaxy"
	"heliosdrift/render"
)

// The yard swaps hulls, keeps the cargo cap in step and charges the difference.
func TestShipyardTrade(t *testing.T) {
	g := newDocked()
	g.Cmd.Cash = 500000
	g.Cmd.Cargo = map[int]int{}
	if !g.BuyShip("python") {
		t.Fatal("buy failed")
	}
	if g.Cmd.Ship != "python" {
		t.Fatalf("ship = %q, want python", g.Cmd.Ship)
	}
	if g.Cmd.CargoCap != 40 {
		t.Fatalf("cargo cap = %d, want 40", g.Cmd.CargoCap)
	}
	want := 500000 - (260000 - 72000)
	if int(g.Cmd.Cash) != want {
		t.Fatalf("cash = %v, want %d", g.Cmd.Cash, want)
	}
	if g.Player.Energy != ShipFor("python").Energy {
		t.Fatalf("energy = %v, want %v", g.Player.Energy, ShipFor("python").Energy)
	}
}

func TestShipyardRefusals(t *testing.T) {
	g := newDocked()
	g.Cmd.Cash = 100
	if g.BuyShip("anaconda") {
		t.Fatal("bought a ship with no money")
	}
	if g.BuyShip("cobra3") {
		t.Fatal("allowed buying the hull already fitted")
	}
	// cargo that would not fit must block the trade
	g.Cmd.Cash = 500000
	g.Cmd.Ship = "cobra3"
	g.Cmd.Cargo[0] = 20
	if g.BuyShip("adder") {
		t.Fatal("sold the hull with cargo that cannot fit")
	}
}

// Jettison drops a canister that can be scooped back up.
func TestJettison(t *testing.T) {
	g := New()
	g.NewGame("DUMP")
	g.Launch()
	g.World.Entities = nil
	g.Cmd.Cargo[0] = 2
	if !g.Jettison() {
		t.Fatal("jettison refused")
	}
	if g.Cmd.Cargo[0] != 1 {
		t.Fatalf("cargo = %d, want 1", g.Cmd.Cargo[0])
	}
	found := false
	for _, e := range g.World.Entities {
		if e.Kind == "canister" && e.CargoType == 0 {
			found = true
		}
	}
	if !found {
		t.Fatal("no canister spawned")
	}
}

// Flying close to the star heats the cabin.
func TestCabinHeat(t *testing.T) {
	g := New()
	g.NewGame("ROAST")
	g.Launch()
	w := g.World
	if w.Sun.Radius <= 0 {
		t.Skip("no star in this system")
	}
	g.Player.Pos = w.Sun.Pos.Add(render.V3(w.Sun.Radius*1.3, 0, 0))
	for i := 0; i < 240; i++ {
		if g.State != "flight" {
			break
		}
		g.Update(0.05)
	}
	if g.Player.Cabin < 0.3 {
		t.Fatalf("cabin stayed cold at %.2f", g.Player.Cabin)
	}
}

// A misjump drops the commander in witchspace: no station, no planet, company.
func TestWitchspace(t *testing.T) {
	g := New()
	g.NewGame("LOST")
	g.EnterSystem(g.Cmd.System, "arrive", 9000, true)
	w := g.World
	if !w.Witchspace {
		t.Fatal("not flagged as witchspace")
	}
	if w.Station != nil {
		t.Fatal("a station appeared in witchspace")
	}
	if w.Planet.Radius != 0 {
		t.Fatal("a planet appeared in witchspace")
	}
	thargoids := 0
	for _, e := range w.Entities {
		if e.Type == "thargoid" || e.Type == "thargon" {
			thargoids++
		}
	}
	if thargoids == 0 {
		t.Fatal("witchspace was empty")
	}
	// and you can leave: pick any system inside jump range
	g.Cmd.Fuel = 7
	gal := galaxy.Get(g.Cmd.Galaxy)
	from := gal.Systems[g.Cmd.System]
	target := -1
	for _, s := range gal.Systems {
		if d := galaxy.Distance(from, s); d > 0.5 && d <= 7 {
			target = s.ID
			break
		}
	}
	if target < 0 {
		t.Fatal("no reachable system")
	}
	if !g.Hyperspace(target) {
		t.Fatal("could not hyperspace out of witchspace")
	}
}

// Three slots, each holding its own commander.
func TestSaveSlots(t *testing.T) {
	dir := t.TempDir()
	SaveDirOverride = dir
	defer func() { SaveDirOverride = "" }()

	g := New()
	g.NewGame("ALPHA")
	g.Cmd.Cash = 1111
	if !g.SaveSlot(1) {
		t.Fatal("save to slot 2 failed")
	}
	g.Cmd.Name = "BETA"
	g.Cmd.Cash = 2222
	if !g.SaveSlot(2) {
		t.Fatal("save to slot 3 failed")
	}
	if name, cash, _, ok := SlotInfo(1); !ok || name != "ALPHA" || cash != 1111 {
		t.Fatalf("slot 2 = %q %v %v", name, cash, ok)
	}
	if name, cash, _, ok := SlotInfo(2); !ok || name != "BETA" || cash != 2222 {
		t.Fatalf("slot 3 = %q %v %v", name, cash, ok)
	}
	if SlotExists(0) {
		t.Fatal("empty slot reports a commander")
	}
	if !g.LoadSlot(1) || g.Cmd.Name != "ALPHA" {
		t.Fatalf("load gave %q", g.Cmd.Name)
	}
	// the loaded commander is docked and flyable
	if g.State != "docked" || g.World == nil {
		t.Fatal("load did not restore a docked world")
	}
	_ = os.Remove(dir)
}
