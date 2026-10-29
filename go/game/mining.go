package game

import (
	"math"
	"math/rand"

	"heliosdrift/models"
	"heliosdrift/render"
)

// breakRock splits a rock into smaller pieces. Asteroids become boulders,
// boulders become splinters, and splinters are what the fuel scoop picks up.
func (g *Game) breakRock(e *Entity, kind string, n int) {
	mesh := models.Get(kind)
	if mesh == nil || g.World == nil {
		return
	}
	r := e.Radius * 0.45
	if kind == "splinter" {
		r = math.Max(3, e.Radius*0.55)
	}
	for i := 0; i < n; i++ {
		pos := e.Pos.Add(randDir().Mul(e.Radius * 0.45))
		// ore chips carry minerals, with the occasional precious find
		ct := 12
		if kind == "splinter" && rand.Float64() < 0.12 {
			ct = 13 + rand.Intn(3)
		}
		meshFor := mesh
		if kind == "splinter" {
			if ore := models.Get("splinter_ore"); ore != nil && ct != 12 {
				meshFor = ore
			}
		}
		life := 240.0
		g.World.Entities = append(g.World.Entities, &Entity{
			Kind: kind, Pos: pos, Mat: render.LookMatrix(randDir()),
			Speed: randRange(15, 55), Radius: r, CargoType: ct,
			Mesh: meshFor, Life: life,
			Spin: render.V3(randRange(-1.2, 1.2), randRange(-1.2, 1.2), randRange(-1.2, 1.2)),
		})
	}
	if kind == "boulder" {
		g.Message("ROCK BREAKS UP", "info")
	}
}
