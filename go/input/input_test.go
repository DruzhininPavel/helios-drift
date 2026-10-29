package input

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestDefaultsCoverEveryAction(t *testing.T) {
	d := Defaults()
	if len(d) != int(NumActions) {
		t.Fatalf("%d defaults for %d actions", len(d), NumActions)
	}
	for a := Action(0); a < NumActions; a++ {
		if _, ok := d[a]; !ok {
			t.Errorf("no default for %s", Name(a))
		}
		if Name(a) == "" || Name(a) == "?" {
			t.Errorf("action %d has no name", a)
		}
	}
}

func TestRebindAndReset(t *testing.T) {
	Reset()
	before := KeyDesc(Fire)
	SetKey(Fire, int32(ebiten.KeyQ))
	if KeyDesc(Fire) == before {
		t.Fatal("rebind did not change the key")
	}
	if KeyDesc(Fire) == "-" {
		t.Fatal("key description is empty after rebind")
	}
	SetPad(Fire, 3)
	if PadDesc(Fire) != "Y" {
		t.Fatalf("pad = %q, want Y", PadDesc(Fire))
	}
	Reset()
	if KeyDesc(Fire) != before {
		t.Fatalf("reset did not restore: %q vs %q", KeyDesc(Fire), before)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	DirOverride = t.TempDir()
	defer func() { DirOverride = "" }()
	Reset()
	SetKey(Missile, int32(ebiten.KeyQ))
	SetPad(Missile, 3)
	if !Save() {
		t.Fatal("save failed")
	}
	Reset()
	if KeyDesc(Missile) == "Q" {
		t.Fatal("reset did not clear the binding")
	}
	if !Load() {
		t.Fatal("load failed")
	}
	if KeyDesc(Missile) != "Q" {
		t.Fatalf("key after load = %q, want Q", KeyDesc(Missile))
	}
	if PadDesc(Missile) != "Y" {
		t.Fatalf("pad after load = %q, want Y", PadDesc(Missile))
	}
}

func TestKeyNames(t *testing.T) {
	if KeyName(int32(0)) == "" {
		t.Fatal("empty name for key 0")
	}
	for a := Action(0); a < NumActions; a++ {
		if KeyDesc(a) == "" {
			t.Errorf("%s has empty key description", Name(a))
		}
	}
}
