// Package lang holds the UI translations and the player's preferences.
package lang

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var current = "en"

// Config is everything the settings screen can change.
type Config struct {
	Lang   string  `json:"lang"`
	Muted  bool    `json:"muted"`
	Volume float64 `json:"volume"` // 0..1
}

var Cfg = Config{Lang: "en", Volume: 0.8}

// ReadOnly suppresses writing the config file (used by tests).
var ReadOnly bool

func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "heliosdrift-config.json"
	}
	return filepath.Join(home, ".heliosdrift-config.json")
}

// Load reads the saved preferences (best effort).
func Load() {
	raw, err := os.ReadFile(configPath())
	if err != nil {
		return
	}
	var c Config
	if json.Unmarshal(raw, &c) == nil {
		if c.Lang == "ru" || c.Lang == "en" {
			Cfg.Lang = c.Lang
		}
		Cfg.Muted = c.Muted
		if c.Volume > 0 || c.Volume == 0 {
			Cfg.Volume = c.Volume
		}
		current = Cfg.Lang
	}
}

// Save writes the preferences back.
func Save() bool {
	if ReadOnly {
		return true
	}
	Cfg.Lang = current
	raw, err := json.MarshalIndent(Cfg, "", "  ")
	if err != nil {
		return false
	}
	return os.WriteFile(configPath(), raw, 0o644) == nil
}

// Set switches the language ("en" or "ru").
// Set switches the language in memory. Use SetSaved to persist it.
func Set(l string) {
	if l != "ru" && l != "en" {
		return
	}
	current = l
	Cfg.Lang = l
}

// SetSaved switches the language and writes the config.
func SetSaved(l string) {
	Set(l)
	Save()
}

func Current() string { return current }

// IsRU reports whether Russian is active.
func IsRU() bool { return current == "ru" }

// T translates a UI string. English strings are the keys; anything not in the
// table falls through unchanged, so a missing translation degrades to English.
func T(s string) string {
	if current != "ru" {
		return s
	}
	if v, ok := ru[s]; ok {
		return v
	}
	return s
}

// N translates a proper noun (ship or commodity name) if it has an entry.
func N(s string) string { return T(s) }

// TranslateLine localises a whole line: first the exact phrase, then each word
// so composed messages ("ARRIVED: LAVE") come out translated too.
// TranslateLine returns the Russian text for a UI string.
//
// Only whole strings and whole leading phrases are translated, never single
// words: data labels, units and names stay in English ("HOLD 12/20 T",
// "STAGE 3/8", "KESTREL III"), while sentences and briefings read in Russian.
// A phrase match must end at a separator, so "ON" never eats the start of
// "ONLY".
func TranslateLine(s string) string {
	if current != "ru" || s == "" {
		return s
	}
	dumpOnce(s)
	if v, ok := ru[s]; ok {
		return v
	}
	// menu rows are drawn as "3  SAVE COMMANDER": translate the text, keep the
	// number and its padding
	if i := leadingNumberLen(s); i > 0 && i < len(s) {
		return s[:i] + TranslateLine(s[i:])
	}
	// same for a leading run of spaces
	if i := leadingSpaceLen(s); i > 0 && i < len(s) {
		return s[:i] + TranslateLine(s[i:])
	}
	best, bestV := "", ""
	for k, v := range ru {
		if len(k) <= 1 || len(k) >= len(s) || !strings.HasPrefix(s, k) {
			continue
		}
		if !isSeparator(s[len(k)]) || len(k) <= len(best) {
			continue
		}
		best, bestV = k, v
	}
	if best != "" {
		return bestV + TranslateLine(s[len(best):])
	}
	return s
}

// leadingNumberLen returns the length of a leading "12  " style prefix, or 0.
func leadingNumberLen(s string) int {
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	if i == 0 {
		return 0
	}
	for i < len(s) && s[i] == ' ' {
		i++
	}
	return i
}

// leadingSpaceLen returns the length of a leading run of spaces.
func leadingSpaceLen(s string) int {
	i := 0
	for i < len(s) && s[i] == ' ' {
		i++
	}
	return i
}

// dumpOnce records every string the UI asks to translate when
// HELIOS_DUMP_TEXT names a file. It exists so a translation pass can be checked
// against what the game actually draws instead of guesswork.
var (
	dumpMu   sync.Mutex
	dumpSeen = map[string]bool{}
)

func dumpOnce(s string) {
	path := os.Getenv("HELIOS_DUMP_TEXT")
	if path == "" {
		return
	}
	dumpMu.Lock()
	defer dumpMu.Unlock()
	if dumpSeen[s] {
		return
	}
	dumpSeen[s] = true
	if f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
		fmt.Fprintf(f, "%s\n", s)
		f.Close()
	}
}

// isSeparator reports whether a phrase may end here: messages are built as
// "PREFIX: data" or "PREFIX  +1200 CR", so a phrase never ends mid-word.
func isSeparator(c byte) bool {
	switch c {
	case ' ', ':', '+', '-', '.', '(', '/', '\t':
		return true
	}
	return c >= '0' && c <= '9'
}

func isWordRune(r rune) bool {
	switch {
	case r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z', r >= '0' && r <= '9':
		return true
	case r == '-' || r == '.' || r == '\'':
		return true
	}
	return false
}
