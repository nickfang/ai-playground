package app

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nfang/markdown-viewer/internal/config"
	"github.com/nfang/markdown-viewer/internal/markdown"
	"github.com/nfang/markdown-viewer/internal/theme"
)

func setupTestApp(t *testing.T) (Model, string) {
	t.Helper()
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "alpha.md"), []byte("# Alpha\n\n- [ ] Todo\n- [x] Done"), 0644)
	os.WriteFile(filepath.Join(dir, "bravo.md"), []byte("# Bravo\n\n- [!] Urgent task"), 0644)
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "subdir", "nested.md"), []byte("# Nested"), 0644)

	cfg := config.DefaultConfig()
	th := theme.EmojiTheme()
	m := New(dir, cfg, th)

	// Simulate window size
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	return updated.(Model), dir
}

func TestNew_WithDirectory(t *testing.T) {
	m, _ := setupTestApp(t)

	if m.focus != FocusBrowser {
		t.Error("initial focus should be on browser")
	}
	if m.mode != ModeNormal {
		t.Error("initial mode should be normal")
	}
}

func TestNew_WithFilePath(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.md")
	os.WriteFile(filePath, []byte("# Test"), 0644)

	cfg := config.DefaultConfig()
	th := theme.EmojiTheme()
	m := New(filePath, cfg, th)

	// Should open the directory containing the file
	if m.browser.Dir() != dir {
		t.Errorf("expected dir %q, got %q", dir, m.browser.Dir())
	}
}

func TestUpdate_WindowSize(t *testing.T) {
	m, _ := setupTestApp(t)

	if m.width != 120 {
		t.Errorf("width should be 120, got %d", m.width)
	}
	if m.height != 40 {
		t.Errorf("height should be 40, got %d", m.height)
	}
}

func TestUpdate_QuitKey(t *testing.T) {
	m, _ := setupTestApp(t)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Error("q key should return a quit command")
	}
}

func TestUpdate_TabTogglesFocus(t *testing.T) {
	m, _ := setupTestApp(t)

	if m.focus != FocusBrowser {
		t.Error("should start with browser focus")
	}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)

	if m.focus != FocusViewer {
		t.Error("Tab should switch to viewer focus")
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)

	if m.focus != FocusBrowser {
		t.Error("Tab again should switch back to browser focus")
	}
}

func TestUpdate_ColonEntersCommandMode(t *testing.T) {
	m, _ := setupTestApp(t)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{':'}})
	m = updated.(Model)

	if m.mode != ModeCommand {
		t.Error(": should enter command mode")
	}
}

func TestUpdate_CommandMode_EscCancels(t *testing.T) {
	m, _ := setupTestApp(t)

	// Enter command mode
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{':'}})
	m = updated.(Model)

	// Type something
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	m = updated.(Model)

	// Esc
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEscape})
	m = updated.(Model)

	if m.mode != ModeNormal {
		t.Error("Esc should return to normal mode")
	}
	if m.commandInput != "" {
		t.Error("command input should be cleared")
	}
}

func TestUpdate_CommandMode_BackspaceDeletes(t *testing.T) {
	m, _ := setupTestApp(t)

	// Enter command mode and type
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{':'}})
	m = updated.(Model)

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	m = updated.(Model)

	if m.commandInput != "ab" {
		t.Errorf("command input should be 'ab', got %q", m.commandInput)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	m = updated.(Model)

	if m.commandInput != "a" {
		t.Errorf("after backspace, command input should be 'a', got %q", m.commandInput)
	}
}

func TestUpdate_CommandMode_BackspaceEmptyExits(t *testing.T) {
	m, _ := setupTestApp(t)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{':'}})
	m = updated.(Model)

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	m = updated.(Model)

	if m.mode != ModeNormal {
		t.Error("backspace on empty command should exit command mode")
	}
}

