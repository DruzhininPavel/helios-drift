package game

import (
	"testing"

	"heliosdrift/render"
)

// Yaw right must put the nose to the ship's right (+X in ship space). The
// first cut had the sign flipped, so the stick felt inverted.
func TestYawTurnsRight(t *testing.T) {
	g := newGameMode(ModeAdvanced)
	g.Launch()
	g.Player.Mat = render.RotX(0) // a known heading, so the test is deterministic
	before := g.Player.Mat.Col(2).Mul(-1)[0]
	g.In.Yaw = 1
	for i := 0; i < 20; i++ {
		g.Update(1.0 / 60.0)
	}
	after := g.Player.Mat.Col(2).Mul(-1)[0]
	if after <= before {
		t.Fatalf("yaw right moved the nose to x=%.3f from %.3f: it turns left", after, before)
	}
}

// Yaw left is the mirror case.
func TestYawTurnsLeft(t *testing.T) {
	g := newGameMode(ModeAdvanced)
	g.Launch()
	g.Player.Mat = render.RotX(0) // a known heading, so the test is deterministic
	before := g.Player.Mat.Col(2).Mul(-1)[0]
	g.In.Yaw = -1
	for i := 0; i < 20; i++ {
		g.Update(1.0 / 60.0)
	}
	after := g.Player.Mat.Col(2).Mul(-1)[0]
	if after >= before {
		t.Fatalf("yaw left moved the nose to x=%.3f from %.3f", after, before)
	}
}
