package bot

import (
	"strconv"

	"heliosdrift/game"
)

// Small helpers the human-style policy needs: the game keeps the yard and the
// trade-in value, the bot just asks for them.

func itoa(n int) string { return strconv.Itoa(n) }

// shipIndex is the position of a hull in the yard list.
func shipIndex(id string) int {
	for i := range game.ShipYard {
		if game.ShipYard[i].ID == id {
			return i
		}
	}
	return -1
}

// ShipName is the display name of a hull.
func ShipName(id string) string { return game.ShipFor(id).Name }

// ShipPrice is what the yard asks for a hull.
func ShipPrice(id string) int { return game.ShipFor(id).Price }

// tradeIn is what the yard gives for the ship being flown.
func tradeIn(g *game.Game) int { return g.TradeInValue() }

// contractTarget is the nearest system with a job to finish: a courier drop or
// a bounty to collect. Jobs are worth a detour, the campaign can wait a stop.
func (t *Trader) contractTarget(g *game.Game) (int, bool) {
	best, bestPay := -1, 0
	for _, c := range g.Cmd.Contracts {
		if c.Galaxy != g.Cmd.Galaxy || c.System == g.Cmd.System {
			continue
		}
		if c.Pay > bestPay {
			best, bestPay = c.System, c.Pay
		}
	}
	if best < 0 {
		return 0, false
	}
	return best, true
}
