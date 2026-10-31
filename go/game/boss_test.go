package game

import (
	"testing"

	"heliosdrift/render"
)

// TestFlagshipNeedsRealGuns answers the design question directly: how much
// fire does the last stage take with each laser class? The gap is what makes
// it worth earning a better gun before the finale.
func TestFlagshipNeedsRealGuns(t *testing.T) {
	for lvl, name := range LaserName {
		g := New()
		g.NewGameMode(ModeAdvanced, "BOT")
		g.Cmd.Laser = lvl
		g.Launch()
		g.Cmd.Campaign.Stage = CampaignTotal() - 1
		g.Cmd.Campaign.Active = true
		g.Cmd.Campaign.Target = g.Cmd.System
		g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
		g.campaignArrive()

		var boss *Entity
		escorts := 0
		for _, e := range g.World.Entities {
			switch {
			case e.CampaignTag == "FLAGSHIP":
				boss = e
			case e.Behaviour == "pirate" && e.CampaignTag == "":
				escorts++
			}
		}
		if boss == nil {
			t.Fatalf("%s: no flagship spawned", name)
		}
		if escorts < 2 {
			t.Errorf("%s: %d escorts, want at least 2", name, escorts)
		}

		shots, quiet := 0, 0
		for boss.Energy > 0 {
			g.Player.Pos = boss.Pos.Add(render.V3(0, 0, 600))
			g.Player.Mat = render.LookMatrix(render.V3(0, 0, -1))
			g.Player.Cooldown = 0
			g.Player.LaserTemp = 0
			before := boss.Energy + boss.Shield
			g.FireLaser()
			shots++
			if boss.Energy+boss.Shield == before {
				quiet++
			} else {
				quiet = 0
			}
			if quiet > 10 {
				t.Fatalf("%s: flagship took no damage after %d shots (hp %.0f/%.0f)",
					name, shots, boss.Energy, boss.Shield)
			}
			if shots > 5000 {
				t.Fatalf("%s: flagship would not die", name)
			}
		}
		t.Logf("%-20s %4d shots, %.1f s of fire", name, shots, float64(shots)*LaserCooldown[lvl])
	}
}
