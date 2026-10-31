package ui

import (
	"testing"

	"heliosdrift/lang"
	"heliosdrift/render"
)

func TestHelpLinesFit(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	for _, p := range helpPages {
		for _, l := range p.Lines {
			if l == "" {
				continue
			}
			w := render.TextWidth(l, 1)
			if w > 270 {
				t.Errorf("%s: %q is %d px wide", p.Title, l, w)
			}
		}
	}
}
