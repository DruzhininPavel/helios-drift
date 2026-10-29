package game

import "testing"

func newGameMode(mode string) *Game {
	g := New()
	g.NewGameMode(mode, "TEST")
	return g
}

// Classic is the 1984 rule set: no contract board, no shipyard, no extras.
func TestClassicBlocksModernMechanics(t *testing.T) {
	g := newGameMode(ModeClassic)
	if !g.Classic() {
		t.Fatal("mode not set")
	}
	g.dockedAt()
	if len(g.Offers) != 0 {
		t.Fatalf("classic offered %d contracts", len(g.Offers))
	}
	g.Offers = []Contract{{Kind: ContractCourier, Good: 0, Qty: 1, System: 3, Galaxy: 0, Pay: 200}}
	if g.AcceptOffer(0) {
		t.Fatal("classic accepted a contract")
	}
	g.Cmd.Cash = 999999
	if g.BuyShip("python") {
		t.Fatal("classic bought a ship")
	}
	if g.Cmd.Ship != "cobra3" {
		t.Fatalf("ship changed to %q", g.Cmd.Ship)
	}
	g.Launch()
	g.Cmd.Cargo[0] = 1
	if g.Jettison() {
		t.Fatal("classic jettisoned cargo")
	}
	g.JumpInSystem()
	if g.State != "flight" {
		t.Fatalf("classic in-system jump changed state to %q", g.State)
	}
}

// Advanced keeps everything on.
func TestAdvancedAllowsModernMechanics(t *testing.T) {
	g := newGameMode(ModeAdvanced)
	if g.Classic() {
		t.Fatal("advanced flagged as classic")
	}
	g.dockedAt()
	g.Offers = []Contract{{Kind: ContractCourier, Good: 0, Qty: 1, System: 3, Galaxy: 0, Pay: 200}}
	if !g.AcceptOffer(0) {
		t.Fatal("advanced refused a contract")
	}
	g.Cmd.Cash = 999999
	g.Cmd.Cargo = map[int]int{}
	if !g.BuyShip("adder") {
		t.Fatal("advanced refused a ship purchase")
	}
}

// The story missions exist in both modes.
func TestClassicHasStoryMissions(t *testing.T) {
	g := newGameMode(ModeClassic)
	g.Cmd.Kills = 40
	g.dockedAt()
	if got := g.MissionOffer(); got != MissionConstrictor {
		t.Fatalf("classic offer = %q, want the constrictor contract", got)
	}
	if !g.AcceptMission(MissionConstrictor) {
		t.Fatal("classic refused the story mission")
	}
}

// The mode travels with the save file and cannot be switched later.
func TestModeSurvivesSave(t *testing.T) {
	SaveDirOverride = t.TempDir()
	defer func() { SaveDirOverride = "" }()
	g := newGameMode(ModeClassic)
	if !g.SaveSlot(0) {
		t.Fatal("save failed")
	}
	g2 := New()
	if !g2.LoadSlot(0) {
		t.Fatal("load failed")
	}
	if !g2.Classic() {
		t.Fatalf("loaded mode = %q, want classic", g2.Cmd.Mode)
	}
	// no API switches it: the only way in is NewGameMode
	g2.NewGameMode(ModeAdvanced, "OTHER")
	if g2.Classic() {
		t.Fatal("mode did not change on a new game")
	}
}

// Chart markers need to know where a job points.
func TestMissionTargets(t *testing.T) {
	g := newGameMode(ModeClassic)
	g.Cmd.Kills = 40
	g.dockedAt()
	if _, _, ok := g.MissionTarget(); ok {
		t.Fatal("no mission is active yet")
	}
	sys, gal, ok := g.MissionOfferTarget()
	if !ok || gal != g.Cmd.Galaxy {
		t.Fatalf("offer target %d/%d ok=%v", sys, gal, ok)
	}
	if !g.AcceptMission(g.MissionOffer()) {
		t.Fatal("accept failed")
	}
	sys, gal, ok = g.MissionTarget()
	if !ok || gal != g.Cmd.Galaxy || sys != g.Cmd.Mission.System {
		t.Fatalf("active target %d/%d, want %d/%d", sys, gal, g.Cmd.Mission.System, g.Cmd.Galaxy)
	}
	// contract board only exists in Advanced
	a := newGameMode(ModeAdvanced)
	a.Offers = []Contract{{Kind: ContractCourier, Good: 0, Qty: 1, System: 7, Galaxy: a.Cmd.Galaxy, Pay: 100}}
	if !a.AcceptOffer(0) {
		t.Fatal("contract accept failed")
	}
	targets := a.ContractTargets()
	if len(targets) != 1 || targets[0] != 7 {
		t.Fatalf("contract targets = %v", targets)
	}
}

