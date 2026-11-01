package game

import (
	"strings"
	"testing"

	"heliosdrift/lang"
	"heliosdrift/render"
)

// Every line the player can read in a conversation must exist in both
// languages, and must be drawable: the font has no lower case and only its own
// glyph set.
func TestDialogueTextIsComplete(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	for _, id := range DialogueNodeIDs() {
		n := DialogueNode(id)
		texts := append([]string{}, n.Lines...)
		if n.Speaker != "" {
			texts = append(texts, n.Speaker)
		}
		for _, o := range n.Options {
			texts = append(texts, o.Text)
		}
		for _, s := range texts {
			if strings.TrimSpace(s) == "" {
				continue
			}
			if lang.T(s) == s {
				t.Errorf("%s: %q has no Russian", id, s)
			}
			if strings.ToUpper(s) != s {
				t.Errorf("%s: %q is not capitals", id, s)
			}
			// the talk screen wraps, so a long line is fine; a line that wraps
			// into fifty rows is not
			if w := render.TextWidth(s, 1); w > 276*6 {
				t.Errorf("%s: %q is absurdly long (%d px)", id, s, w)
			}
		}
	}
}
