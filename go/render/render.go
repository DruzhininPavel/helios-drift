// Package render implements the software framebuffer, palette, rasterizer,
// bitmap font and 3D mesh renderer used by ELITE 16.
package render

import (
	"math"
	"unicode"
)

// ---------------------------------------------------------------- vectors

type Vec3 [3]float64

func V3(x, y, z float64) Vec3     { return Vec3{x, y, z} }
func (a Vec3) Add(b Vec3) Vec3    { return Vec3{a[0] + b[0], a[1] + b[1], a[2] + b[2]} }
func (a Vec3) Sub(b Vec3) Vec3    { return Vec3{a[0] - b[0], a[1] - b[1], a[2] - b[2]} }
func (a Vec3) Mul(s float64) Vec3 { return Vec3{a[0] * s, a[1] * s, a[2] * s} }
func (a Vec3) Dot(b Vec3) float64 { return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] }
func (a Vec3) Cross(b Vec3) Vec3 {
	return Vec3{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}
func (a Vec3) Len() float64 { return math.Sqrt(a.Dot(a)) }
func (a Vec3) Norm() Vec3 {
	l := a.Len()
	if l < 1e-9 {
		return Vec3{0, 0, -1}
	}
	return Vec3{a[0] / l, a[1] / l, a[2] / l}
}

// Mat3 stores the local axes in its columns (column 0 = right, 1 = up,
// 2 = backwards), so Apply is local->world and ApplyT is world->local.
type Mat3 [3][3]float64

func MatIdentity() Mat3 { return Mat3{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}} }

func (m Mat3) Apply(v Vec3) Vec3 {
	return Vec3{
		m[0][0]*v[0] + m[0][1]*v[1] + m[0][2]*v[2],
		m[1][0]*v[0] + m[1][1]*v[1] + m[1][2]*v[2],
		m[2][0]*v[0] + m[2][1]*v[1] + m[2][2]*v[2],
	}
}

func (m Mat3) ApplyT(v Vec3) Vec3 {
	return Vec3{
		m[0][0]*v[0] + m[1][0]*v[1] + m[2][0]*v[2],
		m[0][1]*v[0] + m[1][1]*v[1] + m[2][1]*v[2],
		m[0][2]*v[0] + m[1][2]*v[1] + m[2][2]*v[2],
	}
}

func (m Mat3) Mul(b Mat3) Mat3 {
	var r Mat3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			var s float64
			for k := 0; k < 3; k++ {
				s += m[i][k] * b[k][j]
			}
			r[i][j] = s
		}
	}
	return r
}

func (m Mat3) Col(i int) Vec3 { return Vec3{m[0][i], m[1][i], m[2][i]} }

func RotX(a float64) Mat3 {
	c, s := math.Cos(a), math.Sin(a)
	return Mat3{{1, 0, 0}, {0, c, -s}, {0, s, c}}
}
func RotY(a float64) Mat3 {
	c, s := math.Cos(a), math.Sin(a)
	return Mat3{{c, 0, s}, {0, 1, 0}, {-s, 0, c}}
}
func RotZ(a float64) Mat3 {
	c, s := math.Cos(a), math.Sin(a)
	return Mat3{{c, -s, 0}, {s, c, 0}, {0, 0, 1}}
}

// Orthonormalize rebuilds a rotation matrix from its first two columns.
func (m Mat3) Orthonormalize() Mat3 {
	x := Vec3{m[0][0], m[1][0], m[2][0]}.Norm()
	y := Vec3{m[0][1], m[1][1], m[2][1]}
	z := x.Cross(y).Norm()
	y = z.Cross(x)
	return Mat3{
		{x[0], y[0], z[0]},
		{x[1], y[1], z[1]},
		{x[2], y[2], z[2]},
	}
}

