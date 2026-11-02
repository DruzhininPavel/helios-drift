// Campaign: an eight-stage story chain for Advanced mode. The shape is the one
// space-trading games have used since Galaxy on Fire - a freelancer, a
// corporation, a betrayal, an artefact - but the text and the beats are written
// for this game and use only mechanics it actually has.
package game

import (
	"heliosdrift/galaxy"
	"heliosdrift/lang"
	"heliosdrift/models"
	"sort"
)

// TalkVariant picks a different opening for a stage when its conditions hold.
type TalkVariant struct {
	If   []Cond
	Node string
}

// Quest is one campaign stage.
type Quest struct {
	Title string
	Kind  string // deliver | kill | travel
	Good  int
	Qty   int
	Pay   int
	// Advance is paid when the job is taken, Papers is the paperwork that
	// comes with it: the board hands over the cargo and the forms together.
	Advance int
	Papers  int
	Tag     string // entity tag that has to die
	// Talk is a dialogue node opened on arrival instead of paying out on the
	// spot: the story happens at the counter, not in the mission screen.
	Talk string
	// Variants are alternative openings in order of priority: the same stage
	// tells a different story depending on what the player did last time.
	Variants []TalkVariant
	Spawn    string // ship type spawned on arrival
	Count    int
	Brief    []string
}

const (
	QuestDeliver = "deliver"
	QuestKill    = "kill"
	QuestTravel  = "travel"
)

// campaignQuests is the whole chain, in order.
var campaignQuests = []Quest{
	{
		Title: "MEDICAL RUN",
		Kind:  QuestDeliver, Good: 0, Qty: 4, Pay: 1200, Talk: "c1_disp",
		Advance: 400, Papers: 1,
		Brief: []string{
			"THE FRONTIER CLINIC IS OUT OF SUPPLIES.",
			"TAKE FOUR TONNES OF FOOD OUT TO THEM.",
			"",
			"IT IS A SIMPLE RUN. THAT IS THE POINT.",
		},
	},
	{
		Title: "AMBUSH",
		Kind:  QuestKill, Tag: "AMBUSH", Spawn: "krait", Count: 3, Pay: 2600,
		Talk: "c2_cmdA",
		Variants: []TalkVariant{
			{If: []Cond{{Key: "REFUSED_SEAL", Op: "flag"}}, Node: "c2_cmdB"},
			{If: []Cond{{Key: "HAS_RECEIPT", Op: "flag"}}, Node: "c2_cmdAR"},
		},
		Brief: []string{
			"THE CLINIC WAS A DECOY. THEY WERE WAITING",
			"FOR WHOEVER CAME WITH THE SUPPLIES.",
			"",
			"THREE OF THEM. CLEAR THE SKY AND LEAVE.",
		},
	},
	{
		Title: "BLACK BOX",
		Kind:  QuestKill, Tag: "HULK", Spawn: "python", Count: 1, Pay: 3400,
		Talk: "c3_roster",
		Brief: []string{
			"THEIR SHIP WENT DOWN WITH A FLIGHT",
			"RECORDER. FIND THE WRECK AND BURN IT -",
			"WHOEVER READS THAT BOX WILL COME LOOKING.",
		},
	},
	{
		Title: "THE COURIER",
		Kind:  QuestKill, Tag: "COURIER", Spawn: "cobra1", Count: 1, Pay: 5200,
		Talk: "c4_broadcast",
		Brief: []string{
			"THE BOX NAMED A NAVY LIAISON SELLING",
			"JUMP COORDINATES. HIS COURIER RUNS THE",
			"DATA OUT EVERY SIX DAYS.",
			"",
			"STOP THE COURIER.",
		},
	},
	{
		Title: "THE CONVOY",
		Kind:  QuestKill, Tag: "CONVOY", Spawn: "boa", Count: 3, Pay: 14000,
		Talk: "c5_captain",
		Brief: []string{
			"WITH THE LIAISON EXPOSED, HELIOS COMBINE",
			"IS MOVING SOMETHING BIG. THREE HAULERS,",
			"NO ESCORT WORTH THE NAME.",
			"",
			"BREAK THE CONVOY.",
		},
	},
	{
		Title: "FRAMED",
		Kind:  QuestKill, Tag: "NAVY", Spawn: "viper", Count: 2, Pay: 12000,
		Talk: "c6_brief",
		Brief: []string{
			"THEY HAVE PUT YOUR NAME ON THE CONVOY",
			"MASSACRE. NAVY VIPERS ARE ON YOUR TAIL",
			"AND YOUR RECORD IS ALREADY FILTHY.",
			"",
			"SURVIVE. TWO OF THEM.",
		},
	},
	{
		Title: "THE HIVE",
		Kind:  QuestKill, Tag: "MOTHER", Spawn: "thargoid", Count: 1, Pay: 25000,
		Talk: "c7_hive",
		Brief: []string{
			"THE ARTEFACT IN THAT CONVOY WAS NOT",
			"CARGO. IT WAS A CALL, AND SOMETHING",
			"ANSWERED IT.",
			"",
			"KILL THE MOTHERSHIP BEFORE IT WAKES.",
		},
	},
	{
		Title: "THE FLAGSHIP",
		Kind:  QuestKill, Tag: "FLAGSHIP", Spawn: "anaconda", Count: 1, Pay: 50000,
		Talk: "c8_director",
		Brief: []string{
			"HELIOS COMBINE IS RUNNING. THEIR FLAGSHIP",
			"IS LOADED AND ALREADY CLEARED TO JUMP.",
			"",
			"END IT.",
		},
	},
}