// The campaign chain: every stage must be offered, spawned and advanced.
func TestCampaignChain(t *testing.T) {
	g := newGameMode(ModeAdvanced)
	for i := 0; i < CampaignTotal(); i++ {
		if !g.CampaignOffer() {
			t.Fatalf("stage %d is not offered", i)
		}
		q := g.CampaignQuest()
		if !g.AcceptCampaign() {
			t.Fatalf("stage %d (%s) not accepted", i, q.Title)
		}
		if g.Cmd.Campaign.Target < 0 || g.Cmd.Campaign.Target > 255 {
			t.Fatalf("stage %d has no target", i)
		}
		g.EnterSystem(g.Cmd.Campaign.Target, "arrive", 9000, false)

		switch q.Kind {
		case QuestDeliver:
			g.Dock()
			for j := 0; j < 200 && g.State != "docked"; j++ {
				g.Update(1.0 / 60.0)
			}
		case QuestKill:
			killed := 0
			for _, e := range append([]*Entity(nil), g.World.Entities...) {
				if e.CampaignTag == q.Tag {
					g.DestroyEntity(e, "laser")
					killed++
				}
			}
			if killed == 0 {
				t.Fatalf("stage %d spawned no %q contacts", i, q.Tag)
			}
		}
		if g.Cmd.Campaign.Stage != i+1 {
			t.Fatalf("stage %d advanced to %d", i, g.Cmd.Campaign.Stage)
		}
		if g.Cmd.Campaign.Active {
			t.Fatalf("stage %d still active after completion", i)
		}
	}
	if g.CampaignOffer() {
		t.Fatal("campaign offered a ninth stage")
	}
	if g.CampaignStatus() != "CAMPAIGN COMPLETE" {
		t.Fatalf("status = %q", g.CampaignStatus())
	}
}

// Advanced no longer offers the 1984 contracts; classic still does.
func TestModeDecidesStorySource(t *testing.T) {
	g := newGameMode(ModeAdvanced)
	g.Cmd.Kills = 40
	g.dockedAt()
	if g.MissionOffer() != "" {
		t.Fatal("advanced still offers the constrictor contract")
	}
	if !g.CampaignOffer() {
		t.Fatal("advanced has no campaign to offer")
	}

	c := newGameMode(ModeClassic)
	c.Cmd.Kills = 40
	c.dockedAt()
	if c.MissionOffer() != MissionConstrictor {
		t.Fatalf("classic offer = %q", c.MissionOffer())
	}
	if c.CampaignOffer() {
		t.Fatal("classic offered the campaign")
	}
}

// Campaign progress must survive a save and reload.
func TestCampaignSurvivesSave(t *testing.T) {
	SaveDirOverride = t.TempDir()
	defer func() { SaveDirOverride = "" }()

	g := newGameMode(ModeAdvanced)
	g.Cmd.Campaign = CampaignState{Stage: 3, Active: true, Target: 77, Galaxy: 1}
	g.Cmd.MissionFlags = 0
	if !g.SaveSlot(1) {
		t.Fatal("save failed")
	}
	other := New()
	if !other.LoadSlot(1) {
		t.Fatal("load failed")
	}
	if other.Cmd.Campaign.Stage != 3 || !other.Cmd.Campaign.Active ||
		other.Cmd.Campaign.Target != 77 || other.Cmd.Campaign.Galaxy != 1 {
		t.Fatalf("campaign after load = %+v", other.Cmd.Campaign)
	}
}

// Yaw is an Advanced-mode addition: it turns the ship there and is ignored in
// Classic, where the 1984 flight model had only pitch and roll.
func TestYawAdvancedOnly(t *testing.T) {
	turn := func(mode string) float64 {
		g := newGameMode(mode)
		g.Launch()
		before := g.Player.Mat.Col(2)
		g.In.Yaw = 1
		for i := 0; i < 30; i++ {
			g.Update(1.0 / 60.0)
		}
		after := g.Player.Mat.Col(2)
		return before.Sub(after).Len()
	}
	if d := turn(ModeAdvanced); d < 0.01 {
		t.Fatalf("advanced did not yaw (moved %.4f)", d)
	}
	if d := turn(ModeClassic); d > 0.001 {
		t.Fatalf("classic yawed by %.4f", d)
	}
}
