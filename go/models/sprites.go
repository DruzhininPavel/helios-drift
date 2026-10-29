package models

import (
	"math"

	"heliosdrift/render"
)

// Rng is a small deterministic PRNG (mulberry32) matching the JS build.
type Rng struct{ a uint32 }

func NewRng(seed uint32) *Rng { return &Rng{a: seed} }

func (r *Rng) Float() float64 {
	r.a += 0x6D2B79F5
	t := r.a
	t = (t ^ (t >> 15)) * (t | 1)
	t ^= t + (t^(t>>7))*(t|61)
	return float64((t^(t>>14))>>0) / 4294967296.0
}

func hash3(x, y, z float64) float64 {
	n := math.Sin(x*127.1+y*311.7+z*74.7) * 43758.5453
	return n - math.Floor(n)
}

func vnoise(x, y, z float64) float64 {
	xi, yi, zi := math.Floor(x), math.Floor(y), math.Floor(z)
	xf, yf, zf := x-xi, y-yi, z-zi
	s := func(t float64) float64 { return t * t * (3 - 2*t) }
	u, v, w := s(xf), s(yf), s(zf)
	c := 0.0
	for dx := 0; dx <= 1; dx++ {
		for dy := 0; dy <= 1; dy++ {
			for dz := 0; dz <= 1; dz++ {
				h := hash3(xi+float64(dx), yi+float64(dy), zi+float64(dz))
				wx, wy, wz := 1-u, 1-v, 1-w
				if dx == 1 {
					wx = u
				}
				if dy == 1 {
					wy = v
				}
				if dz == 1 {
					wz = w
				}
				c += h * wx * wy * wz
			}
		}
	}
	return c
}

func fbm(x, y, z float64, oct int) float64 {
	s, a, f := 0.0, 0.5, 1.0
	for i := 0; i < oct; i++ {
		s += a * vnoise(x*f, y*f, z*f)
		a *= 0.5
		f *= 2.05
	}
	return s
}

type planetTex struct {
	w, h  int
	lum   []uint8
	ramp  []uint8
	ramps []string
}

var planetCache = map[string][]*render.FB{}

var rampOrder = []string{"blue", "teal", "green", "brown", "ash", "steel", "orange", "yellow", "red"}

func buildPlanetTexture(seed uint32, kind int) *planetTex {
	const W, H = 256, 128
	r := NewRng(seed)
	ox, oy, oz := r.Float()*100, r.Float()*100, r.Float()*100
	tex := &planetTex{w: W, h: H, lum: make([]uint8, W*H), ramp: make([]uint8, W*H), ramps: rampOrder}
	for j := 0; j < H; j++ {
		lat := (float64(j)/float64(H-1) - 0.5) * math.Pi
		cl, sl := math.Cos(lat), math.Sin(lat)
		for i := 0; i < W; i++ {
			lon := float64(i) / W * 2 * math.Pi
			nx, ny, nz := math.Sin(lon)*cl, sl, math.Cos(lon)*cl
			h := fbm(nx*2.1+ox, ny*2.1+oy, nz*2.1+oz, 3)
			al := math.Abs(ny)
			var ramp string
			var shade float64
			switch kind {
			case 3:
				band := math.Sin(ny*11+h*3.0)*0.5 + 0.5
				if band > 0.66 {
					ramp = "orange"
				} else if band > 0.36 {
					ramp = "brown"
				} else {
					ramp = "yellow"
				}
				shade = 1.6 + band*2.4
			case 2:
				if h > 0.5 {
					ramp = "steel"
				} else {
					ramp = "blue"
				}
				shade = 1.8 + h*2.2
			case 1:
				if h > 0.58 {
					ramp = "brown"
				} else if h > 0.44 {
					ramp = "orange"
				} else {
					ramp = "yellow"
				}
				shade = 1.6 + h*2.4
			case 4:
				if h > 0.6 {
					ramp = "red"
				} else if h > 0.4 {
					ramp = "orange"
				} else {
					ramp = "brown"
				}
				shade = 1.3 + h*2.8
			default:
				switch {
				case h < 0.46:
					ramp, shade = "blue", 1.1+h*3.2
				case h < 0.51:
					ramp, shade = "teal", 2.3
				case h < 0.63:
					ramp, shade = "green", 1.7+(h-0.51)*7
				case h < 0.71:
					ramp, shade = "brown", 2.1+(h-0.63)*6
				default:
					ramp, shade = "ash", 3.6
				}
				if al > 0.87 {
					ramp, shade = "steel", 3.3+h
				}
			}
			si := int(math.Round(shade))
			if si < 0 {
				si = 0
			}
			if si > 5 {
				si = 5
			}
			tex.lum[j*W+i] = uint8(si)
			rid := 0
			for k, name := range rampOrder {
				if name == ramp {
					rid = k
					break
				}
			}
			tex.ramp[j*W+i] = uint8(rid)
		}
	}
	return tex
}

