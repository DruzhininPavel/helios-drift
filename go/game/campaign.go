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
	// BriefVariants rewrite the briefing when something earlier changed what
	// the job actually is: a convoy you let pass is not a convoy you broke.
	BriefVariants []BriefVariant
	// TitleVariants rename it for the same reason.
	TitleVariants []TitleVariant
	// Requirements: a stage may ask for more than the previous one being over.
	// NeedEquip is equipment that has to be bolted on, NeedLaser is the gun
	// class the job expects you to bring, NeedRank is the standing.
	NeedEquip []string
	NeedLaser int
	NeedRank  int
	// EnsureSpawn puts the stage's ships in the system on arrival, even when the
	// stage opens with a conversation: the talk explains, it does not spawn.
	EnsureSpawn bool
}

// BriefVariant is an alternative briefing, first match wins.
type BriefVariant struct {
	If    []Cond
	Lines []string
}

// TitleVariant renames the stage when the story took a different turn.
type TitleVariant struct {
	If    []Cond
	Title string
}

// MissingFor says what the pilot lacks for a stage: the shop says why before it
// charges, and the board says why before it offers.
func (g *Game) MissingFor(q Quest) []string {
	var out []string
	for _, id := range q.NeedEquip {
		if !g.Cmd.Equip[id] {
			name := id
			if e := galaxy.EquipByID(id); e != nil {
				name = e.Name
			}
			out = append(out, name)
		}
	}
	if q.NeedLaser > g.Cmd.Laser {
		out = append(out, LaserName[q.NeedLaser])
	}
	if q.NeedRank > 0 && galaxy.RankIndexFor(g.Cmd.Kills) < q.NeedRank {
		out = append(out, galaxy.Ranks[q.NeedRank])
	}
	return out
}

// TitleFor is the name of the job as this pilot earned it.
func (g *Game) TitleFor(q Quest) string {
	for _, v := range q.TitleVariants {
		if g.conditionsMet(v.If) {
			return v.Title
		}
	}
	return q.Title
}

// BriefFor is the briefing as this pilot has earned it.
func (g *Game) BriefFor(q Quest) []string {
	for _, v := range q.BriefVariants {
		if g.conditionsMet(v.If) {
			return v.Lines
		}
	}
	return q.Brief
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
		// two navy vipers is not a job for a stock pulse laser
		NeedLaser: 1,
		Brief: []string{
			"THEY HAVE PUT YOUR NAME ON THE CONVOY",
			"MASSACRE. NAVY VIPERS ARE ON YOUR TAIL",
			"AND YOUR RECORD IS ALREADY FILTHY.",
			"",
			"SURVIVE. TWO OF THEM.",
		},
		TitleVariants: []TitleVariant{{
			If:    []Cond{{Key: "CONVOY_SPARED", Op: "flag"}},
			Title: "THE PAPERS",
		}},
		BriefVariants: []BriefVariant{{
			If: []Cond{{Key: "CONVOY_SPARED", Op: "flag"}},
			Lines: []string{
				"THE CONVOY DOCKED AND TOLD THEM WHAT",
				"YOU DID. NOBODY IS PUTTING THE DEAD ON",
				"YOUR NAME: THE NAVY WANTS THE PAPERS READ,",
				"THE COMBINE WANTS THEM BURNED.",
				"",
				"TWO VIPERS COME TO TAKE THEM OFF YOU.",
			},
		}},
	},
	{
		Title: "THE HIVE",
		Kind:  QuestKill, Tag: "MOTHER", Spawn: "thargoid", Count: 1, Pay: 25000,
		Talk:        "c7_hive",
		EnsureSpawn: true,
		// something that big does not care about a pulse laser, and a pilot
		// without boosted shields has no business in the same system
		NeedLaser: 1, NeedEquip: []string{"shieldboost1"},
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
		TitleVariants: []TitleVariant{
			{If: []Cond{{Key: "AUDIT_WALKED", Op: "flag"}}, Title: "THE LEDGER"},
		},
		BriefVariants: []BriefVariant{
			{If: []Cond{{Key: "AUDIT_WALKED", Op: "flag"}}, Lines: []string{
				"THE HIVE IS STILL OUT THERE AND THE",
				"CLOSURES ARE ALREADY FILED. YOU DID NOT",
				"SIGN THEM AND YOU CANNOT UNSIGN THEM.",
				"",
				"THE COMBINE'S FLAGSHIP IS LOADED AND",
				"CLEARED TO JUMP. IT LEAVES WITH THE",
				"PAPERWORK THAT BURNS SEVENTEEN STATIONS.",
				"",
				"STOP IT, OR LET IT GO.",
			}},
		},
		// the last job in the chain: a gun worth the name and shielding that
		// will survive being looked at by an anaconda
		NeedLaser: 2, NeedEquip: []string{"shieldboost1", "shieldboost2"},
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
	// Spawned says the stage's ships have been put into the world once: a
	// system rebuild (docking and launching again) loses them, and a stage that
	// only spawns from a conversation could then never be finished.
	Spawned bool `json:"spawned,omitempty"`
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
	if miss := g.MissingFor(q); len(miss) > 0 {
		// the board said what it wants; say it again where the pilot pressed
		g.Message(lang.T("NOT READY")+": "+lang.T(miss[0]), "bad")
		return false
	}
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
	g.Message("CONTRACT ACCEPTED: "+g.TitleFor(q), "good")
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
		return lang.T(g.TitleFor(q)) + lang.T(" TO ") + sys
	}
	return lang.T(g.TitleFor(q)) + lang.T(" IN ") + sys
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

