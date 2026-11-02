package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// Every hint line has to fit the 320 pixel footer on its own: the talk screen
// once carried a 354 pixel one and it was cut at both edges.
func TestHintsFitTheFooter(t *testing.T) {
	for _, l := range []string{"ru", "en"} {
		lang.Set(l)
		render.Translate = lang.TranslateLine
		fb := render.NewFB(320, 224)
		g := game.New()
		g.NewGame("HINT")
		g.State = "docked"
		g.Screen = "launch"
		g.Docked()
		check := func(what string) {
			text := lang.TranslateLine(Hint)
			if w := render.TextWidth(text, 1); w > 314 {
				t.Errorf("%s (%s): hint %q is %d px wide", what, l, text, w)
			}
		}
		check("station")
		for _, sc := range []string{"market", "contracts", "equip", "status", "log", "journal", "shipyard", "slots"} {
			g.Screen = sc
			Draw(fb, 1.0/60)
			check(sc)
		}
		for i := range helpPages {
			S.Help = i
			S.HelpScroll = 99
			g.Screen = "help"
			Draw(fb, 1.0/60)
			check("help")
		}
		g.Screen = "launch"
		g.Launch()
		g.Overlay = "pause"
		Draw(fb, 1.0/60)
		check("pause")
		g.Overlay = ""
		if g.TalkStart("c1_disp") {
			Draw(fb, 1.0/60)
			check("talk")
			g.EndTalk()
		}
		g.State = "docked"
		g.Screen = "launch"
		g.StartEnding("combine")
		Draw(fb, 1.0/60)
		check("ending")
		g.Player.Invulnerable = 0
		g.DamagePlayer(9999, true, "laser")
		Draw(fb, 1.0/60)
		check("death")
	}
	render.Translate = nil
	lang.Set("ru")
}
