package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// The prologue's last line was cut in half by a fixed six-row window while half
// the panel stood empty. Every line of every intro screen must be on the first
// page when the text fits the panel.
func TestIntroTextIsNotCutOff(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	fb := render.NewFB(320, 224)
	g := game.New()
	g.NewGame("FIT")
	for _, id := range []string{"intro_who", "intro_where", "intro_what"} {
		g.TalkStart(id)
		S.TalkScroll = 0
		Draw(fb, 1.0/60.0)
		n := game.DialogueNode(id)
		rows := 0
		for _, ln := range g.TalkText(n) {
			rows += len(wrapText(lang.T(ln), 48)) + 1
		}
		opts := len(g.VisibleOptions(n))
		visible := (186 - opts*11 - 6 - 32) / 10
		if rows > visible {
			t.Fatalf("%s: %d rows in a %d-row window", id, rows, visible)
		}
		g.EndTalk()
	}
}
