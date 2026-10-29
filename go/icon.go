package main

import (
	"fmt"
	"math"

	"heliosdrift/game"
	"heliosdrift/models"
	"heliosdrift/render"
	"heliosdrift/ui"
)

// runIcon renders the 1024x1024 application icon.
func runIcon(path string) {
	const S = 1024
	g := game.New()
	g.NewGame("ICON")
	g.State = "title"

	fb := render.NewFB(S, S)
	fb.ClearGradientStops([]uint32{0xff05070f, 0xff0a1422, 0xff122238, 0xff0a1422, 0xff05070f})

	// starfield
	for i := 0; i < 260; i++ {
		x := (i*137 + 41) % S
		y := (i*251 + 17) % S
		b := render.Color("steel", 2+i%4)
		fb.FillRect(x, y, x+1, y+1, b)
		if i%9 == 0 {
			fb.FillRect(x-1, y-1, x+2, y+2, render.Color("steel", 1))
		}
	}

	// cobra, three-quarter view
	cam := render.NewCamera(S, S, 55)
	cam.CX = S / 2
	cam.CY = S/2 - 40
	mat := render.RotX(0.30).Mul(render.RotY(math.Pi * 0.82))
	render.DrawMesh(fb, cam, models.Get("cobra3"), render.V3(0, 0, -190), mat, 90000)

	// engine flare
	ui.GlowAt(fb, S/2-58, S/2+52, 26, "orange")
	ui.GlowAt(fb, S/2+58, S/2+52, 26, "orange")

	// wordmark
	fb.TextC(game.ProductName, S/2+6, 780, render.Color("steel", 0), 12)
	fb.TextC(game.ProductName, S/2+3, 776, render.Color("steel", 2), 12)
	fb.TextC(game.ProductName, S/2, 772, render.Color("steel", 5), 12)
	fb.TextC(game.ProductSub, S/2+2, 872, render.Color("steel", 2), 6)
	fb.TextC(game.ProductSub, S/2, 870, render.Color("steel", 4), 6)

	if err := writePNG(fb, path); err != nil {
		fmt.Println("icon error:", err)
		return
	}
	fmt.Println("wrote", path)
}