// LookMatrix builds an orientation whose -Z (forward) points along dir.
func LookMatrix(dir Vec3) Mat3 {
	zAxis := dir.Norm().Mul(-1)
	up := Vec3{0, 1, 0}
	if math.Abs(zAxis.Dot(up)) > 0.98 {
		up = Vec3{1, 0, 0}
	}
	x := up.Cross(zAxis).Norm()
	y := zAxis.Cross(x)
	return Mat3{
		{x[0], y[0], zAxis[0]},
		{x[1], y[1], zAxis[1]},
		{x[2], y[2], zAxis[2]},
	}
}

// ---------------------------------------------------------------- palette

var RampNames = []string{"steel", "ash", "red", "orange", "yellow", "green",
	"teal", "blue", "purple", "brown", "pink", "alien"}

var Ramps = map[string][6]uint32{
	"steel":  {0xff100a07, 0xff30221a, 0xff5a4133, 0xff886a55, 0xffc0a38f, 0xfff7e6db},
	"ash":    {0xff0c0a0a, 0xff2a2322, 0xff4f4543, 0xff7c6e6b, 0xffb2a29e, 0xfff7ece8},
	"red":    {0xff060417, 0xff140c40, 0xff261a7a, 0xff3434b8, 0xff4e6aee, 0xffa0c0ff},
	"orange": {0xff020a1a, 0xff06214d, 0xff10438c, 0xff1c76c9, 0xff44aef5, 0xffa8e1ff},
	"yellow": {0xff02151c, 0xff083a4d, 0xff186d8a, 0xff3aa6c9, 0xff6ad9f2, 0xffc0f8ff},
	"green":  {0xff0a1203, 0xff1c330d, 0xff32601a, 0xff4e942c, 0xff7cd25e, 0xffcfffc0},
	"teal":   {0xff1a1402, 0xff423407, 0xff74620f, 0xffa8931e, 0xffddcf3e, 0xfffff6b6},
	"blue":   {0xff1a0804, 0xff521d0e, 0xffa0371a, 0xffd8602f, 0xffff9b5f, 0xffffdcc4},
	"purple": {0xff20040d, 0xff4a0e26, 0xff861d48, 0xffc23671, 0xfff066a5, 0xffffc8e2},
	"brown":  {0xff050d15, 0xff112237, 0xff1e3c5e, 0xff31608b, 0xff528cb9, 0xffa0cae8},
	"pink":   {0xff12051a, 0xff300f4a, 0xff5c1f8a, 0xff8c38c8, 0xffb86af5, 0xffe6c6ff},
	"alien":  {0xff0c1406, 0xff2a4012, 0xff487a22, 0xff60a83a, 0xff7ae07c, 0xffb0ffe0},
}

// Palette holds every colour flattened as ramp*6+shade.
var Palette [72]uint32
var rampIndex = map[string]int{}

func init() {
	for i, name := range RampNames {
		rampIndex[name] = i
		for j := 0; j < 6; j++ {
			Palette[i*6+j] = Ramps[name][j]
		}
	}
}

// Color returns a palette entry for a ramp name and shade 0..5.
func Color(ramp string, shade int) uint32 {
	i, ok := rampIndex[ramp]
	if !ok {
		i = 0
	}
	if shade < 0 {
		shade = 0
	}
	if shade > 5 {
		shade = 5
	}
	return Palette[i*6+shade]
}

func ColorIdx(rampID, shade int) uint32 {
	if shade < 0 {
		shade = 0
	}
	if shade > 5 {
		shade = 5
	}
	if rampID < 0 || rampID > 11 {
		rampID = 0
	}
	return Palette[rampID*6+shade]
}

// RGB packs an 0xRRGGBB colour into the framebuffer's 0xAABBGGRR layout.
func RGB(hex uint32) uint32 {
	r := (hex >> 16) & 255
	g := (hex >> 8) & 255
	b := hex & 255
	return 0xff000000 | (b << 16) | (g << 8) | r
}

var (
	Space    = RGB(0x05070f)
	Bayer    = [16]int{0, 8, 2, 10, 12, 4, 14, 6, 3, 11, 1, 9, 15, 7, 13, 5}
	LightDir = Vec3{-0.40, 0.50, 0.77}.Norm()
	// HalfDir is the halfway vector between the light and the viewer, used for
	// the specular term.
	HalfDir = LightDir.Add(Vec3{0, 0, 1}).Norm()
)

