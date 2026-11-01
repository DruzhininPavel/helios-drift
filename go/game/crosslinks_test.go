package game

import "testing"

// Flags that were set and never read. Each of these is a place where an earlier
// decision has to come back.
func TestEarlierDecisionsComeBack(t *testing.T) {
	t.Run("paid ambush spares a commander", func(t *testing.T) {
		g := openStage(t, 6)
		g.Cmd.Flags["PAID_PERMIT"] = true
		chooseText(t, g, "KILL IT")
		chooseText(t, g, "CUT THE STATION SHIELDS")
		quiet := 0
		for _, e := range g.World.Entities {
			if e.Behaviour == "depart" {
				quiet++
			}
		}
		if quiet == 0 {
			t.Fatal("nobody turned away after the ambush was paid off")
		}
	})

	t.Run("the hive remembers the question", func(t *testing.T) {
		for _, key := range []string{"KNOWS_HIVE", "KNOWS_BOX", "ASKED_PRICE"} {
			g := New()
			g.NewGame("Q")
			g.Cmd.Flags = map[string]bool{key: true}
			lines := g.NodeLines(DialogueNode("c7_hive"))
			if len(lines) == 0 || lines[0] == DialogueNode("c7_hive").Lines[0] {
				t.Fatalf("%s did not change what the hive says", key)
			}
		}
	})

	t.Run("a forge closes the quiet endings", func(t *testing.T) {
		g := openStage(t, 7)
		g.Cmd.Flags["FORGED"] = true
		n := DialogueNode(g.Talk)
		vis := g.VisibleOptions(n)
		if len(vis) != 1 || n.Options[vis[0]].Text != "KILL HIM AND TAKE IT" {
			t.Fatalf("the director offers %d answers to a forger", len(vis))
		}
	})

	t.Run("a broken audit stops delivery", func(t *testing.T) {
		g := New()
		g.NewGame("AUDIT")
		g.State = "docked"
		g.Cmd.Flags["AUDIT_BROKEN"] = true
		g.rollContracts()
		for _, c := range g.Offers {
			if c.Kind == ContractCourier || c.Kind == ContractMail || c.Kind == ContractPassenger {
				t.Fatalf("a %s contract is still on the board", c.Kind)
			}
		}
		// and the market is worth thirty per cent less
		g.dockedAt()
		plain := New()
		plain.NewGame("AUDIT")
		plain.EnterSystem(g.Cmd.System, "arrival", 0, false)
		plain.dockedAt()
		if g.Market[0].Price >= plain.Market[0].Price {
			t.Fatalf("panic prices %d vs %d", g.Market[0].Price, plain.Market[0].Price)
		}
	})

	t.Run("the recording is evidence", func(t *testing.T) {
		g := New()
		g.NewGame("EV")
		g.Cmd.Flags = map[string]bool{"COPIED_RECORDER": true}
		g.Cmd.Campaign.Stage = 5
		g.Cmd.Campaign.Active = true
		g.EnterSystem(7, "launch", 0, false)
		g.Cmd.Campaign.Target = g.Cmd.System
		g.Cmd.Campaign.Galaxy = g.Cmd.Galaxy
		g.campaignArrive()
		g.Cmd.Papers = 1
		chooseText(t, g, "SHOW THE PAPERS")
		chooseText(t, g, "CONTINUE")
		chooseText(t, g, "NO")
		if !hasOption(g, "PLAY THE RECORDING") {
			t.Fatal("the recorder is not offered as evidence")
		}
	})
}

func hasOption(g *Game, text string) bool {
	n := DialogueNode(g.Talk)
	if n == nil {
		return false
	}
	for _, idx := range g.VisibleOptions(n) {
		if n.Options[idx].Text == text {
			return true
		}
	}
	return false
}
