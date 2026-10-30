package bot

import (
	"fmt"
	"math"
	"strconv"

	"heliosdrift/galaxy"
	"heliosdrift/game"
)

var knownStates = map[string]bool{
	"title": true, "docked": true, "flight": true, "hyper": true,
	"dead": true, "docking": true, "witchspace": true,
}

// check reports the first impossible state of the frame, or nil.
func check(g *game.Game, frame int) *Violation {
	bad := func(what, value string) *Violation {
		return &Violation{Frame: frame, What: what, Value: value}
	}
	p := &g.Player
	c := &g.Cmd

	if !knownStates[g.State] {
		return bad("unknown state", g.State)
	}
	if nan := math.IsNaN(p.Pos[0]) || math.IsNaN(p.Pos[1]) || math.IsNaN(p.Pos[2]); nan {
		return bad("player position is NaN", fmt.Sprintf("%v", p.Pos))
	}
	for col := 0; col < 3; col++ {
		v := p.Mat.Col(col)
		if math.IsNaN(v[0]) || math.IsNaN(v[1]) || math.IsNaN(v[2]) {
			return bad("orientation is NaN", "matrix")
		}
	}
	if c.Cash < 0 {
		return bad("negative cash", ftoa(c.Cash))
	}
	if c.Fuel < -0.001 || c.Fuel > game.FuelMax+0.001 {
		return bad("fuel out of range", ftoa(c.Fuel))
	}
	if game.CargoUsed(*c) > c.CargoCap {
		return bad("cargo over capacity", strconv.Itoa(game.CargoUsed(*c))+"/"+strconv.Itoa(c.CargoCap))
	}
	if g.State != "docked" && g.State != "dead" {
		if p.Speed < -0.001 {
			return bad("negative speed", ftoa(p.Speed))
		}
		if p.Energy < -0.001 {
			return bad("negative energy", ftoa(p.Energy))
		}
		if p.ShieldF < -0.001 || p.ShieldA < -0.001 {
			return bad("negative shield", ftoa(math.Min(p.ShieldF, p.ShieldA)))
		}
	}
	// stranded: no fuel, no money to buy any, nothing to sell
	if g.State == "docked" && c.Fuel < 0.1 && c.Cash < float64(galaxy.FuelPricePerLY) && game.CargoUsed(*c) == 0 {
		return bad("stranded: no fuel and no way to pay for it", ftoa(c.Cash))
	}
	if p.Cabin < -0.001 || p.Cabin > 101 {
		return bad("cabin temperature out of range", ftoa(p.Cabin))
	}
	return nil
}

func ftoa(v float64) string { return strconv.FormatFloat(v, 'f', 3, 64) }
