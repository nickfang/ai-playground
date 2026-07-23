package search

import (
	"testing"
)

func TestNew_Search(t *testing.T) {
	m := New()
	if m.Visible() {
		t.Error("should not be visible by default")
	}
	if m.Selected() != nil {
		t.Error("Selected should be nil when no results")
	}
}

func TestShow(t *testing.T) {
	m := New()
	results := []Result{
		{FilePath: "/a.md", FileName: "a.md", MatchType: "file"},
		{FilePath: "/b.md", FileName: "b.md", LineNum: 5, LineContent: "hello", MatchType: "content"},
	}

	m.Show("test", results)

	if !m.Visible() {
		t.Error("should be visible after Show")
	}
	if m.Selected() == nil {
		t.Error("Selected should not be nil with results")
	}
	if m.Selected().FileName != "a.md" {
		t.Errorf("Selected should be first result, got %q", m.Selected().FileName)
	}
}

func TestHide(t *testing.T) {
	m := New()
	m.Show("test", []Result{{FilePath: "/a.md", FileName: "a.md", MatchType: "file"}})

	m.Hide()

	if m.Visible() {
		t.Error("should not be visible after Hide")
	}
	if m.Selected() != nil {
		t.Error("Selected should be nil after Hide")
	}
}

func TestMoveDown_Search(t *testing.T) {
	m := New()
	results := []Result{
		{FilePath: "/a.md", FileName: "a.md", MatchType: "file"},
		{FilePath: "/b.md", FileName: "b.md", MatchType: "file"},
		{FilePath: "/c.md", FileName: "c.md", MatchType: "file"},
	}
	m.Show("test", results)

	m.MoveDown()
	if m.Selected().FileName != "b.md" {
		t.Errorf("after MoveDown, expected b.md, got %q", m.Selected().FileName)
	}

	m.MoveDown()
	if m.Selected().FileName != "c.md" {
		t.Errorf("after second MoveDown, expected c.md, got %q", m.Selected().FileName)
	}

	// At bottom, should not go further
	m.MoveDown()
	if m.Selected().FileName != "c.md" {
		t.Error("MoveDown at bottom should be no-op")
	}
}

func TestMoveUp_Search(t *testing.T) {
	m := New()
	results := []Result{
		{FilePath: "/a.md", FileName: "a.md", MatchType: "file"},
		{FilePath: "/b.md", FileName: "b.md", MatchType: "file"},
	}
	m.Show("test", results)

	m.MoveDown()
	m.MoveUp()
	if m.Selected().FileName != "a.md" {
		t.Errorf("after MoveUp, expected a.md, got %q", m.Selected().FileName)
	}

	// At top, should not go further
	m.MoveUp()
	if m.Selected().FileName != "a.md" {
		t.Error("MoveUp at top should be no-op")
	}
}

func TestMoveDown_EmptyResults(t *testing.T) {
	m := New()
	m.Show("test", nil)

	m.MoveDown() // should not panic
	m.MoveUp()   // should not panic

	if m.Selected() != nil {
		t.Error("Selected should be nil with no results")
	}
}

func TestView_NotVisible(t *testing.T) {
	m := New()
	m.SetSize(80, 24)

	view := m.View()
	if view != "" {
		t.Error("View should return empty string when not visible")
	}
}

func TestView_Visible(t *testing.T) {
	m := New()
	m.SetSize(80, 24)

	results := []Result{
		{FilePath: "/test.md", FileName: "test.md", MatchType: "file"},
	}
	m.Show("test", results)

	view := m.View()
	if view == "" {
		t.Error("View should not be empty when visible")
	}
}

func TestView_ZeroSize(t *testing.T) {
	m := New()
	m.Show("test", []Result{{FilePath: "/a.md", FileName: "a.md", MatchType: "file"}})

	view := m.View()
	if view != "" {
		t.Error("View with zero size should return empty string")
	}
}

func TestView_NoResults(t *testing.T) {
	m := New()
	m.SetSize(80, 24)
	m.Show("noresults", nil)

	view := m.View()
	if view == "" {
		t.Error("View should render even with no results")
	}
}

func TestSetSize_Search(t *testing.T) {
	m := New()
	m.SetSize(100, 50)
	// Should not panic
	m.SetSize(0, 0)
}
