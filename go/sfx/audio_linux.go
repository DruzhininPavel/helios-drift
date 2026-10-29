//go:build linux

package sfx

import (
	"os"
	"path/filepath"
	"strconv"
)

// AudioAvailable reports whether the host has a sound server (PulseAudio or
// PipeWire's pulse layer) or an ALSA device. Ebiten treats a missing backend
// as fatal, so the game checks first and runs silent instead of dying.
func AudioAvailable() bool {
	if os.Getenv("PULSE_SERVER") != "" {
		return true
	}
	uid := strconv.Itoa(os.Getuid())
	paths := []string{
		"/run/user/" + uid + "/pulse/native",
		"/dev/snd",
	}
	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
		paths = append(paths, filepath.Join(dir, "pulse/native"))
	}
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".pulse/native"))
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}
