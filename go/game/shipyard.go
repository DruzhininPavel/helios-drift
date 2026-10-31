package game

// Shipyard: the commander can trade the Cobra in for something that suits a
// different trade - the later Elites let you do this, the 1984 original did not.
type ShipClass struct {
	ID       string
	Name     string
	Price    int
	Cargo    int
	Speed    float64
	Energy   float64
	Missiles int
	// Pitch, Roll and Yaw scale the base turn rates: a fighter turns far
	// quicker than a hauler, and that is what makes a hull worth its price in
	// a fight rather than only in a cargo hold.
	Pitch, Roll, Yaw float64
	Shield           float64 // shield strength, per face
	Accel            float64 // units per second, added to speed
	MaxLaser         int     // highest laser class the hull can mount
	Blurb            string
}

var ShipYard = []ShipClass{
	// hull, name, price, cargo, vmax, energy, missiles, pitch, roll, yaw,
	// shield, accel, max laser, blurb
	{"cobra3", "KESTREL III", 120000, 20, 420, 100, 4, 1.00, 1.00, 1.00, 100, 260, 2, "THE ALL ROUNDER"},
	{"adder", "ADDER", 34000, 12, 300, 80, 2, 1.15, 1.20, 1.10, 60, 330, 1, "CHEAP AND CHEERFUL"},
	{"python", "PYTHON", 260000, 40, 250, 160, 4, 0.60, 0.55, 0.60, 140, 180, 1, "SLOW, HUGE HOLD"},
	{"ferdelance", "FER-DE-LANCE", 420000, 14, 500, 130, 4, 1.45, 1.40, 1.30, 120, 300, 2, "FAST AND LETHAL"},
	{"asp", "ASP MK II", 300000, 18, 460, 120, 4, 1.25, 1.20, 1.05, 90, 340, 2, "NAVAL SCOUT"},
}

// ShipFor returns the class with the given id, defaulting to the Cobra.
func ShipFor(id string) ShipClass {
	for _, s := range ShipYard {
		if s.ID == id {
			return s
		}
	}
	return ShipYard[0]
}

// Ship is the class the commander is flying.
func (g *Game) Ship() ShipClass { return ShipFor(g.Cmd.Ship) }

// CargoCapacity is the hold size of the current hull plus any cargo bay.
func (g *Game) CargoCapacity() int {
	n := g.Ship().Cargo
	if g.Cmd.Equip["largebay"] {
		n += 15
	}
	return n
}

// TradeInValue is what the yard pays for the current hull.
func (g *Game) TradeInValue() int {
	return int(float64(g.Ship().Price) * 0.6)
}

// BuyShip trades the current hull for another one.
func (g *Game) BuyShip(id string) bool {
	if g.Classic() {
		g.Message("NOT AVAILABLE IN CLASSIC MODE", "bad")
		return false
	}
	cmd := &g.Cmd
	cur := g.Ship()
	want := ShipFor(id)
	if want.ID == cur.ID {
		g.Message("YOU ARE ALREADY FLYING ONE", "bad")
		return false
	}
	net := want.Price - g.TradeInValue()
	if net > 0 && float64(net) > cmd.Cash {
		g.Message("NOT ENOUGH CREDITS", "bad")
		return false
	}
	if CargoUsed(*cmd) > want.Cargo+boolInt(cmd.Equip["largebay"], 15) {
		g.Message("SELL CARGO FIRST", "bad")
		return false
	}
	cmd.Cash -= float64(net)
	cmd.Ship = want.ID
	cmd.CargoCap = g.CargoCapacity()
	if cmd.Missiles > want.Missiles {
		cmd.Missiles = want.Missiles
	}
	// a big laser in a small hull has to come out
	if cmd.Laser > want.MaxLaser {
		cmd.Laser = want.MaxLaser
		g.Message("LASER DOWNGRADED FOR THIS HULL", "bad")
	}
	g.Player.Energy = want.Energy
	g.Player.ShieldF = want.Shield
	g.Player.ShieldA = want.Shield
	g.Message("NOW FLYING "+want.Name, "good")
	return true
}

func boolInt(b bool, n int) int {
	if b {
		return n
	}
	return 0
}
