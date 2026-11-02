package ui

import (
	"testing"
	"unicode/utf8"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// The handbook wraps to the plate and scrolls, so a long line is fine - what is
// not fine is a wrapped row that does not fit, or a page whose tail cannot be
// reached.
func TestHelpPagesWrapAndScroll(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	fb := render.NewFB(320, 224)
	for _, p := range helpPages {
		rows := 0
		for _, l := range p.Lines {
			if l == "" {
				rows++
				continue
			}
			wrapped := wrapText(lang.T(l), 44)
			rows += len(wrapped)
			for _, w := range wrapped {
				if utf8.RuneCountInString(w) > 44 {
					t.Errorf("%s: %q wraps to %d characters", p.Title, w, utf8.RuneCountInString(w))
				}
				if width := render.TextWidth(w, 1); width > 276 {
					t.Errorf("%s: %q is %d px wide", p.Title, w, width)
				}
			}
		}
		if rows == 0 {
			t.Errorf("%s: empty page", p.Title)
		}
	}
	// and every page draws, scrolled to its end
	g := game.New()
	g.NewGame("HELP")
	g.State = "docked"
	g.Screen = "help"
	for i := range helpPages {
		S.Help = i
		S.HelpScroll = 0
		Draw(fb, 1.0/60.0)
		for k := 0; k < 20; k++ {
			Key("PageDown")
		}
		Draw(fb, 1.0/60.0)
	}
}