// MakePlanetFrames renders (and caches) rotation frames for a planet sprite.
func MakePlanetFrames(size int, seed uint32, kind, frames int) []*render.FB {
	key := itoa(int(seed)) + ":" + itoa(kind) + ":" + itoa(size)
	if v, ok := planetCache[key]; ok {
		return v
	}
	if frames < 1 {
		frames = 8
	}
	tex := buildPlanetTexture(seed, kind)
	out := make([]*render.FB, frames)
	R := float64(size) / 2
	sun := render.V3(-0.55, 0.42, 0.72).Norm()
	for fr := 0; fr < frames; fr++ {
		fb := render.NewFB(size, size)
		rot := float64(fr) / float64(frames) * 2 * math.Pi
		for y := 0; y < size; y++ {
			ny := -(float64(y) - R) / R
			for x := 0; x < size; x++ {
				nx := (float64(x) - R) / R
				d2 := nx*nx + ny*ny
				if d2 > 1 {
					continue
				}
				nz := math.Sqrt(1 - d2)
				lam := nx*sun[0] + ny*sun[1] + nz*sun[2]
				if lam <= -0.03 {
					continue
				}
				lat := math.Asin(math.Max(-1, math.Min(1, ny)))
				lon := math.Atan2(nx, nz) + rot
				u := math.Mod(lon/(2*math.Pi), 1)
				if u < 0 {
					u += 1
				}
				v := lat/math.Pi + 0.5
				ti := int(u * float64(tex.w))
				if ti >= tex.w {
					ti = tex.w - 1
				}
				tj := int(v * float64(tex.h-1))
				if tj >= tex.h-1 {
					tj = tex.h - 1
				}
				idx := tj*tex.w + ti
				shade := float64(tex.lum[idx])*(0.20+0.80*math.Max(0, lam)) + 0.15
				i0 := int(math.Floor(shade))
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
				rampID, _ := render.RampIndex(tex.ramps[tex.ramp[idx]])
				b := render.Bayer[((y&3)*4 + (x & 3))]
				ci := i0
				if (shade-float64(i0))*16 > float64(b) {
					ci = i1
				}
				fb.Buf[y*size+x] = render.ColorIdx(rampID, ci)
			}
		}
		out[fr] = fb
	}
	planetCache[key] = out
	return out
}

var nebulaCache = map[string]*render.FB{}

// MakeNebula builds a 320x75 dithered cloud backdrop that pans with yaw.
func MakeNebula(seed uint32, galaxy int) *render.FB {
	key := itoa(int(seed)) + ":" + itoa(galaxy)
	if v, ok := nebulaCache[key]; ok {
		return v
	}
	const W, H = 320, 75
	fb := render.NewFB(W, H)
	r := NewRng(seed*7 + uint32(galaxy)*131)
	ox, oy := r.Float()*40, r.Float()*40
	ramp := "purple"
	if galaxy%2 == 0 {
		ramp = "purple"
	} else {
		ramp = "blue"
	}
	ramp2 := "blue"
	if galaxy%3 == 0 {
		ramp2 = "teal"
	}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			n := fbm(float64(x)*0.015+ox, float64(y)*0.032+oy, 1.7, 2)
			n2 := fbm(float64(x)*0.028+9, float64(y)*0.06+4, 4.1, 2)
			d := n*0.75 + n2*0.25
			if d < 0.46 {
				continue
			}
			t := (d - 0.46) / 0.22
			useRamp := ramp2
			if n2 > 0.5 {
				useRamp = ramp
			}
			shade := math.Min(2.2, 0.35+t*1.9)
			ci := int(shade)
			frac := shade - float64(ci)
			b := render.Bayer[((y&3)*4 + (x & 3))]
			if frac*16 < float64(b) {
				continue
			}
			rid, _ := render.RampIndex(useRamp)
			fb.Buf[y*W+x] = render.ColorIdx(rid, ci)
		}
	}
	nebulaCache[key] = fb
	return fb
}

// PlanetKindFor picks a world type from a system seed.
func PlanetKindFor(seed uint32) int {
	r := NewRng(seed).Float()
	switch {
	case r < 0.45:
		return 0
	case r < 0.62:
		return 1
	case r < 0.74:
		return 2
	case r < 0.92:
		return 3
	}
	return 4
}