// CampaignState rides along with the commander.
type CampaignState struct {
	Stage  int  `json:"stage"`
	Active bool `json:"active"`
	Target int  `json:"target"`
	Galaxy int  `json:"galaxy"`
}

// CampaignTotal is the number of stages.
func CampaignTotal() int { return len(campaignQuests) }

// Quest returns the stage definition.
func QuestAt(i int) Quest {
	if i < 0 || i >= len(campaignQuests) {
		return Quest{}
	}
	return campaignQuests[i]
}

// CampaignOffer reports whether a stage can be taken at this station.
func (g *Game) CampaignOffer() bool {
	return !g.Classic() && !g.Cmd.Campaign.Active && g.Cmd.Campaign.Stage < CampaignTotal()
}

// CampaignQuest is the stage currently on offer or in progress.
func (g *Game) CampaignQuest() Quest {
	return QuestAt(g.Cmd.Campaign.Stage)
}

// AcceptCampaign starts the offered stage.
func (g *Game) AcceptCampaign() bool {
	if !g.CampaignOffer() {
		return false
	}
	q := g.CampaignQuest()
	if q.Kind == QuestDeliver {
		if CargoUsed(g.Cmd)+q.Qty > g.Cmd.CargoCap {
			g.Message("NOT ENOUGH CARGO SPACE", "bad")
			return false
		}
		g.Cmd.Cargo[q.Good] += q.Qty
	}
	if q.Advance > 0 {
		g.Cmd.Cash += float64(q.Advance)
	}
	if q.Papers > 0 {
		if g.Cmd.Flags == nil {
			g.Cmd.Flags = map[string]bool{}
		}
		g.Cmd.Papers += q.Papers
	}
	g.Cmd.Campaign.Active = true
	g.Cmd.Campaign.Target = g.missionTarget()
	g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
	g.Message("CONTRACT ACCEPTED: "+q.Title, "good")
	return true
}

// CampaignStatus is the journal line for the current stage.
func (g *Game) CampaignStatus() string {
	c := g.Cmd.Campaign
	if c.Stage >= CampaignTotal() {
		return "CAMPAIGN COMPLETE"
	}
	q := QuestAt(c.Stage)
	if !c.Active {
		return ""
	}
	sys := galaxy.Get(c.Galaxy).Systems[c.Target].Name
	// the stage name and the connective read in the player's language; the
	// system name stays as it is on the chart
	switch q.Kind {
	case QuestDeliver:
		return lang.T(q.Title) + lang.T(" TO ") + sys
	}
	return lang.T(q.Title) + lang.T(" IN ") + sys
}

// campaignTarget returns the system the current stage points at.
func (g *Game) campaignTarget() (int, int, bool) {
	c := g.Cmd.Campaign
	if c.Stage >= CampaignTotal() {
		return 0, 0, false
	}
	if !c.Active {
		return g.OfferTarget, g.Cmd.Galaxy, true
	}
	return c.Target, c.Galaxy, true
}

// campaignArrive spawns the stage opposition when the target system is reached.
func (g *Game) campaignArrive() {
	c := &g.Cmd.Campaign
	if !c.Active || c.Galaxy != g.Cmd.Galaxy || c.Target != g.Cmd.System {
		return
	}
	q := QuestAt(c.Stage)
	if q.Kind != QuestKill {
		return
	}
	// a stage with a conversation speaks first: the shooting is the answer to
	// a question, not the arrival itself
	if q.Talk != "" {
		if g.Cmd.Flags["SPOKEN_"+q.Tag] {
			return
		}
		if g.Cmd.Flags == nil {
			g.Cmd.Flags = map[string]bool{}
		}
		g.Cmd.Flags["SPOKEN_"+q.Tag] = true
		for _, v := range q.Variants {
			if g.conditionsMet(v.If) {
				g.TalkStart(v.Node)
				return
			}
		}
		g.TalkStart(q.Talk)
		return
	}
	g.spawnCampaignGroup(q)
}

