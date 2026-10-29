package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"heliosdrift/galaxy"
)

// Slots is the number of save files a commander can keep.
const Slots = 3

type saveData struct {
	Version int       `json:"version"`
	Cmd     Commander `json:"cmd"`
}

func home() string {
	h, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return h
}

// SaveDirOverride lets tests redirect saves away from the real home directory.
var SaveDirOverride string

func savePath(slot int) string {
	if slot < 0 {
		slot = 0
	}
	if slot >= Slots {
		slot = Slots - 1
	}
	dir := home()
	if SaveDirOverride != "" {
		dir = SaveDirOverride
	}
	return filepath.Join(dir, fmt.Sprintf(".heliosdrift-save%d.json", slot+1))
}

// legacyPath is the single-slot file used by earlier builds.
func legacyPath() string {
	dir := home()
	if SaveDirOverride != "" {
		dir = SaveDirOverride
	}
	return filepath.Join(dir, ".heliosdrift-save.json")
}

// oldPath maps a new save file back to the name used before the rename, so a
// commander saved by an older build still loads.
func oldPath(p string) string {
	dir, base := filepath.Split(p)
	switch {
	case strings.HasPrefix(base, ".heliosdrift-save"):
		return filepath.Join(dir, strings.Replace(base, ".heliosdrift-", ".elite16-", 1))
	}
	return ""
}

// readSave loads p, falling back to the pre-rename file name.
func readSave(p string) ([]byte, error) {
	raw, err := os.ReadFile(p)
	if err == nil {
		return raw, nil
	}
	if old := oldPath(p); old != "" {
		return os.ReadFile(old)
	}
	return nil, err
}

// SaveSlot writes the commander into a slot.
func (g *Game) SaveSlot(slot int) bool {
	data := saveData{Version: 2, Cmd: g.Cmd}
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return false
	}
	return os.WriteFile(savePath(slot), raw, 0o644) == nil
}

// LoadSlot restores a commander and docks at their last station.
func (g *Game) LoadSlot(slot int) bool {
	raw, err := readSave(savePath(slot))
	if err != nil && slot == 0 {
		raw, err = readSave(legacyPath())
	}
	if err != nil {
		return false
	}
	var data saveData
	if err := json.Unmarshal(raw, &data); err != nil {
		return false
	}
	g.Cmd = data.Cmd
	if g.Cmd.Cargo == nil {
		g.Cmd.Cargo = map[int]int{}
	}
	if g.Cmd.Equip == nil {
		g.Cmd.Equip = map[string]bool{}
	}
	if g.Cmd.Ship == "" {
		g.Cmd.Ship = "cobra3"
	}
	g.Cmd.CargoCap = g.CargoCapacity()
	g.EscapeUsed = false
	g.EnterSystem(g.Cmd.System, "arrive", 9000, false)
	g.State = "docked"
	g.dockedAt()
	g.Message("COMMANDER RESTORED", "good")
	return true
}

// SlotExists reports whether a slot holds a commander.
func SlotExists(slot int) bool {
	if _, err := os.Stat(savePath(slot)); err == nil {
		return true
	}
	if slot == 0 {
		_, err := os.Stat(legacyPath())
		return err == nil
	}
	return false
}

// SlotInfo describes a slot for the save/load screen.
func SlotInfo(slot int) (name string, cash float64, where string, ok bool) {
	raw, err := readSave(savePath(slot))
	if err != nil && slot == 0 {
		raw, err = readSave(legacyPath())
	}
	if err != nil {
		return "", 0, "", false
	}
	var data saveData
	if err := json.Unmarshal(raw, &data); err != nil {
		return "", 0, "", false
	}
	sys := "?"
	if gal := galaxy.Get(data.Cmd.Galaxy); gal != nil && data.Cmd.System >= 0 && data.Cmd.System < len(gal.Systems) {
		sys = gal.Systems[data.Cmd.System].Name
	}
	name = data.Cmd.Name
	if name == "" {
		name = "JAMESON"
	}
	return name, data.Cmd.Cash, sys, true
}

// SaveExists reports whether any slot holds a commander.
func SaveExists() bool {
	for i := 0; i < Slots; i++ {
		if SlotExists(i) {
			return true
		}
	}
	return false
}

// Save keeps the old single-slot call sites working.
func (g *Game) Save() bool { return g.SaveSlot(0) }

// Load keeps the old single-slot call sites working.
func (g *Game) Load() bool { return g.LoadSlot(0) }
