package game

import "testing"

// A record can be bought clean, and the price has to scale with the damage:
// shooting two vipers must cost more than one bounced cheque.
func TestPardonClearsRecordForMoney(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.World.Sys.Government = 3 // dictatorship: someone can make a call
	g.World.Sys.Tech = 6
	g.Cmd.Legal = 120
	g.Cmd.Cash = 100000

	if !g.PardonAvailable() {
		t.Fatal("a dictator's station should have a fixer")
	}
	small := g.PardonPrice()
	g.Cmd.Legal = 240
	if g.PardonPrice() <= small {
		t.Fatal("price does not grow with the record")
	}
	g.Cmd.Legal = 120
	if !g.BuyPardon() {
		t.Fatal("pardon refused with enough money on a station that offers it")
	}
	if g.Cmd.Legal != 0 {
		t.Fatalf("legal is %d after the pardon", g.Cmd.Legal)
	}
	if g.Cmd.Cash != 100000-float64(small) {
		t.Fatalf("cash %.0f, want %d", g.Cmd.Cash, 100000-small)
	}
}

// A clean pilot has nothing to buy, and a quiet agricultural world has nobody
// selling it.
func TestPardonOnlyWhereOffered(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.Cmd.Legal = 0
	g.Cmd.Cash = 100000
	if g.PardonAvailable() {
		t.Fatal("nothing to pardon")
	}
	g.Cmd.Legal = 300
	g.World.Sys.Government = 6 // democracy
	g.World.Sys.Tech = 5
	if g.PardonAvailable() {
		t.Fatal("a democracy at tech 5 should not sell pardons")
	}
	g.World.Sys.Tech = 12
	if !g.PardonAvailable() {
		t.Fatal("a high tech world should have lawyers")
	}
}

// Without the money nothing happens, and the record stays.
func TestPardonNeedsMoney(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	g.World.Sys.Government = 7
	g.Cmd.Legal = 100
	g.Cmd.Cash = 10
	if g.BuyPardon() {
		t.Fatal("pardon bought with ten credits")
	}
	if g.Cmd.Legal != 100 {
		t.Fatal("record changed without paying")
	}
}
