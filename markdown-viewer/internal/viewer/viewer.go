package viewer

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nfang/markdown-viewer/internal/markdown"
)

// Model is the markdown viewer pane model.
type Model struct {
	icons        map[markdown.CheckboxState]string
	content      string   // raw markdown
	rendered     string   // rendered output
	lines        []string // rendered lines for scrolling
	renderer     *markdown.Renderer
	offset       int
	width        int
	height       int
	focused      bool
	filePath     string
	filterStates map[markdown.CheckboxState]bool
}

// New creates a new viewer model.
func New(icons map[markdown.CheckboxState]string) Model {
	return Model{
		icons: icons,
	}
}

// SetSize sets the viewer pane dimensions and recreates the renderer.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.rebuildRenderer()
	m.rerender()
}

// SetFocused sets whether this pane has focus.
func (m *Model) SetFocused(focused bool) {
	m.focused = focused
}

// Focused returns whether this pane has focus.
func (m Model) Focused() bool {
	return m.focused
}

// FilePath returns the currently loaded file path.
func (m Model) FilePath() string {
	return m.filePath
}

// SetFilter sets the active checkbox state filters.
func (m *Model) SetFilter(states map[markdown.CheckboxState]bool) {
	m.filterStates = states
	m.rerender()
}

// ClearFilter removes all filters.
func (m *Model) ClearFilter() {
	m.filterStates = nil
	m.rerender()
}

// LoadFile loads and renders a markdown file.
func (m *Model) LoadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	m.filePath = path
	m.content = string(data)
	m.offset = 0
	m.rerender()
	return nil
}

// LoadContent loads raw markdown content (not from a file).
func (m *Model) LoadContent(content string) {
	m.filePath = ""
	m.content = content
	m.offset = 0
	m.rerender()
}

// ScrollUp scrolls the viewer up by one line.
func (m *Model) ScrollUp() {
	if m.offset > 0 {
		m.offset--
	}
}

// ScrollDown scrolls the viewer down by one line.
func (m *Model) ScrollDown() {
	maxOffset := len(m.lines) - m.viewableHeight()
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.offset < maxOffset {
		m.offset++
	}
}

// ScrollToTop scrolls to the top.
func (m *Model) ScrollToTop() {
	m.offset = 0
}

// ScrollToBottom scrolls to the bottom.
func (m *Model) ScrollToBottom() {
	maxOffset := len(m.lines) - m.viewableHeight()
	if maxOffset < 0 {
		maxOffset = 0
	}
	m.offset = maxOffset
}

// Reload re-reads the current file from disk.
func (m *Model) Reload() error {
	if m.filePath == "" {
		return nil
	}
	return m.LoadFile(m.filePath)
}

// View renders the viewer pane.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	var borderColor lipgloss.Color
	if m.focused {
		borderColor = lipgloss.Color("39")
	} else {
		borderColor = lipgloss.Color("240")
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(m.width - 2).
		Height(m.height - 2)

	if m.content == "" {
		placeholder := lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Render("Select a file to preview")
		return style.Render(placeholder)
	}

	vh := m.viewableHeight()
	endIdx := m.offset + vh
	if endIdx > len(m.lines) {
		endIdx = len(m.lines)
	}

	var visible []string
	if m.offset < len(m.lines) {
		visible = m.lines[m.offset:endIdx]
	}

	return style.Render(strings.Join(visible, "\n"))
}

func (m *Model) rebuildRenderer() {
	if m.width <= 4 {
		return
	}
	r, err := markdown.NewRenderer(m.icons, m.width-4)
	if err != nil {
		return
	}
	m.renderer = r
}

func (m *Model) rerender() {
	if m.renderer == nil || m.content == "" {
		m.rendered = ""
		m.lines = nil
		return
	}

	var rendered string
	var err error

	if len(m.filterStates) > 0 {
		rendered, err = m.renderer.RenderWithFilter(m.content, m.filterStates)
	} else {
		rendered, err = m.renderer.Render(m.content)
	}

	if err != nil {
		m.rendered = "Error rendering: " + err.Error()
		m.lines = []string{m.rendered}
		return
	}

	m.rendered = rendered
	m.lines = strings.Split(rendered, "\n")
}

func (m Model) viewableHeight() int {
	h := m.height - 2
	if h < 1 {
		return 1
	}
	return h
}
