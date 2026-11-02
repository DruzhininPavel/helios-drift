package game

import "heliosdrift/galaxy"

// What is left when the story is over.
//
// The campaign ends on one of three pages, and then the pilot is a free man in a
// sector that keeps trading, keeps arresting and keeps paying bounties. This file
// is the part that gives that time a shape: ranks that announce themselves, a
// chart worth filling in, galaxies worth crossing, and a title that changes how
// the market treats him.

// TitleName is what the sector calls a pilot after the campaign: the endings were
// not just an epitaph, they are a standing.
func (g *Game) TitleName() string {
	switch g.Cmd.Ending {
	case "combine":
		return "OWNER OF THE SECTOR"
	case "chaos":
		return "FREE OF LICENCES"
	case "council":
		return "COUNCIL MEMBER"
	}
	return ""
}

// RankUp is called after a kill: a promotion is announced and written down, and
// the rank that used to mean nothing now hands over a gun.
func (g *Game) RankUp(before int) {
	now := galaxy.RankIndexFor(g.Cmd.Kills)
	if now <= before {
		return
	}
	name := galaxy.RankFor(g.Cmd.Kills)
	g.Message("PROMOTED: "+name, "good")
	g.Note(LogKill, "RANK "+name, g.Cmd.Kills, 0)
	// DANGEROUS is where the navy stops asking questions about where a gun came
	// from: a military laser, if the pilot has not already bought better
	if now >= 6 && g.Cmd.Laser < 2 {
		g.Cmd.Laser = 2
		g.Message("MILITARY LASER FITTED - THE NAVY NOTICED", "good")
		g.Note(LogEquip, "MILITARY LASER", 1, 0)
	}
}

// checkSector reports a charted sector: every system in this galaxy seen at least
// once. It is a long job on purpose - it is the answer to "what now".
func (g *Game) checkSector() {
	if g.Cmd.Flags["SECTOR_CHARTED"] {
		return
	}
	gal := galaxy.Get(g.Cmd.Galaxy)
	seen := 0
	for _, sys := range gal.Systems {
		if g.Cmd.Visited[visitedKey(g.Cmd.Galaxy, sys.ID)] {
			seen++
		}
	}
	if seen < len(gal.Systems) {
		return
	}
	if g.Cmd.Flags == nil {
		g.Cmd.Flags = map[string]bool{}
	}
	g.Cmd.Flags["SECTOR_CHARTED"] = true
	g.Message("SECTOR CHARTED - EVERY SYSTEM VISITED", "good")
	g.Note(LogVisit, "SECTOR CHARTED", seen, 0)
}

// visitedKey is how Visited keys are built, in one place.
func visitedKey(gal, sys int) string {
	return itoa(gal) + ":" + itoa(sys)
}

// noteGalaxy counts the galaxies a pilot has been through: crossing all eight is
// the other long road, and it does not need a rank of its own.
func (g *Game) noteGalaxy() {
	if g.Cmd.Flags == nil {
		g.Cmd.Flags = map[string]bool{}
	}
	key := "GAL_" + itoa(g.Cmd.Galaxy)
	if g.Cmd.Flags[key] {
		return
	}
	g.Cmd.Flags[key] = true
	n := 0
	for i := 0; i < galaxy.Galaxies; i++ {
		if g.Cmd.Flags["GAL_"+itoa(i)] {
			n++
		}
	}
	g.Note(LogGalaxy, "GALAXY "+itoa(g.Cmd.Galaxy+1), n, 0)
	if n >= galaxy.Galaxies {
		g.Message("ALL EIGHT GALAXIES CROSSED", "good")
	}
}

// GalaxiesSeen is how many of the eight a pilot has flown in.
func (g *Game) GalaxiesSeen() int {
	n := 0
	for i := 0; i < galaxy.Galaxies; i++ {
		if g.Cmd.Flags["GAL_"+itoa(i)] {
			n++
		}
	}
	return n
}

// sellBonus and buyBonus are the standing a title buys: an owner takes a cut of
// what leaves his sector, a council member buys at the price the council agreed.
func (g *Game) sellBonus() float64 {
	if g.Cmd.Ending == "combine" {
		return 1.05
	}
	return 1.0
}

func (g *Game) buyBonus() float64 {
	if g.Cmd.Ending == "council" {
		return 0.95
	}
	return 1.0
}

// SystemsSeen is how much of this galaxy the pilot has actually looked at.
func (g *Game) SystemsSeen() int {
	gal := galaxy.Get(g.Cmd.Galaxy)
	n := 0
	for _, sys := range gal.Systems {
		if g.Cmd.Visited[visitedKey(g.Cmd.Galaxy, sys.ID)] {
			n++
		}
	}
	return n
}
