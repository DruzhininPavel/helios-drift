package ui

import (
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
)

// The cursor has to walk the rows in the order they are drawn: it used to move
// through the catalogue, so down jumped from a weapon to a shield two shelves on.
func TestEquipCursorFollowsTheRows(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	g := game.New()
	g.NewGame("EQ")
	g.State = "docked"
	g.Screen = "equip"
	g.Docked()
	g.World.Sys.Tech = 12
	S.Equip = 0
	items := equipList(g)
	if len(items) < 4 {
		t.Fatalf("only %d items on the shelf", len(items))
	}
	// the list is grouped: a shelf never starts twice
	seen := map[string]bool{}
	group := ""
	for _, it := range items {
		gr := equipGroup(it.ID)
		if gr != group {
			if seen[gr] {
				t.Fatalf("shelf %s appears twice: the list is not grouped", gr)
			}
			seen[gr] = true
			group = gr
		}
	}
	// and the cursor visits them in order
	var visited []string
	for i := 0; i < len(items); i++ {
		visited = append(visited, equipGroup(items[S.Equip].ID))
		Key("ArrowDown")
	}
	for i := 1; i < len(visited); i++ {
		if visited[i] != visited[i-1] && i < len(visited)-1 {
			// a shelf may change once, going forward
			for j := i + 1; j < len(visited); j++ {
				if visited[j] == visited[i-1] {
					t.Fatalf("the cursor went back to %s after leaving it: %v", visited[i-1], visited)
				}
			}
		}
	}
}
