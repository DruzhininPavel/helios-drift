package game

import "testing"

// The four station faces: each met once, never in Classic, never twice the same
// day, and none of them has anything to do with the campaign.
func TestAcquaintances(t *testing.T) {
	g := New()
	g.NewGame("MEET")
	g.State = "docked"
	g.Screen = "launch"
	met := 0
	for i := 0; i < 400 && met < 4; i++ {
		g.Launch()
		if g.Screen == "talk" {
			if _, ok := map[string]bool{"a1_bar": true, "a2_dock": true, "a3_nurse": true, "a4_old": true}[g.Talk]; !ok {
				t.Fatalf("launched into %q", g.Talk)
			}
			met++
			g.Cmd.Flags["MET_"+map[string]string{
				"a1_bar": "BAR", "a2_dock": "DOCK", "a3_nurse": "NURSE", "a4_old": "OLD",
			}[g.Talk]] = true
			g.EndTalk()
		}
		g.State = "docked"
		g.Screen = "launch"
		g.Cmd.Day++
	}
	if met == 0 {
		t.Fatal("nobody ever had anything to say")
	}
	// and once everyone has been met, the space stays quiet
	for i := 0; i < 50; i++ {
		g.Launch()
		if g.Screen == "talk" {
			t.Fatalf("met %q twice", g.Talk)
		}
		g.State = "docked"
		g.Screen = "launch"
		g.Cmd.Day++
	}
}

// Classic has no small talk: a 1984 commander left the station in silence.
func TestClassicHasNoSmallTalk(t *testing.T) {
	g := New()
	g.NewGameMode(ModeClassic, "MEET")
	g.State = "docked"
	g.Screen = "launch"
	for i := 0; i < 60; i++ {
		g.Launch()
		if g.Screen == "talk" {
			t.Fatalf("classic opened %q", g.Talk)
		}
		g.State = "docked"
		g.Screen = "launch"
		g.Cmd.Day++
	}
}
