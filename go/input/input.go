// Package input holds the rebindable action map. Every action can carry two
// keyboard keys, one gamepad button and one gamepad axis with a direction, so
// steering can sit on the keyboard, a stick or the triggers.
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

// Binding is two keyboard keys, one gamepad button and one gamepad axis with a
// direction. Buttons are stored in the standard layout rather than as raw
// device indices, so the same binding works on an Xbox pad, a DualSense and
// the Steam Deck.
type Binding struct {
	Keys []int32 `json:"keys"`
	Btn  int32   `json:"btn"`  // standard gamepad button, -1 when unbound
	Axis int32   `json:"axis"` // axis index, -1 when unbound
	Dir  int32   `json:"dir"`  // +1 or -1: the direction that counts as pressed
}

// UnmarshalJSON fills the unbound sentinels first, so a file that predates a
// field does not end up bound to button 0 or axis 0.
func (b *Binding) UnmarshalJSON(data []byte) error {
	type plain Binding
	p := plain{Btn: noBtn, Axis: noAxis}
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	*b = Binding(p)
	return nil
}

// Map is the full action table.
type Map map[Action]Binding

// Unbound sentinels. Keys cannot use 0 for that: ebiten.KeyA is the zero
// value, and treating it as "empty" made A impossible to bind.
const (
	noKey  = int32(-1)
	noBtn  = int32(-1)
	noAxis = int32(-1)

	// AxisThreshold is how far a stick or trigger has to move before a digital
	// action counts as pressed.
	AxisThreshold = 0.6

	// axisDeadzone keeps a resting stick from nudging the ship.
	axisDeadzone = 0.15
)

// Standard gamepad buttons, spelled the way a player sees them.
const (
	PadA     = int32(ebiten.StandardGamepadButtonRightBottom)
	PadB     = int32(ebiten.StandardGamepadButtonRightRight)
	PadX     = int32(ebiten.StandardGamepadButtonRightLeft)
	PadY     = int32(ebiten.StandardGamepadButtonRightTop)
	PadLB    = int32(ebiten.StandardGamepadButtonFrontTopLeft)
	PadRB    = int32(ebiten.StandardGamepadButtonFrontTopRight)
	PadLT    = int32(ebiten.StandardGamepadButtonFrontBottomLeft)
	PadRT    = int32(ebiten.StandardGamepadButtonFrontBottomRight)
	PadBack  = int32(ebiten.StandardGamepadButtonCenterLeft)
	PadStart = int32(ebiten.StandardGamepadButtonCenterRight)
	PadLS    = int32(ebiten.StandardGamepadButtonLeftStick)
	PadRS    = int32(ebiten.StandardGamepadButtonRightStick)
	PadUp    = int32(ebiten.StandardGamepadButtonLeftTop)
	PadDown  = int32(ebiten.StandardGamepadButtonLeftBottom)
	PadLeft  = int32(ebiten.StandardGamepadButtonLeftLeft)
	PadRight = int32(ebiten.StandardGamepadButtonLeftRight)
)

// Axis indices: four sticks, then the two triggers.
const (
	AxisLX = int32(0)
	AxisLY = int32(1)
	AxisRX = int32(2)
	AxisRY = int32(3)
	AxisTL = int32(4)
	AxisTR = int32(5)
)

