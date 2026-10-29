package game

import (
	"testing"

	"heliosdrift/galaxy"
	"heliosdrift/render"
)

// flyIn aims the ship at the slot and runs the real update loop.
func flyIn(t *testing.T, sysPick func(*Game) int) {
	t.Helper()
	g := New()
	g.NewGame("DOCKER")
	g.Cmd.System = sysPick(g)
	g.EnterSystem(g.Cmd.System, "arrive", 9000, false)
	st := g.World.Station
	if st == nil {
		t.Fatal("no station")
	}
	normal := st.Mat.Col(2)
	g.Player.Pos = st.Pos.Add(normal.Mul(420))
	g.Player.Mat = render.LookMatrix(normal.Mul(-1))
	g.Player.Speed = 120
	for i := 0; i < 60*20; i++ {
		if g.State == "docked" {
			return
		}
		if g.State != "flight" && g.State != "docking" {
			t.Fatalf("died while docking: %s", g.LastDamageReason)
		}
		g.Update(1.0 / 60.0)
	}
	t.Fatalf("never docked: pos %v, dist %.0f", g.Player.Pos, st.Pos.Sub(g.Player.Pos).Len())
}

// dockNow runs the glide-in so tests see the docked state.
func dockNow(g *Game) {
	g.Dock()
	for i := 0; i < 200 && g.State != "docked"; i++ {
		g.Update(1.0 / 60.0)
	}
}

func pickStation(g *Game, kind string) int {
	for i := 0; i < 256; i++ {
		if stationMesh(galaxy.Get(g.Cmd.Galaxy).Systems[i]) == kind {
			return i
		}
	}
	return 0
}

func TestDockCoriolis(t *testing.T) { flyIn(t, func(g *Game) int { return pickStation(g, "station") }) }
func TestDockDodeca(t *testing.T)   { flyIn(t, func(g *Game) int { return pickStation(g, "dodeca") }) }
func TestDockHermit(t *testing.T)   { flyIn(t, func(g *Game) int { return pickStation(g, "hermit") }) }

// galaxySeed reads the system seed used to pick the station hull.
func galaxySeed(g *Game, i int) uint32 {
	return galaxy.Get(g.Cmd.Galaxy).Systems[i].Seed
}

// The docking computer is the answer to a fiddly slot: it must fly in and dock.
func TestAutopilotDocks(t *testing.T) {
	g := New()
	g.NewGame("AUTO")
	g.Cmd.Equip["dockcomp"] = true
	g.EnterSystem(g.Cmd.System, "arrive", 9000, false)
	st := g.World.Station
	g.Player.Pos = st.Pos.Add(render.V3(600, 200, 3800))
	g.Player.Speed = 200
	g.ToggleAutopilot()
	for i := 0; i < 60*90; i++ {
		if g.State == "docked" {
			return
		}
		if g.State != "flight" && g.State != "docking" {
			t.Fatalf("died on autopilot: %s", g.LastDamageReason)
		}
		g.Update(1.0 / 60.0)
	}
	t.Fatalf("autopilot did not dock: dist %.0f", st.Pos.Sub(g.Player.Pos).Len())
}

// Docking must not fire while the ship is still outside the mouth.
func TestDockOnlyInsideShaft(t *testing.T) {
	g := New()
	g.NewGame("GATE")
	g.Cmd.System = pickStation(g, "station")
	g.EnterSystem(g.Cmd.System, "arrive", 9000, false)
	st := g.World.Station
	normal := st.Mat.Col(2)

	g.Player.Pos = st.Pos.Add(normal.Mul(StationSize + 120))
	g.Player.Mat = render.LookMatrix(normal.Mul(-1))
	g.Player.Speed = 50
	for i := 0; i < 30; i++ {
		g.Update(1.0 / 60.0)
	}
	if g.State != "flight" {
		t.Fatalf("docked %d units short of the mouth (state %q)", 120, g.State)
	}

	g.Player.Pos = st.Pos.Add(normal.Mul(StationSize - 25))
	g.Update(1.0 / 60.0)
	if g.State == "flight" {
		t.Fatal("did not dock once inside the shaft")
	}
}
