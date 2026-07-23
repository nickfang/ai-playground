package theme

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/nfang/markdown-viewer/internal/markdown"
)

// Theme maps checkbox states to their rendered string representations.
type Theme struct {
	Icons map[markdown.CheckboxState]string
}

// EmojiTheme renders checkboxes as emoji icons.
func EmojiTheme() Theme {
	return Theme{
		Icons: map[markdown.CheckboxState]string{
			markdown.StateTodo:       "⬜",
			markdown.StateDone:       "✅",
			markdown.StateInProgress: "🔄",
			markdown.StateCancelled:  "❌",
			markdown.StateUrgent:     "🔥",
		},
	}
}

// SymbolTheme renders checkboxes as colored text symbols.
func SymbolTheme() Theme {
	gray := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	blue := lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Strikethrough(true)
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)

	return Theme{
		Icons: map[markdown.CheckboxState]string{
			markdown.StateTodo:       gray.Render("□"),
			markdown.StateDone:       green.Render("■"),
			markdown.StateInProgress: blue.Render("▶"),
			markdown.StateCancelled:  dim.Render("–"),
			markdown.StateUrgent:     red.Render("!"),
		},
	}
}

// CustomTheme creates a theme from user-defined glyph mappings.
func CustomTheme(glyphs map[string]string) Theme {
	t := EmojiTheme() // fallback to emoji for any missing states

	if v, ok := glyphs["todo"]; ok {
		t.Icons[markdown.StateTodo] = v
	}
	if v, ok := glyphs["done"]; ok {
		t.Icons[markdown.StateDone] = v
	}
	if v, ok := glyphs["inprogress"]; ok {
		t.Icons[markdown.StateInProgress] = v
	}
	if v, ok := glyphs["cancelled"]; ok {
		t.Icons[markdown.StateCancelled] = v
	}
	if v, ok := glyphs["urgent"]; ok {
		t.Icons[markdown.StateUrgent] = v
	}

	return t
}

// ForConfig returns the appropriate theme based on config values.
func ForConfig(themeName string, customGlyphs map[string]string) Theme {
	switch themeName {
	case "symbols":
		return SymbolTheme()
	case "custom":
		return CustomTheme(customGlyphs)
	default:
		return EmojiTheme()
	}
}

// Render returns the icon string for a given checkbox state.
func (t Theme) Render(state markdown.CheckboxState) string {
	if icon, ok := t.Icons[state]; ok {
		return icon
	}
	return "?"
}
