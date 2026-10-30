package bot

import (
	"testing"

	"heliosdrift/galaxy"
	"heliosdrift/game"
)

// TestCampaignKillStage puts the bot straight into the ambush stage of the
// campaign and watches whether it can clear the tagged contacts.
func TestCampaignKillStage(t *testing.T) {
	game.SetSeed(7)
	g := game.New()
	g.NewGameMode(game.ModeAdvanced, "BOT")
	g.Cmd.Equip["dockcomp"] = true
	g.Launch()

	// jump straight to the second stage: three tagged contacts to destroy
	gal := galaxy.Get(0)
	target := -1
	for _, s := range gal.Systems {
		if d := galaxy.Distance(gal.Systems[0], s); d > 0.5 && d < 6.9 {
			target = s.ID
			break
		}
	}
	if target < 0 {
		t.Fatal("no system in range")
	}
	g.Cmd.Campaign.Stage = 1
	g.Cmd.Campaign.Active = true
	g.Cmd.Campaign.Target = target
	g.Cmd.Campaign.Galaxy = 0
	g.Cmd.Fuel = 7
	if !g.Hyperspace(target) {
		t.Fatal("hyperspace refused")
	}
	for i := 0; i < 600; i++ {
		g.Update(DT)
		if g.Cmd.System == target {
			break
		}
	}

	tagged := 0
	for _, e := range g.World.Entities {
		if e.CampaignTag != "" {
			tagged++
		}
	}
	t.Logf("in system %d with %d tagged contacts, %d entities", g.Cmd.System, tagged, len(g.World.Entities))
	if tagged == 0 {
		t.Fatalf("no campaign contacts spawned in the target system")
	}

	pol := NewTrader(7)
	pol.Ace = true
	fireFrames, shots := 0, 0
	for f := 0; f < 60*60*5; f++ {
		in, ev := pol.Act(g, f)
		for _, e := range ev {
			applyEvent(g, e)
		}
		if in.Fire {
			fireFrames++
		}
		if g.Player.Cooldown <= 0 && in.Fire && g.State == "flight" {
			shots++
		}
		g.Player.Invulnerable = 1e6
		g.In = in
		g.Update(DT)
		if f%600 == 0 {
			left := 0
			for _, e := range g.World.Entities {
				if e.CampaignTag != "" && e.Energy > 0 {
					left++
				}
			}
			foe := pol.nearestFoe(g)
			name := "none"
			if foe != nil {
				name = foe.Type
			}
			hp := 0.0
			if foe != nil {
				hp = foe.Energy
			}
			t.Logf("f%5d stage=%d left=%d foe=%s hp=%.0f dist=%.0f fire=%d/%d state=%s",
				f, g.Cmd.Campaign.Stage, left, name, hp,
				g.Player.Pos.Sub(g.World.Station.Pos).Len(), fireFrames, shots, g.State)
		}
		if g.Cmd.Campaign.Stage > 1 {
			t.Logf("stage cleared at frame %d", f)
			return
		}
	}
	t.Fatalf("stage not cleared after 5 minutes")
}

// TestCampaignReachesFinalStage is the playability check the whole harness
// exists for: can the campaign be walked from the first stage to the last
// without god mode? A player who does not die early must be able to finish the
// story, so the bot has to get there with its own flying and shooting.
func TestCampaignReachesFinalStage(t *testing.T) {
	if testing.Short() {
		t.Skip("long run")
	}
	game.SetSeed(40)
	g := game.New()
	g.NewGameMode(game.ModeAdvanced, "BOT")
	g.Cmd.Equip["dockcomp"] = true
	g.Cmd.Equip["fuelscoop"] = true
	pol := NewTrader(40)
	sc := Run(g, 60*60*10, pol)

	if sc.Report.CampaignDone {
		t.Logf("campaign complete at stage %d, cash %.0f", sc.Report.Stage, g.Cmd.Cash)
		return
	}
	t.Fatalf("campaign stalled at stage %d of %d after %.0f s (death: %s)",
		sc.Report.Stage, game.CampaignTotal(), float64(sc.Report.Frames)/60, sc.Report.Death)
}