func TestUpdate_ResizePanes(t *testing.T) {
	m, _ := setupTestApp(t)

	origWidth := m.cfg.SidebarWidth

	// Increase
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}})
	m = updated.(Model)
	if m.cfg.SidebarWidth != origWidth+5 {
		t.Errorf("+ should increase sidebar width by 5, got %d", m.cfg.SidebarWidth)
	}

	// Decrease
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}})
	m = updated.(Model)
	if m.cfg.SidebarWidth != origWidth {
		t.Errorf("- should decrease sidebar width by 5, got %d", m.cfg.SidebarWidth)
	}
}

func TestUpdate_ResizePanes_Clamped(t *testing.T) {
	m, _ := setupTestApp(t)

	// Max out
	for i := 0; i < 20; i++ {
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}})
		m = updated.(Model)
	}
	if m.cfg.SidebarWidth > 50 {
		t.Error("sidebar width should not exceed 50")
	}

	// Min out
	for i := 0; i < 20; i++ {
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}})
		m = updated.(Model)
	}
	if m.cfg.SidebarWidth < 10 {
		t.Error("sidebar width should not go below 10")
	}
}

func TestUpdate_BrowserNavigation(t *testing.T) {
	m, _ := setupTestApp(t)

	first := m.browser.SelectedEntry()
	if first == nil {
		t.Fatal("should have entries")
	}

	// Move down
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)

	second := m.browser.SelectedEntry()
	if first.Name == second.Name {
		t.Error("j should move cursor down in browser")
	}

	// Move back up
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(Model)

	third := m.browser.SelectedEntry()
	if third.Name != first.Name {
		t.Error("k should move cursor back up")
	}
}

func TestExecuteCommand_Filter(t *testing.T) {
	m, _ := setupTestApp(t)

	m.executeCommand("f todo urgent")

	if len(m.filterNames) != 2 {
		t.Errorf("expected 2 filter names, got %d", len(m.filterNames))
	}
}

func TestExecuteCommand_FilterClear(t *testing.T) {
	m, _ := setupTestApp(t)

	m.executeCommand("f todo")
	m.executeCommand("f clear")

	if len(m.filterNames) != 0 {
		t.Error("filter names should be empty after clear")
	}
}

func TestExecuteCommand_FilterEmpty(t *testing.T) {
	m, _ := setupTestApp(t)

	m.executeCommand("f todo")
	m.executeCommand("f")

	if len(m.filterNames) != 0 {
		t.Error("empty :f should clear filters")
	}
}

func TestExecuteCommand_FilterInvalid(t *testing.T) {
	m, _ := setupTestApp(t)

	m.executeCommand("f bogus")

	if m.statusMsg == "" {
		t.Error("invalid filter should set status message")
	}
}

func TestExecuteCommand_Search(t *testing.T) {
	m, _ := setupTestApp(t)

	m.handleSearchCommandExec([]string{"alpha"})

	if m.mode != ModeSearch {
		t.Error("search should enter search mode")
	}
	if !m.search.Visible() {
		t.Error("search overlay should be visible")
	}
}

func TestExecuteCommand_SearchNoArgs(t *testing.T) {
	m, _ := setupTestApp(t)

	m.handleSearchCommandExec(nil)

	if m.statusMsg == "" {
		t.Error("search with no args should show usage message")
	}
}

func TestExecuteCommand_Unknown(t *testing.T) {
	m, _ := setupTestApp(t)

	m.executeCommand("unknown_cmd")

	if m.statusMsg == "" {
		t.Error("unknown command should set status message")
	}
}

func TestHandleFilterCommand_ValidStates(t *testing.T) {
	m, _ := setupTestApp(t)

	tests := []struct {
		args     []string
		expected int
	}{
		{[]string{"todo"}, 1},
		{[]string{"todo", "done"}, 2},
		{[]string{"todo", "done", "urgent"}, 3},
		{[]string{"todo", "done", "inprogress", "cancelled", "urgent"}, 5},
	}

	for _, tt := range tests {
		m.handleFilterCommand(tt.args)
		if len(m.filterNames) != tt.expected {
			t.Errorf("filter(%v) expected %d names, got %d", tt.args, tt.expected, len(m.filterNames))
		}
	}
}

