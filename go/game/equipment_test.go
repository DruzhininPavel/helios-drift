package game

import "testing"

// TestShieldBoostersRaiseTheCap checks the two shield upgrades: each adds to
// the hull's own shielding, and the second needs the first.
func TestShieldBoostersRaiseTheCap(t *testing.T) {
	g := New()
	g.NewGameMode(ModeAdvanced, "BOT")
	g.Cmd.Ship = "cobra3"
	g.Cmd.Cash = 1e6
	g.World.Sys.Tech = 15
	base := g.ShieldMax()

	g.Cmd.Cash = 1e6
	g.BuyEquipment("shieldboost2") // ordered wrong: refused
	if g.Cmd.Equip["shieldboost2"] {
		t.Fatal("heavy shielding fitted without a booster")
	}

	g.Cmd.Cash = 1e6
	g.BuyEquipment("shieldboost1")
	if got := g.ShieldMax(); got != base+40 {
		t.Fatalf("after booster shield max %.0f, want %.0f", got, base+40)
	}
	g.Cmd.Cash = 1e6
	g.BuyEquipment("shieldboost2")
	if got := g.ShieldMax(); got != base+120 {
		t.Fatalf("after heavy shielding shield max %.0f, want %.0f", got, base+120)
	}
	if g.Player.ShieldF != g.ShieldMax() {
		t.Fatalf("shields not topped up: %.0f of %.0f", g.Player.ShieldF, g.ShieldMax())
	}
}

// TestCargoExpansionAddsHold covers the second bay: it stacks on the first.
func TestCargoExpansionAddsHold(t *testing.T) {
	g := New()
	g.NewGameMode(ModeAdvanced, "BOT")
	g.Cmd.Ship = "cobra3"
	g.Cmd.Cash = 1e6
	g.World.Sys.Tech = 15
	base := g.CargoCapacity()

	g.BuyEquipment("cargobay2") // refused without the first bay
	if g.Cmd.Equip["cargobay2"] {
		t.Fatal("cargo expansion fitted without a large bay")
	}
	g.Cmd.Cash = 1e6
	g.BuyEquipment("largebay")
	first := g.CargoCapacity()
	if first != base+15 {
		t.Fatalf("large bay gave %d, want %d", first, base+15)
	}
	g.Cmd.Cash = 1e6
	g.BuyEquipment("cargobay2")
	if got := g.CargoCapacity(); got != first+25 {
		t.Fatalf("cargo expansion gave %d, want %d", got, first+25)
	}
	if g.Cmd.CargoCap != g.CargoCapacity() {
		t.Fatalf("commander hold %d, capacity %d", g.Cmd.CargoCap, g.CargoCapacity())
	}
}

// TestHeavyMiningRigCutsMore checks that the rig actually yields more ore than
// the plain mining laser.
func TestHeavyMiningRigCutsMore(t *testing.T) {
	count := func(rig bool) int {
		g := New()
		g.NewGameMode(ModeAdvanced, "BOT")
		g.Cmd.Equip["mining"] = true
		if rig {
			g.Cmd.Equip["miningrig"] = true
		}
		g.Launch()
		ast := &Entity{Kind: "asteroid", Radius: 40, Energy: 1}
		g.HitEntity(ast, 100, "mining")
		n := 0
		for _, e := range g.World.Entities {
			if e.Kind == "boulder" {
				n++
			}
		}
		return n
	}
	if plain, rig := count(false), count(true); rig <= plain {
		t.Fatalf("rig yields %d pieces, plain laser %d", rig, plain)
	}
}

// TestPlasmaNeedsARealFighter keeps the best gun out of the trading hulls.
func TestPlasmaNeedsARealFighter(t *testing.T) {
	g := New()
	g.NewGameMode(ModeAdvanced, "BOT")
	g.Cmd.Ship = "python"
	g.Cmd.Cash = 1e7
	g.World.Sys.Tech = 15
	g.BuyEquipment("laserplasma")
	if g.Cmd.Laser == 3 {
		t.Fatal("python mounted a plasma accelerator")
	}

	f := New()
	f.NewGameMode(ModeAdvanced, "BOT")
	f.Cmd.Ship = "ferdelance"
	f.Cmd.Cash = 1e7
	f.World.Sys.Tech = 15
	f.BuyEquipment("laserplasma")
	if f.Cmd.Laser != 3 {
		t.Fatalf("fer-de-lance laser %d, want plasma", f.Cmd.Laser)
	}
}
