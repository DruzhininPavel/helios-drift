package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/input"
	"heliosdrift/lang"
)

func TestMain(m *testing.M) {
	lang.ReadOnly = true
	m.Run()
}

func newTitle() {
	g := game.New()
	g.NewGame("TEST")
	g.State = "title"
	g.Screen = "launch"
	S.Title = 0
	S.Settings = 0
	lang.Set("en")
}

// The settings screen must be reachable and closable from the title screen,
// where the game state stays "title" while it is open.
func TestSettingsFromTitle(t *testing.T) {
	newTitle()
	OpenSettings()
	if game.G.Screen != "settings" {
		t.Fatalf("settings did not open, screen=%q", game.G.Screen)
	}
	if !Key("Escape") {
		t.Fatal("escape was not consumed")
	}
	if game.G.Screen == "settings" {
		t.Fatal("escape did not leave settings")
	}
	if game.G.State != "title" {
		t.Fatalf("state=%q, want title", game.G.State)
	}
}

func TestSettingsFromDock(t *testing.T) {
	newTitle()
	game.G.State = "docked"
	game.G.Screen = "launch"
	OpenSettings()
	Key("Escape")
	if game.G.Screen != "launch" || game.G.State != "docked" {
		t.Fatalf("after escape screen=%q state=%q, want launch/docked", game.G.Screen, game.G.State)
	}
}

func TestSettingsLanguageToggle(t *testing.T) {
	newTitle()
	OpenSettings()
	S.Settings = 0
	Key("Enter")
	if !lang.IsRU() {
		t.Fatal("language did not switch to ru")
	}
	Key("Enter")
	if lang.IsRU() {
		t.Fatal("language did not switch back to en")
	}
	lang.Set("en")
}

func TestTitleHasQuit(t *testing.T) {
	newTitle()
	S.Title = 3
	Key("Enter")
	if !QuitRequested() {
		t.Fatal("QUIT did not ask to leave")
	}
}

func TestSettingsNavigation(t *testing.T) {
	newTitle()
	OpenSettings()
	S.Settings = 0
	Key("ArrowDown")
	if S.Settings != 1 {
		t.Fatalf("down -> %d, want 1", S.Settings)
	}
	Key("ArrowUp")
	if S.Settings != 0 {
		t.Fatalf("up -> %d, want 0", S.Settings)
	}
	Key("ArrowUp")
	if S.Settings != 4 {
		t.Fatalf("up from first -> %d, want 4 (wrap)", S.Settings)
	}
}

func TestSettingsVolumeAndMute(t *testing.T) {
	newTitle()
	OpenSettings()
	lang.Cfg.Volume = 0.8
	S.Settings = 2
	Key("ArrowRight")
	if lang.Cfg.Volume == 0.8 {
		t.Fatal("volume did not change")
	}
	S.Settings = 1
	before := lang.Cfg.Muted
	Key("Enter")
	if lang.Cfg.Muted == before {
		t.Fatal("mute did not toggle")
	}
	lang.Cfg.Muted = false
	lang.Cfg.Volume = 0.8
}

// The station menu offers the contract screen once the rank is high enough.
func TestMissionScreenFromMenu(t *testing.T) {
	newTitle()
	game.G.NewGameMode(game.ModeClassic, "TEST") // 1984 contracts live here
	game.G.State = "docked"
	game.G.Screen = "launch"
	game.G.Cmd.Kills = 40
	game.G.DockedForShots()

	Key("9")
	if game.G.Screen != "mission" {
		t.Fatalf("screen = %q, want mission", game.G.Screen)
	}
	Key("Enter")
	if game.G.Screen != "launch" {
		t.Fatalf("screen = %q, want launch after accepting", game.G.Screen)
	}
	if game.G.Cmd.Mission.Kind != game.MissionConstrictor {
		t.Fatalf("mission = %q, want constrictor", game.G.Cmd.Mission.Kind)
	}
}

// Esc backs out without accepting.
func TestMissionScreenDecline(t *testing.T) {
	newTitle()
	game.G.NewGameMode(game.ModeClassic, "TEST")
	game.G.State = "docked"
	game.G.Screen = "launch"
	game.G.Cmd.Kills = 40
	game.G.DockedForShots()
	Key("9")
	Key("Escape")
	if game.G.Screen != "launch" {
		t.Fatalf("screen = %q, want launch", game.G.Screen)
	}
	if game.G.Cmd.Mission.Kind != "" {
		t.Fatal("mission was accepted on decline")
	}
}

