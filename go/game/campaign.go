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

// Quest is one campaign stage.
type Quest struct {
	Title string
	Kind  string // deliver | kill | travel
	Good  int
	Qty   int
	Pay   int
	Tag   string // entity tag that has to die
	Spawn string // ship type spawned on arrival
	Count int
	Brief []string
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
		Kind:  QuestDeliver, Good: 0, Qty: 4, Pay: 1200,
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
		Brief: []string{
			"THEIR SHIP WENT DOWN WITH A FLIGHT",
			"RECORDER. FIND THE WRECK AND BURN IT -",
			"WHOEVER READS THAT BOX WILL COME LOOKING.",
		},
	},
	{
		Title: "THE COURIER",
		Kind:  QuestKill, Tag: "COURIER", Spawn: "cobra1", Count: 1, Pay: 5200,
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
		Kind:  QuestKill, Tag: "CONVOY", Spawn: "boa", Count: 3, Pay: 8800,
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
		Kind:  QuestKill, Tag: "NAVY", Spawn: "viper", Count: 2, Pay: 0,
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
			e.Energy, e.MaxEnergy = 520, 520
			e.Shield, e.MaxShield = 90, 90
			if mesh := models.Get("mothership"); mesh != nil {
				e.Mesh = mesh
			}
		case "FLAGSHIP":
			e.Radius = 70
			e.Energy, e.MaxEnergy = 900, 900
			e.Shield, e.MaxShield = 120, 120
			e.Speed = 120
		}
		g.World.Entities = append(g.World.Entities, e)
	}
	if q.Tag == "MOTHER" {
		for i := 0; i < 4; i++ {
			g.World.Entities = append(g.World.Entities, g.spawnAtDistance("thargon", randRange(3000, 5000), "thargoid"))
		}
	}
	g.Message(q.Title+" - CONTACTS", "bad")
	g.Emit("warning")
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
		g.Message("STAGE COMPLETE  +"+itoa(q.Pay)+" CR", "good")
	} else {
		g.Message("STAGE COMPLETE", "good")
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
	g.Cmd.Cash += float64(q.Pay)
	c.Active = false
	c.Stage++
	g.Message("DELIVERED  +"+itoa(q.Pay)+" CR", "good")
}
