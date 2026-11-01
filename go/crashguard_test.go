package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"heliosdrift/game"
)

// A frame that panics must not take the window with it: the stack goes to a
// log, the run goes to a crash save, and the frame is skipped.
func TestCrashGuardCatchesAFrameAndWritesTheLog(t *testing.T) {
	crashDirOverride = t.TempDir()
	defer func() { crashDirOverride = "" }()
	os.Setenv("HOME", crashDirOverride) // keep the crash save out of the real home
	defer os.Unsetenv("HOME")
	game.G = game.New()
	game.G.NewGame("GUARD")
	before := listCrashes(t)
	frameCrashes = 0

	fatal := guardFrame("update", func() { panic("boom in the frame") })
	if fatal {
		t.Fatal("one crash was treated as fatal")
	}
	after := listCrashes(t)
	if len(after) <= len(before) {
		t.Fatal("no crash log was written")
	}
	log := filepath.Join(crashDir(), after[len(after)-1])
	data, err := os.ReadFile(log)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, want := range []string{"boom in the frame", "where:  update", "goroutine"} {
		if !strings.Contains(text, want) {
			t.Errorf("log has no %q", want)
		}
	}
	home, _ := os.UserHomeDir()
	if _, err := os.Stat(filepath.Join(home, ".heliosdrift-crash.json")); err != nil {
		t.Errorf("no crash save: %v", err)
	}

	// three in a row is the point where the frame is not trusted any more
	guardFrame("update", func() { panic("second") })
	if !guardFrame("update", func() { panic("third") }) {
		t.Fatal("the third crash in a row did not ask for the window to close")
	}
	frameCrashes = 0
}

func listCrashes(t *testing.T) []string {
	t.Helper()
	ents, err := os.ReadDir(crashDir())
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range ents {
		if strings.HasPrefix(e.Name(), "crash-") {
			out = append(out, e.Name())
		}
	}
	return out
}