// spawnQuestGroup releases the ships of a tagged stage group in the current
// system: used when a conversation ends in a fight.
func (g *Game) spawnQuestGroup(tag string) {
	c := &g.Cmd.Campaign
	q := QuestAt(c.Stage)
	if q.Tag != tag {
		return
	}
	g.spawnCampaignGroup(q)
}

// spawnCampaignGroup puts the stage's opposition into the system.
func (g *Game) spawnCampaignGroup(q Quest) {
	for i := 0; i < q.Count; i++ {
		e := g.spawnAtDistance(q.Spawn, randRange(3200, 5200), "pirate")
		e.CampaignTag = q.Tag
		if q.Tag == "NAVY" {
			e.Behaviour = "police"
			e.Hostile = true
		}
		switch q.Tag {
		case "HULK":
			e.Mesh = models.Get("asteroid")
			e.Radius = 40
			e.Behaviour = "wreck"
			e.Speed = 4
		case "MOTHER":
			e.Radius = 62
			e.Energy, e.MaxEnergy = 700, 700
			e.Shield, e.MaxShield = 110, 110
			if mesh := models.Get("mothership"); mesh != nil {
				e.Mesh = mesh
			}
		case "FLAGSHIP":
			// the last fight is meant to be the one you buy a warship for: a
			// pulse laser needs minutes, a military one seconds
			e.Radius = 74
			e.Energy, e.MaxEnergy = 1400, 1400
			e.Shield, e.MaxShield = 180, 180
			e.Speed = 120
		}
		g.World.Entities = append(g.World.Entities, e)
	}
	if q.Tag == "FLAGSHIP" {
		// two escorts: the flagship alone would just be a long slog
		for i := 0; i < 2; i++ {
			esc := g.spawnAtDistance("viper", randRange(2200, 3200), "pirate")
			esc.Hostile = true
			g.World.Entities = append(g.World.Entities, esc)
		}
	}
	if q.Tag == "MOTHER" {
		for i := 0; i < 4; i++ {
			g.World.Entities = append(g.World.Entities, g.spawnAtDistance("thargon", randRange(3000, 5000), "thargoid"))
		}
	}
	g.Message(q.Title+" - CONTACTS", "bad")
	g.Emit("warning")
}

// advanceCampaign moves to the next stage after a dialogue settles the current
// one. Payment and rewards are the dialogue's business, not the chain's.
func (g *Game) advanceCampaign() {
	c := &g.Cmd.Campaign
	if c.Stage >= CampaignTotal() {
		return
	}
	c.Active = false
	c.Stage++
	if c.Stage >= CampaignTotal() {
		g.Message("HELIOS COMBINE IS FINISHED", "good")
		g.Cmd.Legal = 0
		return
	}
	g.Message("STAGE COMPLETE", "good")
}

// campaignKill advances the chain when the tagged target dies.
func (g *Game) campaignKill(tag string) {
	c := &g.Cmd.Campaign
	if !c.Active || tag == "" {
		return
	}
	q := QuestAt(c.Stage)
	if q.Kind != QuestKill || q.Tag != tag {
		return
	}
	c.Active = false
	c.Stage++
	if q.Pay > 0 {
		g.Cmd.Cash += float64(q.Pay)
		g.Note(LogDone, q.Title, c.Stage, float64(q.Pay))
		g.Message("STAGE COMPLETE  +"+itoa(q.Pay)+" CR", "good")
	} else {
		g.Message("STAGE COMPLETE", "good")
	}
	// an audit that ends in an explosion does not end: it breaks
	if tag == "MOTHER" && g.Cmd.Flags["HIVE_ATTACKED"] {
		if g.Cmd.Flags == nil {
			g.Cmd.Flags = map[string]bool{}
		}
		g.Cmd.Flags["AUDIT_BROKEN"] = true
		g.Message("AUDIT INTERRUPTED", "bad")
		g.Note(LogKill, "AUDIT", 0, 0)
	}
	// killing the courier puts a price on the pilot's hull: the board is not
	// shown, but the sky gets busier
	if tag == "COURIER" && g.Cmd.Flags["COURIER_DEAD"] {
		g.Message("BOUNTY POSTED: 20000 ON YOUR HULL", "bad")
		g.Note(LogKill, "BOUNTY", 20000, 0)
	}
	// a broken convoy floods the local market and leaves crates in the debris
	if tag == "CONVOY" && g.Cmd.Flags["CONVOY_ATTACKED"] {
		g.startPanic()
		if g.Cmd.Boxes < 3 {
			g.Cmd.Boxes += 3
			if g.Cmd.Boxes > 3 {
				g.Cmd.Boxes = 3
			}
			g.Message("3 CRATES RECOVERED FROM THE WRECKAGE", "good")
			g.Note(LogEquip, "CRATES", 3, 0)
		}
	}
	// the wreck paid for a burned recorder: the company does not say thank you,
	// it just takes your name off one list and puts it on another
	if tag == "HULK" && g.Cmd.Flags["BURNED_RECORDER"] {
		g.Cmd.Papers++
		g.Message("EXECUTED CORRECTLY. THAT IS RARE.", "good")
		g.Note(LogDone, "CLEAN WORK", c.Stage, 0)
	}
	if c.Stage >= CampaignTotal() {
		// the last stage has no "well done": it has three ways to end, and with
		// the director dead nobody countersigns the seal but the pilot
		g.Cmd.Cash += 25000
		if g.Cmd.Flags["DIRECTOR_DEAD"] {
			g.TalkStart("c8_after")
			return
		}
		g.StartEnding("combine")
	}
}

