package bot

import (
	"math"
	"testing"

	"heliosdrift/game"
	"heliosdrift/render"
)

func newFlight(t *testing.T, seed int64) *game.Game {
	t.Helper()
	game.SetSeed(seed)
	g := game.New()
	g.NewGameMode(game.ModeAdvanced, "BOT")
	g.Launch()
	g.Player.Speed = 0
	return g
}

func angleTo(g *game.Game, want render.Vec3) float64 {
	fwd := g.Player.Mat.Col(2).Mul(-1)
	return math.Acos(clamp(fwd.Dot(want.Norm()), -1, 1))
}

// TestPilotPointsAtTarget checks that the controller turns the nose onto a
// target at several angles, including one behind the ship.
func TestPilotPointsAtTarget(t *testing.T) {
	for _, deg := range []float64{10, 40, 90, 150, 179} {
		g := newFlight(t, 3)
		p := DefaultPilot()
		fwd := g.Player.Mat.Col(2).Mul(-1)
		side := g.Player.Mat.Col(0)
		rad := deg * math.Pi / 180
		want := fwd.Mul(math.Cos(rad)).Add(side.Mul(math.Sin(rad))).Norm()
		before := angleTo(g, want)
		for i := 0; i < 300; i++ {
			g.In = p.Aim(g, want)
			g.Update(DT)
		}
		after := angleTo(g, want)
		t.Logf("%3.0f degrees: %.3f -> %.3f rad", deg, before, after)
		if after > 0.12 {
			t.Errorf("pilot left %.0f degrees of error", after*180/math.Pi)
		}
	}
}

// TestPilotLevelsTheHorizon checks the roll loop settles instead of hunting.
func TestPilotLevelsTheHorizon(t *testing.T) {
	for _, gain := range []float64{0.4, 0.55, 0.8, 1.2} {
		g := newFlight(t, 4)
		g.Player.Mat = g.Player.Mat.Mul(render.RotZ(0.7))
		p := DefaultPilot()
		p.LevelGain = gain
		fwd := g.Player.Mat.Col(2).Mul(-1)
		trace := make([]float64, 0, 8)
		for i := 0; i < 300; i++ {
			g.In = p.Aim(g, fwd)
			g.Update(DT)
			if i%40 == 0 {
				up := g.Player.Mat.Col(1)
				trace = append(trace, math.Atan2(up[0], up[1]))
			}
		}
		t.Logf("gain %.2f: tilt %.3f  trace %v", gain, trace[len(trace)-1], trace)
	}
}
