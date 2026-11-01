package ui

import (
	"testing"
	"time"

	"heliosdrift/game"
	"heliosdrift/lang"
	"heliosdrift/render"
)

// A real pilot has a full log, not an empty one.
func TestFullLogDoesNotHang(t *testing.T) {
	lang.Set("ru")
	defer lang.Set("en")
	fb := render.NewFB(320, 224)
	g := game.New()
	g.NewGame("HANG")
	g.State = "docked"
	for i := 0; i < 300; i++ {
		g.Note(game.LogBuy, "FOOD", 1+i%5, float64(-20*i))
	}
	g.Screen = "log"
	done := make(chan struct{})
	go func() {
		for i := 0; i < 40; i++ {
			Draw(fb, 1.0/60.0)
			for _, k := range []string{"ArrowDown", "ArrowDown", "ArrowUp", "PageDown", "PageUp", "ArrowRight", "ArrowLeft"} {
				Key(k)
			}
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("the log screen hangs with a full log")
	}
}
