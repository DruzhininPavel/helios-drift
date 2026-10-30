package bot

import (
	"math"

	"heliosdrift/game"
	"heliosdrift/render"
)

// Pilot flies the ship with a hand-like controller. It works out the rotation
// that would put the nose on a target and asks for it through the same three
// axes a player has: pitch, roll and yaw. Nothing here moves the ship
// directly, so the bot flies no better than the flight model allows.
type Pilot struct {
	Gain      float64 // how hard to chase the target, 1 = flat out
	LevelGain float64 // how hard to roll the cabin level
	Deadband  float64 // ignore errors smaller than this, in radians
	UseYaw    bool    // yaw only exists in Advanced mode
}

// DefaultPilot returns the calibrated controller.
func DefaultPilot() *Pilot {
	return &Pilot{Gain: 0.9, LevelGain: 1.0, Deadband: 0.01, UseYaw: true}
}

// Aim returns stick input that turns the nose towards the given world
// direction.
func (p *Pilot) Aim(g *game.Game, want render.Vec3) game.InputState {
	in := game.InputState{}
	m := g.Player.Mat
	fwd := m.Col(2).Mul(-1)
	want = want.Norm()

	angle := math.Acos(clamp(fwd.Dot(want), -1, 1))
	if angle < p.Deadband {
		return p.level(g)
	}
	// the rotation that carries the nose onto the target
	axis := fwd.Cross(want)
	if axis.Len() < 1e-9 {
		// exactly opposite: any perpendicular axis will do
		axis = m.Col(1)
	}
	axis = axis.Norm()

	// ask for a share of that rotation, capped at what the ship can do
	step := math.Min(angle*p.Gain, 1.2) // radians this tick at most
	omega := axis.Mul(step)

	// into the ship's own axes: x right, y up (yaw), z back (roll)
	ox := omega.Dot(m.Col(0))
	oy := omega.Dot(m.Col(1))
	oz := omega.Dot(m.Col(2))

	in.Pitch = clamp(ox/game.PitchRate, -1, 1)
	if p.UseYaw {
		in.Yaw = clamp(-oy/game.YawRate, -1, 1)
	} else {
		// no yaw: bank so the target sits overhead, then pull
		if math.Abs(oy) > 0.05 {
			in.Roll = clamp(-math.Copysign(1, oy)*0.9, -1, 1)
		}
	}
	in.Roll = clamp(in.Roll-oz/game.RollRate, -1, 1)

	// once the nose is close, roll the cabin level: the docking slot is a
	// narrow horizontal slit and a rolled ship cannot fit it
	if angle < 0.35 {
		up := m.Col(1)
		level := math.Atan2(up[0], up[1])
		if math.Abs(level) > p.Deadband {
			in.Roll = clamp(level*1.4/game.RollRate+g.In.Roll*0, -1, 1)
			in.Roll = clamp(-level*p.LevelGain, -1, 1)
		}
	}
	return in
}

// Arrive flies to a point instead of past it: the wanted speed falls with the
// distance left, and a nose that is off the target cuts the throttle.
func Arrive(g *game.Game, target render.Vec3, cap float64) float64 {
	to := target.Sub(g.Player.Pos)
	d := to.Len()
	if d < 1e-6 {
		return 0
	}
	fwd := g.Player.Mat.Col(2).Mul(-1)
	aim := clamp((fwd.Dot(to.Norm())-0.5)/0.4, 0, 1)
	want := math.Min(cap, d*0.25) * aim
	switch {
	case g.Player.Speed > want*1.15:
		return -1
	case g.Player.Speed < want*0.85:
		return 1
	}
	return 0
}

// level rolls the cabin upright without turning the nose.
func (p *Pilot) level(g *game.Game) game.InputState {
	in := game.InputState{}
	up := g.Player.Mat.Col(1)
	lv := math.Atan2(up[0], up[1])
	if math.Abs(lv) > p.Deadband {
		in.Roll = clamp(-lv*p.LevelGain, -1, 1)
	}
	return in
}

func clamp(v, lo, hi float64) float64 {
	if math.IsNaN(v) {
		return 0
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
