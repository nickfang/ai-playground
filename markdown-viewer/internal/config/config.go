package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Theme        string            `yaml:"theme"`
	CustomTheme  map[string]string `yaml:"custom_theme"`
	Keybindings  string            `yaml:"keybindings"`
	Sort         string            `yaml:"sort"`
	SidebarWidth int               `yaml:"sidebar_width"`
}

func DefaultConfig() Config {
	return Config{
		Theme:       "emoji",
		CustomTheme: nil,
		Keybindings: "vim",
		Sort:        "alpha",
		SidebarWidth: 30,
	}
}

func Load() Config {
	cfg := DefaultConfig()

	home, err := os.UserHomeDir()
	if err != nil {
		return cfg
	}

	path := filepath.Join(home, ".mdview.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}

	_ = yaml.Unmarshal(data, &cfg)

	// Clamp sidebar width
	if cfg.SidebarWidth < 10 {
		cfg.SidebarWidth = 10
	}
	if cfg.SidebarWidth > 50 {
		cfg.SidebarWidth = 50
	}

	// Validate theme
	switch cfg.Theme {
	case "emoji", "symbols", "custom":
	default:
		cfg.Theme = "emoji"
	}

	// Validate keybindings
	switch cfg.Keybindings {
	case "vim", "arrows":
	default:
		cfg.Keybindings = "vim"
	}

	// Validate sort
	switch cfg.Sort {
	case "alpha", "modified":
	default:
		cfg.Sort = "alpha"
	}

	return cfg
}
