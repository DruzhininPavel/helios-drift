package game

import (
	"testing"

	"heliosdrift/galaxy"
)

func findOffer(g *Game, kind string) int {
	for i, c := range g.Offers {
		if c.Kind == kind {
			return i
		}
	}
	return -1
}

// A courier run loads cargo on accept and pays at the destination.
func TestCourierContract(t *testing.T) {
	g := newDocked()
	// force a courier onto the board so the test is deterministic
	g.Offers = []Contract{{Kind: ContractCourier, Good: 0, Qty: 3,
		System: 42, Galaxy: g.Cmd.Galaxy, Pay: 500}}
	i := findOffer(g, ContractCourier)
	if i < 0 {
		t.Fatal("no courier offer")
	}
	if !g.AcceptOffer(i) {
		t.Fatal("accept failed")
	}
	if g.Cmd.Cargo[0] != 3 {
		t.Fatalf("cargo = %d, want 3", g.Cmd.Cargo[0])
	}
	if len(g.Offers) != 0 {
		t.Fatal("offer still on the board")
	}

	cash := g.Cmd.Cash
	g.Cmd.Galaxy = g.Cmd.Contracts[0].Galaxy
	g.EnterSystem(42, "arrive", 9000, false)
	dockNow(g)

	if g.Cmd.Cash != cash+500 {
		t.Fatalf("cash = %v, want %v", g.Cmd.Cash, cash+500)
	}
	if g.Cmd.Cargo[0] != 0 {
		t.Fatalf("cargo not unloaded: %d", g.Cmd.Cargo[0])
	}
	if len(g.Cmd.Contracts) != 0 {
		t.Fatal("contract not cleared after delivery")
	}
}

// A bounty job pays after the last pirate dies in the target system.
func TestBountyContract(t *testing.T) {
	g := newDocked()
	g.Offers = []Contract{{Kind: ContractBounty, Qty: 2, System: g.Cmd.System,
		Galaxy: g.Cmd.Galaxy, Pay: 300}}
	if !g.AcceptOffer(0) {
		t.Fatal("accept failed")
	}
	cash := g.Cmd.Cash

	kill := func() {
		e := &Entity{Kind: "ship", Type: "krait", Behaviour: "pirate", Radius: 20, Energy: 1}
		g.World.Entities = append(g.World.Entities, e)
		g.DestroyEntity(e, "laser")
	}
	kill()
	if len(g.Cmd.Contracts) != 1 || g.Cmd.Contracts[0].Done != 1 {
		t.Fatalf("progress = %+v, want 1 kill", g.Cmd.Contracts)
	}
	kill()
	if len(g.Cmd.Contracts) != 0 {
		t.Fatal("contract not cleared after the last kill")
	}
	if g.Cmd.Cash < cash+300 {
		t.Fatalf("cash = %v, want at least %v", g.Cmd.Cash, cash+300)
	}
}

// The board refills on docking and never holds more than three jobs.
func TestContractBoard(t *testing.T) {
	g := newDocked()
	g.Cmd.Contracts = nil
	g.Cmd.Cargo = map[int]int{}
	g.dockedAt()
	if len(g.Offers) == 0 {
		t.Fatal("board is empty after docking")
	}
	for len(g.Offers) > 0 {
		if !g.AcceptOffer(0) {
			break
		}
	}
	if len(g.Cmd.Contracts) > 3 {
		t.Fatalf("%d contracts accepted, want at most 3", len(g.Cmd.Contracts))
	}
}

// Contraband: loads the hold, pays on delivery, and only moves through anarchy
// and feudal space.
func TestSmuggleContract(t *testing.T) {
	g := newDocked()
	g.Offers = []Contract{{Kind: ContractSmuggle, Good: 6, Qty: 2, System: 42,
		Galaxy: g.Cmd.Galaxy, Pay: 900, Due: g.Cmd.Day + 3}}
	if !g.AcceptOffer(0) {
		t.Fatal("accept failed")
	}
	if g.Cmd.Cargo[6] != 2 {
		t.Fatalf("cargo = %d, want 2", g.Cmd.Cargo[6])
	}
	if !g.CarryingContraband() {
		t.Fatal("contraband not flagged")
	}
	cash := g.Cmd.Cash
	g.Cmd.Galaxy = g.Cmd.Contracts[0].Galaxy
	g.EnterSystem(42, "arrive", 9000, false)
	dockNow(g)
	if g.Cmd.Cash != cash+900 {
		t.Fatalf("cash = %v, want %v", g.Cmd.Cash, cash+900)
	}
	if g.Cmd.Cargo[6] != 0 || g.CarryingContraband() {
		t.Fatal("contraband still aboard after delivery")
	}
}

func TestSmuggleOnlyInRoughSpace(t *testing.T) {
	g := newGameMode(ModeAdvanced)
	gal := galaxy.Get(g.Cmd.Galaxy)
	rough, lawful := -1, -1
	for _, s := range gal.Systems {
		if rough < 0 && s.Government <= 1 {
			rough = s.ID
		}
		if lawful < 0 && s.Government >= 5 {
			lawful = s.ID
		}
	}
	if rough < 0 || lawful < 0 {
		t.Skip("no suitable systems in this galaxy")
	}
	countSmuggle := func(sys int) int {
		n := 0
		for i := 0; i < 200; i++ {
			g.Cmd.System = sys
			g.EnterSystem(sys, "arrive", 9000, false)
			g.Cmd.Contracts = nil
			g.Offers = nil
			g.rollContracts()
			for _, c := range g.Offers {
				if c.Kind == ContractSmuggle {
					n++
				}
			}
		}
		return n
	}
	if countSmuggle(rough) == 0 {
		t.Fatal("no smuggling jobs in anarchy space")
	}
	if got := countSmuggle(lawful); got != 0 {
		t.Fatalf("%d smuggling jobs in lawful space", got)
	}
}

// TestContractFineCannotGoNegative keeps a missed deadline from bankrupting a
// commander into a negative balance, which the game cannot recover from.
func TestContractFineCannotGoNegative(t *testing.T) {
	g := New()
	g.NewGameMode(ModeAdvanced, "BOT")
	g.Cmd.Cash = 10
	g.Cmd.Contracts = []Contract{{
		Kind: ContractCourier, Pay: 2000, Due: 0, System: 5, Galaxy: 0,
	}}
	g.Cmd.Day = 5 // well past the deadline
	g.expireContracts()
	if g.Cmd.Cash < 0 {
		t.Fatalf("fine drove cash to %.0f", g.Cmd.Cash)
	}
	if len(g.Cmd.Contracts) != 0 {
		t.Fatal("expired contract stayed in the list")
	}
}