// campaignDock advances delivery stages.
func (g *Game) campaignDock() {
	c := &g.Cmd.Campaign
	if !c.Active || c.Galaxy != g.Cmd.Galaxy || c.Target != g.Cmd.System {
		return
	}
	q := QuestAt(c.Stage)
	if q.Kind != QuestDeliver {
		return
	}
	if g.Cmd.Cargo[q.Good] < q.Qty {
		g.Message("CONTRACT FAILED: CARGO MISSING", "bad")
		c.Active = false
		return
	}
	g.Cmd.Cargo[q.Good] -= q.Qty
	if q.Talk != "" {
		// the counter is where the story happens: the payment is settled in
		// the conversation, not on the spot
		g.TalkStart(q.Talk)
		return
	}
	g.Cmd.Cash += float64(q.Pay)
	c.Active = false
	c.Stage++
	g.Message("DELIVERED  +"+itoa(q.Pay)+" CR", "good")
}

// SpawnMissionGroupForTest exposes the stage spawn to the UI tests: the screens
// have to survive whatever the campaign puts in the sky.
func (g *Game) SpawnMissionGroupForTest(tag string) { g.spawnQuestGroup(tag) }

// PanicHere reports whether this system is in the day-long panic that follows a
// broken convoy: prices halve, and the station knows why.
func (g *Game) PanicHere() bool {
	if !g.Cmd.Flags["PANIC_MARKET"] {
		return false
	}
	if g.Cmd.PanicSys != g.Cmd.System || g.Cmd.PanicGalaxy != g.Cmd.Galaxy {
		return false
	}
	return g.Cmd.Day-g.Cmd.PanicDay <= 1
}

// startPanic marks the system for the day: the convoy did not arrive, and
// everyone who bought a licence on credit now wants to be paid before the rumour
// becomes a fact.
func (g *Game) startPanic() {
	if g.Cmd.Flags == nil {
		g.Cmd.Flags = map[string]bool{}
	}
	g.Cmd.Flags["PANIC_MARKET"] = true
	g.Cmd.PanicDay = g.Cmd.Day
	g.Cmd.PanicSys = g.Cmd.System
	g.Cmd.PanicGalaxy = g.Cmd.Galaxy
	g.Note(LogVisit, "PANIC", 0, 0)
	g.Message("MARKET PANIC - PRICES HALVED FOR A DAY", "bad")
}

// dangerNow scales the chance of a pirate wave: a price on the pilot's hull
// draws them out of the rocks, and an interrupted audit empties the lanes.
func (g *Game) dangerNow(base float64) float64 {
	if g.Cmd.Flags["COURIER_DEAD"] {
		base *= 1.6
	}
	if g.Cmd.Flags["AUDIT_BROKEN"] {
		base *= 1.3
	}
	if g.Cmd.Ending == "chaos" {
		// no licences means no licensed cargo worth robbing
		base *= 0.7
	}
	return base
}

