package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"heliosdrift/galaxy"
	"heliosdrift/game"
	"heliosdrift/input"
	"heliosdrift/lang"
	"heliosdrift/models"
	"heliosdrift/render"
	"heliosdrift/sfx"
	"heliosdrift/ui"
)

const (
	ScreenW = 320
	ScreenH = 224
)

var bank *sfx.Bank

type App struct {
	fb    *render.FB
	img   *ebiten.Image
	rgba  []byte
	scale float64
	mouse bool
}

func newApp() *App {
	a := &App{
		fb:   render.NewFB(ScreenW, ScreenH),
		img:  ebiten.NewImage(ScreenW, ScreenH),
		rgba: make([]byte, ScreenW*ScreenH*4),
	}
	return a
}

func (a *App) Update() error {
	g := game.G
	dt := 1.0 / 60.0

	// rebinding prompt: while it is up, swallow everything except Escape
	if ui.IsCapturing() {
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			ui.CancelCapture()
		} else {
			ui.PollCapture(dt)
		}
		g.Time += dt
		if bank != nil {
			bank.SetMaster(lang.Cfg.Volume)
			bank.Muted = lang.Cfg.Muted
		}
		return nil
	}

	// F11 goes fullscreen: handy on a TV or a handheld
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// menus are fixed; only flying uses the rebindable layout
	freeFlight := g.State == "flight" && g.Overlay == ""
	if !freeFlight {
		for _, k := range []ebiten.Key{
			ebiten.KeyEnter, ebiten.KeyEscape, ebiten.KeyTab,
			ebiten.KeyArrowUp, ebiten.KeyArrowDown, ebiten.KeyArrowLeft, ebiten.KeyArrowRight,
			ebiten.KeySpace, ebiten.KeyC, ebiten.KeyV, ebiten.KeyM, ebiten.KeyE,
			ebiten.KeyL, ebiten.KeyJ, ebiten.KeyF, ebiten.KeyB, ebiten.KeyS, ebiten.KeyX, ebiten.KeyZ,
			ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5, ebiten.Key6,
			ebiten.Key7, ebiten.Key8, ebiten.Key9, ebiten.Key0, ebiten.KeyG, ebiten.KeyQ, ebiten.KeyH, ebiten.KeyK,
			ebiten.KeyBackspace,
		} {
			if inpututil.IsKeyJustPressed(k) {
				name := keyName(k)
				if name != "" {
					ui.Key(name)
					if ui.QuitRequested() {
						return ebiten.Termination
					}
				}
			}
		}
		// the pad drives the menus as well, so the game is playable on a
		// handheld with no keyboard attached; D-pad acts exactly like the
		// arrow keys here
		if pad := firstGamepad(); pad >= 0 {
			for _, m := range []struct {
				btn int32
				key ebiten.Key
			}{
				{input.PadA, ebiten.KeyEnter},
				{input.PadB, ebiten.KeyEscape},
				{input.PadStart, ebiten.KeyEscape},
				{input.PadUp, ebiten.KeyArrowUp},
				{input.PadDown, ebiten.KeyArrowDown},
				{input.PadLeft, ebiten.KeyArrowLeft},
				{input.PadRight, ebiten.KeyArrowRight},
			} {
				if inpututil.IsStandardGamepadButtonJustPressed(pad, ebiten.StandardGamepadButton(m.btn)) {
					if name := keyName(m.key); name != "" {
						ui.Key(name)
					}
				}
			}
		}
	}

	// flight input
	if g.State == "flight" && g.Overlay == "" {
		pitch, roll := 0.0, 0.0
		// digital part (keys, buttons) and analogue part (bound axes) add up
		if input.Digital(input.PitchUp) {
			pitch += 1
		}
		if input.Digital(input.PitchDown) {
			pitch -= 1
		}
		if input.Digital(input.RollLeft) {
			roll -= 1
		}
		if input.Digital(input.RollRight) {
			roll += 1
		}
		yaw := 0.0
		if input.Digital(input.YawLeft) {
			yaw -= 1
		}
		if input.Digital(input.YawRight) {
			yaw += 1
		}
		yaw += input.AxisValue(input.YawRight) - input.AxisValue(input.YawLeft)
		pitch += input.AxisValue(input.PitchUp) - input.AxisValue(input.PitchDown)
		roll += input.AxisValue(input.RollRight) - input.AxisValue(input.RollLeft)
		if pitch == 0 && roll == 0 && a.mouse {
			mx, my := ebiten.CursorPosition()
			w, h := ebiten.WindowSize()
			nx := (float64(mx)/float64(w) - 0.5) * 2
			ny := (float64(my)/float64(h) - 0.5) * 2
			if math.Abs(nx) > 0.06 {
				roll = clampF((nx-math.Copysign(0.06, nx))*1.6, -1, 1)
			}
			if math.Abs(ny) > 0.06 {
				pitch = -clampF((ny-math.Copysign(0.06, ny))*1.6, -1, 1)
			}
		}
		accel := 0.0
		if input.Digital(input.Accel) {
			accel += 1
		}
		if input.Digital(input.Decel) {
			accel -= 1
		}
		accel += input.AxisValue(input.Accel) - input.AxisValue(input.Decel)
		fire := input.Held(input.Fire) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

		dispatchFlightActions()
		g.In.Pitch = pitch
		g.In.Roll = roll
		g.In.Yaw = yaw
		g.In.Accel = accel
		g.In.Fire = fire
	} else {
		g.In.Pitch, g.In.Roll, g.In.Yaw, g.In.Accel, g.In.Fire = 0, 0, 0, 0, false
	}

	if g.Overlay == "" {
		g.Update(dt)
	} else {
		g.Time += dt
	}
	if bank != nil {
		bank.SetMaster(lang.Cfg.Volume)
		if bank.Muted != lang.Cfg.Muted {
			bank.Muted = lang.Cfg.Muted
		}
		lvl := 0.0
		if g.State == "flight" && !g.Autopilot {
			lvl = 0.04 + 0.18*(g.Player.Speed/game.MaxSpeed)
		} else if g.State == "flight" {
			lvl = 0.10
		}
		bank.Update(lvl)
	}
	ui.Draw(a.fb, dt)

	// copy framebuffer into the ebiten image
	buf := a.fb.Buf
	for i, c := range buf {
		a.rgba[i*4+0] = byte(c >> 0)
		a.rgba[i*4+1] = byte(c >> 8)
		a.rgba[i*4+2] = byte(c >> 16)
		a.rgba[i*4+3] = 0xff
	}
	a.img.WritePixels(a.rgba)
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	scale := math.Min(float64(sw)/ScreenW, float64(sh)/ScreenH)
	if scale < 1 {
		scale = 1
	}
	scale = math.Floor(scale)
	if scale < 1 {
		scale = 1
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.Filter = ebiten.FilterNearest
	offX := (float64(sw) - ScreenW*scale) / 2
	offY := (float64(sh) - ScreenH*scale) / 2
	op.GeoM.Translate(offX, offY)
	screen.DrawImage(a.img, op)
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func clampF(v, a, b float64) float64 {
	if v < a {
		return a
	}
	if v > b {
		return b
	}
	return v
}

func keyName(k ebiten.Key) string {
	shift := ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight)
	switch k {
	case ebiten.KeyEnter:
		return "Enter"
	case ebiten.KeyBackspace:
		return "Backspace"
	case ebiten.KeyEscape:
		return "Escape"
	case ebiten.KeyTab:
		return "Tab"
	case ebiten.KeyArrowUp:
		return "ArrowUp"
	case ebiten.KeyArrowDown:
		return "ArrowDown"
	case ebiten.KeyArrowLeft:
		return "ArrowLeft"
	case ebiten.KeyArrowRight:
		return "ArrowRight"
	case ebiten.Key1:
		return "1"
	case ebiten.Key2:
		return "2"
	case ebiten.Key3:
		return "3"
	case ebiten.Key4:
		return "4"
	case ebiten.Key5:
		return "5"
	case ebiten.Key6:
		return "6"
	case ebiten.Key7:
		return "7"
	case ebiten.Key8:
		return "8"
	case ebiten.Key9:
		return "9"
	case ebiten.Key0:
		return "0"
	case ebiten.KeyK:
		return "k"
	case ebiten.KeyG:
		return "g"
	case ebiten.KeyH:
		return "h"
	case ebiten.KeyQ:
		return "q"
	case ebiten.KeyB:
		if shift {
			return "B"
		}
		return "b"
	case ebiten.KeyS:
		if shift {
			return "S"
		}
		return "s"
	case ebiten.KeyC:
		return "c"
	case ebiten.KeyV:
		return "v"
	case ebiten.KeyM:
		return "m"
	case ebiten.KeyE:
		return "e"
	case ebiten.KeyL:
		return "l"
	case ebiten.KeyJ:
		return "j"
	case ebiten.KeyF:
		return "f"
	}
	return ""
}