// ExtraDetail switches the metal sheen and panel seams on; -plain turns it off
// so the two looks can be compared side by side.
var ExtraDetail = true

var (
	_ = ExtraDetail
)

// ------------------------------------------------------------------- font

var glyphs = map[rune][7]uint8{
	'0':  {0x0E, 0x11, 0x13, 0x15, 0x19, 0x11, 0x0E},
	'1':  {0x04, 0x0C, 0x04, 0x04, 0x04, 0x04, 0x0E},
	'2':  {0x0E, 0x11, 0x01, 0x02, 0x04, 0x08, 0x1F},
	'3':  {0x1E, 0x01, 0x01, 0x0E, 0x01, 0x01, 0x1E},
	'4':  {0x02, 0x06, 0x0A, 0x12, 0x1F, 0x02, 0x02},
	'5':  {0x1F, 0x10, 0x1E, 0x01, 0x01, 0x11, 0x0E},
	'6':  {0x06, 0x08, 0x10, 0x1E, 0x11, 0x11, 0x0E},
	'7':  {0x1F, 0x01, 0x02, 0x04, 0x08, 0x08, 0x08},
	'8':  {0x0E, 0x11, 0x11, 0x0E, 0x11, 0x11, 0x0E},
	'9':  {0x0E, 0x11, 0x11, 0x0F, 0x01, 0x02, 0x0C},
	'A':  {0x0E, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11},
	'B':  {0x1E, 0x11, 0x11, 0x1E, 0x11, 0x11, 0x1E},
	'C':  {0x0E, 0x11, 0x10, 0x10, 0x10, 0x11, 0x0E},
	'D':  {0x1C, 0x12, 0x11, 0x11, 0x11, 0x12, 0x1C},
	'E':  {0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x1F},
	'F':  {0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x10},
	'G':  {0x0E, 0x11, 0x10, 0x17, 0x11, 0x11, 0x0F},
	'H':  {0x11, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11},
	'I':  {0x0E, 0x04, 0x04, 0x04, 0x04, 0x04, 0x0E},
	'J':  {0x07, 0x02, 0x02, 0x02, 0x02, 0x12, 0x0C},
	'K':  {0x11, 0x12, 0x14, 0x18, 0x14, 0x12, 0x11},
	'L':  {0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x1F},
	'M':  {0x11, 0x1B, 0x15, 0x15, 0x11, 0x11, 0x11},
	'N':  {0x11, 0x11, 0x19, 0x15, 0x13, 0x11, 0x11},
	'O':  {0x0E, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E},
	'P':  {0x1E, 0x11, 0x11, 0x1E, 0x10, 0x10, 0x10},
	'Q':  {0x0E, 0x11, 0x11, 0x11, 0x15, 0x12, 0x0D},
	'R':  {0x1E, 0x11, 0x11, 0x1E, 0x14, 0x12, 0x11},
	'S':  {0x0F, 0x10, 0x10, 0x0E, 0x01, 0x01, 0x1E},
	'T':  {0x1F, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04},
	'U':  {0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E},
	'V':  {0x11, 0x11, 0x11, 0x11, 0x11, 0x0A, 0x04},
	'W':  {0x11, 0x11, 0x11, 0x15, 0x15, 0x1B, 0x11},
	'X':  {0x11, 0x11, 0x0A, 0x04, 0x0A, 0x11, 0x11},
	'Y':  {0x11, 0x11, 0x0A, 0x04, 0x04, 0x04, 0x04},
	'Z':  {0x1F, 0x01, 0x02, 0x04, 0x08, 0x10, 0x1F},
	' ':  {0, 0, 0, 0, 0, 0, 0},
	'.':  {0, 0, 0, 0, 0, 0x0C, 0x0C},
	',':  {0, 0, 0, 0, 0x06, 0x04, 0x08},
	':':  {0, 0x04, 0x04, 0, 0x04, 0x04, 0},
	';':  {0, 0x04, 0x04, 0, 0x04, 0x04, 0x08},
	'-':  {0, 0, 0, 0x1F, 0, 0, 0},
	'_':  {0, 0, 0, 0, 0, 0, 0x1F},
	'+':  {0, 0x04, 0x04, 0x1F, 0x04, 0x04, 0},
	'/':  {0x01, 0x02, 0x02, 0x04, 0x08, 0x08, 0x10},
	'\\': {0x10, 0x08, 0x08, 0x04, 0x02, 0x02, 0x01},
	'%':  {0x19, 0x1A, 0x02, 0x04, 0x08, 0x0B, 0x13},
	'(':  {0x02, 0x04, 0x08, 0x08, 0x08, 0x04, 0x02},
	')':  {0x08, 0x04, 0x02, 0x02, 0x02, 0x04, 0x08},
	'[':  {0x0E, 0x08, 0x08, 0x08, 0x08, 0x08, 0x0E},
	']':  {0x0E, 0x02, 0x02, 0x02, 0x02, 0x02, 0x0E},
	'*':  {0, 0x04, 0x15, 0x0E, 0x15, 0x04, 0},
	'\'': {0x04, 0x04, 0, 0, 0, 0, 0},
	'!':  {0x04, 0x04, 0x04, 0x04, 0x04, 0, 0x04},
	'?':  {0x0E, 0x11, 0x01, 0x06, 0x04, 0, 0x04},
	'=':  {0, 0, 0x1F, 0, 0x1F, 0, 0},
	'>':  {0x10, 0x08, 0x04, 0x02, 0x04, 0x08, 0x10},
	'<':  {0x01, 0x02, 0x04, 0x08, 0x04, 0x02, 0x01},
	'^':  {0x04, 0x0A, 0x11, 0, 0, 0, 0},
	'@':  {0x0E, 0x11, 0x17, 0x15, 0x17, 0x10, 0x0E},
}

