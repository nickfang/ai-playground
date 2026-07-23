package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nfang/markdown-viewer/internal/browser"
	"github.com/nfang/markdown-viewer/internal/config"
	"github.com/nfang/markdown-viewer/internal/markdown"
	"github.com/nfang/markdown-viewer/internal/search"
	"github.com/nfang/markdown-viewer/internal/theme"
	"github.com/nfang/markdown-viewer/internal/viewer"
	"github.com/nfang/markdown-viewer/internal/watcher"
)

const (
	FocusBrowser = iota
	FocusViewer
)

const (
	ModeNormal = iota
	ModeCommand
	ModeSearch
)

type Model struct {
	browser      browser.Model
	viewer       viewer.Model
	search       search.Model
	watcher      *watcher.Watcher
	cfg          config.Config
	theme        theme.Theme
	focus        int
	mode         int
	commandInput string
	width        int
	height       int
	statusMsg    string
	filterNames  []string
}

func New(root string, cfg config.Config, t theme.Theme) Model {
	// Resolve root path
	absRoot, err := filepath.Abs(root)
	if err != nil {
		absRoot = root
	}

	// If root is a file, open its directory with the file selected
	info, err := os.Stat(absRoot)
	if err == nil && !info.IsDir() {
		absRoot = filepath.Dir(absRoot)
	}

	b := browser.New(absRoot, cfg.Sort)
	v := viewer.New(t.Icons)
	s := search.New()
	w, _ := watcher.New()

	b.SetFocused(true)

	return Model{
		browser: b,
		viewer:  v,
		search:  s,
		watcher: w,
		cfg:     cfg,
		theme:   t,
		focus:   FocusBrowser,
		mode:    ModeNormal,
	}
}

func (m Model) Init() tea.Cmd {
	if m.watcher != nil {
		return m.watcher.ListenCmd()
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updatePaneSizes()
		m.search.SetSize(m.width, m.height)
		m.autoPreview()
		return m, nil

	case watcher.FileChangedMsg:
		// Auto-reload on file change
		m.browser.Refresh()
		_ = m.viewer.Reload()
		// Continue watching
		if m.watcher != nil {
			return m, m.watcher.ListenCmd()
		}
		return m, nil

	case tea.KeyMsg:
		// Search overlay takes priority
		if m.mode == ModeSearch {
			return m.handleSearchInput(msg)
		}
		if m.mode == ModeCommand {
			return m.handleCommandInput(msg)
		}
		return m.handleNormalInput(msg)
	}

	return m, nil
}

func (m Model) handleNormalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Clear status message on any key
	m.statusMsg = ""

	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "tab":
		m.toggleFocus()
		return m, nil

	case ":":
		m.mode = ModeCommand
		m.commandInput = ""
		return m, nil

	case "+", "=":
		if m.cfg.SidebarWidth < 50 {
			m.cfg.SidebarWidth += 5
			m.updatePaneSizes()
		}
		return m, nil

	case "-":
		if m.cfg.SidebarWidth > 10 {
			m.cfg.SidebarWidth -= 5
			m.updatePaneSizes()
		}
		return m, nil

	case "r":
		m.browser.Refresh()
		_ = m.viewer.Reload()
		m.statusMsg = "Refreshed"
		return m, nil
	}

	if m.focus == FocusBrowser {
		return m.handleBrowserInput(msg)
	}
	return m.handleViewerInput(msg)
}

func (m Model) handleBrowserInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		m.browser.MoveDown()
		m.autoPreview()
	case "k", "up":
		m.browser.MoveUp()
		m.autoPreview()
	case "enter", "l", "right":
		m.browser.EnterDir()
		m.autoPreview()
	case "backspace", "h", "left":
		m.browser.GoUp()
		m.autoPreview()
	case "s":
		m.browser.ToggleSort()
		m.autoPreview()
	}
	return m, nil
}

func (m Model) handleViewerInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		m.viewer.ScrollDown()
	case "k", "up":
		m.viewer.ScrollUp()
	case "g":
		m.viewer.ScrollToTop()
	case "G":
		m.viewer.ScrollToBottom()
	}
	return m, nil
}

func (m Model) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.search.Hide()
		m.mode = ModeNormal
		return m, nil

	case "j", "down":
		m.search.MoveDown()
		return m, nil

	case "k", "up":
		m.search.MoveUp()
		return m, nil

	case "enter":
		result := m.search.Selected()
		if result != nil {
			_ = m.viewer.LoadFile(result.FilePath)
			// Switch focus to viewer
			m.focus = FocusViewer
			m.browser.SetFocused(false)
			m.viewer.SetFocused(true)
		}
		m.search.Hide()
		m.mode = ModeNormal
		return m, nil
	}

	return m, nil
}

