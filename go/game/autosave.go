package game

import (
	"encoding/json"

	"heliosdrift/galaxy"
)

// The session autosave.
//
// Every time the pilot docks, the commander is remembered in memory. It is not a
// file and it does not survive quitting: it is the answer to "I just died over a
// station I left two minutes ago". The death screen offers it, along with the
// save slots and the way out.

// snapshot is the state of the commander at the last docking.
type snapshot struct {
	Cmd  Commander
	Sys  int
	Gal  int
	Ship string
}

// TakeSnapshot remembers the pilot where he stands now. It runs on docking.
func (g *Game) TakeSnapshot() {
	data, err := json.Marshal(g.Cmd)
	if err != nil {
		return
	}
	var copyCmd Commander
	if err := json.Unmarshal(data, &copyCmd); err != nil {
		return
	}
	g.auto = &snapshot{Cmd: copyCmd, Sys: g.Cmd.System, Gal: g.Cmd.Galaxy, Ship: g.Cmd.Ship}
}

// HasSnapshot reports whether there is somewhere to come back to.
func (g *Game) HasSnapshot() bool { return g.auto != nil }

// SnapshotWhere names the station the pilot would come back to.
func (g *Game) SnapshotWhere() string {
	if g.auto == nil {
		return ""
	}
	gal := galaxy.Get(g.auto.Gal)
	if gal == nil || g.auto.Sys < 0 || g.auto.Sys >= len(gal.Systems) {
		return ""
	}
	return gal.Systems[g.auto.Sys].Name
}

// Respawn puts the pilot back at the last station he docked at, alive, with
// everything he had then. The death is counted: it happened.
func (g *Game) Respawn() bool {
	if g.auto == nil {
		return false
	}
	deaths := g.Cmd.Deaths
	cmd := g.auto.Cmd
	cmd.Deaths = deaths
	g.Cmd = cmd
	g.Dead = nil
	g.Overlay = ""
	g.EnterSystem(g.Cmd.System, "launch", 0, false)
	// EnterSystem puts the pilot in flight; the point of coming back is that he
	// wakes up at the station, not drifting outside it
	g.State = "docked"
	g.Screen = "launch"
	g.dockedAt()
	g.State = "docked"
	return true
}
