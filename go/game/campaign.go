// Campaign: an eight-stage story chain for Advanced mode. The shape is the one
// space-trading games have used since Galaxy on Fire - a freelancer, a
// corporation, a betrayal, an artefact - but the text and the beats are written
// for this game and use only mechanics it actually has.
package game

import (
	"heliosdrift/galaxy"
	"heliosdrift/lang"
	"heliosdrift/models"
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
		Kind:  QuestKill, Tag: "COURIER", Spawn: "cobra1", Count: 1, Pay: 9000,
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
	// the wreck paid for a burned recorder: the company does not say thank you,
	// it just takes your name off one list and puts it on another
	if tag == "HULK" && g.Cmd.Flags["BURNED_RECORDER"] {
		g.Cmd.Papers++
		g.Message("EXECUTED CORRECTLY. THAT IS RARE.", "good")
		g.Note(LogDone, "CLEAN WORK", c.Stage, 0)
	}
	if c.Stage >= CampaignTotal() {
		g.Message("HELIOS COMBINE IS FINISHED", "good")
		g.Cmd.Cash += 25000
		g.Cmd.Legal = 0
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
