package game

// Work that opens after the campaign, and work that opens because the pilot went
// somewhere and somebody noticed him. A side job is a short chain of steps, each
// one either a visit or a fight, and the last one pays.
//
// The design note is docs/after-story.md; this file is the part that runs.

import (
	"heliosdrift/lang"
	"heliosdrift/models"
)

// SideStep is one thing to do. Visits are matched against the station the pilot
// docks at, fights against the tag of the ship he destroyed.
type SideStep struct {
	Text string // what the summary shows
	Kind string // "visit" or "kill"
	// a visit step can require a station of a certain kind:
	MinTech int // 0 means any
	MaxGov  int // 0 means any, otherwise the worst government allowed
	Tag     string
	Qty     int    // visits or kills needed
	Talk    string // conversation opened when the step completes
	Pay     int    // paid when this step completes
	Flag    string // set when this step completes
}

// SideJob is a chain with a condition for appearing.
type SideJob struct {
	ID     string
	Title  string
	Unlock []Cond
	Steps  []SideStep
}

// sideJobs is the list. Everything here is checked against things the game
// already tracks: where the pilot docked, what he killed, what the campaign left
// behind.
var sideJobs = []SideJob{
	{
		// The clinic from chapter one is still short of supplies, and the
		// dispatcher knows why.
		ID:     "clinic",
		Title:  "CLINIC DEBT",
		Unlock: []Cond{{Key: "ENDING_SHOWN", Op: "flag"}},
		Steps: []SideStep{
			{Text: "DOCK AT A HIGH TECH STATION AND ASK", Kind: "visit", MinTech: 9, Qty: 1, Talk: "sj_clinic_a"},
			{Text: "CLEAR THE TWO WHO COLLECT THE CUT", Kind: "kill", Tag: "SIDE_CLINIC", Qty: 2},
			{Text: "RETURN WITH THE ANSWER", Kind: "visit", MinTech: 9, Qty: 1, Talk: "sj_clinic_done", Pay: 3000, Flag: "CLINIC_PAID"},
		},
	},
	{
		// The council's first session needs its members in one room, and the
		// pilot is the one with a ship.
		ID:     "keys",
		Title:  "EIGHTEEN KEYS",
		Unlock: []Cond{{Key: "END_COUNCIL", Op: "flag"}},
		Steps: []SideStep{
			{Text: "FIND THREE MEMBERS, ONE PER STATION", Kind: "visit", Qty: 3, Talk: "sj_keys_a"},
			{Text: "HOLD OFF THE ONE WHO OBJECTS", Kind: "kill", Tag: "SIDE_COUNCIL", Qty: 2},
			{Text: "BRING THEM TOGETHER", Kind: "visit", MinTech: 8, Qty: 1, Talk: "sj_keys_done", Pay: 5000, Flag: "COUNCIL_SESSION"},
		},
	},
	{
		// The courier's list, if he is alive and the pilot kept it: eighteen
		// stations, three copies to carry.
		ID:     "list",
		Title:  "THE COURIER'S LIST",
		Unlock: []Cond{{Key: "COURIER_ALIVE", Op: "flag"}, {Key: "HAS_LIST", Op: "flag"}},
		Steps: []SideStep{
			{Text: "TAKE A COPY TO A FRONTIER STATION", Kind: "visit", MaxGov: 1, Qty: 1, Talk: "sj_list_a"},
			{Text: "TAKE THE SECOND TO A QUIET ONE", Kind: "visit", MinTech: 10, Qty: 1},
			{Text: "DECIDE WHAT THE THIRD ONE IS WORTH", Kind: "visit", Qty: 1, Talk: "sj_list_done", Pay: 4000, Flag: "LIST_PLACED"},
		},
	},
}

// SideJobByID finds a definition.
func SideJobByID(id string) *SideJob {
	for i := range sideJobs {
		if sideJobs[i].ID == id {
			return &sideJobs[i]
		}
	}
	return nil
}

// SideJobs is the whole list, for tests and for the summary screen.
func SideJobs() []SideJob { return sideJobs }

