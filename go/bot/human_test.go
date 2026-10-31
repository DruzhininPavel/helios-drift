package bot

import (
	"testing"

	"heliosdrift/game"
)

// TestHumanModeDwellsAtStations checks the difference between a speedrun and a
// playthrough: the human-style commander spends real time docked, the speedrun
// almost none.
func TestHumanModeDwellsAtStations(t *testing.T) {
	share := func(human bool) float64 {
		game.SetSeed(7)
		g := game.New()
		g.NewGameMode(game.ModeAdvanced, "BOT")
		g.Cmd.Equip["dockcomp"] = true
		g.Cmd.Equip["fuelscoop"] = true
		pol := NewTrader(7)
		pol.Human = human
		sc := Run(g, 60*60*6, pol)
		docked := float64(sc.Report.States["docked"])
		return docked / float64(sc.Report.Frames)
	}
	fast, human := share(false), share(true)
	t.Logf("docked share: speedrun %.1f%%, human %.1f%%", fast*100, human*100)
	if human < 0.25 {
		t.Errorf("human run spends only %.0f%% docked", human*100)
	}
	if human <= fast*3 {
		t.Errorf("human %.2f vs speedrun %.2f: no difference in pace", human, fast)
	}
}

// TestHumanModeSavingsPlan checks that the commander actually trades up: the
// journal has to show a hull change, and the ship in the end is the better one.
func TestHumanModeSavingsPlan(t *testing.T) {
	if testing.Short() {
		t.Skip("long run")
	}
	game.SetSeed(7)
	g := game.New()
	g.NewGameMode(game.ModeAdvanced, "BOT")
	g.Cmd.Equip["dockcomp"] = true
	g.Cmd.Equip["fuelscoop"] = true
	pol := NewTrader(7)
	pol.Human = true
	sc := Run(g, 60*60*40, pol)

	if g.Cmd.Ship == "cobra3" {
		t.Fatalf("never traded up: still in the starter after %.0f minutes", float64(sc.Report.Frames)/3600)
	}
	if !hasLine(sc.Report.Journal, "now flying") {
		t.Fatal("journal has no hull change")
	}
	t.Logf("finished in a %s, %d journal lines", g.Ship().Name, len(sc.Report.Journal))
}

func hasLine(lines []string, want string) bool {
	for _, l := range lines {
		for i := 0; i+len(want) <= len(l); i++ {
			if l[i:i+len(want)] == want {
				return true
			}
		}
	}
	return false
}
