package game

import (
	"math"
	"testing"

	"heliosdrift/render"
)

// Flying at the back of a station must not carry you through it: the slot is a
// shaft in one face, not a tunnel.
func TestStationBackIsSolid(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	st := g.World.Station
	nrm := st.Mat.Col(2)

	// start behind the station on the axis, the way he came at it
	g.State = "flight"
	g.Player.Invulnerable = 0
	g.Player.Pos = st.Pos.Sub(nrm.Mul(StationSize + 60))
	g.Player.Mat = render.LookMatrix(nrm) // nose towards the hull
	g.Player.Speed = 300

	died := false
	for i := 0; i < 240 && !died; i++ {
		// fly straight at the station
		g.Player.Pos = g.Player.Pos.Add(nrm.Mul(6))
		g.Update(1.0 / 60.0)
		if g.State == "dead" {
			died = true
		}
		if g.World == nil {
			break
		}
		// if he somehow crosses the centre, the test has failed the point
		local := st.Mat.ApplyT(g.Player.Pos.Sub(st.Pos))
		if local[2] > 0.1*StationSize {
			t.Fatalf("ship passed through the station: local z %.0f after %d frames", local[2], i)
		}
	}
	if !died {
		t.Fatal("the back of the station did not stop the ship")
	}
}

// And the mouth still works: a straight run down the axis from the front docks.
func TestStationFrontStillDocks(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	st := g.World.Station
	nrm := st.Mat.Col(2)
	for z := 400.0; z > 60; z -= 1 {
		g.State = "flight"
		g.Player.Pos = st.Pos.Add(nrm.Mul(z))
		g.Player.Mat = render.LookMatrix(nrm.Mul(-1))
		g.Player.Speed = 60
		if g.TryDock() {
			if math.Abs(z-StationSize) > 30 {
				t.Fatalf("docked at %.0f, far from the mouth", z)
			}
			return
		}
	}
	t.Fatal("never docked from the front")
}
