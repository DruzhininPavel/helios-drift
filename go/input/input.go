// Package input holds the rebindable action map: every flying control can be
// pointed at a keyboard key or a gamepad button, and reset to defaults.
package input

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Action int

const (
	PitchUp Action = iota
	PitchDown
	RollLeft
	RollRight
	Fire
	Accel
	Decel
	Missile
	ECM
	Bomb
	Autopilot
	JumpInSystem
	GalacticJump
	Jettison
	Chart
	Status
	Pause
	Sound
	NumActions
)

var actionNames = [NumActions]string{
	"PITCH UP", "PITCH DOWN", "ROLL LEFT", "ROLL RIGHT", "FIRE",
	"ACCELERATE", "DECELERATE", "MISSILE", "E.C.M.", "ENERGY BOMB",
	"DOCKING COMPUTER", "JUMP IN SYSTEM", "GALACTIC JUMP", "JETTISON",
	"CHART", "STATUS", "PAUSE", "SOUND",
}

// Name is the label shown on the controls screen.
func Name(a Action) string {
	if a < 0 || a >= NumActions {
		return "?"
	}
	return actionNames[a]
}

// Actions lists every bindable action in display order.
func Actions() []Action {
	out := make([]Action, 0, NumActions)
	for a := Action(0); a < NumActions; a++ {
		out = append(out, a)
	}
	return out
}

// Binding is up to two keyboard keys plus one gamepad button. Two keys let the
// defaults keep both the arrows and WASD.
type Binding struct {
	Keys []int32 `json:"keys"`
	Pad  int32   `json:"pad"` // -1 when unbound
}

// Map is the full action table.
type Map map[Action]Binding

// noPad and noKey are the "unbound" sentinels. Keys cannot use 0 for that:
// ebiten.KeyA is the zero value, and treating it as "empty" made the A key
// impossible to bind, display or press.
const (
	noPad = int32(-1)
	noKey = int32(-1)
)

// Defaults returns the stock layout: arrows and WASD for steering, the keys the
// game has always used for everything else, and pad buttons for the common acts.
func Defaults() Map {
	k := func(keys ...ebiten.Key) []int32 {
		out := make([]int32, len(keys))
		for i, key := range keys {
			out[i] = int32(key)
		}
		return out
	}
	p := func(b ebiten.GamepadButton) int32 { return int32(b) }
	return Map{
		PitchUp:      {Keys: k(ebiten.KeyArrowUp, ebiten.KeyW), Pad: noPad},
		PitchDown:    {Keys: k(ebiten.KeyArrowDown, ebiten.KeyS), Pad: noPad},
		RollLeft:     {Keys: k(ebiten.KeyArrowLeft, ebiten.KeyA), Pad: noPad},
		RollRight:    {Keys: k(ebiten.KeyArrowRight, ebiten.KeyD), Pad: noPad},
		Fire:         {Keys: k(ebiten.KeySpace), Pad: p(ebiten.GamepadButton0)},
		Accel:        {Keys: k(ebiten.KeyX, ebiten.KeyEqual), Pad: p(ebiten.GamepadButton5)},
		Decel:        {Keys: k(ebiten.KeyZ, ebiten.KeyMinus), Pad: p(ebiten.GamepadButton4)},
		Missile:      {Keys: k(ebiten.KeyM), Pad: p(ebiten.GamepadButton1)},
		ECM:          {Keys: k(ebiten.KeyE), Pad: p(ebiten.GamepadButton2)},
		Bomb:         {Keys: k(ebiten.KeyTab), Pad: p(ebiten.GamepadButton4)},
		Autopilot:    {Keys: k(ebiten.KeyL), Pad: p(ebiten.GamepadButton3)},
		JumpInSystem: {Keys: k(ebiten.KeyJ), Pad: p(ebiten.GamepadButton7)},
		GalacticJump: {Keys: k(ebiten.KeyG), Pad: p(ebiten.GamepadButton6)},
		Jettison:     {Keys: k(ebiten.KeyK), Pad: noPad},
		Chart:        {Keys: k(ebiten.KeyC), Pad: p(ebiten.GamepadButton8)},
		Status:       {Keys: k(ebiten.KeyV), Pad: noPad},
		Pause:        {Keys: k(ebiten.KeyEscape), Pad: p(ebiten.GamepadButton9)},
		Sound:        {Keys: k(ebiten.KeyF10), Pad: noPad},
	}
}

// Cur is the live binding table.
var Cur = Defaults()

// Reset restores the stock layout.
func Reset() { Cur = Defaults() }

// SetKey binds a single key to an action. The key is taken away from any other
// action that had it: without that, changing a binding left the old behaviour
// in place and two actions fired at once.
func SetKey(a Action, key int32) {
	for other, b := range Cur {
		if other == a {
			continue
		}
		kept := b.Keys[:0]
		for _, k := range b.Keys {
			if k != key {
				kept = append(kept, k)
			}
		}
		b.Keys = kept
		Cur[other] = b
	}
	b := Cur[a]
	b.Keys = []int32{key}
	Cur[a] = b
}

// SetPad binds a gamepad button to an action, taking it from other actions.
func SetPad(a Action, btn int32) {
	for other, b := range Cur {
		if other == a {
			continue
		}
		if b.Pad == btn {
			b.Pad = noPad
			Cur[other] = b
		}
	}
	b := Cur[a]
	b.Pad = btn
	Cur[a] = b
}

