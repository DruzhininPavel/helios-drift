// Package ui draws the dashboard, the 3D view and the station screens.
package ui

import (
	"fmt"
	"math"
	"strconv"

	"heliosdrift/galaxy"
	"heliosdrift/game"
	"heliosdrift/models"
	"heliosdrift/render"
)

var (
	ColText    = render.RGB(0x9fe8ff)
	ColDim     = render.RGB(0x3f6f8c)
	ColBg      = render.RGB(0x0a1220)
	ColBg2     = render.RGB(0x111d2e)
	ColFrame   = render.RGB(0x2d5f7a)
	ColFrameHi = render.RGB(0x79c9e6)
	ColYellow  = render.RGB(0xffd45e)
	ColGreen   = render.RGB(0x6ef08a)
	ColRed     = render.RGB(0xff6a4e)
	ColWhite   = render.RGB(0xe8f6ff)
	ColOrange  = render.RGB(0xff9c3a)
)

type Sel struct {
	Market, Equip, Menu int
	Contract            int
	Ship                int
	Slot                int
	Mode                int
	SlotSave            bool
	Control             int
	ChartCursor         int
	Title               int
	Settings            int
	PrevState           string
	PrevScreen          string
	PadTestFrom         string
	HasSession          bool
	SessionState        string
	SessionScreen       string
}

var S = Sel{ChartCursor: -1}

// quitRequested is set by the title screen and read by the main loop, which
// is the only place that can end the run.
var quitRequested bool

// QuitRequested reports and clears a request to leave the game.
func QuitRequested() bool {
	q := quitRequested
	quitRequested = false
	return q
}

var Hint string

func money(v float64) string { return strconv.FormatFloat(math.Round(v*10)/10, 'f', 1, 64) }
func itoa(v int) string      { return strconv.Itoa(v) }

func plate(fb *render.FB, x, y, w, h int, title string) {
	fb.PolyGrad([]render.Pt3{{X: float64(x), Y: float64(y)}, {X: float64(x + w), Y: float64(y)},
		{X: float64(x + w), Y: float64(y + h)}, {X: float64(x), Y: float64(y + h)}}, ColBg2, ColBg)
	fb.Rect(x, y, x+w-1, y+h-1, ColFrame)
	fb.HLine(x+1, x+w-2, y+1, render.RGB(0x1b3348))
	fb.HLine(x, x+5, y, ColFrameHi)
	fb.HLine(x+w-6, x+w-1, y, ColFrameHi)
	if title != "" {
		fb.FillRect(x+4, y, x+4+render.TextWidth(title, 1)+6, y+8, ColBg)
		fb.Text(title, x+7, y-3, ColYellow, 1)
	}
}

func segBar(fb *render.FB, x, y, w, h int, frac float64, ramp string, warn bool) {
	if frac < 0 {
		frac = 0
	}
	if frac > 1 {
		frac = 1
	}
	segs := 8
	gap := 1
	sw := (w - gap*(segs-1)) / segs
	filled := frac * float64(segs)
	for i := 0; i < segs; i++ {
		sx := x + i*(sw+gap)
		f := filled - float64(i)
		if f < 0 {
			f = 0
		}
		if f > 1 {
			f = 1
		}
		if f <= 0 {
			fb.FillRect(sx, y, sx+sw-1, y+h-1, render.RGB(0x0d1a28))
			fb.HLine(sx, sx+sw-1, y, render.RGB(0x16283c))
			continue
		}
		shade := 1.6 + f*2.4
		c := render.Color(ramp, int(shade))
		c2 := render.Color(ramp, int(shade)+1)
		for yy := 0; yy < h; yy++ {
			for xx := 0; xx < sw; xx++ {
				b := render.Bayer[((yy&3)*4 + (xx & 3))]
				if f < 0.99 && b < 4 {
					fb.Px(sx+xx, y+yy, c2)
				} else {
					fb.Px(sx+xx, y+yy, c)
				}
			}
		}
	}
	if warn {
		fb.Rect(x-1, y-1, x+w, y+h, ColRed)
	}
}

