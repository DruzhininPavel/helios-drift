package main

import (
	"fmt"
	"math"

	"heliosdrift/game"
	"heliosdrift/render"
)

func combatProbe() {
	g := game.New()
	g.NewGame("T")
	g.Launch()
	g.World.Entities = nil
	g.SpawnPirates(2)
	for _, e := range g.World.Entities {
		e.Pos = g.Player.Pos.Add(render.V3(0, 0, -1500))
		e.Mat = render.LookMatrix(g.Player.Pos.Sub(e.Pos).Norm())
	}
	for i := 0; i < 600; i++ {
		g.Update(1 / 60.0)
		if i%120 == 0 {
			e := g.World.Entities[0]
			d := e.Pos.Sub(g.Player.Pos).Len()
			local := e.Mat.ApplyT(g.Player.Pos.Sub(e.Pos).Norm())
			ax := math.Atan2(local[0], -local[2])
			ay := math.Atan2(local[1], math.Sqrt(local[0]*local[0]+local[2]*local[2]))
			fmt.Printf("f%d dist %.0f spd %.0f cd %.2f ax %.3f ay %.3f energy %.0f\n", i, d, e.Speed, e.Cooldown, ax, ay, g.Player.Energy)
		}
		if g.State != "flight" {
			fmt.Println("state", g.State)
			break
		}
	}
	fmt.Println("final energy", g.Player.Energy, "shield", g.Player.ShieldF)
}