// ------------------------------------------------------------------- FB

type FB struct {
	W, H           int
	Buf            []uint32
	Z              []float32
	X0, Y0, X1, Y1 int
}

func NewFB(w, h int) *FB {
	fb := &FB{W: w, H: h, Buf: make([]uint32, w*h), Z: make([]float32, w*h)}
	fb.ResetClip()
	return fb
}

func (fb *FB) ResetClip() {
	fb.X0, fb.Y0, fb.X1, fb.Y1 = 0, 0, fb.W-1, fb.H-1
}

func (fb *FB) SetClip(x0, y0, x1, y1 int) {
	fb.X0 = clamp(x0, 0, fb.W-1)
	fb.Y0 = clamp(y0, 0, fb.H-1)
	fb.X1 = clamp(x1, 0, fb.W-1)
	fb.Y1 = clamp(y1, 0, fb.H-1)
}

func clamp(v, a, b int) int {
	if v < a {
		return a
	}
	if v > b {
		return b
	}
	return v
}

func clampf(v, a, b float64) float64 {
	if v < a {
		return a
	}
	if v > b {
		return b
	}
	return v
}

func (fb *FB) Clear(c uint32) {
	for i := range fb.Buf {
		fb.Buf[i] = c
	}
	for i := range fb.Z {
		fb.Z[i] = 0
	}
}

func (fb *FB) ClearGradient(top, bottom uint32) {
	tr := int((top >> 0) & 255)
	tg := int((top >> 8) & 255)
	tb := int((top >> 16) & 255)
	br := int((bottom >> 0) & 255)
	bg := int((bottom >> 8) & 255)
	bb := int((bottom >> 16) & 255)
	for y := 0; y < fb.H; y++ {
		t := float64(y) / float64(fb.H-1)
		dith := (t*8 - math.Floor(t*8))
		for x := 0; x < fb.W; x++ {
			useBottom := float64(Bayer[(y&3)*4+(x&3)]) < dith*16
			var r, g, b int
			if useBottom {
				r, g, b = br, bg, bb
			} else {
				r, g, b = tr, tg, tb
			}
			fb.Buf[y*fb.W+x] = 0xff000000 | uint32(b)<<16 | uint32(g)<<8 | uint32(r)
		}
	}
	for i := range fb.Z {
		fb.Z[i] = 0
	}
}

