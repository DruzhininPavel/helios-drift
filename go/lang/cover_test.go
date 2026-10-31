package lang

import (
	"strings"
	"testing"
)

// TestUiScreensAreTranslated guards the mistake that shipped in v1.14-2: the
// new screens were written in English and the Russian map never got them, so a
// Russian player saw "HOW TO PLAY" and an English handbook.
func TestUiScreensAreTranslated(t *testing.T) {
	Set("ru")
	defer Set("en")
	want := []string{
		"HOW TO PLAY", "HANDBOOK", "FLIGHT CONTROLS", "YAW LEFT / RIGHT",
		"MAIN MENU", "RESUME", "DOCKING", "FIGHTING", "TRAVEL",
		"SPEED MUST BE UNDER 85% OR YOU HIT",
		"BUY WHERE IT IS CHEAP, SELL WHERE IT IS DEAR.",
		"C OPENS THE GALACTIC CHART.",
	}
	for _, s := range want {
		if T(s) == s || strings.Contains(T(s), "  ") {
			t.Errorf("%q is not translated: %q", s, T(s))
		}
	}
}

// TestHintsKeepTheirTail: a hint that ends with ESC BACK in English must still
// say НАЗАД in Russian. Trimming a line to fit the bar once dropped the word.
func TestHintsKeepTheirTail(t *testing.T) {
	Set("ru")
	defer Set("en")
	for k, v := range ru {
		if strings.Contains(k, "ESC BACK") && !strings.Contains(v, "НАЗАД") {
			t.Errorf("%q -> %q lost BACK", k, v)
		}
	}
	for _, k := range []string{"GALAXY", "GALACTIC CHART"} {
		if T(k) == k {
			t.Errorf("%q is not translated, and screens compose it with a number", k)
		}
	}
}
