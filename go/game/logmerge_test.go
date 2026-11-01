package game

import "testing"

// Buying ten tonnes in ten presses must read as one line: the log records what
// happened, not which keys were pressed.
func TestRepeatedTradesMerge(t *testing.T) {
	g := New()
	g.NewGame("TEST")
	g.EnterSystem(3, "launch", 0, false)
	before := len(g.Cmd.Log)
	for i := 0; i < 10; i++ {
		g.Note(LogBuy, "FOOD", 1, -22)
	}
	if got := len(g.Cmd.Log) - before; got != 1 {
		t.Fatalf("ten presses wrote %d lines", got)
	}
	last := g.Cmd.Log[len(g.Cmd.Log)-1]
	if last.N != 10 || last.CR != -220 {
		t.Fatalf("merged line is %d tonnes for %.0f, want 10 for -220", last.N, last.CR)
	}
	// a different commodity is a different line
	g.Note(LogBuy, "TEXTILES", 2, -44)
	if got := len(g.Cmd.Log) - before; got != 2 {
		t.Fatalf("a second commodity merged into the first: %d lines", got)
	}
}