// ------------------------------------------------------------------ tools

func writePNG(fb *render.FB, path string) error {
	img := image.NewRGBA(image.Rect(0, 0, fb.W, fb.H))
	for y := 0; y < fb.H; y++ {
		for x := 0; x < fb.W; x++ {
			c := fb.Buf[y*fb.W+x]
			img.Set(x, y, color.RGBA{uint8(c), uint8(c >> 8), uint8(c >> 16), 255})
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// runShots renders the same reference screens as the JS build (no window).
func runShots() {
	os.MkdirAll("build", 0o755)
	g := game.New()
	fb := render.NewFB(ScreenW, ScreenH)
	prefix := ""
	if lang.IsRU() {
		prefix = "ru-"
	}
	shot := func(name string) {
		ui.Draw(fb, 1.0/60)
		if err := writePNG(fb, filepath.Join("build", prefix+name+".png")); err != nil {
			fmt.Println("write error:", err)
			return
		}
		fmt.Println("wrote", name)
	}
	lookAt := func(from, target render.Vec3) {
		g.Player.Pos = from
		g.Player.Mat = render.LookMatrix(target.Sub(from).Norm())
	}

	g.NewGame("JAMESON")
	g.Launch()
	p := &g.Player
	st := g.World.Station
	normal := st.Mat.Col(2)
	g.Messages = nil
	g.World.Entities = nil
	g.Time = 8
	p.Speed = 180

	lookAt(st.Pos.Add(normal.Mul(1600)), st.Pos)
	g.SpawnTraffic(2)
	for i, e := range g.World.Entities {
		if e.Kind != "ship" {
			continue
		}
		e.Pos = p.Pos.Add(render.V3(0.45-float64(i)*0.9, 0.15, -1).Norm().Mul(900 + float64(i)*500))
		e.Mat = render.LookMatrix(p.Pos.Sub(e.Pos).Norm())
		e.Speed = 120
	}
	g.Message("DOCKING CLEARANCE GRANTED", "good")
	ui.Draw(fb, 1.0/60)
	colorStats(fb, 0, 0, 319, 149, "approach")
	shot("g30-approach")

	lookAt(st.Pos.Add(normal.Mul(700)), st.Pos)
	g.World.Entities = nil
	g.Message("ALIGN WITH THE SLOT", "info")
	ui.Draw(fb, 1.0/60)
	shot("g53-slot-close")

	// the same view turned 30 degrees off the slot axis
	off := render.RotY(0.52).Apply(normal).Norm()
	lookAt(st.Pos.Add(off.Mul(700)), st.Pos)
	ui.Draw(fb, 1.0/60)
	shot("g54-slot-angle")

	// inside the shaft, the view the docking glide shows
	lookAt(st.Pos.Add(normal.Mul(game.StationSize-35)), st.Pos)
	g.Message("DOCKED", "good")
	shot("g58-inside-shaft")

	// final approach distance, where the port is what the player actually sees
	lookAt(st.Pos.Add(render.RotY(0.52).Apply(normal).Mul(260)), st.Pos)
	g.Message("FINAL APPROACH - 30 DEGREES", "info")
	shot("g56-port-close30")

	lookAt(st.Pos.Add(normal.Mul(260)), st.Pos)
	shot("g55-port-close")

	lookAt(st.Pos.Add(normal.Mul(520)), st.Pos)
	ui.Draw(fb, 1.0/60)
	colorStats(fb, 0, 0, 319, 149, "slot view")
	shot("g31-slot")

	lookAt(st.Pos.Add(normal.Mul(2200)), st.Pos)
	g.SpawnPirates(3)
	offs := [][3]float64{{-260, 40, -900}, {180, -60, -1300}, {60, 120, -700}}
	for i, e := range g.World.Entities {
		if e.Kind != "ship" {
			continue
		}
		o := offs[i%3]
		e.Pos = p.Pos.Add(render.V3(o[0], o[1], o[2]))
		e.Mat = render.LookMatrix(p.Pos.Sub(e.Pos).Norm())
		e.Speed = 260
		e.Hostile = true
	}
	g.Messages = nil
	g.Message("HOSTILE CONTACT", "bad")
	shot("g32-combat")

	pl := g.World.Planet
	lookAt(pl.Pos.Add(render.V3(0, 0, pl.Radius*3.4)), pl.Pos)
	g.World.Entities = nil
	shot("g33-planet")

	sun := g.World.Sun
	lookAt(sun.Pos.Add(render.V3(sun.Radius*2.2, 0, sun.Radius*2.2)), sun.Pos)
	shot("g34-sun")

	g.State = "docked"
	g.Docked()
	g.Screen = "market"
	shot("g35-market")
	g.Screen = "chart"
	ui.S.ChartCursor = g.Cmd.System
	shot("g36-chart")
	g.Screen = "launch"
	shot("g37-station-menu")
	g.Screen = "equip"
	shot("g39-equip")
	g.Screen = "status"
	shot("g40-status")

	// mining: break a rock and show the debris
	g.World.Entities = nil
	g.SpawnAsteroid(520)
	for _, e := range append([]*game.Entity(nil), g.World.Entities...) {
		if e.Kind == "asteroid" {
			g.HitEntity(e, 500, "mining")
		}
	}
	for _, e := range g.World.Entities {
		if e.Kind == "boulder" {
			e.Pos = g.Player.Pos.Add(e.Pos.Sub(g.Player.Pos).Norm().Mul(380))
			e.Speed = 5
		}
	}
	g.Messages = nil
	g.Message("ROCK BREAKS UP", "info")
	shot("g46-mining")

	// a system whose station is a dodecahedron
	for i := 0; i < 256; i++ {
		if galaxy.Get(0).Systems[i].Seed%4 == 0 {
			g.EnterSystem(i, "arrive", 9000, false)
			break
		}
	}
	g.Messages = nil
	if g.World != nil {
		st2 := g.World.Station
		lookAt(st2.Pos.Add(st2.Mat.Col(2).Mul(1000)), st2.Pos)
		g.Player.Speed = 60
		shot("g44-dodeca")
	}

	g.State = "docked"
	g.Screen = "launch"
	g.Cmd.Kills = 40
	g.Cmd.Mission = game.Mission{}
	g.Cmd.MissionFlags = 0
	g.DockedForShots()
	shot("g42-menu-mission")
	g.Screen = "mission"
	shot("g43-mission-brief")

	g.Offers = []game.Contract{
		{Kind: game.ContractCourier, Good: 0, Qty: 3, System: 17, Galaxy: g.Cmd.Galaxy, Pay: 640},
		{Kind: game.ContractBounty, Qty: 2, System: 123, Galaxy: g.Cmd.Galaxy, Pay: 280},
		{Kind: game.ContractSmuggle, Good: 6, Qty: 3, System: 88, Galaxy: g.Cmd.Galaxy, Pay: 1120, Due: g.Cmd.Day + 5},
	}
	g.Cmd.Contracts = []game.Contract{
		{Kind: game.ContractCourier, Good: 6, Qty: 1, System: 55, Galaxy: g.Cmd.Galaxy, Pay: 900},
	}
	g.Screen = "contracts"
	shot("g45-contracts")

	g.Cmd.Kills = 40
	g.Cmd.Mission = game.Mission{}
	g.Cmd.MissionFlags = 0
	g.DockedForShots()
	g.AcceptMission(game.MissionConstrictor)
	g.Screen = "chart"
	g.Time = 1.2
	shot("g60-chart-mission")

	// journal with a story contract and two accepted jobs
	// campaign briefing (advanced story)
	savedMode := g.Cmd.Mode
	g.Cmd.Mode = game.ModeAdvanced
	g.State = "docked"
	g.Screen = "mission"
	shot("g63-campaign-brief")
	g.Cmd.Campaign.Active = true
	g.Cmd.Campaign.Target = 40
	g.Cmd.Campaign.Galaxy = 0
	g.Cmd.Campaign.Stage = 3
	g.Screen = "journal"
	shot("g64-campaign-journal")
	g.Cmd.Mode = savedMode

	g.Cmd.Contracts = []game.Contract{
		{Kind: game.ContractCourier, Good: 0, Qty: 3, System: 17, Galaxy: g.Cmd.Galaxy, Pay: 640, Due: g.Cmd.Day + 4},
		{Kind: game.ContractBounty, Qty: 4, Done: 1, System: 123, Galaxy: g.Cmd.Galaxy, Pay: 480, Due: g.Cmd.Day + 6},
	}
	g.Cmd.Pax = 2
	g.Screen = "journal"
	shot("g61-journal")

	// witchspace
	g.EnterSystem(g.Cmd.System, "arrive", 9000, true)
	g.Messages = nil
	shot("g51-witchspace")

	// Thargoid mothership in the wild
	g.Cmd.Mission = game.Mission{Kind: game.MissionThargoid, System: g.Cmd.System, Galaxy: g.Cmd.Galaxy, Pay: 25000}
	g.MissionArriveForShots()
	g.Messages = nil
	for _, e := range g.World.Entities {
		if e.Mission {
			d := e.Pos.Sub(g.Player.Pos).Norm()
			g.Player.Pos = e.Pos.Sub(d.Mul(520))
			g.Player.Mat = render.LookMatrix(d)
		}
	}
	shot("g62-mothership")

	// rock hermit and generation ship
	for _, probe := range []struct {
		mod  uint32
		name string
	}{{11, "g49-hermit"}, {13, "g50-genship"}} {
		for i := 0; i < 256; i++ {
			if galaxy.Get(0).Systems[i].Seed%probe.mod == 0 {
				g.EnterSystem(i, "arrive", 9000, false)
				g.Messages = nil
				if probe.mod == 11 {
					st3 := g.World.Station
					lookAt(st3.Pos.Add(st3.Mat.Col(2).Mul(900)), st3.Pos)
				} else {
					for _, e := range g.World.Entities {
						if e.Kind == "genship" {
							d := e.Pos.Sub(g.Player.Pos).Norm()
							g.Player.Pos = e.Pos.Sub(d.Mul(2200))
							g.Player.Mat = render.LookMatrix(d)
						}
					}
				}
				shot(probe.name)
				break
			}
		}
	}

	// a rock hermit, now rare and no longer the starting station
	for i := 0; i < 256; i++ {
		if game.StationMeshForShots(i) == "hermit" {
			g.EnterSystem(i, "arrive", 9000, false)
			g.Messages = nil
			if st4 := g.World.Station; st4 != nil {
				lookAt(st4.Pos.Add(st4.Mat.Col(2).Mul(700)), st4.Pos)
				shot("g57-hermit-port")
			}
			break
		}
	}

	g.State = "docked"
	g.Screen = "shipyard"
	shot("g47-shipyard")

	ui.OpenSlotsForShots(true)
	shot("g48-slots")

	g.Screen = "launch"
	ui.OpenSettings()
	shot("g41-settings")

	g.Screen = "settings"
	ui.S.Settings = 3
	ui.Key("Enter")
	shot("g52-controls")

	ui.DrawPadTestSample(fb)
	shot("g65-padtest")
	g.Screen = "settings"

	g.State = "title"
	g.Screen = "mode"
	ui.S.Mode = 0
	shot("g59-mode")

	g.State = "title"
	g.Screen = "launch"
	g.Time = 3.5
	ui.S.HasSession = true // show the CONTINUE entry in one shot
	shot("g38-title")
	ui.S.HasSession = false
}

// runSim exercises the simulation without any window.
func runSim() {
	g := game.New()
	g.NewGame("SIM")
	fmt.Println("docked at", g.World.Sys.Name, "cash", g.Cmd.Cash)
	g.Buy(0, 3)
	fmt.Println("after buy: cash", g.Cmd.Cash, "cargo", g.Cmd.Cargo)
	g.Launch()
	fmt.Println("launched, entities", len(g.World.Entities))
	g.Cmd.Equip["dockcomp"] = true
	g.ToggleAutopilot()
	docked := false
	for i := 0; i < 7000; i++ {
		g.Update(1 / 60.0)
		if g.State == "docked" {
			docked = true
			fmt.Printf("docked after %.1fs\n", float64(i)/60)
			break
		}
		if g.State == "dead" {
			fmt.Println("died while docking:", g.Dead.Reason)
			break
		}
	}
	if !docked {
		fmt.Println("DOCKING FAILED, pos", g.Player.Pos)
	}
	// hyperspace
	gal := galaxy.Get(0)
	cur := gal.Systems[g.Cmd.System]
	target := -1
	for _, s := range gal.Systems {
		if s.ID != cur.ID && galaxy.Distance(cur, s) < 6.5 {
			target = s.ID
			break
		}
	}
	g.Cmd.Fuel = 7
	if target >= 0 {
		ok := g.Hyperspace(target)
		fmt.Println("hyperspace ok:", ok)
		for i := 0; i < 400; i++ {
			g.Update(1 / 60.0)
			if g.State == "flight" {
				break
			}
		}
		fmt.Println("arrived:", g.State, g.World.Sys.Name)
	}
	// combat
	g.SpawnPirates(3)
	fired := 0
	for i := 0; i < 900; i++ {
		g.In.Fire = true
		g.Update(1 / 60.0)
		if g.State == "dead" {
			break
		}
		fired++
	}
	fmt.Println("combat frames", fired, "state", g.State, "kills", g.Cmd.Kills, "energy", int(g.Player.Energy))
	fmt.Println("planet frames", len(models.MakePlanetFrames(16, 1, 0, 2)))

	// save / load round-trip
	g.State = "docked"
	g.Cmd.Cash = 4242
	g.Cmd.Cargo = map[int]int{3: 2}
	fmt.Println("save:", g.Save(), "exists:", game.SaveExists())
	g.Cmd.Cash = 0
	fmt.Println("load:", g.Load(), "cash:", g.Cmd.Cash, "cargo:", g.Cmd.Cargo, "state:", g.State)
}

func main() {
	shots := flag.Bool("shots", false, "render reference screenshots and exit")
	icon := flag.String("icon", "", "render the app icon PNG to this path and exit")
	sheet := flag.Bool("sheet", false, "render a hull contact sheet to build/sheet.png and exit")
	plain := flag.Bool("plain", false, "disable the extra shading detail (for A/B comparison)")
	langFlag := flag.String("lang", "", "UI language for shots: ru or en")
	sim := flag.Bool("sim", false, "run a headless simulation and exit")
	ver := flag.Bool("version", false, "print the build version and exit")
	probe := flag.Bool("probe", false, "combat probe")
	flag.Parse()

	if *ver {
		fmt.Println(game.ProductFull, game.Version)
		return
	}

	lang.Load()
	input.Load()
	render.Translate = lang.TranslateLine

	game.New()
	if *plain {
		render.ExtraDetail = false
	}
	if *sheet {
		runSheet(*langFlag)
		return
	}
	if *shots {
		if *langFlag != "" {
			lang.Set(*langFlag)
		}
		runShots()
		return
	}
	if *sim {
		runSim()
		return
	}
	if *icon != "" {
		runIcon(*icon)
		return
	}
	if *probe {
		combatProbe()
		return
	}

	// On macOS a bare executable cannot create a window: AppKit aborts with
	// "NSRemoteView encountered app without bundle" as soon as the window is
	// activated. The GUI therefore has to live inside Elite.app.
	if runtime.GOOS == "darwin" && os.Getenv("HELIOS_ALLOW_BARE") == "" {
		exe, _ := os.Executable()
		if !strings.Contains(exe, ".app/Contents/MacOS/") {
			fmt.Println(game.ProductFull + " нужно запускать из бандла:")
			fmt.Println("    open HeliosDrift.app     (или двойной клик по HeliosDrift.app)")
			fmt.Println()
			fmt.Println("macOS не даёт голому бинарю открыть окно и роняет его в AppKit.")
			fmt.Println("Симуляция без окна работает и так:  ./heliosdrift -sim")
			fmt.Println("Если окно всё равно нужно: HELIOS_ALLOW_BARE=1 ./heliosdrift")
			os.Exit(2)
		}
	}

	g := game.G
	g.NewGame("JAMESON")
	g.State = "title"

	bank = sfx.Init()
	bank.SetMaster(lang.Cfg.Volume)
	bank.Muted = lang.Cfg.Muted
	g.OnEvent = func(name string) {
		switch name {
		case "laser":
			bank.Laser.Play()
		case "hit":
			bank.Hit.Play()
		case "explode":
			bank.Explode.Play()
		case "missile":
			bank.Missile.Play()
		case "ecm":
			bank.ECM.Play()
		case "bomb":
			bank.Bomb.Play()
		case "dock":
			bank.Dock.Play()
		case "hyper":
			bank.Hyper.Play()
		case "warning":
			bank.Warning.Play()
		}
	}

	ebiten.SetWindowTitle(game.ProductFull)
	ebiten.SetWindowSize(ScreenW*3, ScreenH*3)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	app := newApp()
	if err := ebiten.RunGame(app); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var _ = strings.ToUpper

// runSheet renders every hull into a labelled contact sheet, so the shapes can
// be checked without launching the game.
func runSheet(langName string) {
	fb := render.NewFB(960, 640)
	fb.Clear(render.RGB(0x0a0e16))
	cam := render.NewCamera(960, 640, 55)
	names := []string{
		"cobra3", "cobra1", "sidewinder", "krait", "mamba", "adder",
		"gecko", "worm", "moray", "python", "boa", "anaconda",
		"transporter", "shuttle", "viper", "asp", "ferdelance",
		"thargoid", "thargon", "constrictor", "canister", "capsule",
		"missile", "asteroid", "station", "dodeca",
	}
	cells := 5
	rows := (len(names) + cells - 1) / cells
	cw, ch := 960/cells, 640/rows
	for i, n := range names {
		cx := (i%cells)*cw + cw/2
		cy := (i/cells)*ch + ch/2
		mesh := models.Get(n)
		if mesh == nil {
			continue
		}
		c := render.NewCamera(cw, ch, 55)
		c.CX, c.CY = float64(cw)/2, float64(ch)/2
		// fit the hull into the cell
		z := -mesh.Radius * 2.0
		cam2 := *c
		_ = cam2
		sub := render.NewFB(cw, ch)
		sub.Clear(render.RGB(0x0a0e16))
		mat := render.RotX(0.42).Mul(render.RotY(0.72))
		render.DrawMesh(sub, c, mesh, render.V3(0, 0, z), mat, 90000)
		render.BlitScaled(fb, sub, float64(cx-cw/2), float64(cy-ch/2), float64(cw), float64(ch), 1e9)
		fb.TextC(n, cx, cy+ch/2-14, render.RGB(0x6f8ba8), 1)
	}
	_ = cam
	_ = langName
	if err := writePNG(fb, "build/sheet.png"); err != nil {
		println("sheet:", err.Error())
	}
}

// firstGamepad returns the first connected pad, or -1.
func firstGamepad() ebiten.GamepadID {
	ids := ebiten.AppendGamepadIDs(nil)
	if len(ids) == 0 {
		return -1
	}
	return ids[0]
}

// dispatchFlightActions turns held keys/buttons into the one-shot commands the
// flying code expects, using the rebindable table.
var prevAct = map[input.Action]bool{}

func dispatchFlightActions() {
	binds := []struct {
		act input.Action
		tok string
	}{
		{input.Missile, "m"},
		{input.ECM, "e"},
		{input.Bomb, "Tab"},
		{input.Autopilot, "l"},
		{input.JumpInSystem, "j"},
		{input.GalacticJump, "g"},
		{input.Jettison, "k"},
		{input.Chart, "c"},
		{input.Status, "v"},
		{input.Pause, "Escape"},
	}
	for _, b := range binds {
		down := input.Held(b.act)
		if down && !prevAct[b.act] {
			ui.Key(b.tok)
		}
		prevAct[b.act] = down
	}
	if down := input.Held(input.Sound); down && !prevAct[input.Sound] {
		toggleMute()
	}
	prevAct[input.Sound] = input.Held(input.Sound)
}

func toggleMute() {
	if bank == nil {
		return
	}
	muted := bank.ToggleMute()
	lang.Cfg.Muted = muted
	lang.Save()
	if muted {
		game.G.Message("SOUND OFF", "info")
	} else {
		game.G.Message("SOUND ON", "info")
	}
}
