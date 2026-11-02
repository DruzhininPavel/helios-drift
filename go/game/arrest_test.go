package game

import "testing"

// Cutting speed for a police stop with slaves aboard is the end of the voyage,
// not a message: the cargo goes, the fine comes out of the purse, the record is
// wiped, and the pilot is set down on the station.
func TestArrestTakesTheCargoAndTheFine(t *testing.T) {
	g := New()
	g.NewGame("ARREST")
	g.EnterSystem(3, "arrive", 9000, false)
	g.State = "flight"
	g.Cmd.Cash = 5000
	g.Cmd.Cargo[3] = 4 // slaves: contraband
	g.Cmd.Legal = 20
	fine := g.PardonPrice()

	g.Player.Speed = g.Ship().Speed * 0.2
	g.spawnPolice(1)
	for _, e := range g.World.Entities {
		if e.Behaviour == "police" {
			e.Pos = g.Player.Pos.Add(e.Mat.Col(2).Mul(-400))
		}
	}
	for i := 0; i < 60*10 && g.State != "docked"; i++ {
		g.Update(1.0 / 60)
	}
	if g.State != "docked" {
		t.Fatalf("state %q after complying with a police stop", g.State)
	}
	if g.Cmd.Cargo[3] != 0 {
		t.Errorf("slaves still aboard: %d", g.Cmd.Cargo[3])
	}
	if g.Cmd.Cash != float64(5000-fine) {
		t.Errorf("cash %.0f after the fine, want %d", g.Cmd.Cash, 5000-fine)
	}
	if g.Cmd.Legal != 0 {
		t.Errorf("record still dirty: %d", g.Cmd.Legal)
	}
	if g.PoliceState != "" {
		t.Errorf("police state %q after the arrest", g.PoliceState)
	}
}

// A clean pilot is searched and released: no cargo, no fine, still flying.
func TestCleanPilotIsReleased(t *testing.T) {
	g := New()
	g.NewGame("CLEAN")
	g.EnterSystem(3, "arrive", 9000, false)
	g.State = "flight"
	g.Cmd.Cash = 1000
	g.Player.Speed = g.Ship().Speed * 0.2
	g.spawnPolice(1)
	for _, e := range g.World.Entities {
		if e.Behaviour == "police" {
			e.Pos = g.Player.Pos.Add(e.Mat.Col(2).Mul(-400))
		}
	}
	for i := 0; i < 60*10 && g.State != "docked"; i++ {
		g.Update(1.0 / 60)
	}
	if g.State == "docked" {
		t.Fatal("a clean pilot was arrested")
	}
	if g.Cmd.Cash != 1000 {
		t.Errorf("cash %.0f after a clean inspection", g.Cmd.Cash)
	}
}