func (m Model) handleCommandInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeNormal
		m.commandInput = ""
		return m, nil

	case "enter":
		m.executeCommand(m.commandInput)
		m.mode = ModeNormal
		m.commandInput = ""
		return m, nil

	case "backspace":
		if len(m.commandInput) > 0 {
			m.commandInput = m.commandInput[:len(m.commandInput)-1]
		} else {
			m.mode = ModeNormal
		}
		return m, nil

	default:
		if len(msg.String()) == 1 {
			m.commandInput += msg.String()
		}
		return m, nil
	}
}

func (m *Model) executeCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "f":
		m.handleFilterCommand(args)
	case "s":
		m.handleSearchCommandExec(args)
	default:
		m.statusMsg = fmt.Sprintf("Unknown command: %s", cmd)
	}
}

func (m *Model) handleFilterCommand(args []string) {
	if len(args) == 0 || (len(args) == 1 && args[0] == "clear") {
		m.viewer.ClearFilter()
		m.filterNames = nil
		m.statusMsg = ""
		return
	}

	states := make(map[markdown.CheckboxState]bool)
	var names []string

	for _, name := range args {
		state, ok := markdown.ParseStateName(name)
		if ok {
			states[state] = true
			names = append(names, name)
		}
	}

	if len(states) > 0 {
		m.viewer.SetFilter(states)
		m.filterNames = names
	} else {
		m.statusMsg = "Unknown filter status"
	}
}

func (m *Model) handleSearchCommandExec(args []string) {
	if len(args) == 0 {
		m.statusMsg = "Usage: :s <query>"
		return
	}

	query := strings.Join(args, " ")
	results := search.Search(m.browser.Dir(), query)
	m.search.Show(query, results)
	m.mode = ModeSearch
}

func (m *Model) toggleFocus() {
	if m.focus == FocusBrowser {
		m.focus = FocusViewer
		m.browser.SetFocused(false)
		m.viewer.SetFocused(true)
	} else {
		m.focus = FocusBrowser
		m.browser.SetFocused(true)
		m.viewer.SetFocused(false)
	}
}

func (m *Model) autoPreview() {
	entry := m.browser.SelectedEntry()
	if entry == nil || entry.IsDir {
		m.viewer.LoadContent("")
		return
	}
	_ = m.viewer.LoadFile(entry.Path)

	// Update watcher to watch the current file and directory
	if m.watcher != nil {
		_ = m.watcher.WatchMultiple(m.browser.Dir(), entry.Path)
	}
}

func (m *Model) updatePaneSizes() {
	if m.width == 0 || m.height == 0 {
		return
	}

	contentHeight := m.height - 1
	browserWidth := m.width * m.cfg.SidebarWidth / 100
	viewerWidth := m.width - browserWidth

	m.browser.SetSize(browserWidth, contentHeight)
	m.viewer.SetSize(viewerWidth, contentHeight)
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Search overlay takes over the screen
	if m.mode == ModeSearch && m.search.Visible() {
		return m.search.View()
	}

	panes := lipgloss.JoinHorizontal(lipgloss.Top, m.browser.View(), m.viewer.View())
	statusBar := m.renderStatusBar()

	return lipgloss.JoinVertical(lipgloss.Left, panes, statusBar)
}

func (m Model) renderStatusBar() string {
	style := lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("252")).
		Width(m.width)

	if m.mode == ModeCommand {
		return style.Render(":" + m.commandInput + "█")
	}

	if m.statusMsg != "" {
		return style.Render(m.statusMsg)
	}

	// Left: path
	dir := m.browser.Dir()
	entry := m.browser.SelectedEntry()
	var pathDisplay string
	if entry != nil {
		pathDisplay = filepath.Join(dir, entry.Name)
	} else {
		pathDisplay = dir
	}

	// Abbreviate home dir
	if home, err := os.UserHomeDir(); err == nil {
		if strings.HasPrefix(pathDisplay, home) {
			pathDisplay = "~" + pathDisplay[len(home):]
		}
	}

	if len(pathDisplay) > m.width/2 {
		pathDisplay = "..." + pathDisplay[len(pathDisplay)-m.width/2+3:]
	}

	// Filter indicator
	filterStr := ""
	if len(m.filterNames) > 0 {
		filterStr = " [filter: " + strings.Join(m.filterNames, ", ") + "]"
	}

	left := pathDisplay + filterStr

	// Right: hotkey hints
	var hints string
	if m.focus == FocusBrowser {
		hints = "q:quit Tab:switch j/k:nav enter:open bs:up r:refresh"
	} else {
		hints = "q:quit Tab:switch j/k:scroll g/G:top/btm r:refresh"
	}

	hintsStyled := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render(hints)
	gap := m.width - lipgloss.Width(left) - lipgloss.Width(hintsStyled)
	if gap < 1 {
		gap = 1
	}

	return style.Render(left + strings.Repeat(" ", gap) + hints)
}
