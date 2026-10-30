// Missions: the two scripted contracts from 1984 Elite. The Constrictor hunt
// gates on combat rank; the Thargoid documents run requires a galactic
// hyperdrive and ends in a scripted ambush.
package game

import (
	"heliosdrift/galaxy"
	"heliosdrift/lang"
	"heliosdrift/models"
)

const (
	MissionConstrictor = "constrictor"
	MissionDocuments   = "documents"
	MissionThargoid    = "thargoid"

	flagConstrictorDone = 1
	flagDocumentsDone   = 2
	flagThargoidDone    = 4
)

type Mission struct {
	Kind   string `json:"kind"`
	System int    `json:"system"`
	Galaxy int    `json:"galaxy"`
	Pay    int    `json:"pay"`
}

// MissionOffer reports what contract the station is currently offering, if any.
func (g *Game) MissionOffer() string {
	// Advanced runs the campaign instead of the two 1984 contracts
	if !g.Classic() {
		return ""
	}
	if g.Cmd.Mission.Kind != "" {
		return ""
	}
	if g.Cmd.MissionFlags&flagConstrictorDone == 0 {
		if g.Cmd.Kills >= 24 {
			return MissionConstrictor
		}
		return ""
	}
	if g.Cmd.MissionFlags&flagDocumentsDone == 0 {
		return MissionDocuments
	}
	return ""
}

// rollOffers fixes the contract destinations for this visit.
func (g *Game) rollOffers() {
	g.OfferTarget = g.missionTarget()
	g.OfferDocGalaxy = (g.Cmd.Galaxy + 1) % galaxy.Galaxies
	g.OfferDocSystem = randInt(256)
}

// missionTarget picks a reachable system for the Constrictor contract.
func (g *Game) missionTarget() int {
	gal := galaxy.Get(g.Cmd.Galaxy)
	cur := gal.Systems[g.Cmd.System]
	var cand []int
	for _, s := range gal.Systems {
		if d := galaxy.Distance(cur, s); d > 2.5 && d <= 7 {
			cand = append(cand, s.ID)
		}
	}
	if len(cand) == 0 {
		return (g.Cmd.System + 37) % 256
	}
	return cand[randInt(len(cand))]
}

// AcceptMission starts a contract. It fails (with a message) when the player
// lacks the hardware the brief demands.
func (g *Game) AcceptMission(kind string) bool {
	cmd := &g.Cmd
	switch kind {
	case MissionConstrictor:
		cmd.Mission = Mission{Kind: kind, System: g.OfferTarget, Galaxy: cmd.Galaxy, Pay: 5000}
		g.Message("CONTRACT ACCEPTED - HUNT THE SHRIKE", "good")
		return true
	case MissionThargoid:
		cmd.Mission = Mission{Kind: kind, System: g.OfferTarget, Galaxy: cmd.Galaxy, Pay: 25000}
		g.Message("STRIKE ORDER ACCEPTED", "good")
		return true
	case MissionDocuments:
		if !cmd.Equip["galdrive"] {
			g.Message("YOU NEED A GALACTIC HYPERDRIVE FIRST", "bad")
			return false
		}
		next := g.OfferDocGalaxy
		ng := galaxy.Get(next)
		tgt := g.OfferDocSystem
		cmd.Mission = Mission{Kind: kind, System: tgt, Galaxy: next, Pay: 15000}
		g.Message("DOCUMENTS ABOARD - DELIVER TO "+ng.Systems[tgt].Name, "good")
		return true
	}
	return false
}

