package input

import "github.com/hajimehoshi/ebiten/v2"

var keyNames = map[int32]string{}

func init() {
	for c := ebiten.KeyA; c <= ebiten.KeyZ; c++ {
		keyNames[int32(c)] = string(rune('A' + int(c-ebiten.KeyA)))
	}
	for c := ebiten.KeyF1; c <= ebiten.KeyF12; c++ {
		keyNames[int32(c)] = "F" + string(rune('1'+int(c-ebiten.KeyF1)))
	}
	for c := ebiten.Key0; c <= ebiten.Key9; c++ {
		keyNames[int32(c)] = string(rune('0' + int(c-ebiten.Key0)))
	}
	special := []struct {
		k ebiten.Key
		n string
	}{
		{ebiten.KeySpace, "SPACE"},
		{ebiten.KeyEnter, "ENTER"},
		{ebiten.KeyEscape, "ESC"},
		{ebiten.KeyTab, "TAB"},
		{ebiten.KeyBackspace, "BKSP"},
		{ebiten.KeyArrowUp, "UP"},
		{ebiten.KeyArrowDown, "DOWN"},
		{ebiten.KeyArrowLeft, "LEFT"},
		{ebiten.KeyArrowRight, "RIGHT"},
		{ebiten.KeyShiftLeft, "LSHIFT"},
		{ebiten.KeyShiftRight, "RSHIFT"},
		{ebiten.KeyControlLeft, "LCTRL"},
		{ebiten.KeyControlRight, "RCTRL"},
		{ebiten.KeyAltLeft, "LALT"},
		{ebiten.KeyMinus, "-"},
		{ebiten.KeyEqual, "="},
		{ebiten.KeyBracketLeft, "["},
		{ebiten.KeyBracketRight, "]"},
		{ebiten.KeySemicolon, ";"},
		{ebiten.KeyApostrophe, "'"},
		{ebiten.KeyComma, ","},
		{ebiten.KeyPeriod, "."},
		{ebiten.KeySlash, "/"},
		{ebiten.KeyBackslash, "\\"},
		{ebiten.KeyPageUp, "PGUP"},
		{ebiten.KeyPageDown, "PGDN"},
		{ebiten.KeyHome, "HOME"},
		{ebiten.KeyEnd, "END"},
		{ebiten.KeyF1, "F1"},
		{ebiten.KeyF2, "F2"},
		{ebiten.KeyF3, "F3"},
		{ebiten.KeyF4, "F4"},
		{ebiten.KeyF5, "F5"},
		{ebiten.KeyF6, "F6"},
	}
	for _, s := range special {
		keyNames[int32(s.k)] = s.n
	}
}

// KeyName renders a key code for the controls screen.
func KeyName(code int32) string {
	if n, ok := keyNames[code]; ok {
		return n
	}
	return "?"
}

// PadName renders a standard gamepad button.
func PadName(b int32) string {
	switch ebiten.GamepadButton(b) {
	case ebiten.GamepadButton0:
		return "A"
	case ebiten.GamepadButton1:
		return "B"
	case ebiten.GamepadButton2:
		return "X"
	case ebiten.GamepadButton3:
		return "Y"
	case ebiten.GamepadButton4:
		return "LB"
	case ebiten.GamepadButton5:
		return "RB"
	case ebiten.GamepadButton6:
		return "LT"
	case ebiten.GamepadButton7:
		return "RT"
	case ebiten.GamepadButton8:
		return "BACK"
	case ebiten.GamepadButton9:
		return "START"
	case ebiten.GamepadButton10:
		return "LS"
	case ebiten.GamepadButton11:
		return "RS"
	case ebiten.GamepadButton12:
		return "D-UP"
	case ebiten.GamepadButton13:
		return "D-DOWN"
	case ebiten.GamepadButton14:
		return "D-LEFT"
	case ebiten.GamepadButton15:
		return "D-RIGHT"
	case ebiten.GamepadButton16:
		return "HOME"
	}
	return "-"
}

// AxisName renders an axis and its direction for the controls screen.
func AxisName(axis, dir int32) string {
	switch axis {
	case 0:
		if dir < 0 {
			return "LSX-"
		}
		return "LSX+"
	case 1:
		if dir < 0 {
			return "LSY-"
		}
		return "LSY+"
	case 2:
		if dir < 0 {
			return "RSX-"
		}
		return "RSX+"
	case 3:
		if dir < 0 {
			return "RSY-"
		}
		return "RSY+"
	case 4:
		return "LT"
	case 5:
		return "RT"
	}
	return "AXIS"
}
