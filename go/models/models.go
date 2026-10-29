// Package models holds the low-poly meshes, the station, and the procedural
// planet and nebula sprites.
package models

import (
	"math"

	"heliosdrift/render"
)

type V = render.Vec3

type FOpts struct {
	Shade  float64
	Tint   float64
	Fixed  bool
	Double bool
}

type Stn struct{ Z, W, H, Y float64 }

type Builder struct {
	Faces  []render.Face
	Radius float64
	center V
	// profile remembers the last lofted hull so details can be placed on its
	// actual surface instead of guessed coordinates.
	profile []Stn
}

func NewBuilder() *Builder { return &Builder{} }

func (b *Builder) Face(pts []V, ramp string, opts *FOpts) {
	if len(pts) < 3 {
		return
	}
	n := pts[1].Sub(pts[0]).Cross(pts[2].Sub(pts[0]))
	if n.Len() < 1e-6 {
		return
	}
	n = n.Norm()
	var c V
	for _, p := range pts {
		c = c.Add(p)
	}
	c = c.Mul(1 / float64(len(pts)))
	if n.Dot(c.Sub(b.center)) < 0 {
		n = n.Mul(-1)
		rev := make([]V, len(pts))
		for i := range pts {
			rev[i] = pts[len(pts)-1-i]
		}
		pts = rev
	}
	rampID, ok := render.RampIndex(ramp)
	if !ok {
		rampID = 0
	}
	f := render.Face{P: pts, N: n, Ramp: rampID}
	if opts != nil {
		if opts.Fixed {
			f.Shade = opts.Shade
			f.Fixed = true
		}
		f.Tint = opts.Tint
		f.Two = opts.Double
	}
	b.Faces = append(b.Faces, f)
	for _, p := range pts {
		if d := p.Len(); d > b.Radius {
			b.Radius = d
		}
	}
}

func (b *Builder) Mesh() *render.Mesh {
	orientShells(b.Faces)
	markCreases(b.Faces)
	return &render.Mesh{Faces: b.Faces, Radius: b.Radius}
}

// markCreases bakes a cheap ambient occlusion term: a face that meets another
// face in a concave crease gets darkened a little. It costs nothing per frame
// and it is what makes stacked plates read as stacked instead of painted.
func markCreases(faces []render.Face) {
	type ikey [3]int64
	key := func(p render.Vec3) ikey {
		return ikey{
			int64(math.Round(p[0] * 20)),
			int64(math.Round(p[1] * 20)),
			int64(math.Round(p[2] * 20)),
		}
	}
	type ekey struct{ a, b ikey }
	mk := func(p, q render.Vec3) ekey {
		a, b := key(p), key(q)
		if a[0] > b[0] || (a[0] == b[0] && (a[1] > b[1] || (a[1] == b[1] && a[2] > b[2]))) {
			a, b = b, a
		}
		return ekey{a, b}
	}
	owners := map[ekey][]int{}
	centres := make([]render.Vec3, len(faces))
	for fi := range faces {
		p := faces[fi].P
		var c render.Vec3
		for _, v := range p {
			c = c.Add(v)
		}
		centres[fi] = c.Mul(1 / float64(len(p)))
		for i := range p {
			owners[mk(p[i], p[(i+1)%len(p)])] = append(owners[mk(p[i], p[(i+1)%len(p)])], fi)
		}
	}
	for fi := range faces {
		p := faces[fi].P
		for i := range p {
			fs := owners[mk(p[i], p[(i+1)%len(p)])]
			if len(fs) != 2 || fs[0] == fs[1] {
				continue
			}
			other := fs[0]
			if other == fi {
				other = fs[1]
			}
			d := faces[fi].N.Dot(faces[other].N)
			ang := math.Acos(math.Max(-1, math.Min(1, d)))
			if ang < 0.25 || ang > 2.6 {
				continue // coplanar or folded back on itself
			}
			// concave when the neighbour sits on the outward side of our plane
			if faces[fi].N.Dot(centres[other].Sub(p[0])) > 0 {
				faces[fi].AO = math.Min(1, faces[fi].AO+ang/1.2)
			}
		}
	}
}

