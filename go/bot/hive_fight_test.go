package bot

import (
	"fmt"
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
)

// Can a commander with the kit the stage asks for actually kill a mothership?
func TestHiveFightIsWinnable(t *testing.T) {
	lang.Cfg.AxisSpeed = 1
	game.SetSeed(11)
	g := game.New()
	g.NewGameMode(game.ModeAdvanced, "BOT")
	g.Cmd.Equip["dockcomp"] = true
	g.Cmd.Equip["shieldboost1"] = true
	g.Cmd.Laser = 2
	g.Cmd.Missiles = 4
	g.EnterSystem(3, "arrive", 9000, false)
	g.State = "flight"
	g.Cmd.Flags["SHIELDS_CUT"] = true
	g.Cmd.Campaign.Active = true
	g.Cmd.Campaign.Stage = 6
	g.SpawnMissionGroupForTest("MOTHER")
	var mother *game.Entity
	for _, e := range g.World.Entities {
		e.Hostile = true
		if e.CampaignTag == "MOTHER" {
			mother = e
		}
	}
	if mother == nil {
		t.Fatal("no mothership spawned")
	}
	pol := NewTrader(11)
	for f := 0; f < 60*180; f++ {
		in, ev := pol.Act(g, f)
		for _, e := range ev {
			applyEvent(g, e)
		}
		g.In = in
		g.Update(DT)
		if f%600 == 0 {
			fmt.Printf("  t=%3.0fs mother hull=%.0f shield=%.0f dead=%v inWorld=%v\n",
				float64(f)/60, mother.Energy, mother.Shield, mother.Dead, inWorld(g, mother))
		}
		if !inWorld(g, mother) {
			fmt.Printf("mothership down after %.0f s, hull left %.0f\n", float64(f)/60, g.Player.Energy)
			return
		}
		if g.State == "dead" {
			t.Fatalf("died to the hive: %s after %.0f s", g.Dead.Reason, float64(f)/60)
		}
	}
	t.Fatalf("mothership survived three minutes, hull %.0f shield %.0f", mother.Energy, mother.Shield)
}

func inWorld(g *game.Game, e *game.Entity) bool {
	for _, x := range g.World.Entities {
		if x == e {
			return true
		}
	}
	return false
}
