// Package input holds the rebindable action map. Every action can carry two
// keyboard keys, one gamepad button and one gamepad axis with a direction, so
// steering can sit on the keyboard, a stick or the triggers.
package input

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Action int

const (
	PitchUp Action = iota
	PitchDown
	RollLeft
	RollRight
	YawLeft
	YawRight
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
	"PITCH UP", "PITCH DOWN", "ROLL LEFT", "ROLL RIGHT", "YAW LEFT", "YAW RIGHT", "FIRE",
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

// Axes are the four standard stick axes. Triggers are not axes here: the raw
// SDL order differs between controllers (some pads report LT right after LY),
// so they are bound as the standard LT/RT buttons instead. A stick can be put
// on any action when an analogue value is wanted.
const (
	AxisLX  = int32(ebiten.StandardGamepadAxisLeftStickHorizontal)
	AxisLY  = int32(ebiten.StandardGamepadAxisLeftStickVertical)
	AxisRX  = int32(ebiten.StandardGamepadAxisRightStickHorizontal)
	AxisRY  = int32(ebiten.StandardGamepadAxisRightStickVertical)
	AxisMax = AxisRY
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
		PitchUp:   {Keys: k(ebiten.KeyArrowUp, ebiten.KeyW), Btn: noBtn, Axis: AxisLY, Dir: -1},
		PitchDown: {Keys: k(ebiten.KeyArrowDown, ebiten.KeyS), Btn: noBtn, Axis: AxisLY, Dir: 1},
		RollLeft:  {Keys: k(ebiten.KeyArrowLeft, ebiten.KeyA), Btn: noBtn, Axis: AxisLX, Dir: -1},
		RollRight: {Keys: k(ebiten.KeyArrowRight, ebiten.KeyD), Btn: noBtn, Axis: AxisLX, Dir: 1},
		// yaw is an Advanced-mode control; Classic ignores it, so it starts
		// unbound and the player decides what to put on it
		YawLeft:      {Keys: k(ebiten.KeyQ), Btn: noBtn, Axis: noAxis},
		YawRight:     {Keys: k(ebiten.KeyR), Btn: noBtn, Axis: noAxis},
		Fire:         {Keys: k(ebiten.KeySpace), Btn: PadA, Axis: noAxis},
		Accel:        {Keys: k(ebiten.KeyX, ebiten.KeyEqual), Btn: PadRT, Axis: noAxis},
		Decel:        {Keys: k(ebiten.KeyZ, ebiten.KeyMinus), Btn: PadLT, Axis: noAxis},
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

// SetPad binds a standard gamepad button to an action. A pad input has one
// slot: binding a button replaces an axis binding, so the old control stops
// working instead of quietly doing both.
func SetPad(a Action, btn int32) {
	steal(Binding{Btn: btn}, a)
	b := Cur[a]
	b.Btn = btn
	b.Axis, b.Dir = noAxis, 0
	Cur[a] = b
}

// SetAxis binds an axis direction to an action, replacing its button binding.
func SetAxis(a Action, axis, dir int32) {
	steal(Binding{Axis: axis, Dir: dir}, a)
	b := Cur[a]
	b.Axis, b.Dir = axis, dir
	b.Btn = noBtn
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
	if axis < 0 || axis > AxisMax || dir == 0 {
		return 0
	}
	id := firstPad()
	if id < 0 {
		return 0
	}
	v := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxis(axis)) * float64(dir)
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
	case b.Axis >= 0 && b.Axis <= AxisMax:
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

// Capture looks for an input to bind. Keys and buttons must be pressed in this
// frame; an axis only has to be held past the threshold. There is no "crossing"
// logic: a stick that is already pushed binds immediately after a short settle
// delay, which is what a player expects and what keeps the binding from firing
// later, behind their back.
func Capture() (kind string, code, dir int32, ok bool) {
	settle := captureSettling()
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
	if settle {
		return "", 0, 0, false
	}
	for ax := int32(0); ax <= AxisMax; ax++ {
		v := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxis(ax))
		if v <= -AxisThreshold {
			return "axis", ax, -1, true
		}
		if v >= AxisThreshold {
			return "axis", ax, 1, true
		}
	}
	return "", 0, 0, false
}

// AxisPick reports an axis held past the threshold, for the prompt read-out.
func AxisPick() (int32, int32, bool) {
	id := firstPad()
	if id < 0 {
		return 0, 0, false
	}
	for ax := int32(0); ax <= AxisMax; ax++ {
		v := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxis(ax))
		if v <= -AxisThreshold {
			return ax, -1, true
		}
		if v >= AxisThreshold {
			return ax, 1, true
		}
	}
	return 0, 0, false
}

// captureOpened records when the prompt appeared, so the Enter that opens it
// cannot also land as the bound key.
var captureOpened time.Time

// BeginCapture resets the prompt state.
func BeginCapture() { captureOpened = time.Now() }

// captureSettling reports whether the prompt is still too fresh to accept an
// axis that was already deflected when it opened.
func captureSettling() bool { return time.Since(captureOpened) < 180*time.Millisecond }

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

// PeekAxis names a stick or trigger currently deflected, so the rebinding
// prompt can show that the push is being seen. Triggers only count upwards:
// some drivers rest them at -1.
func PeekAxis() string {
	id := firstPad()
	if id < 0 {
		return ""
	}
	for ax := int32(0); ax <= AxisMax; ax++ {
		v := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxis(ax))
		if v <= -0.35 {
			return AxisName(ax, -1)
		}
		if v >= 0.35 {
			return AxisName(ax, 1)
		}
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
		if b.Axis >= 0 && b.Btn >= 0 {
			// a pad input is one slot; files written before that rule can
			// carry both, and the button would keep firing behind the axis
			b.Btn = noBtn
		}
		if b.Axis > AxisMax {
			// a raw axis index from an older build: not comparable to the
			// standard axes, so fall back to the default binding
			b.Axis, b.Dir = noAxis, 0
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
