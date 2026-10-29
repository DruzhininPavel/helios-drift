package input

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// Rebound keys must survive a restart: the file is what the player expects to
// carry their layout, not just the running process.
func TestBindRoundTrip(t *testing.T) {
	DirOverride = t.TempDir()
	defer func() { DirOverride = "" }()

	Reset()
	SetKey(PitchUp, int32(ebiten.KeyW))
	SetKey(PitchDown, int32(ebiten.KeyS))
	SetPad(Fire, int32(ebiten.GamepadButton3))
	if !Save() {
		t.Fatal("save failed")
	}

	Reset()
	if !Load() {
		t.Fatal("load failed")
	}
	if got := KeyDesc(PitchUp); got != "W" {
		t.Fatalf("pitch up = %q, want W", got)
	}
	if got := KeyDesc(PitchDown); got != "S" {
		t.Fatalf("pitch down = %q, want S", got)
	}
	if got := PadDesc(Fire); got == "-" {
		t.Fatalf("pad binding lost")
	}
}

// Rebinding one action must not disturb the others.
func TestBindIsolation(t *testing.T) {
	Reset()
	SetKey(RollLeft, int32(ebiten.KeyQ))
	if got := KeyDesc(RollLeft); got != "Q" {
		t.Fatalf("roll left = %q", got)
	}
	if got := KeyDesc(PitchUp); got != "UP/W" {
		t.Fatalf("pitch up changed to %q", got)
	}
	Reset()
}

// Every action needs a distinct name: bindings are stored by name, and a
// duplicate would silently overwrite another row on load.
func TestActionNamesUnique(t *testing.T) {
	seen := map[string]Action{}
	for _, a := range Actions() {
		n := Name(a)
		if n == "?" || n == "" {
			t.Fatalf("action %d has no name", a)
		}
		if prev, ok := seen[n]; ok {
			t.Fatalf("%q used by actions %d and %d", n, prev, a)
		}
		seen[n] = a
	}
}

// Binding a key must take it from whoever had it, otherwise "change the
// controls to WASD" silently leaves the old action on the same key and both
// fire at once.
func TestBindStealsKey(t *testing.T) {
	Reset()
	if got := KeyDesc(PitchUp); got != "UP/W" {
		t.Fatalf("default pitch up = %q", got)
	}
	SetKey(Accel, int32(ebiten.KeyW))
	if got := KeyDesc(PitchUp); got == "UP/W" || got == "W" {
		t.Fatalf("pitch up still holds W: %q", got)
	}
	if got := KeyDesc(Accel); got != "W" {
		t.Fatalf("accel = %q", got)
	}
	Reset()
}

// The mute key must not sit on the missile key: pressing M used to fire a
// missile and silence the game.
func TestSoundAndMissileDoNotShareAKey(t *testing.T) {
	Reset()
	if KeyDesc(Sound) == KeyDesc(Missile) {
		t.Fatalf("sound and missile share %q", KeyDesc(Sound))
	}
}

// A keys file from before the fix has SOUND on M; loading it must repair the
// clash instead of firing a missile and muting the game at once.
func TestLoadRepairsSoundClash(t *testing.T) {
	DirOverride = t.TempDir()
	defer func() { DirOverride = "" }()

	Reset()
	if !Save() {
		t.Fatal("save")
	}
	Cur[Sound] = Binding{Keys: []int32{int32(ebiten.KeyM)}, Pad: noPad}
	if !Save() {
		t.Fatal("save clash")
	}
	Load()
	if KeyDesc(Sound) == KeyDesc(Missile) {
		t.Fatalf("sound still on %q", KeyDesc(Sound))
	}
	Reset()
}

// Enter starts the rebinding prompt, so it must never be captured by its own
// keypress; modifiers and unnamed platform keys are not bindable either.
func TestBindableKeys(t *testing.T) {
	for _, k := range []ebiten.Key{ebiten.KeyEnter, ebiten.KeyEscape, ebiten.KeyShift, ebiten.KeyAltLeft} {
		if Bindable(k) {
			t.Fatalf("%v should not be bindable", k)
		}
	}
	for _, k := range []ebiten.Key{ebiten.KeyA, ebiten.KeyW, ebiten.KeyS, ebiten.KeyD, ebiten.KeySpace, ebiten.KeyF10} {
		if !Bindable(k) {
			t.Fatalf("%v should be bindable", k)
		}
	}
}