// KeysPressed reports whether any key bound to the action is held.
func KeysPressed(a Action) bool {
	b, ok := Cur[a]
	if !ok {
		return false
	}
	for _, k := range b.Keys {
		if k >= 0 && ebiten.IsKeyPressed(ebiten.Key(k)) {
			return true
		}
	}
	return false
}

// PadButton returns the bound pad button, or -1.
func PadButton(a Action) int32 {
	if b, ok := Cur[a]; ok {
		return b.Pad
	}
	return noPad
}

// Held reports whether the action is active on keyboard or pad.
func Held(a Action) bool {
	if KeysPressed(a) {
		return true
	}
	if id := firstPad(); id >= 0 {
		if b := PadButton(a); b >= 0 && ebiten.IsGamepadButtonPressed(ebiten.GamepadID(id), ebiten.GamepadButton(b)) {
			return true
		}
	}
	return false
}

// PadDesc describes the bound pad button for the UI.
func PadDesc(a Action) string {
	b := PadButton(a)
	if b < 0 {
		return "-"
	}
	return PadName(b)
}

// KeyDesc describes the bound keys for the UI, e.g. "UP/W".
func KeyDesc(a Action) string {
	b, ok := Cur[a]
	if !ok || len(b.Keys) == 0 {
		return "-"
	}
	out := ""
	for i, k := range b.Keys {
		if k < 0 {
			continue
		}
		if i > 0 && out != "" {
			out += "/"
		}
		out += KeyName(k)
	}
	if out == "" {
		return "-"
	}
	return out
}

func firstPad() ebiten.GamepadID {
	ids := ebiten.AppendGamepadIDs(nil)
	if len(ids) == 0 {
		return -1
	}
	return ids[0]
}

// Capture returns the first key or pad button pressed this frame, for the
// rebinding prompt. ok is false when nothing is down.
// CaptureDown reports a bindable key that is simply held. It is the fallback
// for keys the platform never reports as a fresh press (a key that went down
// before the window had focus, or a sticky/dead key on an unusual layout).
func CaptureDown() (kind string, code int32, ok bool) {
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if !Bindable(k) {
			continue
		}
		if ebiten.IsKeyPressed(k) {
			return "key", int32(k), true
		}
	}
	return "", 0, false
}

// PeekDown names a bindable key currently held, for the rebinding prompt.
func PeekDown() string {
	if _, code, ok := CaptureDown(); ok {
		return KeyName(code)
	}
	return ""
}

// Capture looks for a key or pad button pressed in this frame. It deliberately
// watches for a fresh press, not for a held key: the Enter that starts the
// prompt is often still down, and a held-key check bound Enter instead of
// whatever the player was aiming for.
func Capture() (kind string, code int32, ok bool) {
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if !Bindable(k) {
			continue
		}
		if inpututil.IsKeyJustPressed(k) {
			return "key", int32(k), true
		}
	}
	if id := firstPad(); id >= 0 {
		for b := ebiten.GamepadButton(0); b <= ebiten.GamepadButtonMax; b++ {
			if inpututil.IsGamepadButtonJustPressed(id, b) {
				return "pad", int32(b), true
			}
		}
	}
	return "", 0, false
}

// Bindable reports whether a key may be assigned to an action. Modifiers are
// excluded because people hold Shift while rebinding; Enter and Escape are
// excluded because Enter starts the prompt and would be captured by its own
// keypress, and Escape cancels it. Everything else must have a display name,
// which also filters out platform keys we know nothing about.
func Bindable(k ebiten.Key) bool {
	switch k {
	case ebiten.KeyShift, ebiten.KeyShiftLeft, ebiten.KeyShiftRight,
		ebiten.KeyControl, ebiten.KeyControlLeft, ebiten.KeyControlRight,
		ebiten.KeyAlt, ebiten.KeyAltLeft, ebiten.KeyAltRight,
		ebiten.KeyMeta, ebiten.KeyMetaLeft, ebiten.KeyMetaRight,
		ebiten.KeyEnter, ebiten.KeyEscape, ebiten.KeyCapsLock:
		return false
	}
	_, named := keyNames[int32(k)]
	return named
}

// repair fixes bindings files written by older builds: SOUND used to sit on
// the missile key, so pressing M fired a missile and muted the game.
func repair() {
	if KeyDesc(Sound) == KeyDesc(Missile) {
		Cur[Sound] = Binding{Keys: []int32{int32(ebiten.KeyF10)}, Pad: noPad}
	}
}

// ---------------------------------------------------------------- persistence

type keyFile struct {
	Binds map[string]Binding `json:"binds"`
}

// DirOverride lets tests keep bindings out of the real home directory.
var DirOverride string

func keyPath() string {
	if DirOverride != "" {
		return filepath.Join(DirOverride, ".heliosdrift-keys.json")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "heliosdrift-keys.json"
	}
	return filepath.Join(home, ".heliosdrift-keys.json")
}

// Save writes the current bindings.
func Save() bool {
	f := keyFile{Binds: map[string]Binding{}}
	for a, b := range Cur {
		f.Binds[Name(a)] = b
	}
	raw, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return false
	}
	return os.WriteFile(keyPath(), raw, 0o644) == nil
}

// Load reads saved bindings, falling back to defaults per action.
func Load() bool {
	raw, err := os.ReadFile(keyPath())
	if err != nil {
		return false
	}
	var f keyFile
	if err := json.Unmarshal(raw, &f); err != nil {
		return false
	}
	Cur = Defaults()
	for a := Action(0); a < NumActions; a++ {
		if b, ok := f.Binds[Name(a)]; ok {
			Cur[a] = b
		}
	}
	repair()
	return true
}