// MissionBrief returns the lines shown on the briefing screen.
func (g *Game) MissionBrief(kind string) []string {
	cmd := &g.Cmd
	switch kind {
	case MissionConstrictor:
		gal := galaxy.Get(cmd.Galaxy)
		name := "?"
		if cmd.Mission.Kind == MissionConstrictor {
			name = gal.Systems[cmd.Mission.System].Name
		} else {
			name = gal.Systems[g.OfferTarget].Name
		}
		return []string{
			"NAVAL INTELLIGENCE HAS LOST A PROTOTYPE",
			"SHRIKE CLASS FIGHTER. ARMED AND",
			"DANGEROUS. FIND IT AND DESTROY IT.",
			"",
			"LAST SEEN: " + name,
			"FEE: 5000 CR",
		}
	case MissionThargoid:
		name := "?"
		if cmd.Mission.Kind == MissionThargoid {
			name = galaxy.Get(cmd.Mission.Galaxy).Systems[cmd.Mission.System].Name
		} else {
			name = galaxy.Get(cmd.Galaxy).Systems[g.OfferTarget].Name
		}
		return []string{
			"NAVAL COMMAND HAS TRACKED A SWARM",
			"MOTHERSHIP. IT SITS DEEP IN THE SYSTEM",
			"WITH ITS DRONES AROUND IT.",
			"",
			"LAST SEEN: " + name,
			"FEE: 25000 CR",
		}
	case MissionDocuments:
		next := (cmd.Galaxy + 1) % galaxy.Galaxies
		name := "?"
		if cmd.Mission.Kind == MissionDocuments {
			name = galaxy.Get(cmd.Mission.Galaxy).Systems[cmd.Mission.System].Name
		} else {
			name = galaxy.Get(g.OfferDocGalaxy).Systems[g.OfferDocSystem].Name
		}
		return []string{
			"COURIER RUN. SWARM DOCUMENTS MUST",
			"REACH NAVAL COMMAND IN THE NEXT GALAXY.",
			"THEY WILL TRY TO STOP YOU.",
			"",
			"DESTINATION: " + name,
			"GALAXY " + itoa(next+1) + "   FEE: 15000 CR",
		}
	}
	return nil
}

// MissionOfferLabel is the board label for the story contract.
func (g *Game) MissionOfferLabel() string {
	switch g.MissionOffer() {
	case MissionConstrictor:
		return "STORY: HUNT THE SHRIKE"
	case MissionDocuments:
		return "STORY: SWARM DOCUMENTS"
	}
	return ""
}

// missionStatus is the one-line summary shown on the status screen.
func (g *Game) MissionStatus() string {
	m := g.Cmd.Mission
	switch m.Kind {
	case MissionConstrictor:
		return lang.T("SHRIKE IN") + " " + galaxy.Get(m.Galaxy).Systems[m.System].Name
	case MissionDocuments:
		return lang.T("DOCUMENTS FOR") + " " + galaxy.Get(m.Galaxy).Systems[m.System].Name
	case MissionThargoid:
		return lang.T("MOTHERSHIP IN") + " " + galaxy.Get(m.Galaxy).Systems[m.System].Name
	}
	return ""
}

// missionArrive spawns the scripted encounter when jumping into the target
// system, and rolls the Thargoid ambush while documents are aboard.
func (g *Game) missionArrive() {
	m := &g.Cmd.Mission
	if m.Kind == "" {
		return
	}
	if m.Kind == MissionConstrictor && m.Galaxy == g.Cmd.Galaxy && m.System == g.Cmd.System {
		e := g.spawnAtDistance("constrictor", 4200, "pirate")
		e.Shield = 60
		e.Energy = 120
		e.MaxEnergy = 120
		e.Bounty = 0 // the fee is the contract, not a bounty
		e.Mission = true
		g.World.Entities = append(g.World.Entities, e)
		g.Message("SHRIKE DETECTED", "bad")
		g.Emit("warning")
		return
	}
	if m.Kind == MissionDocuments && randFloat() < 0.35 {
		g.spawnThargoids()
	}
	if m.Kind == MissionThargoid && m.Galaxy == g.Cmd.Galaxy && m.System == g.Cmd.System {
		mother := g.spawnAtDistance("thargoid", 5200, "thargoid")
		mother.Mission = true
		mother.Radius = 62
		mother.Energy, mother.MaxEnergy = 520, 520
		mother.Shield, mother.MaxShield = 90, 90
		mother.Bounty = 0
		if mesh := models.Get("mothership"); mesh != nil {
			mother.Mesh = mesh
		}
		g.World.Entities = append(g.World.Entities, mother)
		for i := 0; i < 4; i++ {
			g.World.Entities = append(g.World.Entities, g.spawnAtDistance("thargon", randRange(3000, 5000), "thargoid"))
		}
		g.Message("SWARM MOTHERSHIP DETECTED", "bad")
		g.Emit("warning")
	}
}

