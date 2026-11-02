package lang

import (
	"regexp"
	"strings"
	"testing"

	"heliosdrift/render"
)

// Every character in a translation must exist in the 5x7 font, otherwise it
// renders as '?'.
func TestTranslationsRenderable(t *testing.T) {
	brace := regexp.MustCompile(`\{[A-Z]+\}`)
	for k, v := range ru {
		// {PLACEHOLDER} is filled with the pilot's key before it is ever drawn
		v = brace.ReplaceAllString(v, "ESC")
		for _, r := range v {
			if r == ' ' || r == '\n' {
				continue
			}
			if !render.HasGlyph(r) {
				t.Errorf("%q -> %q: no glyph for %q", k, v, r)
			}
		}
	}
}

// The UI font is uppercase only; lowercase would fall back to '?'.
func TestTranslationsUppercase(t *testing.T) {
	brace := regexp.MustCompile(`\{[A-Z]+\}`)
	for k, v := range ru {
		// {PLACEHOLDER} is filled with the pilot's key before it is ever drawn
		v = brace.ReplaceAllString(v, "ESC")
		if v != strings.ToUpper(v) {
			t.Errorf("%q -> %q is not uppercase", k, v)
		}
	}
}

var galaxyEconomies = []string{"RICH INDUSTRIAL", "AVERAGE INDUSTRIAL", "POOR INDUSTRIAL",
	"MAINLY INDUSTRIAL", "MAINLY AGRICULTURAL", "RICH AGRICULTURAL",
	"AVERAGE AGRICULTURAL", "POOR AGRICULTURAL"}

var galaxyGovernments = []string{"ANARCHY", "FEUDAL", "MULTI-GOVERNMENT", "DICTATORSHIP",
	"COMMUNIST", "CONFEDERACY", "DEMOCRACY", "CORPORATE STATE"}

// Long labels must still fit the panels they are drawn in.
func TestTranslationsFit(t *testing.T) {
	// 17 glyphs is the widest value the chart's system panel can show.
	for _, key := range append([]string{"CONFEDERACY", "COMMANDER STATUS"},
		append(galaxyEconomies, galaxyGovernments...)...) {
		if w := render.TextWidth(ru[key], 1); w > 17*6 {
			t.Errorf("%q too wide: %dpx", key, w)
		}
	}
}

func TestTranslateLineFallback(t *testing.T) {
	Set("ru")
	defer Set("en")
	if got := TranslateLine("ARRIVED: ESALETRA"); got != "ПРИБЫТИЕ: ESALETRA" {
		t.Errorf("word-wise translation failed: %q", got)
	}
	if got := TranslateLine("ZZZ UNKNOWN"); got != "ZZZ UNKNOWN" {
		t.Errorf("unknown text should pass through, got %q", got)
	}
	Set("en")
	if got := TranslateLine("CASH"); got != "CASH" {
		t.Errorf("english should be identity, got %q", got)
	}
}
