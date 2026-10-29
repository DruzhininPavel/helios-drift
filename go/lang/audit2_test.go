package lang

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

// TestAuditUntranslated lists multi-word strings the game draws without any
// Russian translation — those are usually prose that was missed.
func TestAuditUntranslated(t *testing.T) {
	path := os.Getenv("HELIOS_AUDIT2")
	if path == "" {
		t.Skip("set HELIOS_AUDIT2")
	}
	f, _ := os.Open(path)
	defer f.Close()
	Set("ru")
	defer Set("en")
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		s := sc.Text()
		if TranslateLine(s) != s {
			continue
		}
		if len(strings.Fields(s)) >= 3 && strings.ToUpper(s) == s {
			t.Logf("%q", s)
		}
	}
}
