package game

import (
	"os"
	"path/filepath"
	"testing"
)

// A commander saved by a build from before the rename must still load.
func TestLegacySaveFileLoads(t *testing.T) {
	dir := t.TempDir()
	SaveDirOverride = dir
	defer func() { SaveDirOverride = "" }()

	g := newGameMode(ModeAdvanced)
	g.Cmd.Name = "OLD TIMER"
	g.Cmd.Cash = 4321
	if !g.SaveSlot(1) {
		t.Fatal("save failed")
	}
	newName := savePath(1)
	oldName := filepath.Join(dir, ".elite16-save2.json")
	if err := os.Rename(newName, oldName); err != nil {
		t.Fatal(err)
	}

	other := New()
	if !other.LoadSlot(1) {
		t.Fatal("legacy save did not load")
	}
	if other.Cmd.Name != "OLD TIMER" || other.Cmd.Cash != 4321 {
		t.Fatalf("loaded %q %v", other.Cmd.Name, other.Cmd.Cash)
	}
}
