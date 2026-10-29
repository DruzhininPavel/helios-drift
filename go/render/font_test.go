package render

import "testing"

// Every Cyrillic capital used in the UI must have a real glyph, not '?'.
func TestCyrillicCoverage(t *testing.T) {
	alphabet := "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
	q := glyphs['?']
	for _, r := range alphabet {
		gl, ok := glyphs[r]
		if !ok {
			t.Errorf("missing glyph for %c", r)
			continue
		}
		nonEmpty := false
		for _, row := range gl {
			if row != 0 {
				nonEmpty = true
			}
		}
		if !nonEmpty {
			t.Errorf("blank glyph for %c", r)
		}
		if gl == q {
			t.Errorf("%c falls back to '?'", r)
		}
	}
	// no two distinct letters may be pixel-identical to a different Latin letter
	// unless intended (the aliased set)
	seen := map[[7]uint8]rune{}
	for _, r := range alphabet {
		gl := glyphs[r]
		if prev, ok := seen[gl]; ok && r != prev {
			t.Logf("%c and %c share a bitmap", prev, r)
		}
		seen[gl] = r
	}
}

func TestCyrillicTextWidth(t *testing.T) {
	if w := TextWidth("ПРИВЕТ", 1); w != 6*6-1 {
		t.Errorf("width of ПРИВЕТ = %d, want %d", w, 6*6-1)
	}
}
