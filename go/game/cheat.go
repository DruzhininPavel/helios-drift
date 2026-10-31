package game

import "heliosdrift/lang"

// Cheats exist for testing, not for playing: a tester who wants to see chapter
// eight should not trade food for an hour first. Both are plain keys and both
// say what they did, so a cheat cannot be mistaken for a bug.

// CheatMoney adds credits.
func (g *Game) CheatMoney(amount float64) {
	g.Cmd.Cash += amount
	g.Message(lang.T("CHEAT")+": +"+itoa(int(amount))+" CR", "good")
}

// ToggleCheatShield turns off damage from guns, missiles and debris. Ramming a
// station, a planet, a sun or a ship still kills: the point of the cheat is to
// skip firefights, not to fly through walls.
func (g *Game) ToggleCheatShield() bool {
	g.CheatShield = !g.CheatShield
	if g.CheatShield {
		g.Message(lang.T("CHEAT SHIELD ON"), "good")
	} else {
		g.Message(lang.T("CHEAT SHIELD OFF"), "info")
	}
	return g.CheatShield
}
