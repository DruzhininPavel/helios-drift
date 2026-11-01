package ui

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// The pilot's own save, with whatever it accumulated.
func TestSavedCommanderJournal(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	home, _ := os.UserHomeDir()
	dir := t.TempDir()
	for slot := 0; slot < game.Slots; slot++ {
		src := filepath.Join(home, ".heliosdrift-save"+strconv.Itoa(slot+1)+".json")
		raw, err := os.ReadFile(src)
		if err != nil {
			continue
		}
		os.WriteFile(filepath.Join(dir, ".heliosdrift-save"+strconv.Itoa(slot+1)+".json"), raw, 0o644)
	}
	game.SaveDirOverride = dir
	defer func() { game.SaveDirOverride = "" }()

	fb := render.NewFB(320, 224)
	for slot := 0; slot < game.Slots; slot++ {
		if !game.SlotExists(slot) {
			continue
		}
		g := game.New()
		if !g.LoadSlot(slot) {
			continue
		}
		done := make(chan string, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					done <- "panic"
				}
			}()
			for _, screen := range []string{"journal", "log", "launch", "status", "contracts"} {
				g.Screen = screen
				Draw(fb, 1.0/60.0)
				Key("ArrowDown")
				Draw(fb, 1.0/60.0)
			}
			done <- "ok"
		}()
		select {
		case res := <-done:
			if res != "ok" {
				t.Fatalf("slot %d: %s", slot, res)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("slot %d hangs", slot)
		}
	}
}
