package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
)

// The station board lists things you do at a station. Settings belong to the
// title screen, where a player goes to change how the game looks and sounds.
func TestSettingsStayOffTheStationBoard(t *testing.T) {
	lang.Set("en")
	defer lang.Set("ru")
	g := game.New()
	g.NewGame("TEST")
	g.State = "docked"
	g.Screen = "launch"
	for _, it := range menuItems() {
		if it == "SETTINGS" {
			t.Fatal("settings are still on the station board")
		}
	}
	found := false
	for _, it := range TitleItems() {
		if it == "SETTINGS" {
			found = true
		}
	}
	if !found {
		t.Fatal("settings vanished from the title menu too")
	}
}
