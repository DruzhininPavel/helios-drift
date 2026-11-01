package game

import "testing"

// Chapter three: board the wreck, read the roster, decide about the recorder,
// ask the company one question, then burn it. Three decisions, one wreck.
func TestChapterThreeWalksItsScreens(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.Cmd.Campaign.Stage = 2
	g.Cmd.Campaign.Active = true
	g.EnterSystem(5, "launch", 0, false)
	g.Cmd.Campaign.Target = g.Cmd.System
	g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
	g.campaignArrive()
	if g.Screen != "talk" || g.Talk != "c3_roster" {
		t.Fatalf("arrival opened %q/%q, want the roster", g.Screen, g.Talk)
	}
	// the roster is a screen with one way forward
	if n := DialogueNode(g.Talk); len(n.Options) != 1 {
		t.Fatalf("roster has %d options", len(n.Options))
	}
	g.TalkChoose(0)
	if g.Talk != "c3_recorder" {
		t.Fatalf("after the roster: %q", g.Talk)
	}
	// copying the recording is a flag, and it leads on
	n := DialogueNode(g.Talk)
	vi := -1
	for i, idx := range g.VisibleOptions(n) {
		if n.Options[idx].Text == "COPY THE RECORDING" {
			vi = i
		}
	}
	g.TalkChoose(vi)
	if !g.Cmd.Flags["COPIED_RECORDER"] {
		t.Fatal("the copy was not recorded")
	}
	if g.Talk != "c3_manifest" {
		t.Fatalf("after copying: %q", g.Talk)
	}
	g.TalkChoose(0) // the menu
	if g.Talk != "c3_holo" {
		t.Fatalf("after the menu: %q", g.Talk)
	}
	g.TalkChoose(0) // one question
	if !g.Cmd.Flags["KNOWS_BOX"] {
		t.Fatal("the question did not set a flag")
	}
	if g.Talk != "c3_ans_box" {
		t.Fatalf("after the question: %q", g.Talk)
	}
	g.TalkChoose(0)
	if g.Talk != "c3_end" {
		t.Fatalf("at the end: %q", g.Talk)
	}
	g.TalkChoose(0) // fire
	if g.Screen == "talk" {
		t.Fatal("the conversation did not close")
	}
	// and the wreck is now in the sky
	found := false
	for _, e := range g.World.Entities {
		if e.CampaignTag == "HULK" {
			found = true
		}
	}
	if !found {
		t.Fatal("no wreck spawned after agreeing to burn it")
	}
}

// Burning the recorder without copying earns the company's approval, which in
// this story means papers and a colder tone.
func TestBurningTheRecorderIsPraised(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.Cmd.Campaign.Stage = 2
	g.Cmd.Campaign.Active = true
	g.EnterSystem(5, "launch", 0, false)
	g.Cmd.Campaign.Target = g.Cmd.System
	g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
	g.campaignArrive()
	g.Cmd.Flags["BURNED_RECORDER"] = true
	papers := g.Cmd.Papers
	g.campaignKill("HULK")
	if g.Cmd.Papers != papers+1 {
		t.Fatalf("papers %d -> %d, want one more", papers, g.Cmd.Papers)
	}
	if g.Cmd.Campaign.Stage != 3 {
		t.Fatalf("stage %d after the wreck", g.Cmd.Campaign.Stage)
	}
}
