package game

import "heliosdrift/lang"

// Arrest is what the hail is for.
//
// Cutting speed for a police stop and being found with slaves or narcotics used
// to mean a message saying "you are under arrest" and a fight you could simply
// fly away from, docking at the station with the cargo still aboard. Now the
// stop has a consequence: the hold is emptied, the fine comes out of the purse,
// the record is wiped clean for it, and the pilot wakes up on the station he was
// flying to anyway. Running from the hail - or shooting at the vipers - is still
// the other way out, and it costs the record instead of the cargo.
func (g *Game) Arrest() bool {
	if g.World == nil {
		return false
	}
	fine := g.PardonPrice()
	taken := fine
	if g.Cmd.Cash < float64(taken) {
		taken = int(g.Cmd.Cash)
	}
	g.Cmd.Cash -= float64(taken)
	g.Cmd.Legal = 0

	// the cargo they came for does not leave with the pilot
	seized := g.Cmd.Cargo[3] + g.Cmd.Cargo[6]
	g.Cmd.Cargo[3] = 0
	g.Cmd.Cargo[6] = 0

	g.PoliceState = ""
	for _, e := range g.World.Entities {
		if e.Behaviour == "police" {
			e.Hostile = false
			e.Behaviour = "flee"
		}
	}
	msg := lang.T("ARRESTED") + " - " + lang.T("CARGO SEIZED") + ", " +
		lang.T("FINE") + " " + itoa(taken) + " CR"
	if seized > 0 {
		msg = lang.T("ARRESTED") + " - " + itoa(seized) + " " +
			lang.T("T SEIZED") + ", " + lang.T("FINE") + " " + itoa(taken) + " CR"
	}
	g.Message(msg, "bad")
	sys := g.World.Sys.Name
	g.Note(LogDone, lang.T("ARRESTED AT")+" "+sys, 0, float64(-taken))

	// and the pilot is set down at the station of this system
	if g.World.Station != nil {
		g.EnterSystem(g.Cmd.System, "arrive", 9000, false)
		g.State = "docked"
		g.Screen = "launch"
		g.dockedAt()
		g.State = "docked"
		g.Message(lang.T("RELEASED AT")+" "+sys, "info")
	}
	return true
}
