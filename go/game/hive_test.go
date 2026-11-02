package game

import "testing"

// "LEAVE" in the hive conversation means leave: no fight, nothing spawned, and
// the conversation opens again next time. It used to go to the same node as
// "KILL IT", so the pilot who walked away was dragged into a fight anyway.
func TestLeavingTheHiveIsLeaving(t *testing.T) {
	g := New()
	g.NewGame("HIVE")
	g.EnterSystem(3, "arrive", 9000, false)
	g.State = "docked"
	g.Cmd.Campaign.Active = true
	g.Cmd.Campaign.Stage = 6
	g.Cmd.Flags["SPOKEN_MOTHER"] = true
	g.World.Entities = nil

	if !g.TalkStart("c7_hive") {
		t.Fatal("the hive does not talk")
	}
	vis := g.VisibleOptions(DialogueNode("c7_hive"))
	pick := -1
	for i, idx := range vis {
		if DialogueNode("c7_hive").Options[idx].Next == "c7_leave" {
			pick = i
		}
	}
	if pick < 0 {
		t.Fatal("no option leaves")
	}
	if g.TalkChoose(pick) {
		g.EndTalk()
	}
	if n := len(g.World.Entities); n != 0 {
		t.Errorf("%d ships spawned after leaving, want none", n)
	}
	if g.Cmd.Flags["HIVE_ATTACKED"] {
		t.Error("leaving counted as an attack")
	}
	if g.Cmd.Flags["SPOKEN_MOTHER"] {
		t.Error("the stage marked itself spoken: the talk would never open again")
	}
	if !g.Cmd.Campaign.Active || g.Cmd.Campaign.Stage != 6 {
		t.Errorf("campaign is at %d active=%v after leaving", g.Cmd.Campaign.Stage, g.Cmd.Campaign.Active)
	}
}

// Cutting the station's shields is a choice with a price, and the mothership is
// only killable once the field is down.
func TestTheFieldDecidesWhetherTheHiveCanBeHit(t *testing.T) {
	g := New()
	g.NewGame("FIELD")
	g.EnterSystem(3, "arrive", 9000, false)
	g.State = "flight"
	mother := g.spawnAtDistance("thargoid", 900, "thargoid")
	mother.CampaignTag = "MOTHER"
	mother.Energy = 50
	g.World.Entities = append(g.World.Entities, mother)

	g.Cmd.Flags["SHIELDS_UP"] = true
	g.HitEntity(mother, 100, "laser")
	if !mother.Dead && mother.Energy < 50 {
		t.Error("the field let damage through")
	}
	g.Cmd.Flags["SHIELDS_CUT"] = true
	g.HitEntity(mother, 500, "laser")
	if !mother.Dead && mother.Energy >= 50 {
		t.Error("cutting the field did not expose the mothership")
	}
}
