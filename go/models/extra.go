// Extra hulls that were still missing from the 1984 roster, plus the
// dodecahedron station. Registered from init so models.go stays untouched.
package models

import (
	"math"

	"heliosdrift/render"
)

func init() {
	addExtraShips(Meshes)
}

func addExtraShips(m map[string]*render.Mesh) {
	// Gecko - small delta pirate
	{
		b := NewBuilder()
		loft(b, []Stn{{-20, 0.7, 0.6, 0}, {-10, 4.5, 1.8, 0}, {4, 8.0, 2.8, 0}, {15, 5.5, 2.2, 0}, {20, 2.6, 1.4, 0}}, "teal", 6)
		prism(b, [][2]float64{{4, -2}, {18, 14}, {17, 19}, {5, 15}}, 1.2, "teal")
		prism(b, [][2]float64{{-4, -2}, {-18, 14}, {-17, 19}, {-5, 15}}, 1.2, "teal")
		box3(b, 0, 2.2, 0, 2.6, 1.4, 9, "alien", &FOpts{Shade: 4.0, Fixed: true})
		mesh := b.Mesh()
		mesh.Exhaust = []V{{0, 0, 20}}
		m["gecko"] = mesh
	}
	// Cobra Mk I - the older, smaller Cobra
	{
		b := NewBuilder()
		loft(b, []Stn{{-26, 1.0, 0.9, 0}, {-18, 6.0, 2.6, 0}, {-4, 10.0, 3.4, 0}, {12, 9.0, 3.0, 0}, {22, 5.0, 2.0, 0}}, "steel", 8)
		prism(b, [][2]float64{{5, -10}, {24, 16}, {24, 22}, {6, 20}}, 1.6, "ash")
		prism(b, [][2]float64{{-5, -10}, {-24, 16}, {-24, 22}, {-6, 20}}, 1.6, "ash")
		box3(b, 0, 2.6, -4, 3.0, 1.6, 10, "teal", &FOpts{Shade: 3.6, Fixed: true})
		mesh := b.Mesh()
		mesh.Exhaust = []V{{5, 0, 22}, {-5, 0, 22}}
		m["cobra1"] = mesh
	}
	// Worm - bulbous light trader
	{
		b := NewBuilder()
		loft(b, []Stn{{-22, 1.0, 1.0, 0}, {-16, 4.0, 3.0, 0}, {-4, 6.0, 4.0, 0}, {10, 5.0, 3.2, 0}, {20, 2.5, 1.6, 0}}, "brown", 8)
		box(b, 0, 3.0, -10, 3.0, 1.6, 8, "steel")
		box3(b, 0, -3.0, 14, 6, 1.6, 4, "orange", &FOpts{Shade: 2.4, Fixed: true})
		mesh := b.Mesh()
		mesh.Exhaust = []V{{0, 0, 20}}
		m["worm"] = mesh
	}
	// Boa - big freighter
	{
		b := NewBuilder()
		loft(b, []Stn{{-40, 3.0, 3.0, 0}, {-32, 10.0, 7.0, 0}, {-12, 14.0, 9.0, 0}, {16, 13.0, 8.0, 0}, {32, 8.0, 5.0, 0}}, "green", 8)
		box(b, 18, -1, 6, 9, 7, 26, "brown")
		box(b, -18, -1, 6, 9, 7, 26, "brown")
		box3(b, 0, 4, -24, 4, 2, 8, "teal", &FOpts{Shade: 3.6, Fixed: true})
		mesh := b.Mesh()
		mesh.Exhaust = []V{{0, 0, 34}}
		m["boa"] = mesh
	}
	// Anaconda - the big one
	{
		b := NewBuilder()
		loft(b, []Stn{{-56, 2.0, 2.0, 0}, {-46, 9.0, 7.0, 0}, {-30, 14.0, 11.0, 0}, {10, 16.0, 11.0, 0}, {40, 12.0, 8.0, 0}, {52, 6.0, 4.0, 0}}, "ash", 8)
		box(b, 0, 10, -40, 12, 6, 16, "steel")
		box(b, 22, -2, 0, 10, 8, 30, "brown")
		box(b, -22, -2, 0, 10, 8, 30, "brown")
		box3(b, 0, 0, 52, 14, 6, 3, "orange", &FOpts{Shade: 2.6, Fixed: true})
		mesh := b.Mesh()
		mesh.Exhaust = []V{{10, 0, 52}, {-10, 0, 52}}
		m["anaconda"] = mesh
	}
	// Transporter - boxy orbital shuttle
	{
		b := NewBuilder()
		loft(b, []Stn{{-18, 3.0, 3.5, 0}, {-14, 7.0, 6.0, 0}, {10, 7.0, 6.0, 0}, {16, 4.0, 4.0, 0}}, "yellow", 8)
		box(b, 0, 6.4, -8, 4, 2, 8, "steel")
		box3(b, 0, 0, 15, 10, 8, 4, "orange", &FOpts{Shade: 2.4, Fixed: true})
		mesh := b.Mesh()
		mesh.Exhaust = []V{{0, 0, 17}}
		m["transporter"] = mesh
	}
	// Thargon - Thargoid drone
	{
		b := NewBuilder()
		loft(b, []Stn{{-7, 0.4, 0.4, 0}, {0, 4.5, 4.5, 0}, {7, 0.4, 0.4, 0}}, "alien", 6)
		box3(b, 0, 0, 0, 2.0, 2.0, 9.0, "teal", &FOpts{Shade: 3.8, Fixed: true})
		m["thargon"] = b.Mesh()
	}
	// Constrictor - the mission ship
	{
		b := NewBuilder()
		loft(b, []Stn{{-34, 0.8, 0.8, 0}, {-22, 4.5, 2.6, 0}, {-6, 8.0, 4.0, 0}, {14, 7.0, 3.2, 0}, {30, 3.0, 2.0, 0}}, "purple", 6)
		prism(b, [][2]float64{{6, -6}, {26, 18}, {25, 24}, {7, 20}}, 1.4, "ash")
		prism(b, [][2]float64{{-6, -6}, {-26, 18}, {-25, 24}, {-7, 20}}, 1.4, "ash")
		box3(b, 0, 3.0, -6, 2.4, 1.8, 12, "teal", &FOpts{Shade: 4.2, Fixed: true})
		mesh := b.Mesh()
		mesh.Exhaust = []V{{0, 0, 30}}
		m["constrictor"] = mesh
	}
	// alloy plate: what is left of a hull after a clean kill
	{
		b := NewBuilder()
		box3(b, 0, 0, 0, 15, 1.6, 11, "ash", &FOpts{Shade: 3.2, Fixed: true})
		box(b, 0, 1.2, 0, 6, 0.8, 4, "steel")
		m["plate"] = b.Mesh()
	}
	// rock hermit: a hollowed asteroid with a docking port on +Z
	{
		b := NewBuilder()
		t := (1 + math.Sqrt(5)) / 2
		ico := [][3]float64{
			{-1, t, 0}, {1, t, 0}, {-1, -t, 0}, {1, -t, 0},
			{0, -1, t}, {0, 1, t}, {0, -1, -t}, {0, 1, -t},
			{t, 0, -1}, {t, 0, 1}, {-t, 0, -1}, {-t, 0, 1},
		}
		idx := [][3]int{
			{0, 11, 5}, {0, 5, 1}, {0, 1, 7}, {0, 7, 10}, {0, 10, 11},
			{1, 5, 9}, {5, 11, 4}, {11, 10, 2}, {10, 7, 6}, {7, 1, 8},
			{3, 9, 4}, {3, 4, 2}, {3, 2, 6}, {3, 6, 8}, {3, 8, 9},
			{4, 9, 5}, {2, 4, 11}, {6, 2, 10}, {8, 6, 7}, {9, 8, 1},
		}
		r := NewRng(97)
		pts := make([]V, 12)
		for i, v := range ico {
			pts[i] = V{v[0], v[1], v[2]}.Norm().Mul(200 * (0.86 + r.Float()*0.24))
		}
		for _, f := range idx {
			b.Face([]V{pts[f[0]], pts[f[1]], pts[f[2]]}, "brown", nil)
		}
		collarSlot(b, 150, 250)
		m["hermit"] = b.Mesh()
	}
	// generation ship: a slow city in space
	{
		b := NewBuilder()
		loft(b, []Stn{{-900, 40, 40, 0}, {-820, 180, 180, 0}, {-600, 210, 210, 0},
			{600, 210, 210, 0}, {820, 180, 180, 0}, {900, 40, 40, 0}}, "ash", 10)
		for _, z := range []float64{-560, -180, 200, 580} {
			for a := 0; a < 10; a++ {
				ang := float64(a) / 10 * math.Pi * 2
				b.Face([]V{
					{math.Cos(ang) * 210, math.Sin(ang) * 210, z - 30},
					{math.Cos(ang+0.6) * 210, math.Sin(ang+0.6) * 210, z - 30},
					{math.Cos(ang+0.6) * 210, math.Sin(ang+0.6) * 210, z + 30},
					{math.Cos(ang) * 210, math.Sin(ang) * 210, z + 30},
				}, "yellow", &FOpts{Shade: 4.2, Fixed: true, Double: true})
			}
		}
		m["genship"] = b.Mesh()
	}
	{
		b := NewBuilder()
		r := NewRng(63)
		for i := 0; i < 6; i++ {
			sz := 3 + r.Float()*3
			box3(b, r.Float()*8-4, r.Float()*6-3, r.Float()*8-4, sz, sz*0.8, sz, "yellow",
				&FOpts{Shade: 3.6, Fixed: true})
		}
		m["splinter_ore"] = b.Mesh()
	}
	// Thargoid mothership: a big alien hull with a ring of red lights
	{
		b := NewBuilder()
		loft(b, []Stn{{-46, 2, 2, 0}, {-20, 30, 22, 0}, {0, 44, 30, 0}, {20, 30, 22, 0}, {46, 2, 2, 0}}, "alien", 10)
		for i := 0; i < 10; i++ {
			ang := float64(i) / 10 * math.Pi * 2
			x := math.Cos(ang) * 40
			y := math.Sin(ang) * 26
			box3(b, x, y, 0, 7, 7, 16, "red", &FOpts{Shade: 4.4, Fixed: true})
		}
		box3(b, 0, 0, -44, 16, 16, 10, "teal", &FOpts{Shade: 3.8, Fixed: true})
		m["mothership"] = b.Mesh()
	}
	m["boulder"] = asteroid(41)
	m["splinter"] = asteroid(73)
	m["dodeca"] = dodeca()
}

