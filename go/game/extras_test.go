package game

import "testing"

// Breaking the convoy halves the local market for a day and leaves three crates
// in the debris. Killing the courier puts a price on the hull and brings pirates.
func TestConsequencesOfTheMiddleChapters(t *testing.T) {
	g := New()
	g.NewGame("EXTRAS")
	g.EnterSystem(4, "arrival", 0, false)
	g.Cmd.Campaign.Stage = 4
	g.Cmd.Campaign.Active = true
	if g.Cmd.Flags == nil {
		g.Cmd.Flags = map[string]bool{}
	}
	g.Cmd.Flags["CONVOY_ATTACKED"] = true
	g.campaignKill("CONVOY")

	if !g.PanicHere() {
		t.Fatal("the market is not in a panic after the convoy")
	}
	if g.Cmd.Boxes != 3 {
		t.Fatalf("%d crates recovered, want three", g.Cmd.Boxes)
	}
	// the day after, prices are back
	g.Cmd.Day += 2
	if g.PanicHere() {
		t.Fatal("the panic outlived its day")
	}
	// and another system never had one
	g.Cmd.Day = g.Cmd.PanicDay
	g.Cmd.System = g.Cmd.PanicSys + 1
	if g.PanicHere() {
		t.Fatal("the panic followed the pilot to another system")
	}
}

func TestBountyMakesTheSkyBusier(t *testing.T) {
	g := New()
	g.NewGame("BOUNTY")
	base := g.dangerNow(1.0)
	if g.Cmd.Flags == nil {
		g.Cmd.Flags = map[string]bool{}
	}
	g.Cmd.Flags["COURIER_DEAD"] = true
	if got := g.dangerNow(1.0); got <= base {
		t.Fatalf("danger %.2f after the bounty, was %.2f", got, base)
	}
}

// The closing notice is the same list every time it is read: it is a notice.
func TestClosingStationsIsStable(t *testing.T) {
	g := New()
	g.NewGame("LIST")
	a := g.ClosingStations()
	b := g.ClosingStations()
	if len(a) != 18 {
		t.Fatalf("%d stations on the notice, want eighteen", len(a))
	}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("the notice changed between readings at %d", i)
		}
	}
}
