package game

import (
	"testing"

	"heliosdrift/models"
)

func countKind(g *Game, kind string) int {
	n := 0
	for _, e := range g.World.Entities {
		if e.Kind == kind {
			n++
		}
	}
	return n
}

// Mining: asteroid -> boulders -> splinters -> scoop into the hold.
func TestMiningChain(t *testing.T) {
	g := New()
	g.NewGame("MINER")
	g.Launch()
	g.Cmd.Equip["mining"] = true
	g.Cmd.Equip["fuelscoop"] = true
	g.World.Entities = nil

	fwd := g.Player.Mat.Col(2).Mul(-1)
	ast := &Entity{
		Kind: "asteroid", Pos: g.Player.Pos.Add(fwd.Mul(600)), Radius: 30,
		Energy: 20, Mesh: models.Get("asteroid"),
	}
	g.World.Entities = append(g.World.Entities, ast)

	g.HitEntity(ast, 100, "mining")
	if n := countKind(g, "boulder"); n != 3 {
		t.Fatalf("asteroid broke into %d boulders, want 3", n)
	}

	for _, e := range append([]*Entity(nil), g.World.Entities...) {
		if e.Kind == "boulder" {
			g.HitEntity(e, 100, "mining")
		}
	}
	if n := countKind(g, "splinter"); n == 0 {
		t.Fatal("boulders did not break into splinters")
	}
	if n := countKind(g, "boulder"); n != 0 {
		t.Fatalf("%d boulders left after breaking", n)
	}

	var sp *Entity
	for _, e := range g.World.Entities {
		if e.Kind == "splinter" {
			sp = e
			break
		}
	}
	before := CargoUsed(g.Cmd)
	g.Player.Pos = sp.Pos
	for i := 0; i < 20 && CargoUsed(g.Cmd) == before; i++ {
		g.Update(0.02)
	}
	if CargoUsed(g.Cmd) <= before {
		t.Fatal("splinter was not scooped")
	}
}

// Without a scoop nothing is collected, and the hold limit still applies.
func TestMiningNeedsScoop(t *testing.T) {
	g := New()
	g.NewGame("MINER")
	g.Launch()
	g.World.Entities = nil
	sp := &Entity{Kind: "splinter", Pos: g.Player.Pos.Add(g.Player.Mat.Col(2).Mul(-10)),
		Radius: 4, CargoType: 12, Mesh: models.Get("splinter"), Life: 60}
	g.World.Entities = append(g.World.Entities, sp)
	for i := 0; i < 10; i++ {
		g.Update(0.02)
	}
	if CargoUsed(g.Cmd) != 0 {
		t.Fatal("cargo collected without a fuel scoop")
	}
}
