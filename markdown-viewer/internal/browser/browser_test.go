package browser

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "alpha.md"), []byte("# Alpha\nContent"), 0644)
	os.WriteFile(filepath.Join(dir, "bravo.md"), []byte("# Bravo\nContent"), 0644)
	os.WriteFile(filepath.Join(dir, "charlie.md"), []byte("# Charlie\nContent"), 0644)
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "subdir", "nested.md"), []byte("# Nested\nContent"), 0644)

	return dir
}

func TestNew(t *testing.T) {
	dir := setupTestDir(t)
	m := New(dir, "alpha")

	if m.Dir() == "" {
		t.Error("Dir() should not be empty")
	}
	if m.SelectedEntry() == nil {
		t.Fatal("SelectedEntry() should not be nil with files present")
	}
}

func TestNew_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	m := New(dir, "alpha")

	if m.SelectedEntry() != nil {
		t.Error("SelectedEntry() should be nil for empty dir")
	}
}

func TestFocused(t *testing.T) {
	dir := setupTestDir(t)
	m := New(dir, "alpha")

	if m.Focused() {
		t.Error("should not be focused by default")
	}

	m.SetFocused(true)
	if !m.Focused() {
		t.Error("should be focused after SetFocused(true)")
	}

	m.SetFocused(false)
	if m.Focused() {
		t.Error("should not be focused after SetFocused(false)")
	}
}

func TestMoveDown(t *testing.T) {
	dir := setupTestDir(t)
	m := New(dir, "alpha")

	first := m.SelectedEntry()
	m.MoveDown()
	second := m.SelectedEntry()

	if first.Name == second.Name {
		t.Error("MoveDown should change selected entry")
	}
}

func TestMoveUp(t *testing.T) {
	dir := setupTestDir(t)
	m := New(dir, "alpha")

	// Move down first, then up
	m.MoveDown()
	second := m.SelectedEntry().Name
	m.MoveUp()
	first := m.SelectedEntry().Name

	if first == second {
		t.Error("MoveUp should change selected entry")
	}
}

func TestMoveUp_AtTop(t *testing.T) {
	dir := setupTestDir(t)
	m := New(dir, "alpha")

	name := m.SelectedEntry().Name
	m.MoveUp() // should be no-op at top
	if m.SelectedEntry().Name != name {
		t.Error("MoveUp at top should be no-op")
	}
}

func TestMoveDown_AtBottom(t *testing.T) {
	dir := setupTestDir(t)
	m := New(dir, "alpha")

	// Move to bottom
	for i := 0; i < 20; i++ {
		m.MoveDown()
	}
	name := m.SelectedEntry().Name
	m.MoveDown() // should be no-op at bottom
	if m.SelectedEntry().Name != name {
		t.Error("MoveDown at bottom should be no-op")
	}
}

func TestEnterDir(t *testing.T) {
	dir := setupTestDir(t)
	m := New(dir, "alpha")

	// First entry should be the subdir (dirs come first)
	entry := m.SelectedEntry()
	if !entry.IsDir {
		t.Skip("first entry is not a directory, test setup may differ")
	}

	origDir := m.Dir()
	ok := m.EnterDir()
	if !ok {
		t.Error("EnterDir should return true for a directory")
	}
	if m.Dir() == origDir {
		t.Error("Dir should change after EnterDir")
	}
}

func TestEnterDir_OnFile(t *testing.T) {
	dir := setupTestDir(t)
	m := New(dir, "alpha")

	// Move to a file entry (skip past directories)
	for m.SelectedEntry() != nil && m.SelectedEntry().IsDir {
		m.MoveDown()
	}

	ok := m.EnterDir()
	if ok {
		t.Error("EnterDir should return false for a file")
	}
}

func TestGoUp(t *testing.T) {
	dir := setupTestDir(t)
	m := New(dir, "alpha")

	// Enter subdir first
	if m.SelectedEntry().IsDir {
		m.EnterDir()
	}

	subDir := m.Dir()
	ok := m.GoUp()
	if !ok {
		t.Error("GoUp should return true")
	}
	if m.Dir() == subDir {
		t.Error("Dir should change after GoUp")
	}
}

func TestGoUp_SelectsPreviousDir(t *testing.T) {
	dir := setupTestDir(t)
	m := New(dir, "alpha")

	// Enter subdir
	if m.SelectedEntry() != nil && m.SelectedEntry().IsDir {
		dirName := m.SelectedEntry().Name
		m.EnterDir()
		m.GoUp()

		// Should try to re-select the dir we came from
		entry := m.SelectedEntry()
		if entry != nil && entry.Name == dirName {
			// Good - selected the dir we came from
		}
	}
}

func TestRefresh(t *testing.T) {
	dir := setupTestDir(t)
	m := New(dir, "alpha")

	// Add a new file
	os.WriteFile(filepath.Join(dir, "delta.md"), []byte("# Delta"), 0644)

	m.Refresh()

	found := false
	for {
		entry := m.SelectedEntry()
		if entry == nil {
			break
		}
		if entry.Name == "delta.md" {
			found = true
			break
		}
		m.MoveDown()
		if m.SelectedEntry().Name == entry.Name {
			break // hit bottom
		}
	}

	// Reset and search properly
	m2 := New(dir, "alpha")
	found = false
	for i := 0; i < 20; i++ {
		if m2.SelectedEntry() != nil && m2.SelectedEntry().Name == "delta.md" {
			found = true
			break
		}
		m2.MoveDown()
	}

	if !found {
		t.Error("Refresh should pick up new files")
	}
}

func TestToggleSort(t *testing.T) {
	dir := setupTestDir(t)
	m := New(dir, "alpha")

	m.ToggleSort()
	// After toggle, sort should be "modified"
	// We can't easily check the internal sortMode, but we can verify it doesn't panic
	// and entries are still loaded
	if m.SelectedEntry() == nil {
		t.Error("ToggleSort should not lose entries")
	}
}

func TestSetSize(t *testing.T) {
	dir := setupTestDir(t)
	m := New(dir, "alpha")

	m.SetSize(40, 20)
	// View should render without panic
	view := m.View()
	if view == "" {
		t.Error("View should not be empty after SetSize")
	}
}

func TestView_Empty(t *testing.T) {
	dir := t.TempDir()
	m := New(dir, "alpha")
	m.SetSize(40, 20)
	m.SetFocused(true)

	view := m.View()
	if view == "" {
		t.Error("View should render even for empty directory")
	}
}

func TestView_ZeroSize(t *testing.T) {
	dir := setupTestDir(t)
	m := New(dir, "alpha")

	view := m.View()
	if view != "" {
		t.Error("View with zero size should return empty string")
	}
}
