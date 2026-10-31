package game

import "testing"

// Chapter one, all three branches: the money is the same, the crate and the
// paperwork are not. This is the smallest test of the whole premise.
func TestChapterOneBranches(t *testing.T) {
	cases := []struct {
		name   string
		opt    int
		boxes  int
		papers int
		flag   string
	}{
		{"take", 0, 1, 1, "HAS_SEAL"},
		{"refuse", 1, 0, 0, "REFUSED_SEAL"},
		{"receipt", 2, 1, 2, "HAS_RECEIPT"},
	}
	for _, c := range cases {
		g := New()
		g.NewGame("TEST")
		g.EnterSystem(3, "launch", 0, false)
		before := g.Cmd.Cash

		if !g.TalkStart("c1_disp") {
			t.Fatalf("%s: no such dialogue", c.name)
		}
		if !g.TalkChoose(c.opt) && g.Talk != "" {
			t.Fatalf("%s: conversation did not end", c.name)
		}
		if g.Cmd.Cash != before+1200 {
			t.Errorf("%s: cash %.0f, want %.0f", c.name, g.Cmd.Cash, before+1200)
		}
		if g.Cmd.Boxes != c.boxes {
			t.Errorf("%s: boxes %d, want %d", c.name, g.Cmd.Boxes, c.boxes)
		}
		if g.Cmd.Papers != c.papers {
			t.Errorf("%s: papers %d, want %d", c.name, g.Cmd.Papers, c.papers)
		}
		if !g.Cmd.Flags[c.flag] {
			t.Errorf("%s: flag %s not set", c.name, c.flag)
		}
	}
}

// The crate takes two tonnes of hold and cannot be dumped: it is bound to the
// pilot, which is the whole point of it.
func TestCrateRules(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	free := g.Cmd.CargoCap
	if CargoUsed(g.Cmd) != 0 {
		t.Fatalf("fresh hold is not empty: %d", CargoUsed(g.Cmd))
	}
	g.Cmd.Boxes = 1
	if got := CargoUsed(g.Cmd); got != 2 {
		t.Fatalf("crate takes %d tonnes, want 2", got)
	}
	if CargoUsed(g.Cmd) > free {
		t.Fatal("the crate came from nowhere: hold is over capacity on a fresh ship")
	}
	g.State = "flight"
	g.Cmd.Cargo = map[int]int{}
	if g.Jettison() {
		t.Fatal("the crate was jettisoned")
	}
}

// A node's graph must be sound: every option leads somewhere real, and every
// condition names something the game knows.
func TestDialogueGraphIntegrity(t *testing.T) {
	for _, id := range DialogueNodeIDs() {
		n := DialogueNode(id)
		for _, o := range n.Options {
			if o.Next != "" && DialogueNode(o.Next) == nil {
				t.Errorf("%s: option %q leads to missing node %q", id, o.Text, o.Next)
			}
			if o.Text == "" {
				t.Errorf("%s: option without text", id)
			}
			for _, e := range o.Do {
				switch e.Kind {
				case FxFlag, FxUnflag, FxCash, FxPapers, FxBox, FxWanted, FxClock, FxAudit, FxEnding, FxQuest:
				default:
					t.Errorf("%s: unknown effect %q", id, e.Kind)
				}
			}
		}
		if len(n.Options) == 0 && n.Kind == "talk" {
			t.Errorf("%s: talk node with nothing to say", id)
		}
	}
}

// Conditions gate options: a locked branch must not be offered.
func TestOptionConditionsGate(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	registerNodes(&Node{
		ID: "test_gate", Kind: "talk", Speaker: "X",
		Lines: []string{"HELLO."},
		Options: []Option{
			{Text: "ALWAYS", Next: ""},
			{Text: "IF BOX", Next: "", If: []Cond{{Key: "boxes", Op: ">=", N: 1}}},
			{Text: "IF FLAG", Next: "", If: []Cond{{Key: "HAS_SEAL", Op: "flag"}}},
		},
	})
	defer delete(dialogueNodes, "test_gate")
	g.TalkStart("test_gate")
	if vis := g.VisibleOptions(DialogueNode("test_gate")); len(vis) != 1 {
		t.Fatalf("visible %v, want only the unconditional option", vis)
	}
	g.Cmd.Boxes = 1
	if vis := g.VisibleOptions(DialogueNode("test_gate")); len(vis) != 2 {
		t.Fatalf("visible %v, want two", vis)
	}
}

// Delivering the food does not pay out on the spot: it puts you at the counter,
// with the dialogue open. That is the change the whole campaign rests on.
func TestDeliveryOpensTheCounter(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	if !g.AcceptCampaign() {
		t.Fatal("no campaign offer at the start")
	}
	if g.Cmd.Cargo[0] < 4 {
		t.Fatalf("the board did not load the cargo: %d tonnes", g.Cmd.Cargo[0])
	}
	if g.Cmd.Cash < 400 {
		t.Fatalf("no advance paid: %.0f", g.Cmd.Cash)
	}
	g.Cmd.Campaign.Target = g.Cmd.System
	g.dockedAt()
	g.campaignDock()
	if g.Screen != "talk" || g.Talk != "c1_disp" {
		t.Fatalf("screen %q talk %q, want the counter", g.Screen, g.Talk)
	}
}

// Classic never meets the story: no counter, no crate, no stage.
func TestClassicHasNoDialogue(t *testing.T) {
	g := New()
	g.NewGameMode(ModeClassic, "TEST")
	g.EnterSystem(3, "launch", 0, false)
	if g.CampaignOffer() {
		t.Fatal("classic offered a campaign stage")
	}
	g.Cmd.Campaign.Target = g.Cmd.System
	g.dockedAt()
	g.campaignDock()
	if g.Screen == "talk" || g.Talk != "" {
		t.Fatal("classic opened a dialogue")
	}
}

// The crate travels with the pilot: a new hull keeps it and reserves the space.
func TestCrateSurvivesShipChange(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.Cmd.Boxes = 1
	g.Cmd.Cash = 500000
	if !g.BuyShip("adder") {
		t.Fatal("could not buy a hull with a crate aboard")
	}
	if g.Cmd.Boxes != 1 {
		t.Fatalf("crate lost on the way to the new ship: %d", g.Cmd.Boxes)
	}
	if g.Cmd.CargoCap < 2 {
		t.Fatalf("new hull has %d tonnes, no room for the crate", g.Cmd.CargoCap)
	}
}

// And it heats the cabin: a box in the hold raises the resting temperature.
func TestCrateHeatsTheCabin(t *testing.T) {
	run := func(boxes int) float64 {
		g := New()
		g.NewGame("TEST")
		g.EnterSystem(3, "launch", 0, false)
		g.Cmd.Boxes = boxes
		g.Player.Invulnerable = 1e9
		for i := 0; i < 240; i++ {
			g.Update(1.0 / 60.0)
		}
		return g.Player.Cabin
	}
	plain, loaded := run(0), run(1)
	if loaded <= plain {
		t.Fatalf("cabin %.3f with a crate, %.3f without: the box does not warm it", loaded, plain)
	}
}
