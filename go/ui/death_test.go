package ui

import (
	"math/rand"
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// Dying with a conversation open, restarting, dying again: the screens have to
// survive a pilot who is not paying attention.
func TestScreensSurviveDeathAndRestart(t *testing.T) {
	if testing.Short() {
		t.Skip("frame fuzz: run without -short")
	}
	lang.Set("ru")
	defer lang.Set("en")
	fb := render.NewFB(320, 224)
	keys := []string{"Enter", "Escape", " ", "1", "2", "3", "r", "q"}
	for seed := int64(0); seed < 25; seed++ {
		rng := rand.New(rand.NewSource(seed))
		g := game.New()
		g.NewGame("DIE")
		g.Cmd.Campaign.Stage = 2
		g.Cmd.Campaign.Active = true
		g.State = "flight"
		g.EnterSystem(3, "launch", 0, false)
		g.Cmd.Campaign.Target = g.Cmd.System
		g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
		g.TalkStart("c3_roster")
		g.SpawnPirates(3)
		for i := 0; i < 400; i++ {
			if i%40 == 0 {
				g.DamagePlayer(60, true, "laser")
			}
			if rng.Float64() < 0.3 {
				Key(keys[rng.Intn(len(keys))])
			}
			g.Update(1.0 / 60.0)
			Draw(fb, 1.0/60.0)
		}
	}
}