func scanner(fb *render.FB, cx, cy, rx, ry int) {
	g := game.G
	p := &g.Player
	w := g.World
	rng := float64(game.RadarRange)
	fb.EllipseFill(cx, cy, rx, ry, render.RGB(0x04121a))
	fb.EllipseDash(cx, cy, rx, ry, render.RGB(0x1b5a6a), 3)
	fb.EllipseDash(cx, cy, rx*2/3, ry*2/3, render.RGB(0x123f4d), 2)
	fb.EllipseDash(cx, cy, rx/3, ry/3, render.RGB(0x123f4d), 2)
	fb.HLine(cx-rx, cx+rx, cy, render.RGB(0x123f4d))
	fb.VLine(cx, cy-ry, cy+ry, render.RGB(0x123f4d))
	fb.Ellipse(cx, cy, rx, ry, render.RGB(0x3fbcd8))

	plot := func(pos render.Vec3, col uint32, size int, stick bool) {
		v := p.Mat.ApplyT(pos.Sub(p.Pos))
		depth := -v[2]
		if depth < -600 {
			return
		}
		side := v[0] / rng
		fwdN := depth / rng
		if math.Abs(side) > 1 || math.Abs(fwdN) > 1 {
			return
		}
		sx := int(math.Round(float64(cx) + side*float64(rx)))
		sy := int(math.Round(float64(cy) - fwdN*float64(ry)))
		alt := math.Max(-1, math.Min(1, v[1]/(rng*0.45)))
		topY := int(math.Round(float64(sy) - alt*float64(ry)*0.55))
		if stick {
			fb.VLine(sx, minInt(sy, topY), maxInt(sy, topY), render.RGB(0x2b6f80))
		}
		if size == 1 {
			fb.Px(sx, topY, col)
		} else {
			fb.FillRect(sx-1, topY-1, sx, topY, col)
		}
	}

	if w.Station != nil {
		plot(w.Station.Pos, render.Color("teal", 5), 2, true)
	}
	if w.Planet.Radius > 0 {
		plot(w.Planet.Pos, render.Color("blue", 3), 2, false)
	}
	if w.Sun.Radius > 0 {
		plot(w.Sun.Pos, render.Color("yellow", 4), 2, false)
	}

	for _, e := range w.Entities {
		if e.Kind == "explosion" {
			continue
		}
		var col uint32
		switch {
		case e.Kind == "missile":
			col = ColWhite
		case e.Kind == "canister":
			col = render.Color("yellow", 4)
		case e.Kind == "asteroid":
			col = render.Color("brown", 4)
		case e.Hostile:
			col = ColRed
		case e.Behaviour == "police":
			col = render.Color("blue", 5)
		default:
			col = render.Color("green", 4)
		}
		plot(e.Pos, col, 1, true)
	}
}

