package ui

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// One screenshot per screen, so every menu and every page can be looked at in a
// single pass instead of one at a time.
func TestShootEveryScreen(t *testing.T) {
	if os.Getenv("SHOOT") == "" {
		t.Skip("set SHOOT=1")
	}
	if l := os.Getenv("SHOOTLANG"); l != "" {
		lang.Set(l)
	} else {
		lang.Set("ru")
	}
	// the game installs this in main(); without it the screens draw English
	render.Translate = lang.TranslateLine
	defer func() { render.Translate = nil; lang.Set("ru") }()
	dir := os.Getenv("SHOOTDIR")
	if dir == "" {
		dir = "/tmp/hd-shots"
	}
	os.MkdirAll(dir, 0o755)

	fb := render.NewFB(320, 224)
	shot := func(name string) {
		Draw(fb, 1.0/60.0)
		img := image.NewRGBA(image.Rect(0, 0, 320, 224))
		// the buffer is 0xBBGGRR, same as main.go's writePNG
		for i, v := range fb.Buf {
			img.Set(i%320, i/320, color.RGBA{uint8(v), uint8(v >> 8), uint8(v >> 16), 255})
		}
		f, err := os.Create(filepath.Join(dir, name+".png"))
		if err != nil {
			t.Fatal(err)
		}
		png.Encode(f, img)
		f.Close()
	}
	setup := func() *game.Game {
		g := game.New()
		g.NewGame("JAMESON")
		return g
	}
	docked := func() *game.Game {
		g := setup()
		g.State = "docked"
		g.Screen = "launch"
		g.Docked()
		g.Cmd.Cash = 12500
		g.Cmd.Fuel = 5.5
		g.Cmd.Cargo[0] = 4
		g.Cmd.Cargo[3] = 2
		g.Messages = nil
		return g
	}

	g := setup()
	g.State = "title"
	shot("01-title")
	g.Screen = "mode"
	shot("02-mode")
	g.Screen = "settings"
	shot("03-settings")
	g.Screen = "controls"
	shot("04-controls")

	g = docked()
	shot("05-station")
	S.Menu = 1
	g.Screen = "market"
	shot("06-market")
	S.MarketMode = "buy"
	S.MarketQty = 10
	shot("07-market-amount")
	S.MarketMode = ""

	g.Screen = "contracts"
	shot("08-contracts")
	g.Screen = "equip"
	shot("09-equipment")
	g.Screen = "shipyard"
	shot("10-shipyard")
	g.Screen = "status"
	shot("11-status")
	g.Screen = "log"
	shot("12-log")
	g.Cmd.Flags["HAS_SEAL"] = true
	g.Cmd.Flags["PAID_PERMIT"] = true
	g.Screen = "diary"
	S.Diary = 0
	S.DiaryOpen = false
	shot("13-diary")
	S.DiaryOpen = true
	shot("13-diary-page")
	S.DiaryOpen = false
	g.Screen = "slots"
	shot("14-slots")
	for i := range helpPages {
		S.Help = i
		S.HelpScroll = 0
		g.Screen = "help"
		shot("15-help-" + itoa(i+1))
	}

	g = docked()
	g.Launch()
	g.Overlay = "pause"
	shot("16-pause")
	g.Overlay = ""
	shot("17-flight")
	g.Overlay = "chart"
	shot("18-chart")
	g.Overlay = "status"
	shot("19-flight-status")
	g.Overlay = ""

	g = docked()
	if g.TalkStart("c1_disp") {
		shot("20-talk")
	}

	g = docked()
	g.Player.Invulnerable = 0
	g.DamagePlayer(9999, true, "laser")
	shot("21-death")

	g = docked()
	g.StartEnding("combine")
	shot("22-ending")
}

// Contact sheets: every screenshot in one picture, scaled up so the text reads.
func TestMakeSheets(t *testing.T) {
	if os.Getenv("SHOOT") == "" {
		t.Skip("set SHOOT=1")
	}
	dir := os.Getenv("SHOOTDIR")
	if dir == "" {
		dir = "/tmp/hd-shots"
	}
	files, _ := filepath.Glob(filepath.Join(dir, "*.png"))
	sort.Strings(files)
	const cols, scale = 4, 2
	tw, th := 320*scale, 224*scale
	perSheet := cols * 3
	for s := 0; s*perSheet < len(files); s++ {
		end := (s + 1) * perSheet
		if end > len(files) {
			end = len(files)
		}
		part := files[s*perSheet : end]
		rows := (len(part) + cols - 1) / cols
		sheet := image.NewRGBA(image.Rect(0, 0, cols*tw, rows*th))
		draw.Draw(sheet, sheet.Bounds(), &image.Uniform{color.RGBA{20, 20, 20, 255}}, image.Point{}, draw.Src)
		for i, f := range part {
			in, err := os.Open(f)
			if err != nil {
				continue
			}
			src, err := png.Decode(in)
			in.Close()
			if err != nil {
				continue
			}
			cell := image.NewRGBA(image.Rect(0, 0, tw, th))
			for y := 0; y < 224; y++ {
				for x := 0; x < 320; x++ {
					c := src.At(x, y)
					for dy := 0; dy < scale; dy++ {
						for dx := 0; dx < scale; dx++ {
							cell.Set(x*scale+dx, y*scale+dy, c)
						}
					}
				}
			}
			ox, oy := (i%cols)*tw, (i/cols)*th
			draw.Draw(sheet, image.Rect(ox, oy, ox+tw, oy+th), cell, image.Point{}, draw.Src)
		}
		out := filepath.Join(dir, "sheet-"+itoa(s+1)+".png")
		fh, err := os.Create(out)
		if err != nil {
			t.Fatal(err)
		}
		png.Encode(fh, sheet)
		fh.Close()
		t.Log("wrote", out)
	}
}