func (fb *FB) Px(x, y int, c uint32) {
	if x < fb.X0 || x > fb.X1 || y < fb.Y0 || y > fb.Y1 {
		return
	}
	fb.Buf[y*fb.W+x] = c
}

func (fb *FB) HLine(x0, x1, y int, c uint32) {
	if y < fb.Y0 || y > fb.Y1 {
		return
	}
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if x0 < fb.X0 {
		x0 = fb.X0
	}
	if x1 > fb.X1 {
		x1 = fb.X1
	}
	for x := x0; x <= x1; x++ {
		fb.Buf[y*fb.W+x] = c
	}
}

func (fb *FB) VLine(x, y0, y1 int, c uint32) {
	if x < fb.X0 || x > fb.X1 {
		return
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	if y0 < fb.Y0 {
		y0 = fb.Y0
	}
	if y1 > fb.Y1 {
		y1 = fb.Y1
	}
	for y := y0; y <= y1; y++ {
		fb.Buf[y*fb.W+x] = c
	}
}

func (fb *FB) Rect(x0, y0, x1, y1 int, c uint32) {
	fb.HLine(x0, x1, y0, c)
	fb.HLine(x0, x1, y1, c)
	fb.VLine(x0, y0, y1, c)
	fb.VLine(x1, y0, y1, c)
}

func (fb *FB) FillRect(x0, y0, x1, y1 int, c uint32) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	if x0 < fb.X0 {
		x0 = fb.X0
	}
	if x1 > fb.X1 {
		x1 = fb.X1
	}
	if y0 < fb.Y0 {
		y0 = fb.Y0
	}
	if y1 > fb.Y1 {
		y1 = fb.Y1
	}
	for y := y0; y <= y1; y++ {
		row := y * fb.W
		for x := x0; x <= x1; x++ {
			fb.Buf[row+x] = c
		}
	}
}

