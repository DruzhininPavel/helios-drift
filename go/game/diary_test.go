package game

import "testing"

// The diary is the story as it was lived: pages appear when the turn happens,
// in order, and two pilots with different choices end up with different books.
func TestDiaryFollowsTheStory(t *testing.T) {
	g := New()
	g.NewGame("DIARY")
	if n := len(g.DiaryAvailable()); n != 1 {
		t.Fatalf("%d pages at the start, want the prologue only", n)
	}
	prologue := g.DiaryAvailable()[0]
	if prologue.ID != "prologue" {
		t.Fatalf("first page is %q", prologue.ID)
	}

	// take the crate
	g.Cmd.Flags["HAS_SEAL"] = true
	ids := idsOf(g)
	if !contains(ids, "ch1_crate") {
		t.Fatalf("no page for the crate: %v", ids)
	}
	// the other branch stays out of this book: a pilot who refused gets a
	// different page and never this one
	other := New()
	other.NewGame("DIARY")
	other.Cmd.Flags["REFUSED_SEAL"] = true
	for _, id := range idsOf(other) {
		if id == "ch1_crate" {
			t.Fatal("the crate page appeared for a pilot who refused it")
		}
	}
	if !contains(idsOf(other), "ch1_refused") {
		t.Fatal("no page for refusing the crate")
	}

	// order holds even when the flags are set out of sequence
	g.Cmd.Flags["SHIELDS_CUT"] = true
	g.Cmd.Flags["CONVOY_SPARED"] = true
	prev := -1
	for _, e := range g.DiaryAvailable() {
		if e.Order < prev {
			t.Fatalf("page %s is out of order", e.ID)
		}
		prev = e.Order
	}

	// and reading is remembered
	if g.DiaryUnread() == 0 {
		t.Fatal("nothing is unread")
	}
	before := g.DiaryUnread()
	g.DiaryMarkRead("prologue")
	g.DiaryMarkRead("prologue")
	if g.DiaryUnread() != before-1 {
		t.Fatalf("unread %d, want %d", g.DiaryUnread(), before-1)
	}
}

func idsOf(g *Game) []string {
	var out []string
	for _, e := range g.DiaryAvailable() {
		out = append(out, e.ID)
	}
	return out
}

func contains(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}
