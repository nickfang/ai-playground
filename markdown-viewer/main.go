package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nfang/markdown-viewer/internal/app"
	"github.com/nfang/markdown-viewer/internal/config"
	"github.com/nfang/markdown-viewer/internal/theme"
)

func main() {
	cfg := config.Load()
	t := theme.ForConfig(cfg.Theme, cfg.CustomTheme)

	// Determine root path
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	p := tea.NewProgram(
		app.New(root, cfg, t),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