// Defaults returns the stock layout: arrows and WASD plus the left stick for
// steering, triggers for speed, and the usual buttons for everything else. The
// D-pad is left free in flight — there it is a set of four ordinary buttons,
// and in menus it walks the arrows.
func Defaults() Map {
	k := func(keys ...ebiten.Key) []int32 {
		out := make([]int32, len(keys))
		for i, key := range keys {
			out[i] = int32(key)
		}
		return out
	}
	return Map{
		PitchUp:      {Keys: k(ebiten.KeyArrowUp, ebiten.KeyW), Btn: noBtn, Axis: AxisLY, Dir: -1},
		PitchDown:    {Keys: k(ebiten.KeyArrowDown, ebiten.KeyS), Btn: noBtn, Axis: AxisLY, Dir: 1},
		RollLeft:     {Keys: k(ebiten.KeyArrowLeft, ebiten.KeyA), Btn: noBtn, Axis: AxisLX, Dir: -1},
		RollRight:    {Keys: k(ebiten.KeyArrowRight, ebiten.KeyD), Btn: noBtn, Axis: AxisLX, Dir: 1},
		Fire:         {Keys: k(ebiten.KeySpace), Btn: PadA, Axis: noAxis},
		Accel:        {Keys: k(ebiten.KeyX, ebiten.KeyEqual), Btn: noBtn, Axis: AxisTR, Dir: 1},
		Decel:        {Keys: k(ebiten.KeyZ, ebiten.KeyMinus), Btn: noBtn, Axis: AxisTL, Dir: 1},
		Missile:      {Keys: k(ebiten.KeyM), Btn: PadB, Axis: noAxis},
		ECM:          {Keys: k(ebiten.KeyE), Btn: PadX, Axis: noAxis},
		Bomb:         {Keys: k(ebiten.KeyTab), Btn: PadLB, Axis: noAxis},
		Autopilot:    {Keys: k(ebiten.KeyL), Btn: PadY, Axis: noAxis},
		JumpInSystem: {Keys: k(ebiten.KeyJ), Btn: PadStart, Axis: noAxis},
		GalacticJump: {Keys: k(ebiten.KeyG), Btn: PadBack, Axis: noAxis},
		Jettison:     {Keys: k(ebiten.KeyK), Btn: noBtn, Axis: noAxis},
		Chart:        {Keys: k(ebiten.KeyC), Btn: PadRB, Axis: noAxis},
		Status:       {Keys: k(ebiten.KeyV), Btn: PadLS, Axis: noAxis},
		Pause:        {Keys: k(ebiten.KeyEscape), Btn: PadRS, Axis: noAxis},
		Sound:        {Keys: k(ebiten.KeyF10), Btn: noBtn, Axis: noAxis},
	}
}

// Cur is the live binding table.
var Cur = Defaults()

// Reset restores the stock layout.
func Reset() { Cur = Defaults() }

// ResetAction restores the stock binding of a single action.
func ResetAction(a Action) {
	if d, ok := Defaults()[a]; ok {
		Cur[a] = d
	}
}

// steal takes keys, the button and the axis away from every other action: one
// input drives one action only.
func steal(b Binding, except Action) {
	for other, o := range Cur {
		if other == except {
			continue
		}
		if len(b.Keys) > 0 && len(o.Keys) > 0 {
			kept := o.Keys[:0]
			for _, k := range o.Keys {
				taken := false
				for _, n := range b.Keys {
					if k == n {
						taken = true
					}
				}
				if !taken {
					kept = append(kept, k)
				}
			}
			o.Keys = kept
		}
		if b.Btn >= 0 && o.Btn == b.Btn {
			o.Btn = noBtn
		}
		if b.Axis >= 0 && o.Axis == b.Axis && o.Dir == b.Dir {
			o.Axis, o.Dir = noAxis, 0
		}
		Cur[other] = o
	}
}

// SetKey binds the primary key of an action, keeping its second key.
func SetKey(a Action, key int32) {
	b := Cur[a]
	second := noKey
	if len(b.Keys) > 1 {
		second = b.Keys[1]
	}
	steal(Binding{Keys: []int32{key}}, a)
	if second == key {
		second = noKey
	}
	b.Keys = []int32{key}
	if second >= 0 {
		b.Keys = append(b.Keys, second)
	}
	Cur[a] = b
}

// SetSecondKey binds the second key of an action, leaving the primary alone.
func SetSecondKey(a Action, key int32) {
	b := Cur[a]
	if len(b.Keys) == 0 {
		b.Keys = []int32{key}
		Cur[a] = b
		return
	}
	if b.Keys[0] == key {
		return
	}
	steal(Binding{Keys: []int32{key}}, a)
	b.Keys = []int32{b.Keys[0], key}
	Cur[a] = b
}

// ClearSecondKey drops the second key of an action.
func ClearSecondKey(a Action) {
	b := Cur[a]
	if len(b.Keys) > 1 {
		b.Keys = b.Keys[:1]
		Cur[a] = b
	}
}

// SetPad binds a standard gamepad button to an action.
func SetPad(a Action, btn int32) {
	steal(Binding{Btn: btn}, a)
	b := Cur[a]
	b.Btn = btn
	Cur[a] = b
}

