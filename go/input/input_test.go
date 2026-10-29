package input

import (
	"fmt"
	"os"
	"path/filepath"
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
	SetPad(Fire, PadY)
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
	if got := KeyDesc(RollLeft); got != "Q/A" {
		t.Fatalf("roll left = %q, want Q/A (the second key is kept)", got)
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
	if got := KeyDesc(Accel); got != "W/=" {
		t.Fatalf("accel = %q, want W/= (primary replaced, second kept)", got)
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
	Cur[Sound] = Binding{Keys: []int32{int32(ebiten.KeyM)}, Btn: noBtn, Axis: noAxis}
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

// KeyA is the zero value of ebiten.Key, so a scan that starts at Key(1) skips
// it. That is exactly how the A key became the one key nobody could bind.
func TestKeyAIsScanned(t *testing.T) {
	if ebiten.KeyA != 0 {
		t.Skipf("ebiten.KeyA is %d; the scan start needs revisiting", ebiten.KeyA)
	}
	if !Bindable(ebiten.KeyA) {
		t.Fatal("A must be bindable")
	}
	SetKey(RollLeft, int32(ebiten.KeyA))
	if got := KeyDesc(RollLeft); got != "A" {
		t.Fatalf("roll left = %q", got)
	}
	Reset()
}

// The stock layout has to show and honour A: it is the zero key value.
func TestDefaultsIncludeA(t *testing.T) {
	Reset()
	if got := KeyDesc(RollLeft); got != "LEFT/A" {
		t.Fatalf("roll left = %q, want LEFT/A", got)
	}
	if !KeysPressed(RollLeft) {
		// nothing is physically pressed in a test, but the scan must not skip A
		for _, k := range Cur[RollLeft].Keys {
			if k < 0 {
				t.Fatalf("roll left has an unset key: %v", Cur[RollLeft].Keys)
			}
		}
	}
}

// An action can carry a primary and a second key, and both must work.
func TestSecondKey(t *testing.T) {
	Reset()
	SetKey(RollLeft, int32(ebiten.KeyQ))
	SetSecondKey(RollLeft, int32(ebiten.KeyE))
	if got := KeyDesc(RollLeft); got != "Q/E" {
		t.Fatalf("roll left = %q, want Q/E", got)
	}
	// the second key is taken from whoever had it
	SetSecondKey(Fire, int32(ebiten.KeyA))
	if got := KeyDesc(RollLeft); got != "Q/E" {
		t.Fatalf("roll left changed to %q", got)
	}
	SetSecondKey(RollRight, int32(ebiten.KeyE))
	if got := KeyDesc(RollLeft); got != "Q" {
		t.Fatalf("roll left = %q, want Q after E moved away", got)
	}
	ClearSecondKey(RollRight)
	if got := KeyDesc(RollRight); got != "RIGHT" {
		t.Fatalf("roll right = %q", got)
	}
	Reset()
}

// Rebinding the primary must not silently drop the second key, and the second
// key must never duplicate the primary.
func TestSecondKeyEdgeCases(t *testing.T) {
	Reset()
	SetKey(RollLeft, int32(ebiten.KeyA)) // A is already the second key
	if got := KeyDesc(RollLeft); got != "A" {
		t.Fatalf("roll left = %q, want A without a duplicate", got)
	}
	Reset()
	SetSecondKey(RollLeft, int32(ebiten.KeyLeft)) // same as the primary
	if got := KeyDesc(RollLeft); got != "LEFT/A" {
		t.Fatalf("roll left = %q", got)
	}
	Reset()
}

// The D-pad must sit at the standard positions: 12 up, 13 down, 14 left,
// 15 right. Using raw device indices here is what made "D-pad down" dead.
func TestStandardPadIndices(t *testing.T) {
	if PadUp != 12 || PadDown != 13 || PadLeft != 14 || PadRight != 15 {
		t.Fatalf("d-pad positions: up=%d down=%d left=%d right=%d", PadUp, PadDown, PadLeft, PadRight)
	}
	if PadA != 0 || PadB != 1 || PadX != 2 || PadY != 3 {
		t.Fatalf("face buttons: a=%d b=%d x=%d y=%d", PadA, PadB, PadX, PadY)
	}
	if PadLB != 4 || PadRB != 5 || PadStart != 9 || PadBack != 8 {
		t.Fatalf("shoulders/start: lb=%d rb=%d start=%d back=%d", PadLB, PadRB, PadStart, PadBack)
	}
}

// Axes are bindable with a direction and show up in the PAD column.
func TestAxisBinding(t *testing.T) {
	Reset()
	SetAxis(PitchUp, AxisRY, -1)
	if got := PadDesc(PitchUp); got != "RSY-" {
		t.Fatalf("pitch up pad = %q", got)
	}
	if a, d := AxisDir(PitchUp); a != AxisRY || d != -1 {
		t.Fatalf("axis=%d dir=%d", a, d)
	}
	// the same axis direction cannot drive two actions
	SetAxis(RollLeft, AxisRY, -1)
	if a, _ := AxisDir(PitchUp); a != noAxis {
		t.Fatalf("pitch up still holds the axis: %d", a)
	}
	Reset()
}

// The stock layout names the familiar buttons.
func TestDefaultPadNames(t *testing.T) {
	Reset()
	for a, want := range map[Action]string{
		Fire: "A", Missile: "B", ECM: "X", Autopilot: "Y",
		Bomb: "LB", Chart: "RB", JumpInSystem: "START", GalacticJump: "BACK",
		Accel: "RT", Decel: "LT", PitchUp: "LSY-", RollLeft: "LSX-",
	} {
		if got := PadDesc(a); got != want {
			t.Fatalf("%s = %q, want %q", Name(a), got, want)
		}
	}
	Reset()
}

// A bindings file from before buttons were standard keeps its keys and gets
// the default pad bindings, instead of binding some random raw index.
func TestLegacyKeysFileMerges(t *testing.T) {
	DirOverride = t.TempDir()
	defer func() { DirOverride = "" }()

	legacy := fmt.Sprintf(`{"binds":{"FIRE":{"keys":[%d],"pad":-1},"PITCH UP":{"keys":[%d],"pad":-1}}}`,
		int(ebiten.KeySpace), int(ebiten.KeyW))
	if err := os.WriteFile(filepath.Join(DirOverride, ".heliosdrift-keys.json"), []byte(legacy), 0o644); err != nil {
		t.Fatal(err)
	}
	if !Load() {
		t.Fatal("load")
	}
	if got := KeyDesc(Fire); got != "SPACE" {
		t.Fatalf("fire key = %q", got)
	}
	if got := PadDesc(Fire); got != "A" {
		t.Fatalf("fire pad = %q, want the default A", got)
	}
	if got := PadDesc(PitchUp); got != "LSY-" {
		t.Fatalf("pitch up pad = %q", got)
	}
	Reset()
}

// A stick counts when it crosses the threshold, not when it is already held.
func TestAxisCrossing(t *testing.T) {
	cases := []struct {
		prev, cur float64
		dir       int32
		ok        bool
	}{
		{0, 0.8, 1, true},
		{0, -0.9, -1, true},
		{0.7, 0.9, 0, false}, // already past: no second trigger
		{0.5, 0.3, 0, false}, // never crossed
		{-0.8, -0.7, 0, false},
	}
	for _, c := range cases {
		dir, ok := axisCrossed(c.prev, c.cur)
		if ok != c.ok || dir != c.dir {
			t.Fatalf("axisCrossed(%v,%v) = %d,%v want %d,%v", c.prev, c.cur, dir, ok, c.dir, c.ok)
		}
	}
}

// An axis index outside the standard four is ignored rather than read.
func TestOutOfRangeAxisIsIgnored(t *testing.T) {
	Reset()
	Cur[Accel] = Binding{Keys: nil, Btn: noBtn, Axis: 5, Dir: 1}
	if v := AxisValue(Accel); v != 0 {
		t.Fatalf("axis 5 returned %v", v)
	}
	if got := PadDesc(Accel); got != "-" {
		t.Fatalf("pad desc = %q", got)
	}
	Reset()
}
