package lang

import (
	"strings"
	"testing"

	"heliosdrift/render"
)

// Every character in a translation must exist in the 5x7 font, otherwise it
// renders as '?'.
func TestTranslationsRenderable(t *testing.T) {
	for k, v := range ru {
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
	for k, v := range ru {
		if v != strings.ToUpper(v) {
			t.Errorf("%q -> %q is not uppercase", k, v)
		}
	}
}

// Long labels must still fit the panels they are drawn in.
func TestTranslationsFit(t *testing.T) {
	// 17 glyphs is the widest value the chart's system panel can show.
	for _, key := range []string{"MAINLY AGRICULTURAL", "CONFEDERACY", "COMMANDER STATUS"} {
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
