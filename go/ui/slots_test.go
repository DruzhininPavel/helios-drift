package ui

import (
	"os"
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
)

// Overwriting and erasing both take two presses: the first asks, the second does
// it, and moving the cursor cancels.
func TestSlotsOverwriteAndErase(t *testing.T) {
	lang.Set("en")
	defer lang.Set("ru")
	dir := t.TempDir()
	game.SaveDirOverride = dir
	defer func() { game.SaveDirOverride = "" }()

	g := game.New()
	g.NewGame("SLOTS")
	g.State = "docked"
	g.Screen = "launch"
	g.Cmd.Cash = 5000
	if !g.SaveSlot(0) {
		t.Fatal("could not write a slot")
	}

	// saving over it: first Enter only asks
	openSlots(true)
	S.Slot, S.SlotConfirm = 0, false
	Key("Enter")
	if !S.SlotConfirm || g.Screen != "slots" {
		t.Fatal("the first Enter overwrote without asking")
	}
	Key("Enter")
	if g.Screen == "slots" {
		t.Fatal("the second Enter did not save")
	}

	// erasing: same rule, and the file is kept aside
	openSlots(false)
	S.Slot, S.SlotConfirm = 0, false
	Key("Delete")
	if !S.SlotConfirm {
		t.Fatal("Delete did not ask")
	}
	if !game.SlotExists(0) {
		t.Fatal("the slot was erased on the first press")
	}
	Key("Delete")
	if game.SlotExists(0) {
		t.Fatal("the second Delete did not erase")
	}
	if _, err := os.Stat(dir + "/.heliosdrift-save1.json.deleted"); err != nil {
		t.Fatalf("the erased commander was not kept aside: %v", err)
	}

	// and moving the cursor calls the whole thing off
	openSlots(true)
	S.Slot, S.SlotConfirm = 0, false
	Key("Delete")
	Key("ArrowDown")
	if S.SlotConfirm {
		t.Fatal("moving the cursor kept the erase armed")
	}
}
