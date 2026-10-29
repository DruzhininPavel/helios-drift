package lang

import (
	"bufio"
	"os"
	"regexp"
	"testing"
	"unicode"
)

// TestAuditDump lists strings that come out half-translated: a Russian phrase
// with English data stuck in the middle (or the reverse).
func TestAuditDump(t *testing.T) {
	path := os.Getenv("HELIOS_AUDIT")
	if path == "" {
		t.Skip("set HELIOS_AUDIT to the dump file")
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	cap := regexp.MustCompile(`\b[A-Z]{2,}\b`)
	sc := bufio.NewScanner(f)
	Set("ru")
	defer Set("en")
	for sc.Scan() {
		s := sc.Text()
		out := TranslateLine(s)
		if out == s {
			continue
		}
		cyr, lat := false, false
		for _, r := range out {
			if unicode.Is(unicode.Cyrillic, r) {
				cyr = true
			}
			if r >= 'A' && r <= 'Z' {
				lat = true
			}
		}
		if cyr && lat && cap.MatchString(out) {
			t.Logf("%q\n    -> %q", s, out)
		}
	}
}
