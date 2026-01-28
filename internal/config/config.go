package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configDir = ".config/briks"
const configFile = "config.json"

// Config stores all persistent settings.
type Config struct {
	Theme        string `json:"theme"`
	StartLevel   int    `json:"start_level"`
	GhostPiece   bool   `json:"ghost_piece"`
	ShowGrid     bool   `json:"show_grid"`
	PreviewCount int    `json:"preview_count"`
	DAS          int    `json:"das"` // Delayed Auto Shift in ms
	ARR          int    `json:"arr"` // Auto Repeat Rate in ms
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Theme:        "default",
		StartLevel:   1,
		GhostPiece:   true,
		ShowGrid:     false,
		PreviewCount: 5,
		DAS:          170,
		ARR:          50,
	}
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir, configFile), nil
}

// Load reads configuration from disk, falling back to defaults.
func Load() *Config {
	path, err := configPath()
	if err != nil {
		return DefaultConfig()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfig()
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return DefaultConfig()
	}

	cfg.validate()
	return cfg
}

// Save writes the configuration to disk.
func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (c *Config) validate() {
	if c.StartLevel < 1 {
		c.StartLevel = 1
	}
	if c.StartLevel > 20 {
		c.StartLevel = 20
	}
	if c.PreviewCount < 1 {
		c.PreviewCount = 1
	}
	if c.PreviewCount > 5 {
		c.PreviewCount = 5
	}
	if c.DAS < 50 {
		c.DAS = 50
	}
	if c.DAS > 500 {
		c.DAS = 500
	}
	if c.ARR < 0 {
		c.ARR = 0
	}
	if c.ARR > 200 {
		c.ARR = 200
	}
}
