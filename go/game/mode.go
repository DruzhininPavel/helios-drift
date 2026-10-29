package game

// Two rule sets. Classic is the 1984 game: two missions, no contract board, no
// shipyard, no jettison, no in-system jump. Advanced is everything the modern
// build adds. The mode is chosen when a commander is created and never changes.
// Product identity. The game is not called Elite: that word is a registered
// trademark of Frontier Developments, and this build is an independent work.
const (
	ProductName = "HELIOS"
	ProductSub  = "DRIFT"
	ProductFull = "HELIOS DRIFT"
)

const (
	ModeClassic  = "classic"
	ModeAdvanced = "advanced"
)

// Classic reports whether this commander plays by the 1984 rules.
func (g *Game) Classic() bool { return g.Cmd.Mode == ModeClassic }

// ModeName is the label shown on the status screen.
func (g *Game) ModeName() string {
	if g.Classic() {
		return "CLASSIC"
	}
	return "ADVANCED"
}