// The controls screen lists every action plus a reset row.
func TestControlsScreen(t *testing.T) {
	newTitle()
	game.G.State = "docked"
	game.G.Screen = "settings"
	S.Settings = 3

	Key("Enter")
	if game.G.Screen != "controls" {
		t.Fatalf("screen = %q, want controls", game.G.Screen)
	}
	Key("ArrowDown")
	if S.Control != 1 {
		t.Fatalf("selection = %d, want 1", S.Control)
	}
	rows := len(input.Actions()) + 2
	S.Control = 0
	Key("ArrowUp")
	if S.Control != rows-1 {
		t.Fatalf("up from first = %d, want %d", S.Control, rows-1)
	}
	Key("ArrowDown")
	if S.Control != 0 {
		t.Fatalf("down from last = %d, want 0", S.Control)
	}

	// reset row restores the stock layout
	input.Reset()
	want := input.KeyDesc(input.Fire)
	input.SetKey(input.Fire, 999)
	S.Control = len(input.Actions())
	if S.Control != len(input.Actions()) {
		t.Fatalf("selection = %d, want the reset row %d", S.Control, len(input.Actions()))
	}
	Key("Enter")
	if input.KeyDesc(input.Fire) != want {
		t.Fatalf("reset left %q, want %q", input.KeyDesc(input.Fire), want)
	}

	Key("Escape")
	if game.G.Screen != "settings" {
		t.Fatalf("screen = %q, want settings", game.G.Screen)
	}
}

// The rule set is chosen when a commander is created, from the title screen.
func TestModeScreenStartsGame(t *testing.T) {
	newTitle()
	S.Title = 0
	Key("Enter")
	if game.G.Screen != "mode" {
		t.Fatalf("screen = %q, want mode", game.G.Screen)
	}
	Key("ArrowDown")
	if S.Mode != 1 {
		t.Fatalf("selection = %d, want advanced", S.Mode)
	}
	Key("ArrowUp")
	Key("Enter")
	if !game.G.Classic() {
		t.Fatalf("mode = %q, want classic", game.G.Cmd.Mode)
	}
	if game.G.State != "docked" || game.G.World == nil {
		t.Fatal("game did not start docked")
	}
	for _, it := range menuItems() {
		if it == "SHIPYARD" || it == "CONTRACTS" {
			t.Fatalf("classic menu offers %q", it)
		}
	}
}

func TestAdvancedMenuHasShipyard(t *testing.T) {
	newTitle()
	game.G.NewGameMode(game.ModeAdvanced, "TEST")
	found := false
	for _, it := range menuItems() {
		if it == "SHIPYARD" {
			found = true
		}
	}
	if !found {
		t.Fatal("advanced menu lost the shipyard")
	}
}

// The campaign briefing replaces the 1984 contract in Advanced mode.
func TestCampaignScreen(t *testing.T) {
	newTitle()
	game.G.NewGameMode(game.ModeAdvanced, "TEST")
	game.G.State = "docked"
	game.G.Screen = "launch"
	if !game.G.CampaignOffer() {
		t.Fatal("advanced has no campaign stage to offer")
	}
	Key("9") // MISSION row hotkey
	found := false
	for _, it := range menuItems() {
		if it == "MISSION" {
			found = true
		}
	}
	if !found {
		t.Fatal("menu lost the MISSION row")
	}
	openMission()
	if game.G.Screen != "mission" {
		t.Fatalf("screen = %q", game.G.Screen)
	}
	Key("Enter")
	if !game.G.Cmd.Campaign.Active {
		t.Fatal("campaign stage was not accepted")
	}
	if game.G.Cmd.Campaign.Stage != 0 {
		t.Fatalf("stage = %d", game.G.Cmd.Campaign.Stage)
	}
}

// The pad test has to be reachable from the settings screen, and Enter on it
// must actually open it — otherwise it is invisible in the app.
func TestPadTestReachable(t *testing.T) {
	newTitle()
	OpenSettings()
	S.Settings = 3 // CONTROLS
	Key("Enter")
	if game.G.Screen != "controls" {
		t.Fatalf("screen = %q, want controls", game.G.Screen)
	}
	S.Control = len(input.Actions()) + 1 // PAD TEST is the last row
	Key("Enter")
	if game.G.Screen != "padtest" {
		t.Fatalf("screen = %q, want padtest", game.G.Screen)
	}
	Key("Escape")
	if game.G.Screen != "controls" {
		t.Fatalf("after escape screen = %q, want controls", game.G.Screen)
	}
}
