package render

import "math"

// Pt3 is a projected point: screen X, Y and view depth (1/Z is used for the
// depth buffer, Z == 0 means "2D, skip depth testing").
type Pt3 struct {
	X, Y float64
	Z    float64
}

type Face struct {
	P     []Vec3
	N     Vec3
	Ramp  int
	Shade float64
	Tint  float64
	Fixed bool
	Two   bool
	// AO is a baked occlusion factor: faces that meet a neighbour in a concave
	// crease are darkened, which is what makes assembled parts read as parts.
	AO float64
}

// Mesh is a flat-shaded low-poly model.
type Mesh struct {
	Faces  []Face
	Radius float64
	// Exhaust lists local-space engine positions (drawn as flares in flight).
	Exhaust []Vec3
}

type Camera struct {
	W, H   int
	CX, CY float64
	Focal  float64
	Near   float64
	Pos    Vec3
	Mat    Mat3
}

func NewCamera(w, h int, fovDeg float64) *Camera {
	c := &Camera{W: w, H: h, Focal: (float64(h) / 2) / math.Tan(fovDeg*math.Pi/360), Near: 3}
	c.CX = float64(w) / 2
	c.CY = float64(h) / 2
	c.Mat = MatIdentity()
	return c
}

func (c *Camera) ToView(world Vec3) Vec3 {
	return c.Mat.ApplyT(world.Sub(c.Pos))
}

func (c *Camera) ProjectView(v Vec3) Pt3 {
	depth := -v[2]
	if depth < 0.001 {
		depth = 0.001
	}
	return Pt3{c.CX + v[0]*c.Focal/depth, c.CY - v[1]*c.Focal/depth, depth}
}

func (c *Camera) Project(world Vec3) Pt3 { return c.ProjectView(c.ToView(world)) }

// ClipNear keeps the part of the polygon in front of the near plane.
func ClipNear(poly []Vec3, near float64) []Vec3 {
	lim := -near
	out := make([]Vec3, 0, len(poly)+2)
	n := len(poly)
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		ain := a[2] <= lim
		bin := b[2] <= lim
		if ain {
			out = append(out, a)
		}
		if ain != bin {
			t := (lim - a[2]) / (b[2] - a[2])
			out = append(out, Vec3{a[0] + (b[0]-a[0])*t, a[1] + (b[1]-a[1])*t, lim})
		}
	}
	return out
}

type drawFace struct {
	pv   []Vec3
	z    float64
	ramp int
	sh   float64
	rim  bool
	area float64
	spec float64
}

// DrawMesh renders a mesh instance and returns true if anything was on screen.
func DrawMesh(fb *FB, cam *Camera, mesh *Mesh, pos Vec3, mat Mat3, fogDist float64) bool {
	if fogDist <= 0 {
		fogDist = 60000
	}
	viewPos := cam.ToView(pos)
	if viewPos[2] > mesh.Radius+100 {
		return false
	}
	list := make([]drawFace, 0, len(mesh.Faces))
	for fi := range mesh.Faces {
		f := &mesh.Faces[fi]
		pv := make([]Vec3, len(f.P))
		var centre Vec3
		culled := false
		for i, lp := range f.P {
			wv := mat.Apply(lp).Add(pos)
			vv := cam.ToView(wv)
			pv[i] = vv
			centre = centre.Add(vv)
		}
		centre = centre.Mul(1 / float64(len(pv)))
		worldN := mat.Apply(f.N)
		viewN := cam.Mat.ApplyT(worldN)
		toFace := centre.Norm()
		if viewN.Dot(toFace) > 0 && !f.Two {
			culled = true
		}
		if culled {
			continue
		}
		nrm := viewN.Norm()
		lam := nrm.Dot(LightDir)
		sh := 1.15 + math.Max(0, lam)*2.9
		// metal sheen: a narrow highlight where the facet lines up with the
		// light, which is what makes flat-shaded hulls read as metal
		spec := 0.0
		if ExtraDetail {
			if half := nrm.Dot(HalfDir); half > 0 {
				spec = math.Pow(half, 26)
			}
		}
		// per-face value jitter: coplanar neighbours sharing one tone is what
		// makes flat-shaded hulls look like untextured plastic
		if !f.Fixed && ExtraDetail {
			h := (uint32(fi)*2654435761 + 7) >> 8
			sh += (float64(h&255)/255.0)*0.5 - 0.25
		}
		rim := math.Abs(viewN.Dot(toFace)) < 0.3
		if f.AO > 0 && !f.Fixed {
			sh *= 1 - 0.22*f.AO
		}
		if f.Fixed {
			sh = f.Shade
		}
		if f.Tint != 0 {
			sh *= f.Tint
		}
		dist := centre.Len()
		fog := 1 - dist/fogDist
		if fog < 0 {
			fog = 0
		}
		sh *= 0.80 + 0.20*fog
		list = append(list, drawFace{pv: pv, z: centre[2], ramp: f.Ramp, sh: sh, rim: rim, area: faceArea(pv), spec: spec})
	}
	// painter's order: most negative view z (farthest) first
	for i := 1; i < len(list); i++ {
		for j := i; j > 0 && list[j].z < list[j-1].z; j-- {
			list[j], list[j-1] = list[j-1], list[j]
		}
	}
	drawn := false
	for _, it := range list {
		allBehind := true
		for _, p := range it.pv {
			if p[2] <= -cam.Near {
				allBehind = false
				break
			}
		}
		if allBehind {
			continue
		}
		clipped := ClipNear(it.pv, cam.Near)
		if len(clipped) < 3 {
			continue
		}
		pts := make([]Pt3, len(clipped))
		ok := true
		for i, p := range clipped {
			pp := cam.ProjectView(p)
			if math.IsNaN(pp.X) || math.IsNaN(pp.Y) || math.IsInf(pp.X, 0) || math.IsInf(pp.Y, 0) {
				ok = false
				break
			}
			pts[i] = pp
		}
		if !ok {
			continue
		}
		// two-tone gradient across the facet plus a rim light on silhouette faces
		shA := it.sh - 1.15
		shB := it.sh + 0.95
		if it.spec > 0 {
			shA += it.spec * 0.30
			shB += it.spec * 0.85
		}
		if shB > 4.35 {
			shB = 4.35 // leave the top step for lights and the sun
		}
		if shA > shB-0.35 {
			shA = shB - 0.35
		}
		gx, gy := LightDir[0], -LightDir[1]
		gl := math.Hypot(gx, gy)
		if gl > 0.001 {
			gx /= gl
			gy /= gl
		}
		fb.PolyShadedGrad(pts, it.ramp, shA, shB, gx, gy)
		if it.rim {
			edge := ColorIdx(it.ramp, int(math.Max(0, it.sh-1.1)))
			for i := 0; i < len(pts); i++ {
				a := pts[i]
				b := pts[(i+1)%len(pts)]
				fb.Line(int(a.X), int(a.Y), int(b.X), int(b.Y), edge)
			}
		}
		drawn = true
	}
	return drawn
}

