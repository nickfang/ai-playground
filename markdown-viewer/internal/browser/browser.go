package browser

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Model is the file browser pane model.
type Model struct {
	entries   []Entry
	cursor    int
	dir       string
	width     int
	height    int
	sortMode  string
	focused   bool
}

// New creates a new file browser model.
func New(dir string, sortMode string) Model {
	absDir, _ := filepath.Abs(dir)
	entries, _ := LoadEntries(absDir, sortMode)
	return Model{
		entries:  entries,
		cursor:   0,
		dir:      absDir,
		sortMode: sortMode,
	}
}

// SetSize sets the browser pane dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// SetFocused sets whether this pane has focus.
func (m *Model) SetFocused(focused bool) {
	m.focused = focused
}

// Focused returns whether this pane has focus.
func (m Model) Focused() bool {
	return m.focused
}

// Dir returns the current directory.
func (m Model) Dir() string {
	return m.dir
}

// SelectedEntry returns the currently selected entry, or nil if empty.
func (m Model) SelectedEntry() *Entry {
	if len(m.entries) == 0 {
		return nil
	}
	return &m.entries[m.cursor]
}

// MoveUp moves the cursor up.
func (m *Model) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

// MoveDown moves the cursor down.
func (m *Model) MoveDown() {
	if m.cursor < len(m.entries)-1 {
		m.cursor++
	}
}

// EnterDir enters the selected directory.
// Returns true if a directory was entered.
func (m *Model) EnterDir() bool {
	entry := m.SelectedEntry()
	if entry == nil || !entry.IsDir {
		return false
	}

	m.dir = entry.Path
	m.Refresh()
	m.cursor = 0
	return true
}

// GoUp goes to the parent directory.
// Returns true if the directory changed.
func (m *Model) GoUp() bool {
	parent := filepath.Dir(m.dir)
	if parent == m.dir {
		return false
	}

	oldDir := filepath.Base(m.dir)
	m.dir = parent
	m.Refresh()

	// Try to select the directory we came from
	for i, e := range m.entries {
		if e.IsDir && e.Name == oldDir {
			m.cursor = i
			return true
		}
	}
	m.cursor = 0
	return true
}

// Refresh reloads the entries for the current directory.
func (m *Model) Refresh() {
	entries, _ := LoadEntries(m.dir, m.sortMode)
	m.entries = entries
	if m.cursor >= len(m.entries) {
		m.cursor = max(0, len(m.entries)-1)
	}
}

// ToggleSort toggles between alpha and modified sort.
func (m *Model) ToggleSort() {
	if m.sortMode == "alpha" {
		m.sortMode = "modified"
	} else {
		m.sortMode = "alpha"
	}
	m.Refresh()
}

// View renders the file browser pane.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	var borderColor lipgloss.Color
	if m.focused {
		borderColor = lipgloss.Color("39") // blue when focused
	} else {
		borderColor = lipgloss.Color("240") // dim when not focused
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(m.width - 2).
		Height(m.height - 2)

	if len(m.entries) == 0 {
		content := lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Render("No markdown files found")
		return style.Render(content)
	}

	// Calculate visible range (scrolling)
	visibleHeight := m.height - 2 // account for border
	startIdx := 0
	if m.cursor >= visibleHeight {
		startIdx = m.cursor - visibleHeight + 1
	}
	endIdx := startIdx + visibleHeight
	if endIdx > len(m.entries) {
		endIdx = len(m.entries)
	}

	var lines []string
	contentWidth := m.width - 4 // border + padding

	for i := startIdx; i < endIdx; i++ {
		e := m.entries[i]
		line := m.renderEntry(e, i == m.cursor, contentWidth)
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")
	return style.Render(content)
}

func (m Model) renderEntry(e Entry, selected bool, maxWidth int) string {
	var prefix string
	if e.IsDir {
		prefix = "📁 "
	} else {
		prefix = "   "
	}

	name := e.Name
	timeStr := ""
	if !e.IsDir {
		timeStr = e.RelativeTime()
	}

	// Calculate available space for name + heading
	// Format: prefix + name + "  " + heading + "  " + time
	timeWidth := len(timeStr)
	nameArea := maxWidth - len(prefix) - timeWidth - 4
	if nameArea < 10 {
		nameArea = 10
	}

	display := name
	if !e.IsDir && e.Heading != "" {
		heading := e.Heading
		remainingSpace := nameArea - len(name) - 2
		if remainingSpace > 3 {
			if len(heading) > remainingSpace {
				heading = heading[:remainingSpace-2] + ".."
			}
			display = fmt.Sprintf("%s  %s", name, lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render(heading))
		}
	}

	line := prefix + display
	if timeStr != "" {
		// Pad to right-align time
		padLen := maxWidth - lipgloss.Width(line) - len(timeStr)
		if padLen < 1 {
			padLen = 1
		}
		line += strings.Repeat(" ", padLen) + lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render(timeStr)
	}

	if selected {
		return lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Bold(true).
			Width(maxWidth).
			Render(line)
	}

	return line
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
