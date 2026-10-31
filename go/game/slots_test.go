package game

import (
	"testing"

	"heliosdrift/models"
	"heliosdrift/render"
)

// Every station mesh must be dockable on a straight approach down the slot
// axis: the trigger volume has to lie where the visible mouth is, otherwise
// the pilot flies into the opening and nothing happens.
func TestEveryStationIsDockable(t *testing.T) {
	for _, mesh := range []string{"station", "dodeca", "hermit"} {
		g := New()
		g.NewGame("TEST")
		g.EnterSystem(3, "launch", 0, false)
		st := g.World.Station
		st.Mesh = models.Get(mesh)
		nrm := st.Mat.Col(2)
		docked := false
		for z := 400.0; z > 60; z -= 1 {
			g.State = "flight"
			g.Player.Pos = st.Pos.Add(nrm.Mul(z))
			g.Player.Mat = render.LookMatrix(nrm.Mul(-1))
			g.Player.Speed = 60
			if g.TryDock() {
				docked = true
				t.Logf("%s: docks at %.0f (station size %.0f)", mesh, z, StationSize)
				break
			}
		}
		if !docked {
			t.Errorf("%s: never docks on a straight run down the slot axis", mesh)
		}
	}
}