// DrawGlow renders a dithered radial glow (sun, engine flares).
// faceArea is the screen area of a projected polygon, used to decide whether a
// facet is big enough to be worth outlining.
func faceArea(pts []Vec3) float64 {
	if len(pts) < 3 {
		return 0
	}
	a := 0.0
	for i := 0; i < len(pts); i++ {
		j := (i + 1) % len(pts)
		a += pts[i][0]*pts[j][1] - pts[j][0]*pts[i][1]
	}
	return math.Abs(a) * 0.5
}

func DrawGlow(fb *FB, cx, cy, radius float64, ramp string, maxShade float64) {
	steps := 6
	for i := steps; i >= 1; i-- {
		r := radius * float64(i) / float64(steps)
		shade := 0.6 + (maxShade-0.6)*(1-float64(i)/float64(steps))*0.5
		c := Color(ramp, int(clampf(shade+(float64(steps-i))*0.55, 0, 5)))
		icx, icy := int(math.Round(cx)), int(math.Round(cy))
		for y := -int(r); y <= int(r); y++ {
			t := 1 - float64(y*y)/(r*r)
			if t < 0 {
				continue
			}
			half := math.Sqrt(t) * r
			yy := icy + y
			for x := icx - int(half); x <= icx+int(half); x++ {
				b := Bayer[((yy&3)*4 + (x & 3))]
				if i > 1 && b > 5 {
					continue
				}
				fb.Px(x, yy, c)
			}
		}
	}
}

// BlitScaled copies a sprite to a destination rectangle (nearest neighbour).
// If depth > 0 it also fills the depth buffer so ships behind it are hidden.
func BlitScaled(fb, src *FB, dx, dy, dw, dh float64, depth float64) {
	if dw <= 0 || dh <= 0 {
		return
	}
	x0 := int(math.Round(dx))
	y0 := int(math.Round(dy))
	w1 := int(math.Round(dw))
	h1 := int(math.Round(dh))
	sxStep := float64(src.W) / float64(w1)
	syStep := float64(src.H) / float64(h1)
	for y := 0; y < h1; y++ {
		ty := y0 + y
		if ty < fb.Y0 || ty > fb.Y1 {
			continue
		}
		syi := int(float64(y) * syStep)
		if syi >= src.H {
			syi = src.H - 1
		}
		row := syi * src.W
		for x := 0; x < w1; x++ {
			tx := x0 + x
			if tx < fb.X0 || tx > fb.X1 {
				continue
			}
			sxi := int(float64(x) * sxStep)
			if sxi >= src.W {
				sxi = src.W - 1
			}
			c := src.Buf[row+sxi]
			if c>>24 == 0 {
				continue
			}
			di := ty*fb.W + tx
			fb.Buf[di] = c
			if depth > 0 {
				fb.Z[di] = float32(1 / depth)
			}
		}
	}
}

// BlitPart copies a wrapped source rectangle with integer scaling.
func BlitPart(fb, src *FB, sx, sy, dw, dh, dx, dy, scale int) {
	if scale < 1 {
		scale = 1
	}
	for y := 0; y < dh; y++ {
		ty := dy + y
		if ty < fb.Y0 || ty > fb.Y1 {
			continue
		}
		syy := ((sy+y/scale)%src.H + src.H) % src.H
		for x := 0; x < dw; x++ {
			tx := dx + x
			if tx < fb.X0 || tx > fb.X1 {
				continue
			}
			sxx := ((sx+x/scale)%src.W + src.W) % src.W
			c := src.Buf[syy*src.W+sxx]
			if c>>24 == 0 {
				continue
			}
			fb.Buf[ty*fb.W+tx] = c
		}
	}
}