// ClosingStations is the list of eighteen stations the combine intends to close.
// It is generated from the galaxy so it is the same list every time the pilot
// reads it, and it never changes: that is the point of a notice.
func (g *Game) ClosingStations() []string {
	gal := galaxy.Get(0)
	idx := make([]int, 0, len(gal.Systems))
	for i := range gal.Systems {
		idx = append(idx, i)
	}
	// order by a hash of the name, not by position: the list should look picked
	// by an office, not by a map
	sort.Slice(idx, func(a, b int) bool {
		return nameHash(gal.Systems[idx[a]].Name) < nameHash(gal.Systems[idx[b]].Name)
	})
	if len(idx) > 18 {
		idx = idx[:18]
	}
	out := make([]string, 0, len(idx))
	for _, i := range idx {
		out = append(out, gal.Systems[i].Name)
	}
	return out
}

// nameHash is a small FNV-1a: enough to shuffle a list of names into an order
// that does not follow the map.
func nameHash(s string) uint32 {
	var h uint32 = 2166136261
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

// updateBroadcast counts the audience of the courier's show: it climbs while he
// is alive and talking, and stops when the chapter does.
func (g *Game) updateBroadcast(dt float64) {
	if g.Cmd.Campaign.Active && g.Cmd.Campaign.Stage == 3 {
		onAir := false
		for _, e := range g.World.Entities {
			if e.CampaignTag == "COURIER" && !e.Dead {
				onAir = true
			}
		}
		if onAir || g.Talk == "c4_broadcast" || g.Talk == "c4_offer" {
			g.Viewers += dt * 46
			return
		}
	}
	if g.Viewers > 0 {
		g.Viewers -= dt * 900
		if g.Viewers < 0 {
			g.Viewers = 0
		}
	}
}

// updateHive runs the waiting game: with the shields left up the mothership
// cannot be touched, so after ten minutes the field falls on its own and the
// chapter continues - a stalemate is not a failure state.
func (g *Game) updateHive(dt float64) {
	if !g.Cmd.Campaign.Active || g.Cmd.Campaign.Stage != 6 {
		return
	}
	if !g.Cmd.Flags["SHIELDS_UP"] || g.Cmd.Flags["SHIELDS_CUT"] {
		return
	}
	g.HiveWait += dt
	if g.HiveWait > 600 {
		if g.Cmd.Flags == nil {
			g.Cmd.Flags = map[string]bool{}
		}
		g.Cmd.Flags["SHIELDS_CUT"] = true
		g.Message("THE MOTHERSHIP POWERS DOWN AND THE FIELD FALLS", "info")
	}
}

// StartEnding closes the campaign on one of three notes: who holds the seal when
// the shooting stops.
func (g *Game) StartEnding(kind string) {
	if g.Cmd.Flags == nil {
		g.Cmd.Flags = map[string]bool{}
	}
	g.Cmd.Ending = kind
	g.Cmd.Campaign.Active = false
	g.Cmd.Campaign.Stage = CampaignTotal()
	g.Message("THE CAMPAIGN IS OVER", "good")
	g.Note(LogDone, "ENDING: "+kind, 0, 0)
}

// LicensedHere reports the five per cent that comes with buying a licence off
// the convoy - only on the station where it was bought.
func (g *Game) LicensedHere() bool {
	if !g.Cmd.Flags["BOUGHT_LICENSE"] {
		return false
	}
	return g.Cmd.LicenseSys == g.Cmd.System && g.Cmd.LicenseGalaxy == g.Cmd.Galaxy
}

// acquaintanceNodes are the four people a pilot keeps running into. Each is met
// once per commander, at a station, and none of them has anything to do with the
// campaign: they are the sector's furniture, and furniture has stories.
var acquaintanceNodes = []struct {
	Node string
	Flag string
}{
	{"a1_bar", "MET_BAR"},
	{"a2_dock", "MET_DOCK"},
	{"a3_nurse", "MET_NURSE"},
	{"a4_old", "MET_OLD"},
}

// meetSomeone fires a station conversation now and then, so landing somewhere is
// not only a market. Classic mode stays silent: a 1984 commander had no small
// talk.
func (g *Game) meetSomeone() {
	if g.Classic() || g.Screen == "talk" || g.Cmd.Ending != "" {
		return
	}
	if g.Cmd.Flags == nil {
		g.Cmd.Flags = map[string]bool{}
	}
	var open []string
	for _, a := range acquaintanceNodes {
		if !g.Cmd.Flags[a.Flag] {
			open = append(open, a.Node)
		}
	}
	if len(open) == 0 {
		return
	}
	// one in three dockings, and never twice in a row
	if randFloat() > 0.34 || g.Cmd.LastMeet == g.Cmd.Day {
		return
	}
	pick := open[randInt(len(open))]
	g.Cmd.LastMeet = g.Cmd.Day
	g.TalkStart(pick)
}