// missionKill completes the Constrictor contract.
func (g *Game) missionKill(e *Entity) {
	cmd := &g.Cmd
	if !e.Mission {
		return
	}
	cmd.Cash += float64(cmd.Mission.Pay)
	switch cmd.Mission.Kind {
	case MissionThargoid:
		cmd.MissionFlags |= flagThargoidDone
	case MissionConstrictor:
		cmd.MissionFlags |= flagConstrictorDone
	}
	cmd.Mission = Mission{}
	name := "SHRIKE DESTROYED"
	if cmd.Mission.Kind == MissionThargoid {
		name = "MOTHERSHIP DESTROYED"
	}
	g.Message(name+"  +"+itoa(cmd.Mission.Pay)+" CR", "good")
}

// missionDock completes the courier run.
func (g *Game) missionDock() {
	cmd := &g.Cmd
	m := cmd.Mission
	if m.Kind != MissionDocuments {
		return
	}
	if m.Galaxy != cmd.Galaxy || m.System != cmd.System {
		return
	}
	cmd.Cash += float64(m.Pay)
	cmd.MissionFlags |= flagDocumentsDone
	cmd.Mission = Mission{}
	g.Message("DOCUMENTS DELIVERED  +"+itoa(m.Pay)+" CR", "good")
}

// GalacticJump moves the commander to the next galaxy, consuming the drive the
// way the original does.
func (g *Game) GalacticJump() bool {
	if g.State != "flight" {
		return false
	}
	if !g.Cmd.Equip["galdrive"] {
		g.Message("NO GALACTIC HYPERDRIVE", "bad")
		return false
	}
	delete(g.Cmd.Equip, "galdrive")
	g.Cmd.Galaxy = (g.Cmd.Galaxy + 1) % galaxy.Galaxies
	tgt := randInt(256)
	g.Cmd.TargetSystem = tgt
	g.State = "hyper"
	g.Hyper = &Hyper{Target: tgt}
	g.Message("GALACTIC HYPERSPACE - GALAXY "+itoa(g.Cmd.Galaxy+1), "good")
	g.Emit("hyper")
	return true
}

// StationMeshForShots exposes the hull picker to the screenshot pass.
func StationMeshForShots(i int) string {
	return stationMesh(galaxy.Get(0).Systems[i])
}

// DockedForShots re-rolls the station state for the screenshot pass.
func (g *Game) DockedForShots() {
	g.dockedAt()
}

// MissionTarget exposes the active mission's destination for the chart and HUD.
func (g *Game) MissionTarget() (system, galaxy int, ok bool) {
	if !g.Classic() {
		if g.Cmd.Campaign.Active {
			return g.Cmd.Campaign.Target, g.Cmd.Campaign.Galaxy, true
		}
		return 0, 0, false
	}
	m := g.Cmd.Mission
	if m.Kind == "" {
		return 0, 0, false
	}
	return m.System, m.Galaxy, true
}

// MissionOfferTarget is the system a not-yet-accepted story contract points at.
func (g *Game) MissionOfferTarget() (system, galaxy int, ok bool) {
	if !g.Classic() {
		if g.CampaignOffer() {
			return g.campaignTarget()
		}
		return 0, 0, false
	}
	if g.MissionOffer() == "" {
		return 0, 0, false
	}
	kind := g.MissionOffer()
	if kind == MissionDocuments {
		return g.OfferDocSystem, g.OfferDocGalaxy, true
	}
	return g.OfferTarget, g.Cmd.Galaxy, true
}

// ContractTargets lists destinations of accepted non-story contracts.
func (g *Game) ContractTargets() []int {
	out := make([]int, 0, len(g.Cmd.Contracts))
	for _, c := range g.Cmd.Contracts {
		if c.Galaxy == g.Cmd.Galaxy {
			out = append(out, c.System)
		}
	}
	return out
}

// MissionArriveForShots runs the scripted encounter for the screenshot pass.
func (g *Game) MissionArriveForShots() { g.missionArrive() }
