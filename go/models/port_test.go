package models

import (
	"testing"

	"heliosdrift/render"
)

// A docking port that reads as a hole is the bug: render the hull over a key
// colour at several angles and fail if background survives inside the
// silhouette.
func TestPortNotAHole(t *testing.T) {
	key := render.RGB(0xff00ff)
	// the game's own resolution: a crack that cannot be seen there does not count
	const W, H = 320, 224
	for _, yaw := range []float64{0, 0.26, 0.52, -0.52, 1.05} {
		fb := render.NewFB(W, H)
		fb.Clear(key)
		cam := render.NewCamera(W, H, 55)
		cam.CX, cam.CY = W/2, H/2
		m := Meshes["station"]
		render.DrawMesh(fb, cam, m, render.V3(0, 0, -m.Radius*2.2), render.RotY(yaw), 90000)

		// flood fill the background from the border; anything unreachable is a hole
		seen := make([]bool, W*H)
		stack := []int{}
		push := func(x, y int) {
			i := y*W + x
			if x < 0 || y < 0 || x >= W || y >= H || seen[i] || fb.Buf[i] != key {
				return
			}
			seen[i] = true
			stack = append(stack, i)
		}
		for x := 0; x < W; x++ {
			push(x, 0)
			push(x, H-1)
		}
		for y := 0; y < H; y++ {
			push(0, y)
			push(W-1, y)
		}
		for len(stack) > 0 {
			i := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			x, y := i%W, i/W
			push(x+1, y)
			push(x-1, y)
			push(x, y+1)
			push(x, y-1)
		}
		holes := 0
		for i, v := range fb.Buf {
			if v == key && !seen[i] {
				holes++
			}
		}
		// a hairline seam may survive; the bug this guards against filled the
		// whole port and would light up hundreds of pixels
		if holes > 4 {
			t.Errorf("yaw %.2f: %d background pixels enclosed inside the hull", yaw, holes)
		}
	}
}

// The port must read as an interior, not as a hole punched through to space.
// A mouth as dark as the void is what looked like a floating black rectangle.
func TestPortIsNotBlack(t *testing.T) {
	const W, H = 320, 224
	fb := render.NewFB(W, H)
	fb.Clear(0) // pure black stands in for deep space
	cam := render.NewCamera(W, H, 55)
	cam.CX, cam.CY = W/2, H/2
	m := Meshes["station"]
	dist := m.Radius * 2.2
	render.DrawMesh(fb, cam, m, render.V3(0, 0, -dist), render.RotY(0), 90000)

	// project the middle of the port
	v := cam.ToView(render.V3(0, 0, -dist+StationSize+1))
	if v[2] >= 0 {
		t.Fatal("port is behind the camera")
	}
	p := cam.ProjectView(v)
	cx, cy := int(p.X), int(p.Y)

	sum, n := 0, 0
	for y := cy - 3; y <= cy+3; y++ {
		for x := cx - 12; x <= cx+12; x++ {
			if x < 0 || y < 0 || x >= W || y >= H {
				continue
			}
			c := fb.Buf[y*W+x]
			sum += int(c&0xff) + int((c>>8)&0xff) + int((c>>16)&0xff)
			n++
		}
	}
	if n == 0 {
		t.Fatal("port is off screen")
	}
	avg := sum / n / 3
	if avg < 20 {
		t.Fatalf("port renders as black (avg %d): it will read as a hole", avg)
	}
}

// The port must be a shaft, not a decal: there have to be faces that run from
// the hull surface inward.
func TestPortHasShaft(t *testing.T) {
	m := Meshes["station"]
	front := float64(StationSize)
	walls := 0
	deepest := front
	for _, f := range m.Faces {
		minz, maxz := 1e9, -1e9
		for _, p := range f.P {
			if p[2] < minz {
				minz = p[2]
			}
			if p[2] > maxz {
				maxz = p[2]
			}
		}
		if maxz > front-1 && minz < front-20 {
			walls++
		}
		if minz < deepest {
			deepest = minz
		}
	}
	if walls < 4 {
		t.Fatalf("%d shaft walls, want at least 4: the port is flat", walls)
	}
	if front-deepest < 40 {
		t.Fatalf("shaft is only %.0f units deep", front-deepest)
	}
}
