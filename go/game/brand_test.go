package game

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"heliosdrift/galaxy"
)

// The game must not carry names owned by Frontier Developments. This scans
// every Go string literal in the tree, so a stray "Cobra Mk III" in a new
// message fails the build before it ever ships.
func TestNoTrademarkedNames(t *testing.T) {
	bad := regexp.MustCompile(`"[^"\n]*(CONSTRICTOR|THARGOID|THARGON|COBRA|ELITE)[^"\n]*"`)
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	var hits []string
	filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}
		src, err := os.ReadFile(p)
		if err != nil {
			return nil
		}
		for _, m := range bad.FindAllString(string(src), -1) {
			hits = append(hits, filepath.Base(p)+": "+m)
		}
		return nil
	})
	if len(hits) > 0 {
		t.Fatalf("foreign names found:\n%s", strings.Join(hits, "\n"))
	}
}

// The renaming itself: what the player actually reads.
func TestRenamedContent(t *testing.T) {
	if galaxy.Ranks[len(galaxy.Ranks)-1] != "LEGEND" {
		t.Fatalf("top rank = %q", galaxy.Ranks[len(galaxy.Ranks)-1])
	}
	var kestrel bool
	for _, s := range ShipYard {
		if strings.Contains(s.Name, "KESTREL") {
			kestrel = true
		}
	}
	if !kestrel {
		t.Fatal("no KESTREL in the shipyard")
	}
	g := newGameMode(ModeClassic)
	g.Cmd.Kills = 500
	g.Cmd.Mission = Mission{Kind: MissionConstrictor, System: 1, Galaxy: 0}
	s := g.MissionStatus()
	if !strings.Contains(s, "SHRIKE") || strings.Contains(s, "CONSTRICTOR") {
		t.Fatalf("mission status = %q", s)
	}
}
