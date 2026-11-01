package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

func dumpShot(t *testing.T, fb *render.FB, path string) {
	img := image.NewRGBA(image.Rect(0, 0, fb.W, fb.H))
	for i, v := range fb.Buf {
		img.Set(i%fb.W, i/fb.W, color.RGBA{uint8(v >> 16), uint8(v >> 8), uint8(v), 255})
	}
	f, _ := os.Create(path)
	png.Encode(f, img)
	f.Close()
}

func TestShotJournalAndLog(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	home, _ := os.UserHomeDir()
	dir := t.TempDir()
	for slot := 1; slot <= 3; slot++ {
		src := filepath.Join(home, fmt.Sprintf(".heliosdrift-save%d.json", slot))
		raw, err := os.ReadFile(src)
		if err == nil {
			os.WriteFile(filepath.Join(dir, fmt.Sprintf(".heliosdrift-save%d.json", slot)), raw, 0o644)
		}
	}
	game.SaveDirOverride = dir
	defer func() { game.SaveDirOverride = "" }()

	fb := render.NewFB(320, 224)
	for slot := 0; slot < 3; slot++ {
		g := game.New()
		if !g.LoadSlot(slot) {
			continue
		}
		g.Screen = "journal"
		Draw(fb, 1.0/60.0)
		dumpShot(t, fb, fmt.Sprintf("/tmp/save%d-journal.png", slot+1))
		g.Screen = "log"
		Draw(fb, 1.0/60.0)
		dumpShot(t, fb, fmt.Sprintf("/tmp/save%d-log.png", slot+1))
	}
}
