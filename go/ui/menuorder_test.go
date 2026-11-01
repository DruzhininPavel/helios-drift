package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
)

// The board order: sheet, trade, work, ship, log, save, chart, way out.
func TestStationMenuOrder(t *testing.T) {
	lang.Set("en")
	defer lang.Set("ru")
	g := game.New()
	g.NewGame("TEST")
	g.State = "docked"
	g.Screen = "launch"
	g.Cmd.Legal = 400
	g.World.Sys.Government = 3
	g.World.Sys.Tech = 12
	items := menuItems()
	t.Logf("menu: %v", items)
	if items[0] != "COMMANDER STATUS" {
		t.Errorf("first row is %q, want the commander sheet", items[0])
	}
	if items[len(items)-1] != "LEAVE STATION" {
		t.Errorf("last row is %q, want the way out", items[len(items)-1])
	}
	if items[len(items)-2] != "GALACTIC CHART" {
		t.Errorf("second from bottom is %q, want the chart", items[len(items)-2])
	}
}
