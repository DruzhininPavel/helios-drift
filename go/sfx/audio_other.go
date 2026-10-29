//go:build !linux

package sfx

// AudioAvailable is always true off Linux: CoreAudio and WASAPI are part of
// the OS, so there is nothing to probe.
func AudioAvailable() bool { return true }
