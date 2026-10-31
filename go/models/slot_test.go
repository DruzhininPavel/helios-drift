package models

import "testing"

// Every station's mouth has to sit on the docking plane the game maths uses.
// The rock hermit's collar once ran thirty units past the hull, so a pilot
// could fly into the opening and nothing happened: the trigger was behind it.
func TestStationMouthsAreFlush(t *testing.T) {
	w := StationSlotW * StationSize
	h := StationSlotH * StationSize
	for _, name := range []string{"station", "dodeca", "hermit"} {
		m := Meshes[name]
		if m == nil {
			t.Fatalf("%s missing", name)
		}
		mouth := -1e9
		for _, f := range m.Faces {
			for _, v := range f.P {
				if v[0] > -w && v[0] < w && v[1] > -h && v[1] < h && v[2] > mouth {
					mouth = v[2]
				}
			}
		}
		// nothing may stand in the mouth beyond the docking plane: the
		// coriolis slot is recessed to 143, so only the upper bound matters
		if mouth > StationSize+2 {
			t.Errorf("%s: slot geometry reaches %.0f, past the docking plane %.0f",
				name, mouth, StationSize)
		}
	}
}
