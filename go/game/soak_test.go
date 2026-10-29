package game

import (
	"math/rand"
	"testing"

	"heliosdrift/galaxy"
)

// A long simulated flight with random input and a lot of traffic: catches nil
// dereferences and runaway state in the combat, AI and docking paths.
func TestFlightSoak(t *testing.T) {
	g := newGameMode(ModeAdvanced)
	g.Cmd.Equip["ecm"] = true
	g.Cmd.Equip["dockcomp"] = true
	g.Cmd.Missiles = 4
	g.Cmd.Cargo[0] = 3
	g.Launch()

	rng := rand.New(rand.NewSource(7))
	for i := 0; i < 60*90; i++ {
		if i%45 == 0 {
			g.In.Pitch = rng.Float64()*2 - 1
			g.In.Roll = rng.Float64()*2 - 1
			g.In.Accel = rng.Float64()*2 - 1
			g.In.Fire = rng.Intn(3) == 0
		}
		switch {
		case i%300 == 0:
			g.SpawnPirates(1 + rng.Intn(3))
		case i%420 == 0:
			g.SpawnTraffic(1 + rng.Intn(2))
		case i%700 == 0:
			g.SpawnAsteroid(3000)
		case i%900 == 0:
			g.FireMissile()
		case i%1200 == 0:
			g.FireECM()
		case i%1500 == 0:
			g.Jettison()
		}
		g.Update(1.0 / 60.0)
		if g.State == "dead" {
			for j := 0; j < 400 && g.State == "dead"; j++ {
				g.Update(1.0 / 60.0)
			}
			if g.State == "dead" {
				break
			}
		}
		if g.State == "docking" {
			for j := 0; j < 120 && g.State == "docking"; j++ {
				g.Update(1.0 / 60.0)
			}
			if g.State == "docked" {
				g.Launch()
			}
		}
		if g.Cmd.Fuel < 0.3 {
			g.Cmd.Fuel = 7
		}
	}
	if g.World == nil {
		t.Fatal("world disappeared")
	}
	if g.Player.Energy < 0 || g.Player.Energy > g.Ship().Energy+1 {
		t.Fatalf("energy out of range: %v", g.Player.Energy)
	}
}

// Same soak in classic mode, which has a different set of allowed actions.
func TestFlightSoakClassic(t *testing.T) {
	g := newGameMode(ModeClassic)
	g.Launch()
	rng := rand.New(rand.NewSource(11))
	for i := 0; i < 60*45; i++ {
		if i%60 == 0 {
			g.In.Pitch = rng.Float64()*2 - 1
			g.In.Roll = rng.Float64()*2 - 1
			g.In.Fire = rng.Intn(3) == 0
		}
		if i%400 == 0 {
			g.SpawnPirates(2)
		}
		g.Update(1.0 / 60.0)
		if g.State == "dead" {
			break
		}
		if g.Cmd.Fuel < 0.3 {
			g.Cmd.Fuel = 7
		}
	}
}

// Docking, trading, contracts, refits and jumping, over and over.
func TestStationSoak(t *testing.T) {
	g := newGameMode(ModeAdvanced)
	g.Cmd.Cash = 2000000
	rng := rand.New(rand.NewSource(3))
	g.dockedAt()

	for i := 0; i < 120; i++ {
		for j := 0; j < 6; j++ {
			g.Buy(rng.Intn(17), 1+rng.Intn(3))
			g.Sell(rng.Intn(17), 1+rng.Intn(2))
		}
		g.Refuel()
		if len(g.Offers) > 0 && rng.Intn(2) == 0 {
			g.AcceptOffer(rng.Intn(len(g.Offers)))
		}
		if rng.Intn(4) == 0 {
			g.BuyShip(ShipYard[rng.Intn(len(ShipYard))].ID)
		}
		if CargoUsed(g.Cmd) > g.Cmd.CargoCap {
			t.Fatalf("cargo over capacity: %d/%d", CargoUsed(g.Cmd), g.Cmd.CargoCap)
		}
		if g.Cmd.Cash < 0 {
			t.Fatalf("negative cash: %v", g.Cmd.Cash)
		}

		g.Launch()
		gal := galaxy.Get(g.Cmd.Galaxy)
		from := gal.Systems[g.Cmd.System]
		next := -1
		for _, s := range gal.Systems {
			if d := galaxy.Distance(from, s); d > 0.5 && d <= 7 {
				next = s.ID
				break
			}
		}
		if next >= 0 && g.Cmd.Fuel >= 1 {
			if g.Hyperspace(next) {
				for j := 0; j < 400 && g.State == "hyper"; j++ {
					g.Update(1.0 / 60.0)
				}
			}
		}
		g.Cmd.Fuel = 7
		g.State = "docked"
		g.dockedAt()
	}
}
