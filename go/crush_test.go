package main

import (
	"os"
	"path/filepath"
	"testing"

	"heliosdrift/input"
)

// A rebind must survive: saved atomically, with the previous file kept aside, and
// a broken file preserved rather than overwritten.
func TestBindingsSurviveSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	input.DirOverride = dir
	defer func() { input.DirOverride = "" }()

	input.Reset()
	input.SetKey(input.Accel, 42)
	if !input.Save() {
		t.Fatal("save failed")
	}
	path := filepath.Join(dir, ".heliosdrift-keys.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("no keys file: %v", err)
	}
	// a second save keeps the first one as a backup
	input.SetKey(input.Accel, 43)
	if !input.Save() {
		t.Fatal("second save failed")
	}
	if _, err := os.Stat(path + ".bak"); err != nil {
		t.Fatalf("no backup: %v", err)
	}
	// and a broken file does not silently replace the pilot's bindings
	os.WriteFile(path, []byte("{not json"), 0o644)
	input.Load()
	if _, err := os.Stat(path + ".bad"); err != nil {
		t.Fatalf("broken file was not kept: %v", err)
	}
}