// Dodecahedron station. The hull is rotated so one pentagon faces +Z, which
// keeps the docking maths identical to the Coriolis: the slot sits on +Z in
// station-local space.
func dodeca() *render.Mesh {
	b := NewBuilder()
	phi := (1 + math.Sqrt(5)) / 2
	inv := 1 / phi

	var raw []V
	for _, x := range []float64{-1, 1} {
		for _, y := range []float64{-1, 1} {
			for _, z := range []float64{-1, 1} {
				raw = append(raw, V{x, y, z})
			}
		}
	}
	for _, a := range []float64{-inv, inv} {
		for _, c := range []float64{-phi, phi} {
			raw = append(raw, V{0, a, c}, V{a, c, 0}, V{c, 0, a})
		}
	}

	// basis with a face normal facing +Z, rotated into place. The face normals
	// of a dodecahedron are the icosahedron vertices, not the dodeca vertices.
	n := V{0, phi, 1}.Norm()
	ex := V{0, 1, 0}.Cross(n).Norm()
	ey := n.Cross(ex).Norm()
	rot := func(v V) V { return V{v.Dot(ex), v.Dot(ey), v.Dot(n)} }

	// scale so the front pentagon's plane sits exactly at StationSize, which is
	// where the shared docking maths expects the slot
	var front []V
	best := 0.0
	for _, v := range raw {
		if d := v.Dot(n); d > best {
			best = d
		}
	}
	for _, v := range raw {
		// relative tolerance: the five vertices of the +Z pentagon differ only
		// by floating point noise
		if v.Dot(n) > best-1e-3 {
			front = append(front, rot(v))
		}
	}
	var c V
	for _, v := range front {
		c = c.Add(v)
	}
	c = c.Mul(1 / float64(len(front)))
	faceR := 0.0
	for _, v := range front {
		if d := v.Sub(c).Len(); d > faceR {
			faceR = d
		}
	}
	// scale to the same overall radius as the Coriolis so the shared crash
	// radius and docking box still make sense
	maxLen := 1.0
	for _, v := range raw {
		if l := v.Len(); l > maxLen {
			maxLen = l
		}
	}
	scale := StationSize / maxLen
	verts := make([]V, len(raw))
	for i, v := range raw {
		verts[i] = rot(v).Mul(scale)
	}

	// faces: every icosahedron direction owns the 5 vertices closest to it
	// face normals: cyclic permutations of (0, phi, 1), which is the dual
	// icosahedron for this vertex set
	ico := []V{}
	for _, a := range []float64{-1, 1} {
		for _, b := range []float64{-phi, phi} {
			ico = append(ico, V{0, b, a}, V{a, 0, b}, V{b, a, 0})
		}
	}
	for faceIdx, nn := range ico {
		// normals must live in the same rotated frame as the vertices
		nn = rot(nn).Norm()
		var ring []V
		// the five vertices of a dodecahedron face sit 37.4 deg off the face
		// normal, i.e. cos = 0.794, so the cut-off has to sit below that
		for _, v := range verts {
			if v.Norm().Dot(nn) > 0.75 {
				ring = append(ring, v)
			}
		}
		if len(ring) != 5 {
			continue
		}
		// sort around the face normal so the polygon does not self-cross
		var fc V
		for _, v := range ring {
			fc = fc.Add(v)
		}
		fc = fc.Mul(1 / 5)
		u := ring[0].Sub(fc).Norm()
		w := nn.Cross(u)
		sort5(ring, func(a, b V) bool {
			return math.Atan2(a.Sub(fc).Dot(w), a.Sub(fc).Dot(u)) <
				math.Atan2(b.Sub(fc).Dot(w), b.Sub(fc).Dot(u))
		})
		b.Face(ring, "steel", nil)

		// surface detail: lit window strips on every face and a spire on a few,
		// otherwise twelve identical pentagons read as a blank ball
		u3 := ring[0].Sub(fc).Norm()
		w3 := nn.Cross(u3).Norm()
		for _, off := range []float64{-16, 10} {
			c2 := fc.Add(nn.Mul(1.5)).Add(u3.Mul(off))
			a := c2.Sub(u3.Mul(11)).Sub(w3.Mul(4))
			bb := c2.Add(u3.Mul(11)).Sub(w3.Mul(4))
			cc := c2.Add(u3.Mul(11)).Add(w3.Mul(4))
			dd := c2.Sub(u3.Mul(11)).Add(w3.Mul(4))
			b.Face([]V{a, bb, cc, dd}, "yellow", &FOpts{Shade: 4.4, Fixed: true, Double: true})
		}
		if faceIdx%3 == 0 {
			base := fc.Add(nn.Mul(120))
			apex := fc.Add(nn.Mul(190))
			for k := 0; k < 6; k++ {
				a1 := float64(k) * math.Pi / 3
				a2 := float64(k+1) * math.Pi / 3
				p1 := base.Add(u3.Mul(math.Cos(a1) * 7)).Add(w3.Mul(math.Sin(a1) * 7))
				p2 := base.Add(u3.Mul(math.Cos(a2) * 7)).Add(w3.Mul(math.Sin(a2) * 7))
				b.Face([]V{p1, p2, apex}, "red", &FOpts{Shade: 3.4, Fixed: true, Double: true})
			}
		}
	}

	// docking collar on the +Z face
	z := 0.0
	for _, v := range front {
		z = v.Dot(n) * scale
	}
	collarSlot(b, z-70, z+40)
	return b.Mesh()
}