// SideProgress is the step a job is on, and how far into it the pilot is.
func (g *Game) SideProgress(id string) (step, count int, active bool) {
	if g.Cmd.SideStep == nil {
		return 0, 0, false
	}
	step, active = g.Cmd.SideStep[id]
	if !active {
		return 0, 0, false
	}
	return step, g.Cmd.SideCount[id], true
}

// SideActive jobs, in the order they were opened.
func (g *Game) SideActive() []*SideJob {
	var out []*SideJob
	for i := range sideJobs {
		if _, _, active := g.SideProgress(sideJobs[i].ID); active {
			out = append(out, &sideJobs[i])
		}
	}
	return out
}

// sideUnlocked reports whether a job's conditions hold now.
func (g *Game) sideUnlocked(j *SideJob) bool {
	return g.conditionsMet(j.Unlock)
}

// sideReady is the first job that should open now and has not been opened yet.
func (g *Game) sideReady() *SideJob {
	for i := range sideJobs {
		j := &sideJobs[i]
		if g.Cmd.Flags["SIDE_DONE_"+j.ID] || g.Cmd.Flags["SIDE_START_"+j.ID] {
			continue
		}
		if !g.sideUnlocked(j) {
			continue
		}
		return j
	}
	return nil
}

// startStep opens the step the job is on: the talk, and the ships.
func (g *Game) startStep(j *SideJob, step int) {
	if step >= len(j.Steps) {
		return
	}
	s := j.Steps[step]
	if s.Talk != "" && g.Screen != "talk" {
		g.TalkStart(s.Talk)
	}
	if s.Kind == "kill" && s.Tag != "" {
		// the fight the step needs is in the sky the moment it is accepted
		for i := 0; i < s.Qty; i++ {
			e := g.spawnAtDistance("krait", randRange(2600, 4200), "pirate")
			e.CampaignTag = s.Tag
			e.Hostile = true
			e.Behaviour = "pirate"
			e.Mesh = models.Get("krait")
			// spawnAtDistance builds the ship, the caller puts it in the sky
			g.World.Entities = append(g.World.Entities, e)
		}
	}
}

// SideStartIfReady opens the next available job. It runs when the pilot docks,
// which is also where it says so.
func (g *Game) SideStartIfReady() {
	if g.Cmd.SideStep == nil {
		g.Cmd.SideStep = map[string]int{}
		g.Cmd.SideCount = map[string]int{}
	}
	j := g.sideReady()
	if j == nil {
		return
	}
	// everything that is ready opens on the same dock: a pilot who has just become
	// somebody should not have to leave and come back for the second job
	for {
		for j != nil {
			g.openSideJob(j)
			g.SideOnDockOnce(j)
			j = g.sideReady()
		}
		return
	}
}

// openSideJob starts one job at its first step.
func (g *Game) openSideJob(j *SideJob) {
	g.Cmd.Flags["SIDE_START_"+j.ID] = true
	g.Cmd.SideStep[j.ID] = 0
	g.Cmd.SideCount[j.ID] = 0
	g.Note(LogJob, lang.T(j.Title), 1, 0)
	g.Message("NEW WORK: "+lang.T(j.Title), "good")
	// the dock that opened the job counts as the visit that started it: the
	// pilot is already standing there, and asking him to leave and come back
	// would be a joke at his expense
	g.startStep(j, 0)
}

// sideCompleteStep advances a job and pays whatever the step pays.
func (g *Game) sideCompleteStep(j *SideJob, step int) {
	s := j.Steps[step]
	if s.Pay > 0 {
		g.Cmd.Cash += float64(s.Pay)
		g.Message("PAID "+itoa(s.Pay)+" CR", "good")
		g.Note(LogDone, lang.T(j.Title), s.Pay, float64(s.Pay))
	}
	if s.Flag != "" {
		if g.Cmd.Flags == nil {
			g.Cmd.Flags = map[string]bool{}
		}
		g.Cmd.Flags[s.Flag] = true
	}
	next := step + 1
	if next >= len(j.Steps) {
		g.Cmd.Flags["SIDE_DONE_"+j.ID] = true
		delete(g.Cmd.SideStep, j.ID)
		delete(g.Cmd.SideCount, j.ID)
		g.Message(lang.T(j.Title)+" - DONE", "good")
		g.Note(LogDone, lang.T(j.Title), 0, 0)
		return
	}
	g.Cmd.SideStep[j.ID] = next
	g.Cmd.SideCount[j.ID] = 0
	g.startStep(j, next)
}

