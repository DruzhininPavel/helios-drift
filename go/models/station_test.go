package models

import (
	"testing"

	"heliosdrift/render"
)

// The docking port used to read as a hole through the hull: the tunnel walls
// were culled from the mouth side, so stars showed through. Render face-on and
// make sure no background pixel survives inside the silhouette at the port.
func TestStationPortSealed(t *testing.T) {
	key := render.RGB(0xff00ff)
	for _, name := range []string{"station", "dodeca"} {
		const W, H = 400, 400
		fb := render.NewFB(W, H)
		fb.Clear(key)
		cam := render.NewCamera(W, H, 55)
		cam.CX, cam.CY = W/2, H/2
		mesh := Meshes[name]
		if mesh == nil {
			t.Fatalf("%s: missing mesh", name)
		}
		// face-on: the station's +Z face points at the camera
		render.DrawMesh(fb, cam, mesh, render.V3(0, 0, -mesh.Radius*2.2), render.RotY(0), 90000)

		leaks, rows := 0, 0
		for y := H/2 - 12; y <= H/2+12; y++ {
			first, last := -1, -1
			for x := 0; x < W; x++ {
				if fb.Buf[y*W+x] != key {
					if first < 0 {
						first = x
					}
					last = x
				}
			}
			if first < 0 || last <= first {
				continue
			}
			rows++
			for x := first; x <= last; x++ {
				if fb.Buf[y*W+x] == key {
					leaks++
				}
			}
		}
		if rows < 10 {
			t.Fatalf("%s: hull did not render face-on (%d rows)", name, rows)
		}
		if leaks > 0 {
			t.Errorf("%s: %d background pixels inside the hull across the port", name, leaks)
		}
	}
}
