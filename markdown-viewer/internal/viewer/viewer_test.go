package viewer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nfang/markdown-viewer/internal/markdown"
)

func testIcons() map[markdown.CheckboxState]string {
	return map[markdown.CheckboxState]string{
		markdown.StateTodo:       "⬜",
		markdown.StateDone:       "✅",
		markdown.StateInProgress: "🔄",
		markdown.StateCancelled:  "❌",
		markdown.StateUrgent:     "🔥",
	}
}

func TestNew(t *testing.T) {
	m := New(testIcons())
	if m.Focused() {
		t.Error("should not be focused by default")
	}
	if m.FilePath() != "" {
		t.Error("FilePath should be empty initially")
	}
}

func TestSetFocused(t *testing.T) {
	m := New(testIcons())

	m.SetFocused(true)
	if !m.Focused() {
		t.Error("should be focused after SetFocused(true)")
	}

	m.SetFocused(false)
	if m.Focused() {
		t.Error("should not be focused after SetFocused(false)")
	}
}

func TestLoadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	os.WriteFile(path, []byte("# Test\n\n- [ ] Todo"), 0644)

	m := New(testIcons())
	m.SetSize(80, 24)

	err := m.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile error: %v", err)
	}

	if m.FilePath() != path {
		t.Errorf("FilePath should be %q, got %q", path, m.FilePath())
	}
}

func TestLoadFile_NonExistent(t *testing.T) {
	m := New(testIcons())
	m.SetSize(80, 24)

	err := m.LoadFile("/nonexistent/file.md")
	if err == nil {
		t.Error("LoadFile should return error for nonexistent file")
	}
}

func TestLoadContent(t *testing.T) {
	m := New(testIcons())
	m.SetSize(80, 24)

	m.LoadContent("# Hello\n\nWorld")

	if m.FilePath() != "" {
		t.Error("FilePath should be empty after LoadContent")
	}
}

func TestLoadContent_Empty(t *testing.T) {
	m := New(testIcons())
	m.SetSize(80, 24)

	m.LoadContent("")
	// Should not panic, view should show placeholder
}

func TestScrolling(t *testing.T) {
	m := New(testIcons())
	m.SetSize(80, 10) // small height to enable scrolling

	// Create long content
	content := "# Heading\n"
	for i := 0; i < 50; i++ {
		content += "- [ ] Item line\n"
	}
	m.LoadContent(content)

	// Should start at top
	m.ScrollDown()
	m.ScrollDown()
	m.ScrollDown()

	// ScrollUp should work
	m.ScrollUp()

	// ScrollToTop should reset
	m.ScrollToTop()

	// ScrollToBottom
	m.ScrollToBottom()

	// Calling these shouldn't panic
	m.ScrollUp()
	m.ScrollDown()
}

func TestScrollUp_AtTop(t *testing.T) {
	m := New(testIcons())
	m.SetSize(80, 24)
	m.LoadContent("# Short content")

	m.ScrollUp() // should be no-op at top, no panic
}

func TestScrollDown_ShortContent(t *testing.T) {
	m := New(testIcons())
	m.SetSize(80, 24)
	m.LoadContent("# Short")

	m.ScrollDown() // content fits in view, should be no-op
}

func TestSetFilter(t *testing.T) {
	m := New(testIcons())
	m.SetSize(80, 24)
	m.LoadContent("- [ ] Todo\n- [x] Done\n- [!] Urgent")

	filter := map[markdown.CheckboxState]bool{
		markdown.StateUrgent: true,
	}
	m.SetFilter(filter)
	// Should re-render with filter applied
}

func TestClearFilter(t *testing.T) {
	m := New(testIcons())
	m.SetSize(80, 24)
	m.LoadContent("- [ ] Todo\n- [x] Done")

	filter := map[markdown.CheckboxState]bool{markdown.StateTodo: true}
	m.SetFilter(filter)
	m.ClearFilter()
	// Should re-render without filter
}

func TestReload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	os.WriteFile(path, []byte("# Original"), 0644)

	m := New(testIcons())
	m.SetSize(80, 24)
	m.LoadFile(path)

	// Modify file
	os.WriteFile(path, []byte("# Updated"), 0644)

	err := m.Reload()
	if err != nil {
		t.Fatalf("Reload error: %v", err)
	}
}

func TestReload_NoFile(t *testing.T) {
	m := New(testIcons())
	m.SetSize(80, 24)

	err := m.Reload()
	if err != nil {
		t.Error("Reload with no file should return nil")
	}
}

func TestView_WithContent(t *testing.T) {
	m := New(testIcons())
	m.SetSize(80, 24)
	m.LoadContent("# Hello\n\n- [ ] A task")

	view := m.View()
	if view == "" {
		t.Error("View should not be empty with content")
	}
}

func TestView_NoContent(t *testing.T) {
	m := New(testIcons())
	m.SetSize(80, 24)

	view := m.View()
	if view == "" {
		t.Error("View should show placeholder when empty")
	}
}

func TestView_ZeroSize(t *testing.T) {
	m := New(testIcons())
	view := m.View()
	if view != "" {
		t.Error("View with zero size should return empty string")
	}
}

func TestSetSize_Rerenders(t *testing.T) {
	m := New(testIcons())
	m.SetSize(80, 24)
	m.LoadContent("# Test")

	// Resize should re-render
	m.SetSize(40, 12)
	view := m.View()
	if view == "" {
		t.Error("View should render after resize")
	}
}

func TestViewableHeight(t *testing.T) {
	m := New(testIcons())

	m.SetSize(80, 10)
	h := m.viewableHeight()
	if h != 8 { // 10 - 2 for border
		t.Errorf("viewableHeight() = %d, want 8", h)
	}

	m.SetSize(80, 2)
	h = m.viewableHeight()
	if h != 1 { // minimum 1
		t.Errorf("viewableHeight() with small height = %d, want 1", h)
	}

	m.SetSize(80, 1)
	h = m.viewableHeight()
	if h != 1 {
		t.Errorf("viewableHeight() with height 1 = %d, want 1", h)
	}
}