// SideOnDock gives every running job a chance to advance on this station.
func (g *Game) SideOnDock() {
	if g.Cmd.SideStep == nil {
		return
	}
	sys := g.World.Sys
	for i := range sideJobs {
		j := &sideJobs[i]
		step, count, active := g.SideProgress(j.ID)
		if !active || step >= len(j.Steps) {
			continue
		}
		s := j.Steps[step]
		if s.Kind != "visit" {
			continue
		}
		if s.MinTech > 0 && sys.Tech < s.MinTech {
			continue
		}
		if s.MaxGov > 0 && sys.Government > s.MaxGov {
			continue
		}
		count++
		need := s.Qty
		if need < 1 {
			need = 1
		}
		if count < need {
			g.Cmd.SideCount[j.ID] = count
			continue
		}
		// the visit that finished the step might also be the one that needs a
		// conversation: give the sum out after the talk, not before it
		before := g.Screen
		g.sideCompleteStep(j, step)
		if g.Screen == "talk" && before != "talk" {
			continue
		}
	}
}

// SideOnKill reports a destroyed ship to every running job.
func (g *Game) SideOnKill(tag string) {
	if tag == "" || g.Cmd.SideStep == nil {
		return
	}
	for i := range sideJobs {
		j := &sideJobs[i]
		step, count, active := g.SideProgress(j.ID)
		if !active || step >= len(j.Steps) {
			continue
		}
		s := j.Steps[step]
		if s.Kind != "kill" || s.Tag != tag {
			continue
		}
		count++
		need := s.Qty
		if need < 1 {
			need = 1
		}
		if count < need {
			g.Cmd.SideCount[j.ID] = count
			g.Message(lang.T(j.Title)+": "+itoa(count)+"/"+itoa(need), "info")
			continue
		}
		g.sideCompleteStep(j, step)
	}
}

// SideLine is what the summary shows for a running job.
func (g *Game) SideLine(j *SideJob) string {
	step, count, active := g.SideProgress(j.ID)
	if !active || step >= len(j.Steps) {
		return ""
	}
	s := j.Steps[step]
	line := lang.T(s.Text)
	if s.Qty > 1 {
		line += " " + itoa(count) + "/" + itoa(s.Qty)
	}
	return line
}

// noteEndingFlag turns the ending into a flag the side jobs can unlock on.
func (g *Game) noteEndingFlag(kind string) {
	if g.Cmd.Flags == nil {
		g.Cmd.Flags = map[string]bool{}
	}
	switch kind {
	case "council":
		g.Cmd.Flags["END_COUNCIL"] = true
	case "combine":
		g.Cmd.Flags["END_COMBINE"] = true
	case "chaos":
		g.Cmd.Flags["END_CHAOS"] = true
	}
}

// SideOnDockOnce lets a single job see this station, for the dock that opened it.
func (g *Game) SideOnDockOnce(j *SideJob) {
	step, count, active := g.SideProgress(j.ID)
	if !active || step >= len(j.Steps) {
		return
	}
	s := j.Steps[step]
	if s.Kind != "visit" {
		return
	}
	sys := g.World.Sys
	if s.MinTech > 0 && sys.Tech < s.MinTech {
		return
	}
	if s.MaxGov > 0 && sys.Government > s.MaxGov {
		return
	}
	need := s.Qty
	if need < 1 {
		need = 1
	}
	count++
	if count < need {
		g.Cmd.SideCount[j.ID] = count
		return
	}
	g.sideCompleteStep(j, step)
}
