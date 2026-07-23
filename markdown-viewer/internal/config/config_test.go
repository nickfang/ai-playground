package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Theme != "emoji" {
		t.Errorf("expected theme 'emoji', got %q", cfg.Theme)
	}
	if cfg.Keybindings != "vim" {
		t.Errorf("expected keybindings 'vim', got %q", cfg.Keybindings)
	}
	if cfg.Sort != "alpha" {
		t.Errorf("expected sort 'alpha', got %q", cfg.Sort)
	}
	if cfg.SidebarWidth != 30 {
		t.Errorf("expected sidebar_width 30, got %d", cfg.SidebarWidth)
	}
	if cfg.CustomTheme != nil {
		t.Errorf("expected nil custom_theme, got %v", cfg.CustomTheme)
	}
}

func TestLoad_NoFile(t *testing.T) {
	// Load should return defaults when no config file exists
	cfg := Load()
	def := DefaultConfig()

	if cfg.Theme != def.Theme {
		t.Errorf("expected default theme %q, got %q", def.Theme, cfg.Theme)
	}
	if cfg.SidebarWidth != def.SidebarWidth {
		t.Errorf("expected default sidebar_width %d, got %d", def.SidebarWidth, cfg.SidebarWidth)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	// Create a temp config file
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	configContent := `theme: symbols
keybindings: arrows
sort: modified
sidebar_width: 40
custom_theme:
  todo: "○"
  done: "●"
`
	err := os.WriteFile(filepath.Join(tmpDir, ".mdview.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	cfg := Load()

	if cfg.Theme != "symbols" {
		t.Errorf("expected theme 'symbols', got %q", cfg.Theme)
	}
	if cfg.Keybindings != "arrows" {
		t.Errorf("expected keybindings 'arrows', got %q", cfg.Keybindings)
	}
	if cfg.Sort != "modified" {
		t.Errorf("expected sort 'modified', got %q", cfg.Sort)
	}
	if cfg.SidebarWidth != 40 {
		t.Errorf("expected sidebar_width 40, got %d", cfg.SidebarWidth)
	}
	if cfg.CustomTheme["todo"] != "○" {
		t.Errorf("expected custom_theme todo '○', got %q", cfg.CustomTheme["todo"])
	}
}

func TestLoad_ClampsWidth(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Too small
	os.WriteFile(filepath.Join(tmpDir, ".mdview.yaml"), []byte("sidebar_width: 5"), 0644)
	cfg := Load()
	if cfg.SidebarWidth != 10 {
		t.Errorf("expected clamped sidebar_width 10, got %d", cfg.SidebarWidth)
	}

	// Too large
	os.WriteFile(filepath.Join(tmpDir, ".mdview.yaml"), []byte("sidebar_width: 80"), 0644)
	cfg = Load()
	if cfg.SidebarWidth != 50 {
		t.Errorf("expected clamped sidebar_width 50, got %d", cfg.SidebarWidth)
	}
}

func TestLoad_InvalidValues(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	configContent := `theme: invalid_theme
keybindings: invalid_kb
sort: invalid_sort
`
	os.WriteFile(filepath.Join(tmpDir, ".mdview.yaml"), []byte(configContent), 0644)

	cfg := Load()

	if cfg.Theme != "emoji" {
		t.Errorf("expected fallback theme 'emoji', got %q", cfg.Theme)
	}
	if cfg.Keybindings != "vim" {
		t.Errorf("expected fallback keybindings 'vim', got %q", cfg.Keybindings)
	}
	if cfg.Sort != "alpha" {
		t.Errorf("expected fallback sort 'alpha', got %q", cfg.Sort)
	}
}
