package game

import "testing"

// The cheat shield is for skipping firefights, not for flying through walls.
func TestCheatShieldStopsFireNotCollisions(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.Player.Invulnerable = 0
	g.ToggleCheatShield()

	full := g.Player.ShieldF + g.Player.Energy
	g.DamagePlayer(30, true, "laser")
	if got := g.Player.ShieldF + g.Player.Energy; got != full {
		t.Fatalf("laser got through the cheat shield: %.0f -> %.0f", full, got)
	}
	g.DamagePlayer(10, true, "debris")
	if got := g.Player.ShieldF + g.Player.Energy; got != full {
		t.Fatalf("debris got through: %.0f -> %.0f", full, got)
	}

	// a hull you flew into is still a hull you flew into
	g.DamagePlayer(45, true, "ram")
	if got := g.Player.ShieldF + g.Player.Energy; got >= full {
		t.Fatal("ramming should still hurt with the cheat shield on")
	}
}

func TestCheatMoney(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	before := g.Cmd.Cash
	g.CheatMoney(10000)
	if g.Cmd.Cash != before+10000 {
		t.Fatalf("cash %.0f, want %.0f", g.Cmd.Cash, before+10000)
	}
}
