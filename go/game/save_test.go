package game

import (
	"os"
	"path/filepath"
	"testing"
)

// TestThreeSlots checks the save slots: a commander written to one slot loads
// back from that slot only, and the slot describes itself on the load screen.
func TestThreeSlots(t *testing.T) {
	dir := t.TempDir()
	SaveDirOverride = dir
	defer func() { SaveDirOverride = "" }()

	g := newGameMode(ModeAdvanced)
	g.Cmd.Name = "FIRST"
	g.Cmd.Cash = 111
	if !g.SaveSlot(0) {
		t.Fatal("save to slot 1 failed")
	}
	g.Cmd.Name = "SECOND"
	g.Cmd.Cash = 222
	if !g.SaveSlot(1) {
		t.Fatal("save to slot 2 failed")
	}
	if SlotExists(2) {
		t.Fatal("empty slot reports a commander")
	}

	other := New()
	if !other.LoadSlot(1) {
		t.Fatal("slot 2 did not load")
	}
	if other.Cmd.Name != "SECOND" || other.Cmd.Cash != 222 {
		t.Fatalf("loaded %q %.0f from slot 2", other.Cmd.Name, other.Cmd.Cash)
	}

	name, cash, where, ok := SlotInfo(0)
	if !ok || name != "FIRST" || cash != 111 || where == "" {
		t.Fatalf("slot 1 info: %q %.0f %q %v", name, cash, where, ok)
	}

	// the files live under the current name only
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			t.Errorf("unexpected file in the save directory: %s", e.Name())
		}
	}
	if !SaveExists() {
		t.Fatal("SaveExists is false with two slots written")
	}
}
