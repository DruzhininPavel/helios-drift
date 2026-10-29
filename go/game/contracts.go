// Procedural contracts, the way the later Elites did it: a station board with
// courier runs and bounty jobs generated from the local systems.
package game

import (
	"math"
	"math/rand"

	"heliosdrift/galaxy"
)

type Contract struct {
	Kind   string `json:"kind"` // courier | bounty
	Good   int    `json:"good"`
	Qty    int    `json:"qty"`
	System int    `json:"system"`
	Galaxy int    `json:"galaxy"`
	Pay    int    `json:"pay"`
	Done   int    `json:"done"`
	Due    int    `json:"due"`
}

const (
	ContractCourier   = "courier"
	ContractBounty    = "bounty"
	ContractPassenger = "passenger"
	ContractMail      = "mail"
	ContractSmuggle   = "smuggle"
)

// rollContracts refreshes the station board.
func (g *Game) rollContracts() {
	g.Offers = nil
	if g.Classic() {
		return // contract boards came with the later Elites
	}
	if len(g.Cmd.Contracts) >= 3 {
		return
	}
	gal := galaxy.Get(g.Cmd.Galaxy)
	cur := gal.Systems[g.Cmd.System]
	n := 1 + rand.Intn(3)
	for i := 0; i < n; i++ {
		tgt := g.missionTarget()
		d := galaxy.Distance(cur, gal.Systems[tgt])
		due := g.Cmd.Day + 3 + rand.Intn(5)
		switch r := rand.Float64(); {
		case r < 0.42:
			good := rand.Intn(12)
			qty := 1 + rand.Intn(6)
			pay := int(math.Round(float64(qty)*galaxy.Commodities[good].Base*0.55 + d*55 + rand.Float64()*40))
			g.Offers = append(g.Offers, Contract{
				Kind: ContractCourier, Good: good, Qty: qty,
				System: tgt, Galaxy: g.Cmd.Galaxy, Pay: pay, Due: due,
			})
		case r < 0.55 && g.World != nil && g.World.Sys.Government <= 1:
			// contraband only moves through anarchy and feudal space
			good := 3
			if rand.Float64() < 0.5 {
				good = 6
			}
			qty := 2 + rand.Intn(4)
			pay := qty * (240 + rand.Intn(160))
			g.Offers = append(g.Offers, Contract{
				Kind: ContractSmuggle, Good: good, Qty: qty,
				System: tgt, Galaxy: g.Cmd.Galaxy, Pay: pay, Due: due,
			})
		case r < 0.70:
			qty := 1 + rand.Intn(4)
			g.Offers = append(g.Offers, Contract{
				Kind: ContractPassenger, Qty: qty,
				System: tgt, Galaxy: g.Cmd.Galaxy, Pay: qty * (90 + rand.Intn(70)), Due: due + 1,
			})
		case r < 0.78:
			qty := 1 + rand.Intn(5)
			g.Offers = append(g.Offers, Contract{
				Kind: ContractMail, Qty: qty,
				System: tgt, Galaxy: g.Cmd.Galaxy, Pay: qty * (45 + rand.Intn(35)), Due: due,
			})
		default:
			qty := 2 + rand.Intn(4)
			g.Offers = append(g.Offers, Contract{
				Kind: ContractBounty, Qty: qty, Done: 0,
				System: tgt, Galaxy: g.Cmd.Galaxy, Pay: qty * (110 + rand.Intn(60)), Due: due + 2,
			})
		}
	}
}

// ContractLabel renders one board row.
func (g *Game) ContractLabel(c Contract) string {
	sys := galaxy.Get(c.Galaxy).Systems[c.System].Name
	tail := " " + itoa(c.Pay) + " CR " + itoa(maxInt(0, c.Due-g.Cmd.Day)) + " DAYS"
	switch c.Kind {
	case ContractCourier:
		return "COURIER " + galaxy.Commodities[c.Good].Name + " x" + itoa(c.Qty) + " TO " + sys + tail
	case ContractPassenger:
		return "PASSENGERS " + itoa(c.Qty) + " TO " + sys + tail
	case ContractMail:
		return "MAIL " + itoa(c.Qty) + " T TO " + sys + tail
	case ContractSmuggle:
		return "SMUGGLE " + galaxy.Commodities[c.Good].Name + " x" + itoa(c.Qty) + " TO " + sys + tail
	}
	return "BOUNTY " + itoa(c.Qty-c.Done) + " PIRATES IN " + sys + tail
}

// ContractDetail is the journal line for an accepted job, with progress.
func (g *Game) ContractDetail(c Contract) string {
	sys := galaxy.Get(c.Galaxy).Systems[c.System].Name
	days := itoa(maxInt(0, c.Due-g.Cmd.Day)) + "D"
	switch c.Kind {
	case ContractCourier:
		return "COURIER " + galaxy.Commodities[c.Good].Name + " x" + itoa(c.Qty) + " TO " + sys + "  " + days
	case ContractPassenger:
		return "PASSENGERS " + itoa(c.Qty) + " TO " + sys + "  " + days
	case ContractMail:
		return "MAIL " + itoa(c.Qty) + "T TO " + sys + "  " + days
	case ContractSmuggle:
		return "SMUGGLE " + galaxy.Commodities[c.Good].Name + " x" + itoa(c.Qty) + " TO " + sys + "  " + days
	}
	return "BOUNTY " + itoa(c.Done) + "/" + itoa(c.Qty) + " IN " + sys + "  " + days
}

