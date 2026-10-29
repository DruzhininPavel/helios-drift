package game

import (
	"math/rand"

	"heliosdrift/galaxy"
	"heliosdrift/models"
	"heliosdrift/render"
)

// SpawnPlate drops a piece of hull plating from a wreck. Scooping it yields a
// tonne of alloys.
func (g *Game) SpawnPlate(pos render.Vec3) {
	g.World.Entities = append(g.World.Entities, &Entity{
		Kind: "plate", Pos: pos.Add(randDir().Mul(30)),
		Mat: render.LookMatrix(randDir()), Speed: randRange(10, 45),
		Radius: 9, CargoType: 9, Mesh: models.Get("plate"), Life: 180,
	})
}

// spawnDeparture launches a trader out of the station slot - the traffic you
// see in the original comes from the station, not out of thin space.
func (g *Game) spawnDeparture() {
	st := g.World.Station
	if st == nil {
		return
	}
	types := []string{"adder", "python", "cobra3", "shuttle", "moray", "worm", "boa", "anaconda", "transporter"}
	t := types[rand.Intn(len(types))]
	normal := st.Mat.Col(2)
	pos := st.Pos.Add(normal.Mul(StationSize + 70))
	e := g.baseShip(t, pos, "depart")
	e.Mat = render.LookMatrix(normal)
	e.Speed = e.MaxSpeed * 0.65
	e.Age = 0
	e.CargoType = 1 + rand.Intn(6)
	g.World.Entities = append(g.World.Entities, e)
}

// Jettison kicks one tonne of cargo out of the hold; it becomes a scoopable
// canister the way it does in the original.
func (g *Game) Jettison() bool {
	if g.State != "flight" || g.World == nil {
		return false
	}
	if g.Classic() {
		g.Message("NOT AVAILABLE IN CLASSIC MODE", "bad")
		return false
	}
	cmd := &g.Cmd
	for i := 0; i < len(galaxy.Commodities); i++ {
		if cmd.Cargo[i] > 0 {
			cmd.Cargo[i]--
			pos := g.Player.Pos.Add(g.Player.Mat.Col(2).Mul(-40))
			e := g.baseCargo(pos, i)
			g.World.Entities = append(g.World.Entities, e)
			g.Message("JETTISONED "+galaxy.Commodities[i].Name, "info")
			return true
		}
	}
	g.Message("NOTHING TO JETTISON", "bad")
	return false
}

// baseCargo builds a free-floating canister of the given commodity.
func (g *Game) baseCargo(pos render.Vec3, kind int) *Entity {
	return &Entity{
		Kind: "canister", Pos: pos, Mat: render.LookMatrix(randDir()),
		Speed: 20, Radius: 8, CargoType: kind, Mesh: models.Get("canister"), Life: 240,
	}
}
