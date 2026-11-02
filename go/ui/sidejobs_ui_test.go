package ui

import (
	"strings"
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// A running side job shows up in the summary with its current step.
func TestSideJobAppearsInTheSummary(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	g := game.New()
	g.NewGame("SIDEUI")
	g.State = "docked"
	g.Screen = "journal"
	g.Cmd.SideStep = map[string]int{"clinic": 1}
	g.Cmd.SideCount = map[string]int{"clinic": 0}
	g.Cmd.Flags["SIDE_START_clinic"] = true

	fb := render.NewFB(320, 224)
	Draw(fb, 1.0/60.0)
	if Hint == "" {
		t.Fatal("the summary drew nothing")
	}
	if !strings.Contains(lang.T("SIDE WORK"), lang.T("SIDE WORK")) {
		t.Fatal("no heading")
	}
	j := game.SideJobByID("clinic")
	if line := g.SideLine(j); line == "" {
		t.Fatal("a running job has no line")
	}
}
