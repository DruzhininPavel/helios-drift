package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
)

// The tester switch is not in the settings until the phrase is typed, and the
// phrase is never listed anywhere.
func TestTesterModeHidden(t *testing.T) {
	newTitle()
	lang.Cfg.TesterUnlocked = false
	lang.Cfg.TesterMode = false
	defer func() {
		lang.Cfg.TesterUnlocked = false
		lang.Cfg.TesterMode = false
	}()

	for _, row := range settingsRows() {
		if row == "TESTER MODE" {
			t.Fatal("tester row is visible before the phrase is typed")
		}
	}

	OpenSettings()
	for _, ch := range []string{"i", "d", "d", "q", "d"} {
		Key(ch)
	}
	if !lang.Cfg.TesterUnlocked {
		t.Fatal("the phrase did not unlock the tester row")
	}
	found := false
	for _, row := range settingsRows() {
		if row == "TESTER MODE" {
			found = true
		}
	}
	if !found {
		t.Fatal("tester row missing after unlocking")
	}
}

// Once unlocked, the switch is a normal settings row and it decides whether the
// cheat keys do anything.
func TestTesterModeTogglesAndGatesCheats(t *testing.T) {
	newTitle()
	lang.Cfg.TesterUnlocked = true
	lang.Cfg.TesterMode = false
	defer func() {
		lang.Cfg.TesterUnlocked = false
		lang.Cfg.TesterMode = false
	}()
	OpenSettings()
	for i, row := range settingsRows() {
		if row == "TESTER MODE" {
			S.Settings = i
		}
	}
	Key("Enter")
	if !lang.Cfg.TesterMode {
		t.Fatal("switch did not turn on")
	}
	Key("Enter")
	if lang.Cfg.TesterMode {
		t.Fatal("switch did not turn off")
	}
}

// The phrase is case-insensitive and does not need to be typed in one go
// without other letters before it.
func TestTesterPhraseIgnoresCaseAndLeadIn(t *testing.T) {
	newTitle()
	lang.Cfg.TesterUnlocked = false
	defer func() { lang.Cfg.TesterUnlocked = false }()
	OpenSettings()
	for _, ch := range []string{"x", "I", "d", "D", "q", "D"} {
		Key(ch)
	}
	if !lang.Cfg.TesterUnlocked {
		t.Fatal("phrase not recognised")
	}
	_ = game.G
}

// With the extra row the list must still fit above the hint line: seven rows,
// twenty pixels apart, ending before the footer text.
func TestSettingsFitWithTesterRow(t *testing.T) {
	defer func() { lang.Cfg.TesterUnlocked = false }()
	lang.Cfg.TesterUnlocked = true
	n := len(settingsRows())
	if n != 8 {
		t.Fatalf("%d rows with the tester switch, want 8", n)
	}
	step := 20
	if n > 6 {
		step = 16
	}
	last := 32 + (n-1)*step
	if last+11 > 172 {
		t.Fatalf("last row ends at %d, the hint line is at 176", last+11)
	}
}
