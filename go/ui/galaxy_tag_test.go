package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
)

// A job in another galaxy has to say so: seven forward jumps is a long way to
// discover that the contract is behind you.
func TestContractsNameTheirGalaxy(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	g := game.New()
	g.NewGame("GAL")
	g.Cmd.Galaxy = 1
	c := game.Contract{Kind: game.ContractCourier, Good: 0, Qty: 3, System: 2,
		Galaxy: 0, Due: g.Cmd.Day + 4, Pay: 900}
	lines := contractWhat(c)
	joined := ""
	for _, ln := range lines {
		joined += ln + "\n"
	}
	if !containsText(joined, "ГАЛАКТИКА 1") {
		t.Fatalf("the destination does not name its galaxy:\n%s", joined)
	}

	// a local job stays quiet: no clutter where it does not mean anything
	c.Galaxy = 1
	joined = ""
	for _, ln := range contractWhat(c) {
		joined += ln + "\n"
	}
	if containsText(joined, "ГАЛАКТИКА") {
		t.Fatalf("a local job names a galaxy:\n%s", joined)
	}

	// and the story line does the same
	g.Cmd.Campaign.Stage = 3
	g.Cmd.Campaign.Active = true
	g.Cmd.Campaign.Galaxy = 0
	story := ""
	for _, ln := range storyWhat() {
		story += ln + "\n"
	}
	if !containsText(story, "КОНТРАКТ В ГАЛАКТИКЕ 1") {
		t.Fatalf("the story does not say where it went:\n%s", story)
	}
}
