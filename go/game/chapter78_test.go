package game

import "testing"

func openStage(t *testing.T, stage int) *Game {
	t.Helper()
	g := New()
	g.NewGame("CH78")
	g.Cmd.Campaign.Stage = stage
	g.Cmd.Campaign.Active = true
	g.EnterSystem(13, "launch", 0, false)
	g.Cmd.Campaign.Target = g.Cmd.System
	g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
	g.campaignArrive()
	if g.Screen != "talk" {
		t.Fatalf("stage %d opened %q", stage, g.Screen)
	}
	return g
}

func countTag(g *Game, tag string) int {
	n := 0
	for _, e := range g.World.Entities {
		if e.CampaignTag == tag {
			n++
		}
	}
	return n
}

// Chapter seven: the hive talks once, and the only way to hurt it is to cut the
// station's shields.
func TestChapterSevenShields(t *testing.T) {
	g := openStage(t, 6)
	chooseText(t, g, "KILL IT")
	if !g.Cmd.Flags["HIVE_ATTACKED"] {
		t.Fatal("the attack was not recorded")
	}
	chooseText(t, g, "CUT THE STATION SHIELDS")
	if countTag(g, "MOTHER") != 1 {
		t.Fatal("the mothership is not in the sky")
	}
	// two escorts, the drones come with the mothership itself
	escort := 0
	for _, e := range g.World.Entities {
		if e.Behaviour == "pirate" {
			escort++
		}
	}
	if escort < 2 {
		t.Fatalf("%d pirates came for the fallen shields, want at least the two escorts", escort)
	}
	// with the shields cut, damage lands
	for _, e := range g.World.Entities {
		if e.CampaignTag == "MOTHER" {
			before := e.Energy + e.Shield
			g.HitEntity(e, 50, "laser")
			if e.Energy+e.Shield >= before {
				t.Fatal("the mothership is still immune with the shields cut")
			}
		}
	}
}

// Leaving the shields up makes it invulnerable, and the field falls after ten
// minutes rather than the chapter failing.
func TestChapterSevenWaiting(t *testing.T) {
	g := openStage(t, 6)
	chooseText(t, g, "KILL IT")
	chooseText(t, g, "LEAVE THE SHIELDS UP")
	chooseText(t, g, "CONTINUE") // close the scene
	for _, e := range g.World.Entities {
		if e.CampaignTag == "MOTHER" {
			before := e.Energy + e.Shield
			g.HitEntity(e, 400, "laser")
			if e.Energy+e.Shield < before {
				t.Fatal("the mothership took damage through the station field")
			}
		}
	}
	// the rule, not the frame loop: ten minutes of waiting and the field falls
	for i := 0; i < 601; i++ {
		g.updateHive(1.0)
	}
	if !g.Cmd.Flags["SHIELDS_CUT"] {
		t.Fatalf("the field never fell after ten minutes: wait=%.0f stage=%d active=%v state=%q screen=%q",
			g.HiveWait, g.Cmd.Campaign.Stage, g.Cmd.Campaign.Active, g.State, g.Screen)
	}
}

// The auditor path closes the chapter without a fight, and only with a document.
func TestChapterSevenAuditorPath(t *testing.T) {
	g := openStage(t, 6)
	g.Cmd.Campaign.Stage = 6
	g.Cmd.Flags["MET_AUDITOR"] = true
	g.Cmd.Flags["HAS_COPY"] = true
	g.TalkStart("c7_hive")
	chooseText(t, g, "OPEN THE CHANNEL")
	chooseText(t, g, "PRESENT THE COPY")
	if !g.Cmd.Flags["AUDIT_REFORMED"] {
		t.Fatal("the reform was not recorded")
	}
	if g.Cmd.Campaign.Stage != 7 {
		t.Fatalf("stage %d after the hearing", g.Cmd.Campaign.Stage)
	}
	if countTag(g, "MOTHER") != 0 {
		t.Fatal("a ship was spawned on the peaceful path")
	}
}

// Chapter eight: three answers, three endings, and the council only if the
// courier is alive to bring the eighteen copies.
func TestChapterEightEndings(t *testing.T) {
	t.Run("combine", func(t *testing.T) {
		g := openStage(t, 7)
		chooseText(t, g, "TAKE THE SEAL")
		chooseText(t, g, "CONTINUE")
		if g.Cmd.Ending != "combine" {
			t.Fatalf("ending %q", g.Cmd.Ending)
		}
		if g.Screen != "ending" {
			t.Fatalf("screen %q after the last page", g.Screen)
		}
	})
	t.Run("chaos", func(t *testing.T) {
		g := openStage(t, 7)
		chooseText(t, g, "LET HIM GO")
		chooseText(t, g, "CONTINUE")
		if g.Cmd.Ending != "chaos" {
			t.Fatalf("ending %q", g.Cmd.Ending)
		}
	})
	t.Run("council needs the courier", func(t *testing.T) {
		g := openStage(t, 7)
		chooseText(t, g, "KILL HIM AND TAKE IT")
		chooseText(t, g, "CONTINUE")
		if countTag(g, "FLAGSHIP") != 1 {
			t.Fatal("no flagship after choosing to kill the director")
		}
		g.campaignKill("FLAGSHIP")
		if g.Talk != "c8_after" {
			t.Fatalf("after the flagship: %q", g.Talk)
		}
		n := DialogueNode(g.Talk)
		for _, idx := range g.VisibleOptions(n) {
			if n.Options[idx].Text == "GIVE IT TO THE STATIONS" {
				t.Fatal("the council is offered with the courier dead")
			}
		}
		chooseText(t, g, "KEEP IT")
		chooseText(t, g, "CONTINUE")
		if g.Cmd.Ending != "combine" {
			t.Fatalf("ending %q", g.Cmd.Ending)
		}
	})
	t.Run("council with the courier alive", func(t *testing.T) {
		g := openStage(t, 7)
		g.Cmd.Flags["COURIER_ALIVE"] = true
		chooseText(t, g, "KILL HIM AND TAKE IT")
		chooseText(t, g, "CONTINUE")
		g.campaignKill("FLAGSHIP")
		chooseText(t, g, "GIVE IT TO THE STATIONS")
		chooseText(t, g, "CONTINUE")
		if g.Cmd.Ending != "council" {
			t.Fatalf("ending %q", g.Cmd.Ending)
		}
	})
}
