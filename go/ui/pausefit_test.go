package ui

import (
	"testing"

	"heliosdrift/lang"
	"heliosdrift/render"
)

// The pause list has two columns in a 312 pixel plate: a long name used to run
// off the right edge and lose its last letters.
func TestPauseRowsFitThePlate(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	rows := pauseRows()
	perColumn := (len(rows) + 1) / 2
	for i, r := range rows {
		xk := 12 + (i/perColumn)*150
		keyRight := xk + render.TextWidth(r[0], 1)
		valRight := xk + 46 + render.TextWidth(r[1], 1)
		if keyRight > xk+44 {
			t.Errorf("key %q is %d px wide, the gutter is 44", r[0], keyRight-xk)
		}
		if valRight > 306 {
			t.Errorf("%q ends at x=%d, the plate ends at 312", r[1], valRight)
		}
	}
	if perColumn*11+30 > 194 {
		t.Errorf("%d rows per column do not fit the plate", perColumn)
	}
}
