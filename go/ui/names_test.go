package ui

import (
	"testing"

	"heliosdrift/lang"
	"heliosdrift/render"
)

// A system called QUIT used to read "ВЫХОД" on the chart while contracts kept
// the English word: the same place under two names. Names are data, so they
// must look identical in both languages.
func TestNamesLookTheSameInBothLanguages(t *testing.T) {
	defer lang.Set("en")
	for _, name := range []string{"QUIT", "DIST", "FLYING", "THILQU", "ITILIL"} {
		lang.Set("en")
		en := render.NewFB(200, 16)
		drawData(en, name, 2, 2, 0xffffff)
		lang.Set("ru")
		ru := render.NewFB(200, 16)
		drawData(ru, name, 2, 2, 0xffffff)
		for i := range en.Buf {
			if en.Buf[i] != ru.Buf[i] {
				t.Fatalf("name %q is drawn differently in Russian", name)
			}
		}
	}
}

// And the dictionary really would have changed them, otherwise the test above
// proves nothing.
func TestDictionaryWouldHaveRenamedThem(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	for _, name := range []string{"QUIT", "DIST"} {
		if lang.T(name) == name {
			t.Fatalf("%q is not in the dictionary, pick another sample", name)
		}
	}
}
