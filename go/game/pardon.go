package game

// A record can be cleaned, but only where someone has the file and the nerve.
// The flat fee is the lawyer; the rest is per point of heat, so shooting two
// vipers costs more than one bounced cheque.

// PardonPrice is what the fixer wants for wiping the sheet.
func (g *Game) PardonPrice() int {
	if g.Cmd.Legal <= 0 {
		return 0
	}
	price := 500 + g.Cmd.Legal*40
	if price > 60000 {
		price = 60000 // past a certain point nobody wants to know you
	}
	return price
}

// PardonAvailable reports whether this station has anyone who can lose a file.
// A high tech world has lawyers; a dictatorship, a communist state or a
// corporate state just has people who can make a call.
func (g *Game) PardonAvailable() bool {
	if g.Cmd.Legal <= 0 || g.World == nil || g.World.Sys == nil || g.Classic() {
		return false
	}
	switch g.World.Sys.Government {
	case 3, 4, 7: // dictatorship, communist, corporate state
		return true
	}
	return g.World.Sys.Tech >= 10
}

// BuyPardon clears the record for money. It does nothing for a pilot the navy
// is actively hunting: those vipers are already in the air.
func (g *Game) BuyPardon() bool {
	if !g.PardonAvailable() {
		// say the rule, not just the refusal: the row is always there now
		g.Message("NOBODY HERE CAN HELP - TRY A HIGH TECH OR CORPORATE STATION", "bad")
		return false
	}
	price := g.PardonPrice()
	if float64(price) > g.Cmd.Cash {
		g.Message("NOT ENOUGH CREDITS", "bad")
		return false
	}
	was := g.Cmd.Legal
	g.Cmd.Cash -= float64(price)
	g.Cmd.Legal = 0
	g.Note(LogPardon, "", was, -float64(price))
	g.Message("RECORD CLEARED", "good")
	return true
}

// LegalText names the pilot's standing with the law, for the sheet and the log.
func LegalText(legal int) string {
	switch {
	case legal <= 0:
		return "CLEAN"
	case legal < 500:
		return "OFFENDER"
	}
	return "FUGITIVE"
}
