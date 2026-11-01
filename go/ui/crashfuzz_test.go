package ui

import (
	"math/rand"
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// Drive the real screens while pirates shoot at a talking pilot: a crash here is
// a crash the player can reach.
func TestScreensSurviveAPirateAttack(t *testing.T) {
	if testing.Short() {
		t.Skip("frame fuzz: run without -short")
	}
	lang.Set("ru")
	defer lang.Set("en")
	fb := render.NewFB(320, 224)
	keys := []string{"ArrowUp", "ArrowDown", "Enter", "Escape", " ", "b", "s", "PageUp", "PageDown", "1", "2", "3"}
	for seed := int64(0); seed < 40; seed++ {
		rng := rand.New(rand.NewSource(seed))
		g := game.New()
		g.NewGame("FUZZ")
		g.Cmd.Campaign.Stage = 2
		g.Cmd.Campaign.Active = true
		g.State = "flight"
		g.Screen = ""
		g.EnterSystem(3, "launch", 0, false)
		g.Cmd.Campaign.Target = g.Cmd.System
		g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
		g.SpawnMissionGroupForTest("HULK")
		g.SpawnPirates(3)
		for i := 0; i < 600; i++ {
			if rng.Float64() < 0.25 {
				Key(keys[rng.Intn(len(keys))])
			}
			g.Update(1.0 / 60.0)
			Draw(fb, 1.0/60.0)
		}
	}
}
