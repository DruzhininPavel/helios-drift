package game

import (
	"testing"

	"heliosdrift/render"
)

// A fighter carries no freight, but its wreck should still be worth scooping.
func TestFightersLeaveSalvage(t *testing.T) {
	found := 0
	for trial := 0; trial < 40; trial++ {
		g := New()
		SetSeed(int64(trial*7919 + 13))
		g.NewGame("TEST")
		g.EnterSystem(3, "launch", 0, false)
		e := g.baseShip("viper", render.V3(0, 0, -500), "pirate")
		e.CargoType = 0
		g.World.Entities = append(g.World.Entities, e)
		g.DestroyEntity(e, "laser")
		for _, x := range g.World.Entities {
			if x.Kind == "canister" {
				found++
				break
			}
		}
	}
	if found == 0 {
		t.Fatal("no wreck dropped anything in 40 tries")
	}
	if found == 40 {
		t.Fatal("every wreck dropped cargo, the chance is not a chance")
	}
	t.Logf("%d of 40 wrecks left a canister", found)
}