// clearCampaignTag sends a tagged group away: used when a conversation ends a
// stage without a fight.
func (g *Game) clearCampaignTag(tag string) {
	if g.World == nil {
		return
	}
	for i := len(g.World.Entities) - 1; i >= 0; i-- {
		e := g.World.Entities[i]
		if e.CampaignTag != tag {
			continue
		}
		g.World.Entities = append(g.World.Entities[:i], g.World.Entities[i+1:]...)
	}
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
	// one group per stage: the arrival and the conversation both want to spawn
	// the hive, and two motherships is not a boss fight
	if g.World != nil {
		for _, e := range g.World.Entities {
			if e.CampaignTag == q.Tag && !e.Dead {
				return
			}
		}
	}
	g.Cmd.Campaign.Spawned = true
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
			e.Energy, e.MaxEnergy = 480, 480
			e.Shield, e.MaxShield = 80, 80
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
		// two drones, not four: the mothership is the fight, the swarm was a
		// wall in front of it
		for i := 0; i < 2; i++ {
			g.World.Entities = append(g.World.Entities, g.spawnAtDistance("thargon", randRange(3000, 5000), "thargoid"))
		}
	}
	g.Message(g.TitleFor(q)+" - CONTACTS", "bad")
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
	c.Spawned = false
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
	c.Spawned = false
	c.Stage++
	if q.Pay > 0 {
		g.Cmd.Cash += float64(q.Pay)
		g.Note(LogDone, g.TitleFor(q), c.Stage, float64(q.Pay))
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
	c.Spawned = false
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
	g.noteEndingFlag(kind)
	g.Cmd.Campaign.Active = false
	g.Cmd.Campaign.Stage = CampaignTotal()
	g.Message("THE CAMPAIGN IS OVER", "good")
	g.Note(LogDone, "ENDING: "+kind, 0, 0)
	// the epilogue is a page, not a toast: show it where the pilot can read it
	if g.State == "docked" {
		g.Screen = "ending"
	} else {
		g.pendingEnding = true
	}
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

// NextStop names the system the pilot's open work points at: the campaign
// stage's target first, then the first accepted contract, then side work. It is
// what the tasks page puts at the bottom so the page answers "and now?".
func (g *Game) NextStop() string {
	gal := galaxy.Get(g.Cmd.Galaxy)
	if gal == nil {
		return ""
	}
	name := func(sys int) string {
		if sys < 0 || sys >= len(gal.Systems) {
			return ""
		}
		return gal.Systems[((sys%256)+256)%256].Name
	}
	if g.Cmd.Campaign.Active && g.Cmd.Campaign.Target >= 0 {
		if n := name(g.Cmd.Campaign.Target); n != "" {
			return n
		}
	}
	for _, c := range g.Cmd.Contracts {
		if n := name(c.System); n != "" && c.Galaxy == g.Cmd.Galaxy {
			return n
		}
	}
	if g.MissionStatus() != "" {
		if n := name(g.Cmd.Mission.System); n != "" {
			return n
		}
	}
	return ""
}