func compass(fb *render.FB, cx, cy, r int) {
	g := game.G
	fb.EllipseFill(cx, cy, r, r, render.RGB(0x04121a))
	fb.Ellipse(cx, cy, r, r, render.RGB(0x3fbcd8))
	st := g.World.Station
	if st == nil {
		// witchspace: no station to point at
		return
	}
	rel := st.Pos.Sub(g.Player.Pos)
	v := g.Player.Mat.ApplyT(rel)
	depth := -v[2]
	l := math.Sqrt(v[0]*v[0] + v[1]*v[1] + depth*depth)
	if l < 0.001 {
		l = 1
	}
	sx := int(math.Round(float64(cx) + v[0]/l*float64(r-2)))
	sy := int(math.Round(float64(cy) - v[1]/l*float64(r-2)))
	if depth > 0 {
		fb.FillRect(sx-1, sy-1, sx+1, sy+1, ColGreen)
	} else {
		fb.Px(sx, sy, render.Color("teal", 2))
	}
	fb.Px(cx, cy, render.RGB(0x1b5a6a))
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Draw paints the current screen.
func Draw(fb *render.FB, dt float64) {
	g := game.G
	if g.Screen == "settings" {
		chrome(fb, "SETTINGS")
		settingsScreen(fb)
		footer(fb)
		return
	}
	if g.Screen == "controls" {
		chrome(fb, "CONTROLS")
		controlsScreen(fb)
		footer(fb)
		return
	}
	if g.Screen == "padtest" {
		chrome(fb, "PAD TEST")
		padTestScreen(fb)
		footer(fb)
		return
	}
	if g.Screen == "mode" && g.State == "title" {
		chrome(fb, "NEW GAME")
		modeScreen(fb)
		footer(fb)
		return
	}
	if g.Screen == "slots" && g.State == "title" {
		chrome(fb, "LOAD COMMANDER")
		slotsScreen(fb)
		footer(fb)
		return
	}
	if g.State == "flight" && g.Overlay != "" {
		if g.Overlay == "pause" {
			pauseScreen(fb)
			return
		}
		chrome(fb, "GALACTIC CHART - IN FLIGHT")
		if g.Overlay == "chart" {
			chartScreen(fb)
		} else {
			statusScreen(fb)
		}
		footer(fb)
		return
	}
	switch g.State {
	case "title":
		Title(fb)
	case "dead":
		gameOver(fb)
	case "hyper":
		hyperScreen(fb)
	case "docking":
		flight(fb, dt)
	case "docked":
		docked(fb)
	default:
		flight(fb, dt)
	}
}

// ---------------------------------------------------------------- flight

func flight(fb *render.FB, dt float64) {
	g := game.G
	p := &g.Player
	w := g.World
	if w == nil {
		return
	}
	fb.SetClip(0, 0, 319, 149)
	fb.ClearGradientStops([]uint32{0xff0f0705, 0xff160c07, 0xff1c1008, 0xff22140a, 0xff26170b, 0xff2b1a0c})

	// nebula backdrop pans with yaw
	if w.Nebula == nil {
		w.Nebula = models.MakeNebula(w.Sys.Seed, g.Cmd.Galaxy)
	}
	fwdW := p.Mat.Col(2).Mul(-1)
	yaw := math.Atan2(fwdW[0], -fwdW[2])
	pitch := math.Asin(math.Max(-1, math.Min(1, fwdW[1])))
	nebOff := int(math.Round(yaw / (2 * math.Pi) * float64(w.Nebula.W) * 4))
	nebY := int(math.Round(pitch * 9))
	if nebY < -4 {
		nebY = -4
	}
	if nebY > 4 {
		nebY = 4
	}
	render.BlitPart(fb, w.Nebula, nebOff, 0, 320, 150, 0, -nebY, 2)

	shakeX, shakeY := 0, 0
	if g.Shake > 0 {
		shakeX = int(math.Round((rand01() - 0.5) * g.Shake * 6))
		shakeY = int(math.Round((rand01() - 0.5) * g.Shake * 6))
	}
	cam := render.NewCamera(game.ViewW, game.ViewH, 62)
	cam.CX = 160 + float64(shakeX)
	cam.CY = 75 + float64(shakeY)
	cam.Pos = p.Pos
	cam.Mat = p.Mat
	focal := cam.Focal

	// stars
	for _, s := range w.Stars {
		world := p.Pos.Add(s.Dir.Mul(400000))
		vp := cam.ToView(world)
		if vp[2] > -1000 {
			continue
		}
		sp := cam.ProjectView(vp)
		if sp.X < -2 || sp.X > 322 || sp.Y < -2 || sp.Y > 152 {
			continue
		}
		col := render.Color("steel", 3)
		if s.Bright > 0.86 {
			col = render.Color("steel", 5)
		} else if s.Bright > 0.6 {
			col = render.Color("steel", 4)
		}
		fb.Px(int(sp.X), int(sp.Y), col)
		if s.Size == 2 && s.Bright > 0.8 {
			fb.Px(int(sp.X)+1, int(sp.Y), render.Color("steel", 3))
			fb.Px(int(sp.X), int(sp.Y)+1, render.Color("steel", 3))
		}
	}

	// planet (absent in witchspace)
	planet := w.Planet
	if w.PlanetFrames == nil && planet.Radius > 0 {
		w.PlanetFrames = models.MakePlanetFrames(128, planet.Seed, planet.Kind, 8)
	}
	pv := cam.ToView(planet.Pos)
	pdepth := -pv[2]
	if pdepth > planet.Radius*0.5 {
		psp := cam.ProjectView(pv)
		sr := planet.Radius * focal / pdepth
		frameIdx := int(g.Time*0.6) % len(w.PlanetFrames)
		if sr > 3 {
			render.BlitScaled(fb, w.PlanetFrames[frameIdx], psp.X-sr, psp.Y-sr, sr*2, sr*2, math.Max(1, pdepth-planet.Radius*0.6))
			if sr > 7 {
				rimRamps := []string{"teal", "orange", "steel", "yellow", "red"}
				rim := rimRamps[0]
				if planet.Kind >= 0 && planet.Kind < len(rimRamps) {
					rim = rimRamps[planet.Kind]
				}
				// two soft dithered bands hugging the limb
				for band := 0; band < 2; band++ {
					rad := sr + 1.0 + float64(band)*1.4
					col := render.Color(rim, 3-band)
					steps := int(math.Max(90, rad*7))
					thr := 3 + band*6
					for ai := 0; ai < steps; ai++ {
						ang := float64(ai) / float64(steps) * 2 * math.Pi
						qx := int(math.Round(psp.X + math.Cos(ang)*rad))
						qy := int(math.Round(psp.Y + math.Sin(ang)*rad))
						if qx < 0 || qx > 319 || qy < 0 || qy > 149 {
							continue
						}
						if render.Bayer[(qy&3)*4+(qx&3)] > thr {
							fb.Px(qx, qy, col)
						}
					}
				}
			}
		} else {
			fb.FillRect(int(psp.X)-1, int(psp.Y)-1, int(psp.X), int(psp.Y), render.Color("blue", 4))
		}
	}

	// sun (absent in witchspace)
	sv := cam.ToView(w.Sun.Pos)
	sdepth := -sv[2]
	if w.Sun.Radius > 0 && sdepth > w.Sun.Radius {
		ssp := cam.ProjectView(sv)
		sunR := math.Min(160, w.Sun.Radius*focal/sdepth)
		if ssp.X > -sunR && ssp.X < 320+sunR && ssp.Y > -sunR && ssp.Y < 150+sunR {
			drawSun(fb, ssp.X, ssp.Y, sunR)
		}
	}

	// dust
	for i := range w.Dust {
		dp := &w.Dust[i]
		rel := dp.Pos.Sub(p.Pos)
		if rel.Len() > 2600 {
			dir := render.V3(rand01()-0.5, rand01()-0.5, -1).Norm()
			dp.Pos = p.Pos.Add(dir.Mul(400 + rand01()*2200))
			continue
		}
		dv := cam.ToView(dp.Pos)
		if dv[2] > -40 {
			continue
		}
		dsp := cam.ProjectView(dv)
		if dsp.X < 0 || dsp.X > 319 || dsp.Y < 0 || dsp.Y > 149 {
			continue
		}
		frac := 1 - math.Min(1, -dv[2]/2600)
		col := render.Color("steel", 2+int(frac*3))
		fb.Px(int(dsp.X), int(dsp.Y), col)
		if p.Speed > 250 && frac > 0.45 {
			fb.Px(int(dsp.X), int(dsp.Y)+1, render.Color("steel", 2))
		}
	}

	// entities, far to near
	list := make([]*game.Entity, 0, len(w.Entities)+1)
	list = append(list, w.Entities...)
	if w.Station != nil {
		list = append(list, w.Station)
	}
	for i := 1; i < len(list); i++ {
		for j := i; j > 0; j-- {
			d1 := list[j].Pos.Sub(p.Pos).Len()
			d2 := list[j-1].Pos.Sub(p.Pos).Len()
			if d1 <= d2 {
				break
			}
			list[j], list[j-1] = list[j-1], list[j]
		}
	}

	for _, e := range list {
		if e.Kind == "explosion" {
			ev := cam.ToView(e.Pos)
			if ev[2] > 0 {
				continue
			}
			esp := cam.ProjectView(ev)
			t := e.Age / e.Life
			rad := (e.Radius * (0.4 + t*2.2)) * focal / math.Max(20, -ev[2])
			if rad < 1 {
				rad = 1
			}
			for rr := 0; rr < 3; rr++ {
				fr := math.Mod(t+float64(rr)*0.18, 1)
				r2 := rad * (0.4 + fr)
				ramp := "yellow"
				if fr >= 0.4 && fr < 0.7 {
					ramp = "orange"
				} else if fr >= 0.7 {
					ramp = "red"
				}
				col := render.Color(ramp, 5-int(fr*3))
				steps := int(math.Max(10, r2*3))
				for a := 0; a < steps; a++ {
					ang := float64(a) / float64(steps) * 2 * math.Pi
					px := int(math.Round(esp.X + math.Cos(ang)*r2))
					py := int(math.Round(esp.Y + math.Sin(ang)*r2))
					if render.Bayer[(py&3)*4+(px&3)] < 6 {
						fb.Px(px, py, col)
					}
				}
			}
			continue
		}
		if e.Kind == "station" {
			render.DrawMesh(fb, cam, e.Mesh, e.Pos, e.Mat, 60000)
			nrm := e.Mat.Col(2)
			beacon := e.Pos.Add(nrm.Mul(game.StationSize + 45))
			bv := cam.ToView(beacon)
			if bv[2] < 0 {
				bsp := cam.ProjectView(bv)
				bc := render.Color("red", 2)
				if math.Sin(g.Time*4) > 0 {
					bc = ColRed
				}
				fb.FillRect(int(bsp.X)-1, int(bsp.Y)-1, int(bsp.X)+1, int(bsp.Y)+1, bc)
			}
			continue
		}
		dist := e.Pos.Sub(p.Pos).Len()
		if dist > 40000 {
			continue
		}
		if e.Mesh != nil {
			render.DrawMesh(fb, cam, e.Mesh, e.Pos, e.Mat, 30000)
			if e.Speed > 25 && dist < 18000 {
				back := e.Mat.Col(2)
				tail := e.Pos.Add(back.Mul(e.Radius * 0.85))
				tv := cam.ToView(tail)
				if tv[2] < -20 {
					tp := cam.ProjectView(tv)
					glowR := math.Max(1, math.Min(7, e.Radius*0.28*focal/(-tv[2])))
					ramp := "teal"
					if e.Hostile {
						ramp = "orange"
					}
					glowDot(fb, tp.X, tp.Y, glowR, ramp)
				}
			}
		}
	}

	reticle(fb, cam)

	if g.BeamFlash > 0 {
		flick := 0.4 + rand01()*0.6
		fb.Line(90, 149, 157, 79, render.Color("red", int(4*flick)+1))
		fb.Line(230, 149, 163, 79, render.Color("red", int(4*flick)+1))
		fb.Px(156, 75, ColWhite)
		fb.Px(164, 75, ColWhite)
	}
	if g.EcmFlash > 0 {
		for y := 0; y < 150; y += 2 {
			f := math.Sin(float64(y)*0.35+g.Time*40)*0.5 + 0.5
			if f > 0.6 {
				fb.HLine(0, 319, y, render.Color("teal", 2+int(f*2)))
			}
		}
	}
	if g.BombFlash > 0 {
		ba := math.Min(1, g.BombFlash)
		for yy := 0; yy < 150; yy++ {
			f2 := 1 - math.Abs(float64(yy)-75)/75
			if rand01() < f2*ba {
				fb.HLine(0, 319, yy, render.Color("yellow", 4+int(f2)))
			}
		}
	}

	fb.ResetClip()
	cockpit(fb)
	crosshair(fb, 160, 75)

	// messages
	my := 6
	for i := len(g.Messages) - 1; i >= 0; i-- {
		msg := g.Messages[i]
		col := ColYellow
		if msg.Kind == "bad" {
			col = ColRed
		} else if msg.Kind == "good" {
			col = ColGreen
		}
		tw := render.TextWidth(msg.Text, 1)
		fb.FillRect(160-tw/2-3, my-1, 160+tw/2+2, my+7, render.RGB(0x050b12))
		fb.TextC(msg.Text, 160, my, col, 1)
		my += 9
	}

	hud(fb)
}

func rand01() float64 { return float64(randInt()) / 2147483648.0 }

func randInt() int {
	return int(randFloat() * 2147483647)
}

// GlowAt draws a soft dithered glow at a screen position (icon / effects).
func GlowAt(fb *render.FB, x, y int, r float64, ramp string) {
	glowDot(fb, float64(x), float64(y), r, ramp)
}

func glowDot(fb *render.FB, x, y, r float64, ramp string) {
	ir := int(r)
	for yy := -ir; yy <= ir; yy++ {
		for xx := -ir; xx <= ir; xx++ {
			d := math.Sqrt(float64(xx*xx + yy*yy))
			if d > r+0.4 {
				continue
			}
			b := render.Bayer[((int(y)+yy)&3)*4+((int(x)+xx)&3)]
			if d > r*0.55 && b > 9 {
				continue
			}
			shade := 3
			if d < r*0.45 {
				shade = 5
			} else if d < r*0.8 {
				shade = 4
			}
			fb.Px(int(x)+xx, int(y)+yy, render.Color(ramp, shade))
		}
	}
}

func reticle(fb *render.FB, cam *render.Camera) {
	g := game.G
	w := g.World
	var best *game.Entity
	bestDist := math.Inf(1)
	var bs render.Pt3
	for _, e := range w.Entities {
		if e.Kind != "ship" || !e.Hostile {
			continue
		}
		vv := cam.ToView(e.Pos)
		if vv[2] > -60 {
			continue
		}
		s := cam.ProjectView(vv)
		if s.X < 0 || s.X > 319 || s.Y < 0 || s.Y > 149 {
			continue
		}
		d := vv.Len()
		if d < bestDist {
			bestDist = d
			best = e
			bs = s
		}
	}
	if best == nil {
		return
	}
	depth := -cam.ToView(best.Pos)[2]
	rad := math.Max(5, math.Min(40, best.Radius*cam.Focal/depth))
	x0 := int(math.Round(bs.X - rad))
	x1 := int(math.Round(bs.X + rad))
	y0 := int(math.Round(bs.Y - rad*0.8))
	y1 := int(math.Round(bs.Y + rad*0.8))
	col := render.Color("red", 4)
	l := maxInt(3, int(math.Round(rad*0.35)))
	fb.HLine(x0, x0+l, y0, col)
	fb.VLine(x0, y0, y0+l, col)
	fb.HLine(x1-l, x1, y0, col)
	fb.VLine(x1, y0, y0+l, col)
	fb.HLine(x0, x0+l, y1, col)
	fb.VLine(x0, y1-l, y1, col)
	fb.HLine(x1-l, x1, y1, col)
	fb.VLine(x1, y1-l, y1, col)
	label := best.Type + " " + itoa(int(depth))
	fb.Text(label, int(bs.X)-render.TextWidth(label, 1)/2, y1+3, render.Color("red", 3), 1)
}

func cockpit(fb *render.FB) {
	dark := render.RGB(0x040a11)
	beam := render.Color("steel", 1)
	beamHi := render.Color("steel", 3)
	glass := render.Color("teal", 2)
	// left and right structural members of the canopy frame
	fb.Poly([]render.Pt3{{X: 0, Y: 92}, {X: 9, Y: 92}, {X: 74, Y: 150}, {X: 0, Y: 150}}, beam)
	fb.Poly([]render.Pt3{{X: 319, Y: 92}, {X: 310, Y: 92}, {X: 245, Y: 150}, {X: 319, Y: 150}}, beam)
	fb.Line(0, 92, 74, 150, beamHi)
	fb.Line(319, 92, 245, 150, beamHi)
	fb.Line(9, 92, 74, 150, glass)
	fb.Line(310, 92, 245, 150, glass)
	// bottom bezel band
	fb.FillRect(0, 146, 319, 149, render.RGB(0x0b1622))
	fb.HLine(0, 319, 146, beamHi)
	fb.HLine(0, 319, 147, render.Color("steel", 2))
	for x := 6; x < 319; x += 14 {
		fb.Px(x, 148, render.Color("steel", 2))
	}
	// corner rivets
	for _, p := range [][2]int{{6, 96}, {313, 96}, {40, 142}, {280, 142}} {
		fb.Px(p[0], p[1], render.Color("steel", 4))
	}
	// top edge
	fb.HLine(0, 319, 0, dark)
	fb.HLine(0, 319, 1, render.RGB(0x0b1622))
	for yy := 0; yy < 3; yy++ {
		for xx := 0; xx < 320; xx++ {
			if render.Bayer[(yy&3)*4+(xx&3)] < (3-yy)*4 {
				fb.Px(xx, yy, dark)
			}
		}
	}
}

func crosshair(fb *render.FB, cx, cy int) {
	c := render.Color("green", 4)
	c2 := render.Color("green", 3)
	fb.Px(cx, cy, ColWhite)
	fb.HLine(cx-12, cx-6, cy, c)
	fb.HLine(cx+6, cx+12, cy, c)
	fb.VLine(cx, cy-9, cy-4, c2)
	fb.VLine(cx, cy+4, cy+9, c2)
	fb.Px(cx-14, cy, c2)
	fb.Px(cx+14, cy, c2)
	fb.Px(cx, cy-11, c2)
	fb.Px(cx, cy+11, c2)
}

// ------------------------------------------------------------------ HUD

func hud(fb *render.FB) {
	g := game.G
	p := &g.Player
	cmd := &g.Cmd
	y0 := game.ViewH
	fb.FillRect(0, y0, 319, y0+24, render.RGB(0x0c1524))
	fb.FillRect(0, y0+25, 319, y0+49, render.RGB(0x0a1220))
	fb.FillRect(0, y0+50, 319, 223, render.RGB(0x080e19))
	fb.HLine(0, 319, y0, ColFrameHi)
	fb.HLine(0, 319, y0+1, render.RGB(0x1b3348))

	gx, bx, bw := 6, 22, 58
	maxE := g.Ship().Energy
	fb.Text("FS", gx, y0+5, ColText, 1)
	shieldMax := game.G.ShieldMax()
	if shieldMax <= 0 {
		shieldMax = game.MaxShield
	}
	segBar(fb, bx, y0+5, bw, 5, p.ShieldF/shieldMax, "blue", p.ShieldF < 22)
	fb.Text("AS", gx, y0+13, ColText, 1)
	segBar(fb, bx, y0+13, bw, 5, p.ShieldA/shieldMax, "blue", p.ShieldA < 22)
	fb.Text("FU", gx, y0+21, ColText, 1)
	segBar(fb, bx, y0+21, bw, 5, cmd.Fuel/game.FuelMax, "green", cmd.Fuel < 0.8)
	fb.Text("EN", gx, y0+29, ColText, 1)
	segBar(fb, bx, y0+29, bw, 5, p.Energy/maxE, "yellow", p.Energy < 22)
	fb.Text("HT", gx, y0+37, ColText, 1)
	segBar(fb, bx, y0+37, bw, 5, p.LaserTemp, "red", p.LaserTemp > 0.85)
	fb.Text("CT", gx, y0+45, ColText, 1)
	segBar(fb, bx, y0+45, bw, 5, p.Cabin, "red", p.Cabin > 0.75)

	fb.Text("SPD "+itoa(int(p.Speed)), gx, y0+55, ColText, 1)
	fb.Text(strconv.FormatFloat(cmd.Fuel, 'f', 1, 64)+"LY", bx+bw-26, y0+55, ColGreen, 1)
	fb.Text("SYS "+g.World.Sys.Name, gx, y0+64, render.RGB(0x79c9e6), 1)

	scanner(fb, 160, y0+40, 56, 26)

	rx, vx := 220, 238
	fb.Text("MSL", rx, y0+8, ColText, 1)
	for i := 0; i < game.MissilesMax; i++ {
		mx := vx + i*10
		if i < cmd.Missiles {
			fb.FillRect(mx, y0+8, mx+5, y0+13, render.Color("ash", 4))
			fb.VLine(mx+2, y0+9, y0+12, ColRed)
			fb.HLine(mx, mx+5, y0+8, render.Color("ash", 5))
		} else {
			fb.Rect(mx, y0+8, mx+5, y0+13, render.RGB(0x22415a))
		}
	}
	fb.Text("LSR", rx, y0+17, ColText, 1)
	segBar(fb, vx, y0+17, 34, 5, p.LaserTemp, "red", false)
	laserCodes := [4]string{"PLS", "BEA", "MIL", "MIN"}
	fb.TextR(laserCodes[cmd.Laser], 318, y0+17, ColDim, 1)

	lamps := []struct {
		ch string
		on bool
	}{
		{"E", cmd.Equip["ecm"]}, {"S", cmd.Equip["fuelscoop"]},
		{"D", cmd.Equip["dockcomp"]}, {"B", cmd.Bombs > 0},
	}
	for i, l := range lamps {
		lx := rx + i*13
		col := render.RGB(0x4c6a80)
		if l.on {
			col = ColGreen
		}
		fb.Text(l.ch, lx, y0+26, col, 1)
		if l.on {
			fb.Px(lx+5, y0+26, render.Color("green", 5))
		}
	}

	fb.TextR(money(cmd.Cash)+" CR", 318, y0+35, ColYellow, 1)
	legalTxt := "CLEAN"
	legalCol := ColGreen
	if cmd.Legal >= 500 {
		legalTxt = "FUGITIVE"
		legalCol = ColRed
	} else if cmd.Legal > 0 {
		legalTxt = "OFFENDER"
		legalCol = ColRed
	}
	fb.TextR(legalTxt, 318, y0+44, legalCol, 1)
	fb.TextR("HOLD "+itoa(game.CargoUsed(*cmd))+"/"+itoa(cmd.CargoCap), 318, y0+53, ColDim, 1)

	compass(fb, 230, y0+55, 8)
	if st := g.World.Station; st != nil {
		dist := st.Pos.Sub(p.Pos).Len()
		distTxt := itoa(int(dist))
		if dist > 9999 {
			distTxt = strconv.FormatFloat(dist/1000, 'f', 1, 64) + "K"
		}
		// green once the port faces you: the slot is on one face only, and
		// hunting for it without a cue is the main docking frustration
		col := ColText
		if facing := p.Pos.Sub(st.Pos).Norm().Dot(st.Mat.Col(2)); facing > 0.55 {
			col = ColGreen
		}
		fb.Text("STN "+distTxt, 246, y0+62, col, 1)
	}
}

// ------------------------------------------------------------------ title

func Title(fb *render.FB) {
	g := game.G
	fb.ClearGradient(0xff0d0604, 0xff0f0705)
	for i := 0; i < 70; i++ {
		sx := int(math.Mod(float64(i*73)+g.Time*12*(0.3+float64(i%5)*0.2), 320))
		sy := (i * 41) % 224
		fb.Px(sx, sy, render.Color("steel", 2+i%3))
	}
	cam := render.NewCamera(320, 224, 55)
	cam.CX = 160
	cam.CY = 128
	mat := render.RotX(0.35 + math.Sin(g.Time*0.3)*0.15).Mul(render.RotY(g.Time * 0.5))
	render.DrawMesh(fb, cam, models.Get("cobra3"), render.V3(0, 6, -230), mat, 60000)
	fb.TextC(game.ProductName, 162, 16, render.Color("steel", 0), 5)
	fb.TextC(game.ProductName, 161, 15, render.Color("steel", 2), 5)
	fb.TextC(game.ProductName, 160, 14, render.Color("steel", 5), 5)
	fb.TextC(game.ProductSub, 161, 54, render.Color("steel", 2), 3)
	fb.TextC(game.ProductSub, 160, 53, render.Color("steel", 4), 3)
	fb.TextC("THE GALAXY IS YOURS, COMMANDER", 160, 84, ColDim, 1)
	// one left-aligned column, no digits: mixed widths read as a ragged block.
	// CONTINUE appears only when a flight is still in progress.
	const menuX = 116
	const markX = 102
	items := TitleItems()
	if S.Title >= len(items) {
		S.Title = len(items) - 1
	}
	for i, name := range items {
		y := 140 + i*11
		col := ColText
		if name == "LOAD COMMANDER" && !game.SaveExists() {
			col = ColDim
		}
		if i == S.Title {
			col = ColYellow
			fb.Text(">", markX, y, ColYellow, 1)
		}
		fb.Text(name, menuX, y, col, 1)
	}
	fb.TextC("ARROWS / MOUSE STEER   SPACE FIRE   X/Z SPEED", 160, 196, ColDim, 1)
	fb.TextC("C CHART   J JUMP   ESC CONTROLS", 160, 207, render.RGB(0x2d5f7a), 1)
	// build stamp, so a stale app or a downloaded binary is identifiable at a
	// glance — a local build is stamped with its time
	if game.Version != "" {
		fb.TextR(game.Version, 316, 216, render.RGB(0x2b4a5e), 1)
	}
	// the seed is shown so a run can be repeated exactly, and a bug report can
	// say which world it came from
	if seed := game.Seed(); seed > 0 {
		fb.Text(fmt.Sprintf("SEED %d", seed), 4, 216, render.RGB(0x2b4a5e), 1)
	}
}

// TitleItems is the main menu, in order. Continue is offered only while a
// session exists, so leaving to the menu no longer throws the flight away.
func TitleItems() []string {
	out := []string{"NEW GAME"}
	if S.HasSession {
		out = append(out, "CONTINUE FLIGHT")
	}
	return append(out, "LOAD COMMANDER", "SETTINGS", "QUIT")
}

func gameOver(fb *render.FB) {
	g := game.G
	fb.Clear(render.RGB(0x05070f))
	fb.TextC("GAME OVER", 160, 60, ColRed, 3)
	reason := "SHIP LOST"
	if g.Dead != nil {
		reason = g.Dead.Reason
	}
	fb.TextC(reason, 160, 100, ColText, 1)
	fb.TextC("CASH "+money(g.Cmd.Cash)+" CR   KILLS "+itoa(g.Cmd.Kills), 160, 120, ColDim, 1)
	fb.TextC("PRESS ENTER TO CONTINUE", 160, 160, ColYellow, 1)
}

func hyperScreen(fb *render.FB) {
	g := game.G
	fb.Clear(render.RGB(0x02030a))
	t := 1.0
	target := g.Cmd.TargetSystem
	if g.Hyper != nil {
		t = g.Hyper.T
		target = g.Hyper.Target
	}
	cx, cy := 160.0, 112.0
	n := 70
	for i := 0; i < n; i++ {
		a := float64(i)/float64(n)*2*math.Pi + t*0.4
		r0 := 8 + math.Mod(float64(i*37)+t*900, 260)
		r1 := r0 + 10 + t*40
		col := render.Color("steel", 2+i%3)
		ca, sa := math.Cos(a), math.Sin(a)
		fb.Line(int(cx+ca*r0), int(cy+sa*r0*0.75), int(cx+ca*r1), int(cy+sa*r1*0.75), col)
	}
	fb.TextC("HYPERSPACE", 160, 20, ColWhite, 2)
	gal := galaxy.Get(g.Cmd.Galaxy)
	fb.TextC("TO "+gal.Systems[((target%256)+256)%256].Name, 160, 40, ColYellow, 1)
	if g.Hyper != nil {
		fb.TextC(strconv.FormatFloat(g.Hyper.Dist, 'f', 1, 64)+" LIGHT YEARS", 160, 200, ColDim, 1)
	}
}

// drawSun paints a banded star with a soft dithered corona and a flare.
func drawSun(fb *render.FB, cx, cy, radius float64) {
	ir := int(radius*2.6) + 2
	for y := -ir; y <= ir; y++ {
		for x := -ir; x <= ir; x++ {
			px, py := int(cx)+x, int(cy)+y
			if px < 0 || px > 319 || py < 0 || py > 149 {
				continue
			}
			d := math.Hypot(float64(x), float64(y)) / radius
			if d > 2.6 {
				continue
			}
			var ramp string
			var shade float64
			switch {
			case d < 0.42:
				ramp, shade = "yellow", 5
			case d < 0.72:
				ramp, shade = "yellow", 4.4
			case d < 1.0:
				ramp, shade = "orange", 4.0
			case d < 1.35:
				ramp, shade = "orange", 2.8
			default:
				ramp, shade = "orange", 1.6
			}
			// stable hash dither so the bands break up like a 16-bit gradient
			h := ((px*13 + py*7) * 2654435761) >> 16 & 15
			dens := 1.0
			if d >= 1.0 {
				dens = math.Max(0, 1-(d-1.0)/1.6)
			}
			if float64(h) > dens*16 {
				continue
			}
			fb.Px(px, py, render.Color(ramp, int(shade)))
		}
	}
	// anamorphic flare
	for dy := -1; dy <= 1; dy++ {
		py := int(cy) + dy
		for x := -int(radius * 2.6); x <= int(radius*2.6); x++ {
			px := int(cx) + x
			if px < 0 || px > 319 || py < 0 || py > 149 {
				continue
			}
			if dy != 0 && (x&1) != 0 {
				continue
			}
			if dy == 0 || (x*13+py*7)&3 == 0 {
				fb.Px(px, py, render.Color("yellow", 4-dy))
			}
		}
	}
}
