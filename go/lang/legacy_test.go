package lang

import (
	"os"
	"path/filepath"
	"testing"
)

// Preferences written before the rename are picked up by the new build.
func TestLegacyConfigFallback(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	raw := []byte(`{"lang":"ru","muted":true,"volume":0.4}`)
	if err := os.WriteFile(filepath.Join(dir, ".elite16-config.json"), raw, 0o644); err != nil {
		t.Fatal(err)
	}
	defer func() { Set("en"); SetSaved("en") }()

	Load()
	if !IsRU() {
		t.Fatal("legacy config ignored")
	}
}