func TestView_NormalMode(t *testing.T) {
	m, _ := setupTestApp(t)

	view := m.View()
	if view == "" {
		t.Error("View should not be empty in normal mode")
	}
}

func TestView_CommandMode(t *testing.T) {
	m, _ := setupTestApp(t)

	m.mode = ModeCommand
	m.commandInput = "test"

	view := m.View()
	if view == "" {
		t.Error("View should not be empty in command mode")
	}
}

func TestView_LoadingState(t *testing.T) {
	cfg := config.DefaultConfig()
	th := theme.EmojiTheme()
	m := New(".", cfg, th)

	// Before window size msg, width/height are 0
	view := m.View()
	if view != "Loading..." {
		t.Errorf("View before resize should show 'Loading...', got %q", view)
	}
}

func TestSearchMode_EscDismisses(t *testing.T) {
	m, _ := setupTestApp(t)

	// Enter search mode
	m.handleSearchCommandExec([]string{"test"})
	if m.mode != ModeSearch {
		t.Fatal("should be in search mode")
	}

	// Esc
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEscape})
	m = updated.(Model)

	if m.mode != ModeNormal {
		t.Error("Esc should exit search mode")
	}
}

func TestSearchMode_Navigation(t *testing.T) {
	m, _ := setupTestApp(t)

	m.handleSearchCommandExec([]string{"alpha"})

	// Navigate in search results
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	_ = updated.(Model)
	// Should not panic
}

func TestSearchMode_EnterSelectsResult(t *testing.T) {
	m, _ := setupTestApp(t)

	m.handleSearchCommandExec([]string{"alpha"})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)

	if m.mode != ModeNormal {
		t.Error("Enter should exit search mode")
	}
	if m.focus != FocusViewer {
		t.Error("Enter should switch focus to viewer")
	}
}

func TestAutoPreview(t *testing.T) {
	m, _ := setupTestApp(t)

	entry := m.browser.SelectedEntry()
	if entry != nil && !entry.IsDir {
		if m.viewer.FilePath() == "" {
			t.Error("auto-preview should load the selected file")
		}
	}
}

func TestToggleFocus(t *testing.T) {
	m, _ := setupTestApp(t)

	if !m.browser.Focused() {
		t.Error("browser should start focused")
	}
	if m.viewer.Focused() {
		t.Error("viewer should start unfocused")
	}

	m.toggleFocus()

	if m.browser.Focused() {
		t.Error("browser should be unfocused after toggle")
	}
	if !m.viewer.Focused() {
		t.Error("viewer should be focused after toggle")
	}
}

func TestInit(t *testing.T) {
	m, _ := setupTestApp(t)

	cmd := m.Init()
	// With a watcher, Init should return a command
	if m.watcher != nil && cmd == nil {
		t.Error("Init should return watcher listen command")
	}
}

func TestStatusBar_ClearsOnKey(t *testing.T) {
	m, _ := setupTestApp(t)
	m.statusMsg = "some message"

	// Any normal key should clear status
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)

	if m.statusMsg != "" {
		t.Error("status message should be cleared on keypress")
	}
}

// Test that ParseStateName in the app context works with the markdown package
func TestFilterCommand_UsesMarkdownPackage(t *testing.T) {
	// Verify the state mapping is consistent
	state, ok := markdown.ParseStateName("todo")
	if !ok || state != markdown.StateTodo {
		t.Error("markdown.ParseStateName should work for 'todo'")
	}

	state, ok = markdown.ParseStateName("urgent")
	if !ok || state != markdown.StateUrgent {
		t.Error("markdown.ParseStateName should work for 'urgent'")
	}
}
