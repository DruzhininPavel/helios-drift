package game

import (
	"math/rand"
	"testing"
)

func TestChapterThreeWithPiratesDoesNotPanic(t *testing.T) {
	for seed := int64(0); seed < 60; seed++ {
		rng := rand.New(rand.NewSource(seed))
		g := New()
		g.NewGame("CRASH")
		g.Cmd.Campaign.Stage = 2
		g.Cmd.Campaign.Active = true
		g.EnterSystem(3, "launch", 0, false)
		g.Cmd.Campaign.Target = g.Cmd.System
		g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
		g.spawnQuestGroup("HULK")
		g.SpawnPirates(3)
		for i := 0; i < 60*20; i++ {
			if rng.Float64() < 0.02 {
				g.FireLaser()
			}
			if rng.Float64() < 0.01 && g.Screen == "talk" {
				g.TalkChoose(rng.Intn(3))
			}
			g.Update(1.0 / 60.0)
			if g.Screen == "talk" {
				g.Update(1.0 / 60.0)
			}
		}
	}
}
