package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

// Slot is where this commander belongs: the slot he was loaded from or last
// saved to. Quick save writes there.
func (g *Game) Slot() int {
	if g.Cmd.Slot < 0 || g.Cmd.Slot >= Slots {
		return -1
	}
	return g.Cmd.Slot
}

// QuickSave writes to the slot the commander came from, or to the first one if
// he has never been saved.
func (g *Game) QuickSave() (int, bool) {
	slot := g.Slot()
	if slot < 0 {
		slot = 0
	}
	return slot, g.SaveSlot(slot)
}

// QuickLoad reads the slot the commander belongs to.
func (g *Game) QuickLoad() (int, bool) {
	slot := g.Slot()
	if slot < 0 || !SlotExists(slot) {
		return slot, false
	}
	return slot, g.LoadSlot(slot)
}

func (g *Game) SaveSlot(slot int) bool {
	g.Cmd.Slot = slot
	data := saveData{Version: 2, Cmd: g.Cmd}
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return false
	}
	return os.WriteFile(savePath(slot), raw, 0o644) == nil
}

// LoadSlot restores a commander and docks at their last station.
func (g *Game) LoadSlot(slot int) bool {
	raw, err := os.ReadFile(savePath(slot))
	if err != nil {
		return false
	}
	var data saveData
	if err := json.Unmarshal(raw, &data); err != nil {
		return false
	}
	g.Cmd = data.Cmd
	g.Cmd.Slot = slot
	if g.Cmd.Cargo == nil {
		g.Cmd.Cargo = map[int]int{}
	}
	if g.Cmd.Equip == nil {
		g.Cmd.Equip = map[string]bool{}
	}
	if g.Cmd.Flags == nil {
		g.Cmd.Flags = map[string]bool{}
	}
	if g.Cmd.Visited == nil {
		g.Cmd.Visited = map[string]bool{}
	}
	if g.Cmd.SideStep == nil {
		g.Cmd.SideStep = map[string]int{}
	}
	if g.Cmd.SideCount == nil {
		g.Cmd.SideCount = map[string]int{}
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
	_, err := os.Stat(savePath(slot))
	return err == nil
}

// SlotInfo describes a slot for the save/load screen.
func SlotInfo(slot int) (name string, cash float64, where string, ok bool) {
	raw, err := os.ReadFile(savePath(slot))
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

// Save writes the first slot.
func (g *Game) Save() bool { return g.SaveSlot(0) }

// Load reads the first slot.
func (g *Game) Load() bool { return g.LoadSlot(0) }

// DeleteSlot removes a commander from a slot. The file is moved aside rather than
// destroyed, so a wrong press is recoverable until the next save writes over it.
func DeleteSlot(slot int) bool {
	p := savePath(slot)
	if !SlotExists(slot) {
		return false
	}
	return os.Rename(p, p+".deleted") == nil
}