// AcceptOffer takes a board offer. Courier runs load the cargo immediately.
func (g *Game) AcceptOffer(i int) bool {
	if g.Classic() {
		g.Message("NOT AVAILABLE IN CLASSIC MODE", "bad")
		return false
	}
	if i < 0 || i >= len(g.Offers) {
		return false
	}
	c := g.Offers[i]
	if len(g.Cmd.Contracts) >= 3 {
		g.Message("TOO MANY CONTRACTS ALREADY", "bad")
		return false
	}
	switch c.Kind {
	case ContractCourier, ContractSmuggle:
		if CargoUsed(g.Cmd)+c.Qty > g.Cmd.CargoCap {
			g.Message("NOT ENOUGH CARGO SPACE", "bad")
			return false
		}
		g.Cmd.Cargo[c.Good] += c.Qty
	case ContractPassenger:
		g.Cmd.Pax += c.Qty
	case ContractMail:
		g.Cmd.Mail += c.Qty
	}
	g.Cmd.Contracts = append(g.Cmd.Contracts, c)
	g.Offers = append(g.Offers[:i], g.Offers[i+1:]...)
	g.Message("CONTRACT ACCEPTED", "good")
	return true
}

// deliverContracts pays out courier runs at the destination station.
func (g *Game) deliverContracts() {
	cmd := &g.Cmd
	keep := cmd.Contracts[:0]
	for _, c := range cmd.Contracts {
		match := c.Galaxy == cmd.Galaxy && c.System == cmd.System
		if match {
			switch c.Kind {
			case ContractCourier, ContractSmuggle:
				if cmd.Cargo[c.Good] >= c.Qty {
					cmd.Cargo[c.Good] -= c.Qty
					cmd.Cash += float64(c.Pay)
					g.Message("DELIVERED  +"+itoa(c.Pay)+" CR", "good")
					continue
				}
				g.Message("CONTRACT FAILED: CARGO MISSING", "bad")
				continue
			case ContractPassenger:
				cmd.Pax -= c.Qty
				cmd.Cash += float64(c.Pay)
				g.Message("PASSENGERS DELIVERED  +"+itoa(c.Pay)+" CR", "good")
				continue
			case ContractMail:
				cmd.Mail -= c.Qty
				cmd.Cash += float64(c.Pay)
				g.Message("MAIL DELIVERED  +"+itoa(c.Pay)+" CR", "good")
				continue
			}
		}
		keep = append(keep, c)
	}
	cmd.Contracts = keep
}

// expireContracts drops jobs whose deadline has passed, with a fine.
func (g *Game) expireContracts() {
	cmd := &g.Cmd
	keep := cmd.Contracts[:0]
	for _, c := range cmd.Contracts {
		if cmd.Day > c.Due {
			fine := c.Pay / 5
			cmd.Cash -= float64(fine)
			if cmd.Pax > 0 && c.Kind == ContractPassenger {
				cmd.Pax -= c.Qty
			}
			if c.Kind == ContractMail {
				cmd.Mail -= c.Qty
			}
			g.Message("CONTRACT EXPIRED  -"+itoa(fine)+" CR", "bad")
			continue
		}
		keep = append(keep, c)
	}
	cmd.Contracts = keep
}

// countContracts advances bounty jobs when a pirate dies.
func (g *Game) countContracts() {
	cmd := &g.Cmd
	keep := cmd.Contracts[:0]
	for _, c := range cmd.Contracts {
		if c.Kind == ContractBounty && c.Galaxy == cmd.Galaxy && c.System == cmd.System {
			c.Done++
			if c.Done >= c.Qty {
				cmd.Cash += float64(c.Pay)
				g.Message("BOUNTY CONTRACT PAID  +"+itoa(c.Pay)+" CR", "good")
				continue
			}
		}
		keep = append(keep, c)
	}
	cmd.Contracts = keep
}

// CarryingContraband reports whether slaves or narcotics are in the hold.
func (g *Game) CarryingContraband() bool {
	return g.Cmd.Cargo[3] > 0 || g.Cmd.Cargo[6] > 0
}

// ContractSummary is shown on the status screen.
func (g *Game) ContractSummary() string {
	if len(g.Cmd.Contracts) == 0 {
		return ""
	}
	c := g.Cmd.Contracts[0]
	sys := galaxy.Get(c.Galaxy).Systems[c.System].Name
	if c.Kind == ContractCourier {
		return "COURIER TO " + sys
	}
	return "BOUNTY IN " + sys
}