func sort5(v []V, less func(a, b V) bool) {
	for i := 1; i < len(v); i++ {
		for j := i; j > 0 && less(v[j], v[j-1]); j-- {
			v[j], v[j-1] = v[j-1], v[j]
		}
	}
}

// collarSlot builds a docking collar that sticks out of a hull whose surface
// cannot easily be cut: an open-ended tube with lit walls, a back wall and a
// frame. From outside it reads as a real opening because you can see its
// thickness from any angle.
func collarSlot(b *Builder, zBack, zFront float64) {
	w := StationSlotW * StationSize * 0.95
	h := StationSlotH * StationSize * 1.05
	// back wall (inside the rock), lit so it never reads as space
	b.Face([]V{{-w, h, zBack}, {w, h, zBack}, {w, -h, zBack}, {-w, -h, zBack}}, "steel",
		&FOpts{Shade: 1.3, Fixed: true, Double: true})
	// four walls running out to the mouth
	b.Face([]V{{-w, h, zBack}, {-w, h, zFront}, {w, h, zFront}, {w, h, zBack}}, "steel",
		&FOpts{Shade: 2.6, Fixed: true, Double: true})
	b.Face([]V{{-w, -h, zBack}, {w, -h, zBack}, {w, -h, zFront}, {-w, -h, zFront}}, "steel",
		&FOpts{Shade: 1.6, Fixed: true, Double: true})
	b.Face([]V{{-w, h, zBack}, {-w, -h, zBack}, {-w, -h, zFront}, {-w, h, zFront}}, "steel",
		&FOpts{Shade: 2.2, Fixed: true, Double: true})
	b.Face([]V{{w, h, zBack}, {w, h, zFront}, {w, -h, zFront}, {w, -h, zBack}}, "steel",
		&FOpts{Shade: 1.9, Fixed: true, Double: true})
	// lights on the back wall and a frame around the mouth
	for i := -2; i <= 2; i++ {
		lx := float64(i) * (w / 2.4)
		b.Face([]V{{lx - 4, 5, zBack + 1}, {lx + 4, 5, zBack + 1}, {lx + 4, -5, zBack + 1}, {lx - 4, -5, zBack + 1}},
			"yellow", &FOpts{Shade: 5, Fixed: true, Double: true})
	}
	const rim = 8.0
	for _, sg := range []float64{-1, 1} {
		b.Face([]V{{-w - rim, sg * h, zFront}, {w + rim, sg * h, zFront}, {w + rim, sg * (h + rim), zFront}, {-w - rim, sg * (h + rim), zFront}},
			"steel", &FOpts{Shade: 4.6, Fixed: true, Double: true})
		for _, sx := range []float64{-1, 1} {
			b.Face([]V{{sx * w, -h, zFront}, {sx * (w + rim), -h, zFront}, {sx * (w + rim), h, zFront}, {sx * w, h, zFront}},
				"steel", &FOpts{Shade: 4.2, Fixed: true, Double: true})
		}
	}
}
