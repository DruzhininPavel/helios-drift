package game

import (
	"fmt"
	"testing"
)

// Accepting the black-box contract and jumping to the target must open the
// roster, not a firefight.
func TestAcceptThenJumpOpensTheRoster(t *testing.T) {
	g := New()
	g.NewGame("FLOW")
	g.State = "docked"
	g.Screen = "launch"
	g.Cmd.Campaign.Stage = 2
	if !g.AcceptCampaign() {
		t.Fatal("the contract was refused")
	}
	fmt.Printf("target=%d galaxy=%d active=%v\n", g.Cmd.Campaign.Target, g.Cmd.Campaign.Galaxy, g.Cmd.Campaign.Active)
	g.Launch()
	g.EnterSystem(g.Cmd.Campaign.Target, "arrive", 0, false)
	if g.Screen != "talk" {
		t.Fatalf("screen %q, want the roster; talk=%q", g.Screen, g.Talk)
	}
}
