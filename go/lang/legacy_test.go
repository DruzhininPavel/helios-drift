package lang

import (
	"os"
	"path/filepath"
	"testing"
)

// Preferences written before the rename are picked up by the new build.
func TestLegacyConfigFallback(t *testing.T) {
	dir := t.TempDir()
	// os.UserHomeDir reads USERPROFILE on Windows and HOME everywhere else
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)
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

// The interface keeps its data labels, units and names in English; only prose
// is translated. This is the split the player asked for.
func TestLabelsStayEnglish(t *testing.T) {
	Set("ru")
	defer Set("en")
	for _, s := range []string{
		"HOLD 12/20 T", "STAGE 3/8", "SPD 210", "CASH 1200 CR", "FUEL 5.2 LY",
		"KESTREL III", "LARGE CARGO BAY", "FOOD", "LEGEND", "RICH INDUSTRIAL",
		"ANARCHY", "PAX 2   MAIL 1 T", "RADIUS 6400 KM", "TECH LEVEL 9",
	} {
		if got := TranslateLine(s); got != s {
			t.Fatalf("%q was translated to %q", s, got)
		}
	}
}

// Prose, menus and messages still read in Russian.
func TestProseTranslates(t *testing.T) {
	Set("ru")
	defer Set("en")
	for _, s := range []string{
		"BUY / SELL CARGO", "COMMANDER STATUS", "THE FRONTIER CLINIC IS OUT OF SUPPLIES.",
		"CONTRACT ACCEPTED: MEDICAL RUN", "ARRIVED: LAVE",
	} {
		if got := TranslateLine(s); got == s {
			t.Fatalf("%q was left in English", s)
		}
	}
}