// SetAxis binds an axis direction to an action.
func SetAxis(a Action, axis, dir int32) {
	steal(Binding{Axis: axis, Dir: dir}, a)
	b := Cur[a]
	b.Axis, b.Dir = axis, dir
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

// PadButton returns the bound standard button, or -1.
func PadButton(a Action) int32 {
	if b, ok := Cur[a]; ok {
		return b.Btn
	}
	return noBtn
}

// AxisDir returns the bound axis and its direction, or -1 when unbound.
func AxisDir(a Action) (int32, int32) {
	if b, ok := Cur[a]; ok {
		return b.Axis, b.Dir
	}
	return noAxis, 0
}

func firstPad() ebiten.GamepadID {
	ids := ebiten.AppendGamepadIDs(nil)
	if len(ids) == 0 {
		return -1
	}
	return ids[0]
}

// AxisValue returns how far the action's axis is pushed in its own direction,
// 0..1 — this is what keeps stick steering analogue.
func AxisValue(a Action) float64 {
	axis, dir := AxisDir(a)
	if axis < 0 || dir == 0 {
		return 0
	}
	id := firstPad()
	if id < 0 {
		return 0
	}
	v := ebiten.GamepadAxisValue(id, ebiten.GamepadAxisType(axis)) * float64(dir)
	if v <= axisDeadzone {
		return 0
	}
	// rescale what is left so the stick still reaches 1.0 at the edge
	return (v - axisDeadzone) / (1 - axisDeadzone)
}

// Digital reports whether the action is pressed on a key or a button, without
// counting an analogue axis: flight uses this for the on/off part and
// AxisValue for the proportional part, so a bound stick is not counted twice.
func Digital(a Action) bool {
	if KeysPressed(a) {
		return true
	}
	id := firstPad()
	if id < 0 {
		return false
	}
	btn := PadButton(a)
	if btn < 0 {
		return false
	}
	return ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButton(btn))
}

// Held reports whether the action is active on keyboard, pad button or axis.
func Held(a Action) bool {
	if KeysPressed(a) {
		return true
	}
	id := firstPad()
	if id < 0 {
		return false
	}
	if btn := PadButton(a); btn >= 0 {
		if ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButton(btn)) {
			return true
		}
	}
	return AxisValue(a) > AxisThreshold
}

// PadDesc describes the bound button or axis for the UI.
func PadDesc(a Action) string {
	b, ok := Cur[a]
	if !ok {
		return "-"
	}
	switch {
	case b.Btn >= 0:
		return PadName(b.Btn)
	case b.Axis >= 0:
		return AxisName(b.Axis, b.Dir)
	}
	return "-"
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

// Capture looks for an input pressed in this frame: a key, a standard gamepad
// button, or an axis pushed past the threshold. It watches for a fresh press,
// so a key held from before the prompt does not bind itself.
func Capture() (kind string, code, dir int32, ok bool) {
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if !Bindable(k) {
			continue
		}
		if inpututil.IsKeyJustPressed(k) {
			return "key", int32(k), 0, true
		}
	}
	id := firstPad()
	if id < 0 {
		return "", 0, 0, false
	}
	for b := ebiten.StandardGamepadButton(0); b <= ebiten.StandardGamepadButtonMax; b++ {
		if inpututil.IsStandardGamepadButtonJustPressed(id, b) {
			return "pad", int32(b), 0, true
		}
	}
	for ax := int32(0); ax <= AxisTR; ax++ {
		v := ebiten.GamepadAxisValue(id, ebiten.GamepadAxisType(ax))
		if prev[ax] > -AxisThreshold && prev[ax] < AxisThreshold &&
			(v <= -AxisThreshold || v >= AxisThreshold) {
			d := int32(1)
			if v < 0 {
				d = -1
			}
			return "axis", ax, d, true
		}
	}
	return "", 0, 0, false
}

// prev remembers the last axis positions so a stick already pushed when the
// prompt opened is not captured on the first frame.
var prev [6]float64

// Poll records the axis positions; call it once per frame.
func Poll() {
	id := firstPad()
	if id < 0 {
		return
	}
	for ax := int32(0); ax <= AxisTR; ax++ {
		prev[ax] = ebiten.GamepadAxisValue(id, ebiten.GamepadAxisType(ax))
	}
}

// CaptureDown reports a bindable key that is simply held: the fallback for
// keys the platform never reports as a fresh press.
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
		Cur[Sound] = Binding{Keys: []int32{int32(ebiten.KeyF10)}, Btn: noBtn, Axis: noAxis}
	}
}

// ---------------------------------------------------------------- persistence

type keyFile struct {
	Version int                `json:"version"`
	Binds   map[string]Binding `json:"binds"`
}

// keysVersion is bumped when the binding shape changes. Version 1 stored raw
// device button indices, which mean nothing in the standard layout.
const keysVersion = 2

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
	f := keyFile{Version: keysVersion, Binds: map[string]Binding{}}
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
		b, ok := f.Binds[Name(a)]
		if !ok {
			continue
		}
		if f.Version < keysVersion {
			// old files carry keys we still trust and raw pad numbers we do
			// not: keep the keys, take the button and axis from the defaults
			d := Defaults()[a]
			b.Btn, b.Axis, b.Dir = d.Btn, d.Axis, d.Dir
		}
		Cur[a] = b
	}
	repair()
	return true
}
