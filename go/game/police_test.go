package game

import "testing"

// A clean pilot who slows down is inspected and released: the police do not
// open fire just because they are there.
func TestPoliceInspectAcleanPilot(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.Cmd.Legal = 0
	g.Cmd.Cargo = map[int]int{}
	g.Player.Speed = g.Ship().Speed * 0.2 // cut speed, as the hail asks
	g.spawnPolice(1)
	if g.PoliceState != "hail" {
		t.Fatalf("state %q, want a hail", g.PoliceState)
	}
	// the viper has to be near for the inspection, so put it there
	for _, e := range g.World.Entities {
		if e.Behaviour == "police" {
			e.Pos = g.Player.Pos.Add(e.Mat.Col(2).Mul(-400))
			e.Hostile = false
		}
	}
	for i := 0; i < 120 && g.PoliceState == "hail"; i++ {
		g.Update(1.0 / 60.0)
	}
	if g.PoliceState != "" {
		t.Fatalf("state %q after the inspection, want cleared", g.PoliceState)
	}
}

// A wanted pilot who stops gets arrested: the same stop, a different ending.
func TestPoliceArrestAWantedPilot(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.Cmd.Legal = 400
	g.Player.Speed = g.Ship().Speed * 0.2
	g.spawnPolice(1)
	for _, e := range g.World.Entities {
		if e.Behaviour == "police" {
			e.Pos = g.Player.Pos.Add(e.Mat.Col(2).Mul(-400))
		}
	}
	for i := 0; i < 120 && g.PoliceState == "hail"; i++ {
		g.Update(1.0 / 60.0)
	}
	if g.PoliceState != "fight" {
		t.Fatalf("state %q, want a fight", g.PoliceState)
	}
}

// Ignoring the hail for long enough has the same result as being wanted.
func TestPoliceGiveUpWaiting(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.Cmd.Legal = 0
	g.Player.Speed = g.Ship().Speed // keep flying, ignore them
	g.spawnPolice(1)
	for _, e := range g.World.Entities {
		if e.Behaviour == "police" {
			e.Pos = g.Player.Pos.Add(e.Mat.Col(2).Mul(-4000))
		}
	}
	for i := 0; i < 60*20 && g.PoliceState == "hail"; i++ {
		g.Player.Speed = g.Ship().Speed
		g.Update(1.0 / 60.0)
	}
	if g.PoliceState != "fight" {
		t.Fatalf("state %q after ignoring the hail, want a fight", g.PoliceState)
	}
}

// The hail asks you to slow down; it does not ask you to survive a countdown
// while they fly in from two thousand kilometres away. Stop, and they come.
func TestStoppingIsEnoughEvenIfTheyAreFar(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.Cmd.Legal = 0
	g.Cmd.Cargo = map[int]int{}
	g.Player.Speed = 0
	g.spawnPolice(1)
	last := "hail"
	for i := 0; i < 60*45; i++ {
		g.Player.Speed = 0 // the pilot holds still, as asked
		g.Update(1.0 / 60.0)
		last = g.PoliceState
		if last != "hail" {
			break
		}
	}
	if last != "" {
		t.Fatalf("after standing still the state is %q, want the inspection cleared", last)
	}
}
