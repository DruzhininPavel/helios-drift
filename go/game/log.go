package game

import "strings"

// The commander keeps a log. It is written by the simulation as things happen,
// not by the player, and it survives in the save: a run of eight chapters and
// two hundred trades leaves nothing to remember by hand otherwise.
//
// Entries hold data, not sentences: the kind, what it was, how many and how
// much money. The screen builds the sentence, so the log reads Russian or
// English like everything else without a second copy of every line.
const (
	LogKill   = "kill"   // a ship destroyed
	LogBuy    = "buy"    // freight bought
	LogSell   = "sell"   // freight sold
	LogFuel   = "fuel"   // refuelled
	LogEquip  = "equip"  // equipment fitted
	LogShip   = "ship"   // a hull bought
	LogJob    = "job"    // a contract taken
	LogDone   = "done"   // a contract paid or a stage finished
	LogFail   = "fail"   // a contract lapsed
	LogJump   = "jump"   // a hyperspace jump
	LogMine   = "mine"   // ore scooped or a rock broken
	LogDeath  = "death"  // the pilot lost a ship
	LogVisit  = "visit"  // first time in a system
	LogPirate = "pirate" // hostiles dropped out of hyperspace on you
	LogPolice = "police" // a viper came to look at your record
	LogSwarm  = "swarm"  // the swarm, which is a different thing entirely
	LogPardon = "pardon" // a record quietly cleaned
)

// LogLimit is how many entries are kept. Older ones fall off the front: a
// hundred and fifty pages of trade would bury the chapter you are trying to
// remember.
const LogLimit = 300

// LogEntry is one line of the record.
type LogEntry struct {
	Day  int     `json:"day"`
	Sys  string  `json:"sys,omitempty"` // where it happened
	Kind string  `json:"kind"`          // one of the Log* constants
	A    string  `json:"a,omitempty"`   // ship type, commodity, equipment, system
	N    int     `json:"n,omitempty"`   // tonnes, units, stage number
	CR   float64 `json:"cr,omitempty"`  // money moved, signed
}

// Note appends an entry, trimmed to LogLimit.
func (g *Game) Note(kind, a string, n int, cr float64) {
	e := LogEntry{Day: g.Cmd.Day, Kind: kind, A: a, N: n, CR: cr}
	if g.World != nil && g.World.Sys != nil {
		e.Sys = g.World.Sys.Name
	}
	// ten tonnes bought in ten presses is one line in the record, not ten: the
	// log is a record of what happened, not of which keys were pressed
	if kind == LogBuy || kind == LogSell || kind == LogFuel || kind == LogMine {
		for i := len(g.Cmd.Log) - 1; i >= 0 && i >= len(g.Cmd.Log)-3; i-- {
			last := &g.Cmd.Log[i]
			if last.Kind == e.Kind && last.A == e.A && last.Sys == e.Sys && last.Day == e.Day {
				last.N += e.N
				last.CR += e.CR
				return
			}
		}
	}
	g.Cmd.Log = append(g.Cmd.Log, e)
	if len(g.Cmd.Log) > LogLimit {
		g.Cmd.Log = g.Cmd.Log[len(g.Cmd.Log)-LogLimit:]
	}
}

// LogKind labels a category, used by the screen's filter row.
func LogKindLabel(kind string) string {
	switch kind {
	case LogKill:
		return "KILLED"
	case LogBuy:
		return "BOUGHT"
	case LogSell:
		return "SOLD"
	case LogFuel:
		return "REFUELLED"
	case LogEquip:
		return "FITTED"
	case LogShip:
		return "NEW SHIP"
	case LogJob:
		return "CONTRACT"
	case LogDone:
		return "PAID"
	case LogFail:
		return "FAILED"
	case LogJump:
		return "JUMPED"
	case LogMine:
		return "MINED"
	case LogDeath:
		return "LOST SHIP"
	case LogVisit:
		return "FIRST VISIT"
	case LogPirate:
		return "PIRATES"
	case LogPolice:
		return "POLICE"
	case LogSwarm:
		return "SWARM"
	case LogPardon:
		return "PARDON"
	}
	return strings.ToUpper(kind)
}

// LogGroup buckets a kind for the filter row: trade, fighting, jobs, travel.
func LogGroup(kind string) string {
	switch kind {
	case LogBuy, LogSell, LogFuel, LogEquip, LogShip:
		return "TRADE"
	case LogKill, LogDeath, LogPirate, LogPolice, LogSwarm:
		return "FIGHT"
	case LogJob, LogDone, LogFail, LogPardon:
		return "JOBS"
	case LogJump, LogMine, LogVisit:
		return "TRAVEL"
	}
	return "TRADE"
}
