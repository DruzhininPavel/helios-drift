// Package lang holds the UI translations and the player's preferences.
package lang

import (
	"encoding/json"
	"os"
	"path/filepath"
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
		// keep preferences written before the rename
		if old, e1 := os.ReadFile(filepath.Join(filepath.Dir(configPath()), ".elite16-config.json")); e1 == nil {
			raw = old
		} else {
			return
		}
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
func TranslateLine(s string) string {
	if current != "ru" || s == "" {
		return s
	}
	if v, ok := ru[s]; ok {
		return v
	}
	var b []rune
	runes := []rune(s)
	i := 0
	for i < len(runes) {
		j := i
		for j < len(runes) && !isWordRune(runes[j]) {
			j++
		}
		b = append(b, runes[i:j]...)
		if j >= len(runes) {
			break
		}
		k := j
		for k < len(runes) && isWordRune(runes[k]) {
			k++
		}
		word := string(runes[j:k])
		if v, ok := ru[word]; ok {
			b = append(b, []rune(v)...)
		} else {
			b = append(b, runes[j:k]...)
		}
		i = k
	}
	return string(b)
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