func orientShells(faces []render.Face) {
	keyOf := func(p V) string {
		return itoa(int(math.Round(p[0]*100))) + "," + itoa(int(math.Round(p[1]*100))) + "," + itoa(int(math.Round(p[2]*100)))
	}
	parent := map[string]string{}
	var find func(string) string
	find = func(a string) string {
		root := a
		for parent[root] != root {
			root = parent[root]
		}
		for parent[a] != root {
			n := parent[a]
			parent[a] = root
			a = n
		}
		return root
	}
	for i := range faces {
		f := &faces[i]
		k0 := keyOf(f.P[0])
		if _, ok := parent[k0]; !ok {
			parent[k0] = k0
		}
		for j := 1; j < len(f.P); j++ {
			kj := keyOf(f.P[j])
			if _, ok := parent[kj]; !ok {
				parent[kj] = kj
			}
			ra, rb := find(k0), find(kj)
			if ra != rb {
				parent[rb] = ra
			}
		}
	}
	type acc struct {
		sx, sy, sz, n, nx, ny, nz, area float64
	}
	accs := map[string]*acc{}
	for i := range faces {
		f := &faces[i]
		root := find(keyOf(f.P[0]))
		a, ok := accs[root]
		if !ok {
			a = &acc{}
			accs[root] = a
		}
		for _, p := range f.P {
			a.sx += p[0]
			a.sy += p[1]
			a.sz += p[2]
			a.n++
		}
		u := f.P[1].Sub(f.P[0])
		v := f.P[2].Sub(f.P[0])
		area := u.Cross(v).Len() * 0.5
		a.nx += f.N[0] * area
		a.ny += f.N[1] * area
		a.nz += f.N[2] * area
		a.area += area
	}
	for i := range faces {
		f := &faces[i]
		s := accs[find(keyOf(f.P[0]))]
		cx, cy, cz := s.sx/s.n, s.sy/s.n, s.sz/s.n
		nlen := math.Sqrt(s.nx*s.nx + s.ny*s.ny + s.nz*s.nz)
		flatness := 0.0
		if s.area > 0 {
			flatness = nlen / s.area
		}
		ref := V{cx, cy, cz}
		if flatness > 0.6 {
			ref = V{0, 0, 0}
		}
		var fc V
		for _, p := range f.P {
			fc = fc.Add(p)
		}
		fc = fc.Mul(1 / float64(len(f.P)))
		dir := fc.Sub(ref)
		if dir.Len() < 0.001 {
			continue
		}
		if f.N.Dot(dir) < 0 {
			f.N = f.N.Mul(-1)
			rev := make([]V, len(f.P))
			for k := range f.P {
				rev[k] = f.P[len(f.P)-1-k]
			}
			f.P = rev
		}
	}
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var buf [12]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func loftZ(b *Builder, st []Stn, ramp string, sides int, opts *FOpts) {
	rings := make([][]V, len(st))
	for i, s := range st {
		ring := make([]V, sides)
		for j := 0; j < sides; j++ {
			a := math.Pi/float64(sides) + float64(j)/float64(sides)*2*math.Pi
			ring[j] = V{math.Cos(a) * s.W, s.Y + math.Sin(a)*s.H, s.Z}
		}
		rings[i] = ring
	}
	for i := 0; i+1 < len(rings); i++ {
		for j := 0; j < sides; j++ {
			a0 := rings[i][j]
			a1 := rings[i][(j+1)%sides]
			b0 := rings[i+1][j]
			b1 := rings[i+1][(j+1)%sides]
			b.Face([]V{a0, a1, b1, b0}, ramp, opts)
		}
	}
	if opts == nil || !opts.Double {
		first := make([]V, len(rings[0]))
		for i := range rings[0] {
			first[i] = rings[0][len(rings[0])-1-i]
		}
		b.Face(first, ramp, opts)
		last := rings[len(rings)-1]
		b.Face(append([]V{}, last...), ramp, opts)
	}
}

func loft(b *Builder, st []Stn, ramp string, sides int) {
	loftZ(b, st, ramp, sides, nil)
	b.profile = st
}

// SpineDetails stamps small hull plates along the top of the last lofted hull:
// hatches and access panels that break up the big flat facets.
func (b *Builder) SpineDetails(n int, ramp string) {
	if len(b.profile) < 3 || n <= 0 {
		return
	}
	z0 := b.profile[1].Z
	z1 := b.profile[len(b.profile)-2].Z
	for i := 0; i < n; i++ {
		t := (float64(i) + 0.5) / float64(n)
		z := z0 + (z1-z0)*t
		w, top := b.sectionAt(z)
		if w < 1.5 {
			continue
		}
		sz := w * 0.5
		if sz > 7 {
			sz = 7
		}
		b.Face([]V{
			{sz, top + 0.5, z - sz*0.9},
			{sz, top + 0.5, z + sz*0.9},
			{-sz, top + 0.5, z + sz*0.9},
			{-sz, top + 0.5, z - sz*0.9},
		}, ramp, &FOpts{Tint: 0.72})
	}
}

// sectionAt interpolates half-width and top height at a given z.
func (b *Builder) sectionAt(z float64) (w, top float64) {
	p := b.profile
	for i := 0; i+1 < len(p); i++ {
		a, c := p[i], p[i+1]
		lo, hi := a.Z, c.Z
		if lo > hi {
			lo, hi = hi, lo
		}
		if z < lo || z > hi {
			continue
		}
		f := 0.0
		if hi-lo > 1e-6 {
			f = (z - p[i].Z) / (p[i+1].Z - p[i].Z)
		}
		return p[i].W + (p[i+1].W-p[i].W)*f, p[i].Y + p[i].H + ((p[i+1].Y+p[i+1].H)-(p[i].Y+p[i].H))*f
	}
	return 0, 0
}

func loftShade(b *Builder, st []Stn, ramp string, sides int, shade float64) {
	loftZ(b, st, ramp, sides, &FOpts{Shade: shade, Fixed: true})
}

// prism extrudes a 2D polygon given in the XZ plane along Y.
func prism(b *Builder, poly [][2]float64, thick float64, ramp string) {
	y := thick / 2
	top := make([]V, len(poly))
	bot := make([]V, len(poly))
	for i, p := range poly {
		top[i] = V{p[0], y, p[1]}
		bot[i] = V{p[0], -y, p[1]}
	}
	b.Face(top, ramp, nil)
	b.Face(bot, ramp, nil)
	for i := range poly {
		j := (i + 1) % len(poly)
		b.Face([]V{top[i], top[j], bot[j], bot[i]}, ramp, nil)
	}
}

func box3(b *Builder, cx, cy, cz, sx, sy, sz float64, ramp string, opts *FOpts) {
	x, y, z := sx/2, sy/2, sz/2
	v := []V{
		{cx - x, cy + y, cz - z}, {cx + x, cy + y, cz - z},
		{cx + x, cy - y, cz - z}, {cx - x, cy - y, cz - z},
		{cx - x, cy + y, cz + z}, {cx + x, cy + y, cz + z},
		{cx + x, cy - y, cz + z}, {cx - x, cy - y, cz + z},
	}
	b.Face([]V{v[0], v[1], v[2], v[3]}, ramp, opts)
	b.Face([]V{v[5], v[4], v[7], v[6]}, ramp, opts)
	b.Face([]V{v[4], v[0], v[3], v[7]}, ramp, opts)
	b.Face([]V{v[1], v[5], v[6], v[2]}, ramp, opts)
	b.Face([]V{v[4], v[5], v[1], v[0]}, ramp, opts)
	b.Face([]V{v[3], v[2], v[6], v[7]}, ramp, opts)
}

func box(b *Builder, cx, cy, cz, sx, sy, sz float64, ramp string) {
	box3(b, cx, cy, cz, sx, sy, sz, ramp, nil)
}

// ------------------------------------------------------------------ ships

func buildShips() map[string]*render.Mesh {
	m := map[string]*render.Mesh{}

	// Cobra Mk III
	{
		b := NewBuilder()
		loft(b, []Stn{
			{-34, 1.2, 1.0, 0}, {-24, 5.5, 2.6, 0}, {-10, 9.0, 3.4, 0.4},
			{8, 11.5, 3.2, 0.2}, {26, 9.5, 2.6, 0}, {31, 4.0, 1.8, 0},
		}, "steel", 8)
		prism(b, [][2]float64{{4, -18}, {30, 16}, {31, 24}, {9, 26}, {7, 4}}, 1.8, "steel")
		prism(b, [][2]float64{{-4, -18}, {-30, 16}, {-31, 24}, {-9, 26}, {-7, 4}}, 1.8, "steel")
		box(b, 29, 0, 18, 2.0, 2.0, 12, "ash")
		box(b, -29, 0, 18, 2.0, 2.0, 12, "ash")
		loftZ(b, []Stn{
			{-22, 1.0, 0.8, 2.4}, {-14, 4.2, 2.2, 2.8}, {0, 4.0, 1.6, 3.0}, {8, 2.0, 0.8, 2.6},
		}, "teal", 6, &FOpts{Shade: 4.2, Fixed: true})
		box(b, 4.4, 0.2, 30, 4.0, 3.0, 4.0, "ash")
		box(b, -4.4, 0.2, 30, 4.0, 3.0, 4.0, "ash")
		box(b, 0, 3.1, 18, 2.0, 1.4, 12, "steel")
		// hull detail: hatches, intakes and a spine, so the big flat panels
		// are not one blank surface up close
		box(b, 6.5, 3.4, 2, 4.0, 0.9, 7, "ash")
		box(b, -6.5, 3.4, 2, 4.0, 0.9, 7, "ash")
		box(b, 0, 3.6, -14, 2.6, 0.8, 9, "ash")
		box3(b, 9.5, 2.2, -16, 3.0, 1.6, 3.0, "teal", &FOpts{Shade: 4.2, Fixed: true})
		box3(b, -9.5, 2.2, -16, 3.0, 1.6, 3.0, "teal", &FOpts{Shade: 4.2, Fixed: true})
		mesh := b.Mesh()
		mesh.Exhaust = []V{{4.4, 0.2, 32}, {-4.4, 0.2, 32}}
		m["cobra3"] = mesh
	}
	// Sidewinder
	{
		b := NewBuilder()
		loft(b, []Stn{{-28, 1.0, 0.8, 0}, {-14, 6.0, 2.4, 0}, {2, 8.0, 3.4, 0}, {18, 6.0, 3.0, 0}, {24, 3.0, 2.0, 0}}, "ash", 6)
		prism(b, [][2]float64{{5, -6}, {24, 20}, {23, 25}, {7, 20}}, 1.4, "ash")
		prism(b, [][2]float64{{-5, -6}, {-24, 20}, {-23, 25}, {-7, 20}}, 1.4, "ash")
		box(b, 0, 3.2, -6, 2.0, 1.6, 8, "steel")
		box(b, 4.5, 2.6, 2, 3.0, 0.8, 6, "ash")
		box(b, -4.5, 2.6, 2, 3.0, 0.8, 6, "ash")
		box3(b, 0, 0, 23, 5.0, 3.0, 3.0, "red", &FOpts{Shade: 2.2, Fixed: true})
		m["sidewinder"] = b.Mesh()
	}
	// Krait - diamond planform
	{
		b := NewBuilder()
		loft(b, []Stn{{-30, 0.8, 0.7, 0}, {-18, 5.0, 2.4, 0}, {-4, 13.0, 4.0, 0}, {6, 14.0, 4.2, 0.4}, {18, 7.0, 2.6, 0}, {26, 2.5, 1.4, 0}}, "green", 6)
		box3(b, 0, 3.0, -8, 2.6, 2.4, 12, "alien", &FOpts{Shade: 4.0, Fixed: true})
		box(b, 4.0, 0, 23, 3.4, 2.6, 5.0, "ash")
		box(b, -4.0, 0, 23, 3.4, 2.6, 5.0, "ash")
		b.SpineDetails(3, "ash")
		m["krait"] = b.Mesh()
	}
	// Mamba
	{
		b := NewBuilder()
		loft(b, []Stn{{-32, 0.8, 0.7, 0}, {-16, 4.0, 2.0, 0}, {2, 7.5, 3.2, 0}, {20, 6.0, 2.6, 0}, {27, 3.0, 1.8, 0}}, "purple", 6)
		prism(b, [][2]float64{{6, -2}, {22, 22}, {20, 26}, {8, 20}}, 1.2, "purple")
		prism(b, [][2]float64{{-6, -2}, {-22, 22}, {-20, 26}, {-8, 20}}, 1.2, "purple")
		loft(b, []Stn{{-40, 0.4, 0.4, 0}, {-28, 1.6, 1.2, 0}}, "ash", 5)
		box3(b, 0, 2.8, 0, 2.0, 1.6, 12, "teal", &FOpts{Shade: 3.6, Fixed: true})
		b.SpineDetails(3, "ash")
		m["mamba"] = b.Mesh()
	}
	// Adder
	{
		b := NewBuilder()
		loft(b, []Stn{{-24, 2.0, 2.0, 0}, {-18, 7.0, 3.0, 0}, {10, 8.5, 4.0, 0.4}, {22, 7.0, 3.4, 0}}, "brown", 8)
		prism(b, [][2]float64{{5, -8}, {22, 14}, {22, 20}, {5, 18}}, 2.6, "brown")
		prism(b, [][2]float64{{-5, -8}, {-22, 14}, {-22, 20}, {-5, 18}}, 2.6, "brown")
		box(b, 0, 4.2, -8, 3.0, 2.0, 10, "steel")
		box3(b, 0, 0, 23, 8.0, 4.0, 3.0, "orange", &FOpts{Shade: 2.4, Fixed: true})
		b.SpineDetails(3, "ash")
		m["adder"] = b.Mesh()
	}
	// Python
	{
		b := NewBuilder()
		loft(b, []Stn{{-40, 3.0, 3.0, 0}, {-30, 9.0, 7.0, 0}, {-10, 13.0, 9.0, 0}, {20, 13.0, 8.0, 0}, {34, 9.0, 6.0, 0}, {38, 5.0, 4.0, 0}}, "orange", 8)
		box(b, 15, -1, 12, 8.0, 6.0, 24, "brown")
		box(b, -15, -1, 12, 8.0, 6.0, 24, "brown")
		box(b, 0, 8.8, -2, 5.0, 2.6, 16, "steel")
		box(b, 0, 0, 37, 14.0, 8.0, 4.0, "ash")
		b.SpineDetails(3, "ash")
		m["python"] = b.Mesh()
	}
	// Viper
	{
		b := NewBuilder()
		loft(b, []Stn{{-30, 1.6, 1.6, 0}, {-20, 5.0, 3.0, 0}, {-4, 8.0, 4.0, 0}, {14, 7.0, 3.4, 0}, {24, 4.0, 2.4, 0}}, "blue", 8)
		prism(b, [][2]float64{{6, -4}, {26, 16}, {26, 22}, {7, 18}}, 1.6, "blue")
		prism(b, [][2]float64{{-6, -4}, {-26, 16}, {-26, 22}, {-7, 18}}, 1.6, "blue")
		box(b, 0, 4.2, -6, 3.2, 1.2, 14, "ash")
		box3(b, 4.2, 0, 25, 3.0, 2.6, 3.0, "yellow", &FOpts{Shade: 4.4, Fixed: true})
		box3(b, -4.2, 0, 25, 3.0, 2.6, 3.0, "yellow", &FOpts{Shade: 4.4, Fixed: true})
		b.SpineDetails(3, "ash")
		m["viper"] = b.Mesh()
	}
	// Asp Mk II
	{
		b := NewBuilder()
		loft(b, []Stn{{-36, 1.0, 1.0, 0}, {-20, 5.0, 3.4, 0}, {0, 8.5, 4.6, 0}, {20, 9.0, 4.0, 0}, {30, 6.0, 3.0, 0}}, "teal", 8)
		prism(b, [][2]float64{{7, 10}, {20, 30}, {20, 33}, {8, 26}}, 1.4, "teal")
		prism(b, [][2]float64{{-7, 10}, {-20, 30}, {-20, 33}, {-8, 26}}, 1.4, "teal")
		box(b, 0, 4.6, 4, 3.0, 2.0, 14, "ash")
		box3(b, 0, 0, 29.5, 4.0, 3.0, 4.0, "pink", &FOpts{Shade: 4.0, Fixed: true})
		b.SpineDetails(3, "ash")
		m["asp"] = b.Mesh()
	}
	// Moray
	{
		b := NewBuilder()
		loft(b, []Stn{{-28, 2.0, 1.4, 0}, {-12, 7.0, 2.6, 0}, {6, 8.0, 3.0, 0}, {22, 5.0, 2.0, 0}}, "teal", 8)
		prism(b, [][2]float64{{7, -2}, {18, 10}, {17, 18}, {7, 16}}, 1.2, "teal")
		prism(b, [][2]float64{{-7, -2}, {-18, 10}, {-17, 18}, {-7, 16}}, 1.2, "teal")
		box3(b, 0, 2.8, -4, 2.6, 1.8, 12, "green", &FOpts{Shade: 4.2, Fixed: true})
		box(b, 0, -3.2, 2, 3.0, 4.0, 18, "steel")
		b.SpineDetails(3, "ash")
		m["moray"] = b.Mesh()
	}
	// Fer-de-Lance
	{
		b := NewBuilder()
		loft(b, []Stn{{-30, 1.0, 1.0, 0}, {-14, 5.0, 3.0, 0}, {4, 10.0, 4.0, 0}, {18, 8.0, 3.0, 0}, {26, 4.0, 2.0, 0}}, "pink", 8)
		prism(b, [][2]float64{{9, -4}, {30, 14}, {30, 20}, {10, 18}}, 1.4, "pink")
		prism(b, [][2]float64{{-9, -4}, {-30, 14}, {-30, 20}, {-10, 18}}, 1.4, "pink")
		loftZ(b, []Stn{{-16, 1.0, 0.6, 2.6}, {-6, 3.2, 1.6, 3.0}, {6, 2.0, 1.0, 2.8}}, "yellow", 6, &FOpts{Shade: 4.4, Fixed: true})
		box3(b, 5, 0, 25, 3.0, 2.4, 3.0, "yellow", &FOpts{Shade: 3.4, Fixed: true})
		box3(b, -5, 0, 25, 3.0, 2.4, 3.0, "yellow", &FOpts{Shade: 3.4, Fixed: true})
		b.SpineDetails(3, "ash")
		m["ferdelance"] = b.Mesh()
	}
	// Thargoid - flattened octagonal saucer
	{
		b := NewBuilder()
		loft(b, []Stn{{-22, 1.5, 1.2, 0}, {-16, 9.0, 3.0, 0}, {-7, 17.0, 4.2, 0}, {2, 18.5, 4.4, 0}, {11, 14.0, 3.6, 0}, {19, 6.0, 2.4, 0}}, "alien", 8)
		loftZ(b, []Stn{{-6, 5.0, 2.2, 3.4}, {1, 6.5, 3.4, 4.2}, {8, 4.0, 2.0, 3.4}}, "green", 8, &FOpts{Shade: 3.8, Fixed: true})
		for i := 0; i < 6; i++ {
			a := float64(i) * math.Pi / 3
			box3(b, math.Cos(a)*13, -0.5, math.Sin(a)*9+1, 3.6, 3.0, 3.4, "alien", &FOpts{Shade: 4.4, Fixed: true})
		}
		m["thargoid"] = b.Mesh()
	}
	// Shuttle
	{
		b := NewBuilder()
		box(b, 0, 0, 0, 10, 6, 22, "yellow")
		loft(b, []Stn{{-20, 2.0, 2.0, 0}, {-11, 5.0, 3.0, 0}}, "yellow", 6)
		box(b, 8, -1, 3, 8, 2.2, 12, "ash")
		box(b, -8, -1, 3, 8, 2.2, 12, "ash")
		box(b, 0, 3.2, 2, 4.0, 2.0, 10, "steel")
		b.SpineDetails(3, "ash")
		m["shuttle"] = b.Mesh()
	}
	m["asteroid"] = asteroid(11)
	m["asteroid2"] = asteroid(29)
	m["canister"] = canister()
	m["capsule"] = capsule()
	m["missile"] = missile()
	m["station"] = coriolis()
	return m
}

func asteroid(seed uint32) *render.Mesh {
	b := NewBuilder()
	t := (1 + math.Sqrt(5)) / 2
	v := [][3]float64{
		{-1, t, 0}, {1, t, 0}, {-1, -t, 0}, {1, -t, 0},
		{0, -1, t}, {0, 1, t}, {0, -1, -t}, {0, 1, -t},
		{t, 0, -1}, {t, 0, 1}, {-t, 0, -1}, {-t, 0, 1},
	}
	f := [][3]int{
		{0, 11, 5}, {0, 5, 1}, {0, 1, 7}, {0, 7, 10}, {0, 10, 11},
		{1, 5, 9}, {5, 11, 4}, {11, 10, 2}, {10, 7, 6}, {7, 1, 8},
		{3, 9, 4}, {3, 4, 2}, {3, 2, 6}, {3, 6, 8}, {3, 8, 9},
		{4, 9, 5}, {2, 4, 11}, {6, 2, 10}, {8, 6, 7}, {9, 8, 1},
	}
	r := NewRng(seed)
	pts := make([]V, 12)
	for i, p := range v {
		j := 0.66 + r.Float()*0.6
		pts[i] = V{p[0] * j * 11, p[1] * j * 11, p[2] * j * 11}
	}
	for _, face := range f {
		b.Face([]V{pts[face[0]], pts[face[1]], pts[face[2]]}, "brown", nil)
	}
	return b.Mesh()
}

func canister() *render.Mesh {
	b := NewBuilder()
	loft(b, []Stn{{-7, 2.2, 2.2, 0}, {-5, 4.0, 4.0, 0}, {5, 4.0, 4.0, 0}, {7, 2.2, 2.2, 0}}, "yellow", 8)
	box3(b, 0, 0, 0, 8.3, 8.3, 3.0, "orange", &FOpts{Shade: 3.0, Fixed: true})
	return b.Mesh()
}

func capsule() *render.Mesh {
	b := NewBuilder()
	loft(b, []Stn{{-9, 1.4, 1.4, 0}, {-5, 3.4, 3.4, 0}, {3, 4.0, 3.6, 0}, {8, 2.0, 2.0, 0}}, "red", 8)
	box(b, 0, 3.0, -1, 2.0, 1.4, 5.0, "ash")
	return b.Mesh()
}

func missile() *render.Mesh {
	b := NewBuilder()
	loft(b, []Stn{{-10, 0.3, 0.3, 0}, {-4, 1.3, 1.3, 0}, {6, 1.3, 1.3, 0}, {8, 0.9, 0.9, 0}}, "ash", 6)
	box3(b, 0, 0, 7, 0.4, 5.0, 3.0, "red", &FOpts{Shade: 3.4, Fixed: true})
	box3(b, 0, 0, 7, 5.0, 0.4, 3.0, "red", &FOpts{Shade: 3.4, Fixed: true})
	return b.Mesh()
}

// Coriolis station: cuboctahedron with a docking slot through the +Z face.
func coriolis() *render.Mesh {
	b := NewBuilder()
	const S = StationSize
	w := StationSlotW * S
	h := StationSlotH * S
	// five plain diamond faces + eight corner triangles
	for a := 0; a < 3; a++ {
		for _, sgn := range []float64{-1, 1} {
			if a == 2 && sgn == 1 {
				continue
			}
			var others []int
			for q := 0; q < 3; q++ {
				if q != a {
					others = append(others, q)
				}
			}
			ring := [4][2]float64{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
			pts := make([]V, 4)
			for r := 0; r < 4; r++ {
				var q V
				q[a] = sgn * S
				q[others[0]] = ring[r][0] * S
				q[others[1]] = ring[r][1] * S
				pts[r] = q
			}
			tint := 0.9 + float64(a%3)*0.09
			b.Face(pts, "ash", &FOpts{Tint: tint})
			// inset hull plate, slightly proud of the face
			var centre V
			for _, p := range pts {
				centre = centre.Add(p)
			}
			centre = centre.Mul(0.25)
			plate := make([]V, 4)
			for k, p := range pts {
				q := p.Add(centre.Sub(p).Mul(0.22))
				q[a] += sgn * 1.5
				plate[k] = q
			}
			b.Face(plate, "steel", &FOpts{Tint: 0.78 + float64(a%2)*0.1})
			// a second, smaller panel inset so each face reads as stacked plates
			plate2 := make([]V, 4)
			for k, p := range pts {
				q := p.Add(centre.Sub(p).Mul(0.42))
				q[a] += sgn * 2.2
				plate2[k] = q
			}
			b.Face(plate2, "ash", &FOpts{Tint: 0.62 + float64(a%2)*0.12})

			// a pair of lit windows on each plate
			for w := 0; w < 2; w++ {
				var cw V
				for _, p := range plate {
					cw = cw.Add(p)
				}
				cw = cw.Mul(0.25)
				edge := plate[0].Add(plate[1]).Mul(0.5)
				edge = cw.Add(edge.Sub(cw).Mul(0.55 + float64(w)*0.32))
				var q V
				for k := 0; k < 3; k++ {
					q[k] = edge[k] + (cw[k]-edge[k])*0.06
				}
				q[a] += sgn * 1.8
				half := 1.6
				var o1, o2 V
				o1 = q
				o2 = q
				o1[others[0]] += half
				o2[others[0]] -= half
				o1[others[1]] += half * 0.5
				o2[others[1]] -= half * 0.5
				b.Face([]V{q, o1, o2}, "yellow", &FOpts{Shade: 4.6, Fixed: true})
			}
		}
	}
	for _, cx := range []float64{-1, 1} {
		for _, cy := range []float64{-1, 1} {
			for _, cz := range []float64{-1, 1} {
				c := V{cx, cy, cz}
				tri := make([]V, 3)
				for k := 0; k < 3; k++ {
					var pt V
					others := []int{}
					for q := 0; q < 3; q++ {
						if q != k {
							others = append(others, q)
						}
					}
					pt[others[0]] = c[others[0]] * S
					pt[others[1]] = c[others[1]] * S
					tri[k] = pt
				}
				b.Face(tri, "ash", &FOpts{Tint: 1.12})
			}
		}
	}
	// The +Z face is a diamond with a rectangular hole in it, and the hole runs
	// back into the hull as a real shaft. The renderer keeps a z-buffer, so a
	// genuine cavity sorts correctly; the earlier "painted port" sat in front of
	// the surface with no thickness and read as a sticker.
	dh := S - h
	b.Face([]V{{-S, 0, S}, {-dh, h, S}, {-w, h, S}, {-w, -h, S}, {-dh, -h, S}}, "ash", nil)
	b.Face([]V{{S, 0, S}, {dh, -h, S}, {w, -h, S}, {w, h, S}, {dh, h, S}}, "ash", nil)
	b.Face([]V{{-dh, h, S}, {dh, h, S}, {0, S, S}}, "ash", nil)
	b.Face([]V{{-dh, -h, S}, {0, -S, S}, {dh, -h, S}}, "ash", nil)

	// shaft: four lit walls running inward, then a back wall with lights
	const depth = 78.0
	T := S - depth
	b.Face([]V{{-w, h, S}, {-w, h, T}, {w, h, T}, {w, h, S}}, "steel", &FOpts{Shade: 2.4, Fixed: true})
	b.Face([]V{{-w, -h, S}, {w, -h, S}, {w, -h, T}, {-w, -h, T}}, "steel", &FOpts{Shade: 1.5, Fixed: true})
	b.Face([]V{{-w, h, S}, {-w, -h, S}, {-w, -h, T}, {-w, h, T}}, "steel", &FOpts{Shade: 2.0, Fixed: true})
	b.Face([]V{{w, h, S}, {w, h, T}, {w, -h, T}, {w, -h, S}}, "steel", &FOpts{Shade: 1.8, Fixed: true})
	b.Face([]V{{-w, h, T}, {-w, -h, T}, {w, -h, T}, {w, h, T}}, "steel", &FOpts{Shade: 1.1, Fixed: true})
	// landing lights deep in the shaft
	for i := -2; i <= 2; i++ {
		lx := float64(i) * (w / 2.6)
		b.Face([]V{{lx - 3, 4, T + 1}, {lx + 3, 4, T + 1}, {lx + 3, -4, T + 1}, {lx - 3, -4, T + 1}},
			"yellow", &FOpts{Shade: 5, Fixed: true})
	}
	// a lit frame band around the mouth, wide enough to survive at 320x224
	const rim = 7.0
	b.Face([]V{{-w - rim, h, S + 0.4}, {w + rim, h, S + 0.4}, {w + rim, h + rim, S + 0.4}, {-w - rim, h + rim, S + 0.4}},
		"steel", &FOpts{Shade: 4.6, Fixed: true})
	b.Face([]V{{-w - rim, -h - rim, S + 0.4}, {w + rim, -h - rim, S + 0.4}, {w + rim, -h, S + 0.4}, {-w - rim, -h, S + 0.4}},
		"steel", &FOpts{Shade: 4.2, Fixed: true})
	for _, sx := range []float64{-1, 1} {
		b.Face([]V{{sx * w, h, S + 0.4}, {sx * (w + rim), h, S + 0.4}, {sx * (w + rim), -h, S + 0.4}, {sx * w, -h, S + 0.4}},
			"steel", &FOpts{Shade: 4.4, Fixed: true})
	}
	box(b, 0, 140, 140, 60, 22, 60, "steel")
	box(b, 0, -140, 140, 60, 22, 60, "steel")
	box(b, 140, 0, -140, 46, 46, 46, "steel")
	box(b, -140, 0, -140, 46, 46, 46, "steel")
	box(b, 0, 0, -S+12, 60, 60, 18, "steel")
	// masts on the corner towers, each with a beacon: the silhouette is what
	// reads as a station at distance
	for _, sx := range []float64{-1, 1} {
		for _, sy := range []float64{-1, 1} {
			box(b, sx*140, 200, sy*140, 4, 100, 4, "steel")
			box3(b, sx*140, 255, sy*140, 9, 9, 9, "red", &FOpts{Shade: 4.4, Fixed: true})
			box(b, sx*140, 168, sy*140, 26, 8, 26, "ash")
		}
	}
	m := b.Mesh()
	// grazing faces were being culled along the silhouette, leaving a one pixel
	// crack that showed space through the hull. The station is only 14 hull
	// faces, so drawing them two-sided costs nothing and closes the seam.
	for i := range m.Faces {
		m.Faces[i].Two = true
	}
	return m
}

const (
	StationSize  = 220.0
	StationSlotW = 0.33
	StationSlotH = 0.095
)

var Meshes = buildShips()

func Get(name string) *render.Mesh { return Meshes[name] }
