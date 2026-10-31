package bot

import (
	"fmt"
	"sort"

	"heliosdrift/game"
)

// ShipResult is one line of the hull benchmark: what a hull is actually worth
// in a fight and in a cargo run, measured rather than guessed.
type ShipResult struct {
	Name        string
	Price       int
	Cargo       int
	Speed       float64
	Shield      float64
	Energy      float64
	Turn        float64 // pitch multiplier, as the handbook numbers
	KillSeconds float64 // seconds survived against three pirates
	DamageTaken float64 // shot damage absorbed while doing it
}

// BenchmarkShips flies every hull in the yard through the same two trials:
// three pirates in a fresh system, and ten minutes of trading. The bot does
// the flying, so the only thing that changes between runs is the hull.
func BenchmarkShips() []ShipResult {
	var out []ShipResult
	for i := range game.ShipYard {
		h := game.ShipYard[i]
		res := ShipResult{
			Name: h.Name, Price: h.Price, Cargo: h.Cargo,
			Speed: h.Speed, Shield: h.Shield, Energy: h.Energy, Turn: h.Pitch,
		}
		res.KillSeconds, res.DamageTaken = surviveTrial(h.ID)
		out = append(out, res)
	}
	sort.SliceStable(out, func(a, b int) bool { return out[a].Price < out[b].Price })
	return out
}

// surviveTrial drops the hull into a fight it cannot win and measures how
// long it lasts. That is the number a player feels: how much time a hauler
// buys to run, and how long a fighter can stay in the scrap.
func surviveTrial(ship string) (seconds, damage float64) {
	var total, hits float64
	const trials = 5
	for i := 0; i < trials; i++ {
		game.SetSeed(int64(21 + i*13))
		g := game.New()
		g.NewGameMode(game.ModeAdvanced, "BOT")
		g.Cmd.Ship = ship
		g.Cmd.Laser = 2
		if g.Cmd.Laser > g.Ship().MaxLaser {
			g.Cmd.Laser = g.Ship().MaxLaser
		}
		g.Launch()
		g.Player.Invulnerable = 0
		g.SpawnPirates(3)
		pol := NewTrader(int64(21 + i*13))
		pol.Ace = true

		done := 240.0
		for f := 0; f < 60*240; f++ {
			in, ev := pol.Act(g, f)
			for _, e := range ev {
				applyEvent(g, e)
			}
			g.In = in
			g.Update(DT)
			if g.State == "dead" {
				done = float64(f) / 60
				break
			}
		}
		total += done
		hits += g.DamageTaken
	}
	return total / trials, hits / trials
}

// tankTrial measures how much punishment a hull takes in a minute against
// three pirates. The pilot is kept alive so the number is the damage, not the
// moment of death.
func tankTrial(ship string) float64 {
	var total float64
	const trials = 3
	for i := 0; i < trials; i++ {
		game.SetSeed(int64(41 + i*13))
		g := game.New()
		g.NewGameMode(game.ModeAdvanced, "BOT")
		g.Cmd.Ship = ship
		g.Cmd.Laser = 2
		if g.Cmd.Laser > g.Ship().MaxLaser {
			g.Cmd.Laser = g.Ship().MaxLaser
		}
		g.Launch()
		g.Player.Invulnerable = 0
		g.SpawnPirates(3)
		pol := NewTrader(int64(41 + i*13))
		pol.Ace = true
		for f := 0; f < 60*60; f++ {
			in, ev := pol.Act(g, f)
			for _, e := range ev {
				applyEvent(g, e)
			}
			if g.State == "dead" {
				revive(g)
			}
			g.In = in
			g.Update(DT)
		}
		total += g.DamageTaken
	}
	return total / trials
}

// ShipTable renders the benchmark as a plain text table.
func ShipTable(rows []ShipResult) string {
	s := fmt.Sprintf("%-12s %8s %5s %5s %5s %6s %8s %9s\n",
		"hull", "price", "cargo", "vmax", "shield", "turn", "survive", "damage")
	for _, r := range rows {
		kill := fmt.Sprintf("%.1fs", r.KillSeconds)
		s += fmt.Sprintf("%-12s %8d %5d %5.0f %6.0f %6.2f %8s %9.0f\n",
			r.Name, r.Price, r.Cargo, r.Speed, r.Shield, r.Turn, kill, r.DamageTaken)
	}
	return s
}
