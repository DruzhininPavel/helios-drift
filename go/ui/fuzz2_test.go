package ui

import (
	"math/rand"
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// Everything at once: a wreck in the sky, pirates, missiles, ramming, jumps and
// conversations, with every frame drawn.
func TestScreensSurviveEverything(t *testing.T) {
	if testing.Short() {
		t.Skip("frame fuzz: run without -short")
	}
	lang.Set("ru")
	defer lang.Set("en")
	fb := render.NewFB(320, 224)
	keys := []string{"ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight", " ", "b", "s", "m", "j", "1", "2", "3", "Enter", "Escape", "PageUp", "PageDown"}
	for seed := int64(0); seed < 120; seed++ {
		rng := rand.New(rand.NewSource(seed))
		g := game.New()
		g.NewGame("FUZZ2")
		g.Cmd.Campaign.Stage = 2
		g.Cmd.Campaign.Active = true
		g.State = "flight"
		g.EnterSystem(3, "launch", 0, false)
		g.Cmd.Campaign.Target = g.Cmd.System
		g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
		g.SpawnMissionGroupForTest("HULK")
		g.SpawnPirates(3)
		for i := 0; i < 900; i++ {
			if rng.Float64() < 0.3 {
				Key(keys[rng.Intn(len(keys))])
			}
			if rng.Float64() < 0.05 {
				g.FireMissile()
			}
			g.Update(1.0 / 60.0)
			Draw(fb, 1.0/60.0)
		}
	}
}
