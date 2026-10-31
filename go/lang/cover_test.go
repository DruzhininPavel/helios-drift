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
