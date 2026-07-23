package search

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Model is the search overlay model.
type Model struct {
	results  []Result
	cursor   int
	query    string
	visible  bool
	width    int
	height   int
}

// New creates a new search model.
func New() Model {
	return Model{}
}

// SetSize sets overlay dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Visible returns whether the overlay is shown.
func (m Model) Visible() bool {
	return m.visible
}

// Show activates the overlay with search results.
func (m *Model) Show(query string, results []Result) {
	m.query = query
	m.results = results
	m.cursor = 0
	m.visible = true
}

// Hide dismisses the overlay.
func (m *Model) Hide() {
	m.visible = false
	m.results = nil
	m.cursor = 0
	m.query = ""
}

// MoveUp moves cursor up in results.
func (m *Model) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

// MoveDown moves cursor down in results.
func (m *Model) MoveDown() {
	if m.cursor < len(m.results)-1 {
		m.cursor++
	}
}

// Selected returns the currently selected result, or nil.
func (m Model) Selected() *Result {
	if len(m.results) == 0 {
		return nil
	}
	return &m.results[m.cursor]
}

// View renders the search overlay.
func (m Model) View() string {
	if !m.visible || m.width == 0 || m.height == 0 {
		return ""
	}

	// Overlay dimensions: 80% of screen, centered
	overlayW := m.width * 80 / 100
	overlayH := m.height * 80 / 100
	if overlayW < 40 {
		overlayW = 40
	}
	if overlayH < 10 {
		overlayH = 10
	}

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Width(overlayW - 2).
		Height(overlayH - 2)

	// Header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Render(fmt.Sprintf("Search: %s (%d results)", m.query, len(m.results)))

	// Results list
	visibleResults := overlayH - 5 // header + border + padding
	if visibleResults < 1 {
		visibleResults = 1
	}

	startIdx := 0
	if m.cursor >= visibleResults {
		startIdx = m.cursor - visibleResults + 1
	}
	endIdx := startIdx + visibleResults
	if endIdx > len(m.results) {
		endIdx = len(m.results)
	}

	var lines []string
	lines = append(lines, header)
	lines = append(lines, "")

	if len(m.results) == 0 {
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Render("  No results found"))
	}

	for i := startIdx; i < endIdx; i++ {
		r := m.results[i]
		line := m.renderResult(r, i == m.cursor, overlayW-4)
		lines = append(lines, line)
	}

	lines = append(lines, "")
	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Render("  j/k:navigate  enter:open  esc:close")
	lines = append(lines, footer)

	content := strings.Join(lines, "\n")
	overlay := borderStyle.Render(content)

	// Center the overlay
	padLeft := (m.width - overlayW) / 2
	padTop := (m.height - overlayH) / 2

	var centeredLines []string
	for i := 0; i < padTop; i++ {
		centeredLines = append(centeredLines, "")
	}
	for _, line := range strings.Split(overlay, "\n") {
		centeredLines = append(centeredLines, strings.Repeat(" ", padLeft)+line)
	}

	return strings.Join(centeredLines, "\n")
}

func (m Model) renderResult(r Result, selected bool, maxWidth int) string {
	var line string

	if r.MatchType == "file" {
		icon := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render("📄")
		line = fmt.Sprintf("  %s %s", icon, r.FileName)
	} else {
		loc := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render(
			fmt.Sprintf("%s:%d", r.FileName, r.LineNum))
		content := r.LineContent
		if len(content) > maxWidth-len(r.FileName)-10 {
			content = content[:maxWidth-len(r.FileName)-13] + "..."
		}
		line = fmt.Sprintf("  %s  %s", loc, content)
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