func (fb *FB) Line(x0, y0, x1, y1 int, c uint32) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy
	for i := 0; i < 4000; i++ {
		fb.Px(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// Span is one horizontal run of a polygon with linear 1/z.
type Span struct {
	X0, X1   float64
	IZ0, IZ1 float64
}

// Spans walks a polygon and calls cb for every scanline span.
func (fb *FB) Spans(pts []Pt3, cb func(y int, s Span)) {
	n := len(pts)
	if n < 3 {
		return
	}
	minY, maxY := math.Inf(1), math.Inf(-1)
	hasZ := pts[0].Z != 0
	for _, p := range pts {
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	if math.IsInf(minY, 1) || math.IsInf(maxY, -1) {
		return
	}
	y0 := clamp(int(math.Ceil(minY)), fb.Y0, fb.Y1)
	y1 := clamp(int(math.Floor(maxY)), fb.Y0, fb.Y1)
	type isect struct{ x, iz float64 }
	var xs []isect
	for y := y0; y <= y1; y++ {
		sy := float64(y) + 0.5
		xs = xs[:0]
		for i := 0; i < n; i++ {
			a := pts[i]
			b := pts[(i+1)%n]
			if (a.Y <= sy && b.Y > sy) || (b.Y <= sy && a.Y > sy) {
				t := (sy - a.Y) / (b.Y - a.Y)
				iz := 0.0
				if hasZ {
					za, zb := a.Z, b.Z
					if za < 0.001 {
						za = 0.001
					}
					if zb < 0.001 {
						zb = 0.001
					}
					iz = (1-t)/za + t/zb
				}
				xs = append(xs, isect{a.X + t*(b.X-a.X), iz})
			}
		}
		if len(xs) < 2 {
			continue
		}
		for i := 1; i < len(xs); i++ {
			for j := i; j > 0 && xs[j].x < xs[j-1].x; j-- {
				xs[j], xs[j-1] = xs[j-1], xs[j]
			}
		}
		for k := 0; k+1 < len(xs); k += 2 {
			cb(y, Span{xs[k].x, xs[k+1].x, xs[k].iz, xs[k+1].iz})
		}
	}
}

func (fb *FB) Poly(pts []Pt3, c uint32) {
	fb.Spans(pts, func(y int, s Span) {
		xa := clamp(int(math.Ceil(s.X0-0.5)), fb.X0, fb.X1)
		xb := clamp(int(math.Floor(s.X1-0.5)), fb.X0, fb.X1)
		if xb < xa {
			return
		}
		row := y * fb.W
		for x := xa; x <= xb; x++ {
			fb.Buf[row+x] = c
		}
	})
}

// PolyShaded fills with an ordered-dithered two-tone ramp shade, depth tested
// when the points carry Z.
func (fb *FB) PolyShaded(pts []Pt3, rampID int, shade float64) {
	s := clampf(shade, 0, 5)
	i0 := int(math.Floor(s))
	if i0 > 5 {
		i0 = 5
	}
	f := s - float64(i0)
	i1 := i0 + 1
	if i1 > 5 {
		i1 = 5
	}
	c0 := ColorIdx(rampID, i0)
	c1 := ColorIdx(rampID, i1)
	th := f * 16
	flat := th < 0.5 || th > 15.5
	flatCol := c0
	if th > 15.5 {
		flatCol = c1
	}
	depth := pts[0].Z != 0
	buf, z := fb.Buf, fb.Z
	fb.Spans(pts, func(y int, sp Span) {
		xa := clamp(int(math.Ceil(sp.X0-0.5)), fb.X0, fb.X1)
		xb := clamp(int(math.Floor(sp.X1-0.5)), fb.X0, fb.X1)
		if xb < xa {
			return
		}
		row := y * fb.W
		dx := sp.X1 - sp.X0
		step := 0.0
		if dx > 0.0001 {
			step = (sp.IZ1 - sp.IZ0) / dx
		}
		iz := sp.IZ0 + (float64(xa)-sp.X0)*step
		if flat {
			for x := xa; x <= xb; x++ {
				idx := row + x
				if !depth || float32(iz) > z[idx] {
					if depth {
						z[idx] = float32(iz)
					}
					buf[idx] = flatCol
				}
				iz += step
			}
			return
		}
		by := (y & 3) * 4
		for x := xa; x <= xb; x++ {
			idx := row + x
			if !depth || float32(iz) > z[idx] {
				if depth {
					z[idx] = float32(iz)
				}
				if float64(Bayer[by+(x&3)]) < th {
					buf[idx] = c1
				} else {
					buf[idx] = c0
				}
			}
			iz += step
		}
	})
}

func (fb *FB) PolyGrad(pts []Pt3, cTop, cBot uint32) {
	minY, maxY := math.Inf(1), math.Inf(-1)
	for _, p := range pts {
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	if maxY-minY < 1 {
		fb.Poly(pts, cTop)
		return
	}
	tr := int((cTop >> 0) & 255)
	tg := int((cTop >> 8) & 255)
	tb := int((cTop >> 16) & 255)
	br := int((cBot >> 0) & 255)
	bg := int((cBot >> 8) & 255)
	bb := int((cBot >> 16) & 255)
	fb.Spans(pts, func(y int, s Span) {
		xa := clamp(int(math.Ceil(s.X0-0.5)), fb.X0, fb.X1)
		xb := clamp(int(math.Floor(s.X1-0.5)), fb.X0, fb.X1)
		if xb < xa {
			return
		}
		t := (float64(y) - minY) / (maxY - minY)
		dith := math.Mod(t*8, 1)
		row := y * fb.W
		for x := xa; x <= xb; x++ {
			var r, g, b int
			if float64(Bayer[(y&3)*4+(x&3)]) < dith*16 {
				r, g, b = br, bg, bb
			} else {
				r, g, b = tr, tg, tb
			}
			fb.Buf[row+x] = 0xff000000 | uint32(b)<<16 | uint32(g)<<8 | uint32(r)
		}
	})
}

func (fb *FB) EllipseFill(cx, cy, rx, ry int, c uint32) {
	for y := -ry; y <= ry; y++ {
		t := 1 - float64(y*y)/float64(ry*ry)
		if t < 0 {
			continue
		}
		half := math.Sqrt(t) * float64(rx)
		fb.HLine(int(math.Round(float64(cx)-half)), int(math.Round(float64(cx)+half)), cy+y, c)
	}
}

func (fb *FB) Ellipse(cx, cy, rx, ry int, c uint32) {
	steps := 2 * (rx + ry)
	if steps < 24 {
		steps = 24
	}
	for i := 0; i < steps; i++ {
		a := float64(i) / float64(steps) * 2 * math.Pi
		fb.Px(int(math.Round(float64(cx)+math.Cos(a)*float64(rx))),
			int(math.Round(float64(cy)+math.Sin(a)*float64(ry))), c)
	}
}

func (fb *FB) EllipseDash(cx, cy, rx, ry int, c uint32, dash int) {
	steps := 2 * (rx + ry)
	if steps < 32 {
		steps = 32
	}
	for i := 0; i < steps; i++ {
		if i%(dash*2) >= dash {
			continue
		}
		a := float64(i) / float64(steps) * 2 * math.Pi
		fb.Px(int(math.Round(float64(cx)+math.Cos(a)*float64(rx))),
			int(math.Round(float64(cy)+math.Sin(a)*float64(ry))), c)
	}
}

// Translate is installed by the game to localise every drawn string.
var Translate func(string) string

func tr(s string) string {
	if Translate != nil {
		return Translate(s)
	}
	return s
}

func (fb *FB) Text(s string, x, y int, c uint32, scale int) {
	fb.textRaw(tr(s), x, y, c, scale)
}

// TextRaw draws without localisation (used for the ELITE wordmark).
func (fb *FB) TextRaw(s string, x, y int, c uint32, scale int) {
	fb.textRaw(s, x, y, c, scale)
}

func (fb *FB) textRaw(s string, x, y int, c uint32, scale int) {
	if scale < 1 {
		scale = 1
	}
	cx := x
	for _, ch := range s {
		if ch == '\n' {
			y += 8 * scale
			cx = x
			continue
		}
		gl, ok := glyphs[ch]
		if !ok {
			// the font is capitals only: render lower case as capitals rather
			// than dumping a row of question marks into the panel
			gl, ok = glyphs[unicode.ToUpper(ch)]
		}
		if !ok {
			gl = glyphs['?']
		}
		for row := 0; row < 7; row++ {
			bits := gl[row]
			for col := 0; col < 5; col++ {
				if bits&(1<<(4-uint(col))) == 0 {
					continue
				}
				if scale == 1 {
					fb.Px(cx+col, y+row, c)
				} else {
					fb.FillRect(cx+col*scale, y+row*scale, cx+col*scale+scale-1, y+row*scale+scale-1, c)
				}
			}
		}
		cx += 6 * scale
	}
}

func TextWidth(s string, scale int) int { return TextWidthRaw(tr(s), scale) }

// TextWidthRaw measures without localisation.
func TextWidthRaw(s string, scale int) int {
	if scale < 1 {
		scale = 1
	}
	return len([]rune(s))*6*scale - scale
}

func (fb *FB) TextC(s string, cx, y int, c uint32, scale int) {
	fb.Text(s, cx-TextWidth(s, scale)/2, y, c, scale)
}

func (fb *FB) TextR(s string, rx, y int, c uint32, scale int) {
	fb.Text(s, rx-TextWidth(s, scale), y, c, scale)
}

// RampIndex returns the palette ramp id for a ramp name.
func RampIndex(name string) (int, bool) {
	i, ok := rampIndex[name]
	return i, ok
}

// PolyShadedGrad fills a polygon with a dithered linear shade ramp between
// shadeA and shadeB along the screen-space direction (gx, gy).
func (fb *FB) PolyShadedGrad(pts []Pt3, rampID int, shadeA, shadeB, gx, gy float64) {
	minX, maxX := math.Inf(1), math.Inf(-1)
	minY, maxY := math.Inf(1), math.Inf(-1)
	for _, p := range pts {
		minX = math.Min(minX, p.X)
		maxX = math.Max(maxX, p.X)
		minY = math.Min(minY, p.Y)
		maxY = math.Max(maxY, p.Y)
	}
	cx := (minX + maxX) / 2
	cy := (minY + maxY) / 2
	span := math.Max(maxX-minX, maxY-minY) * 0.5
	if span < 0.75 {
		span = 0.75
	}
	depth := pts[0].Z != 0
	buf, z := fb.Buf, fb.Z
	fb.Spans(pts, func(y int, sp Span) {
		xa := clamp(int(math.Ceil(sp.X0-0.5)), fb.X0, fb.X1)
		xb := clamp(int(math.Floor(sp.X1-0.5)), fb.X0, fb.X1)
		if xb < xa {
			return
		}
		row := y*fb.W + xa
		dx := sp.X1 - sp.X0
		step := 0.0
		if dx > 0.0001 {
			step = (sp.IZ1 - sp.IZ0) / dx
		}
		iz := sp.IZ0 + (float64(xa)-sp.X0)*step
		by := (y & 3) * 4
		fy := (float64(y) - cy) * gy
		for x := xa; x <= xb; x++ {
			idx := row + (x - xa)
			if !depth || float32(iz) > z[idx] {
				if depth {
					z[idx] = float32(iz)
				}
				t := 0.5 + ((float64(x)-cx)*gx+fy)/span*0.5
				if t < 0 {
					t = 0
				} else if t > 1 {
					t = 1
				}
				sh := shadeA + (shadeB-shadeA)*t
				i0 := int(sh)
				if i0 > 5 {
					i0 = 5
				}
				if i0 < 0 {
					i0 = 0
				}
				i1 := i0 + 1
				if i1 > 5 {
					i1 = 5
				}
				if (sh-float64(i0))*16 > float64(Bayer[by+(x&3)]) {
					buf[idx] = ColorIdx(rampID, i1)
				} else {
					buf[idx] = ColorIdx(rampID, i0)
				}
			}
			iz += step
		}
	})
}

// ClearGradientStops fills the framebuffer with a multi-stop dithered sky,
// dithering only between neighbouring stops (as the reference build does).
func (fb *FB) ClearGradientStops(stops []uint32) {
	n := len(stops)
	if n < 2 {
		fb.Clear(stops[0])
		return
	}
	for y := 0; y < fb.H; y++ {
		t := float64(y) / float64(fb.H-1) * float64(n-1)
		i0 := int(t)
		if i0 > n-2 {
			i0 = n - 2
		}
		f := t - float64(i0)
		a, b := stops[i0], stops[i0+1]
		ar, ag, ab := int(a&255), int((a>>8)&255), int((a>>16)&255)
		br, bg, bb := int(b&255), int((b>>8)&255), int((b>>16)&255)
		for x := 0; x < fb.W; x++ {
			useB := float64(Bayer[(y&3)*4+(x&3)]) < f*16
			r, g, bl := ar, ag, ab
			if useB {
				r, g, bl = br, bg, bb
			}
			fb.Buf[y*fb.W+x] = 0xff000000 | uint32(bl)<<16 | uint32(g)<<8 | uint32(r)
		}
	}
	for i := range fb.Z {
		fb.Z[i] = 0
	}
}

// LineDash draws a line with every second pixel skipped: a lighter touch than a
// solid stroke when marking plate seams at low resolution.
func (fb *FB) LineDash(x0, y0, x1, y1 int, c uint32) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy
	n := 0
	for {
		if n%2 == 0 {
			fb.Px(x0, y0, c)
		}
		n++
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}
