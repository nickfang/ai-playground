package browser

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRelativeTime(t *testing.T) {
	tests := []struct {
		name     string
		modified time.Time
		expected string
	}{
		{"just now", time.Now().Add(-10 * time.Second), "just now"},
		{"minutes ago", time.Now().Add(-5 * time.Minute), "5m ago"},
		{"hours ago", time.Now().Add(-3 * time.Hour), "3h ago"},
		{"days ago", time.Now().Add(-2 * 24 * time.Hour), "2d ago"},
		{"month ago", time.Now().Add(-45 * 24 * time.Hour), "1mo ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Entry{Modified: tt.modified}
			got := e.RelativeTime()
			if got != tt.expected {
				t.Errorf("RelativeTime() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestLoadEntries_Empty(t *testing.T) {
	dir := t.TempDir()
	entries, err := LoadEntries(dir, "alpha")
	if err != nil {
		t.Fatalf("LoadEntries error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries in empty dir, got %d", len(entries))
	}
}

func TestLoadEntries_OnlyMdFiles(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "readme.md"), []byte("# Readme"), 0644)
	os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("notes"), 0644)
	os.WriteFile(filepath.Join(dir, "code.go"), []byte("package main"), 0644)

	entries, err := LoadEntries(dir, "alpha")
	if err != nil {
		t.Fatalf("LoadEntries error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 .md entry, got %d", len(entries))
	}
	if entries[0].Name != "readme.md" {
		t.Errorf("expected readme.md, got %q", entries[0].Name)
	}
}

func TestLoadEntries_DirsFirst(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "aaa.md"), []byte("# AAA"), 0644)
	os.Mkdir(filepath.Join(dir, "zzz_dir"), 0755)

	entries, err := LoadEntries(dir, "alpha")
	if err != nil {
		t.Fatalf("LoadEntries error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if !entries[0].IsDir {
		t.Error("first entry should be a directory")
	}
	if entries[1].IsDir {
		t.Error("second entry should be a file")
	}
}

func TestLoadEntries_HiddenDirsSkipped(t *testing.T) {
	dir := t.TempDir()

	os.Mkdir(filepath.Join(dir, ".hidden"), 0755)
	os.Mkdir(filepath.Join(dir, "visible"), 0755)

	entries, err := LoadEntries(dir, "alpha")
	if err != nil {
		t.Fatalf("LoadEntries error: %v", err)
	}

	for _, e := range entries {
		if e.Name == ".hidden" {
			t.Error("hidden directory should be skipped")
		}
	}
}

func TestLoadEntries_AlphaSort(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "charlie.md"), []byte("# C"), 0644)
	os.WriteFile(filepath.Join(dir, "alpha.md"), []byte("# A"), 0644)
	os.WriteFile(filepath.Join(dir, "bravo.md"), []byte("# B"), 0644)

	entries, err := LoadEntries(dir, "alpha")
	if err != nil {
		t.Fatalf("LoadEntries error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Name != "alpha.md" {
		t.Errorf("first should be alpha.md, got %q", entries[0].Name)
	}
	if entries[1].Name != "bravo.md" {
		t.Errorf("second should be bravo.md, got %q", entries[1].Name)
	}
	if entries[2].Name != "charlie.md" {
		t.Errorf("third should be charlie.md, got %q", entries[2].Name)
	}
}

func TestLoadEntries_ModifiedSort(t *testing.T) {
	dir := t.TempDir()

	// Create files with different modification times
	os.WriteFile(filepath.Join(dir, "old.md"), []byte("# Old"), 0644)
	time.Sleep(50 * time.Millisecond)
	os.WriteFile(filepath.Join(dir, "new.md"), []byte("# New"), 0644)

	entries, err := LoadEntries(dir, "modified")
	if err != nil {
		t.Fatalf("LoadEntries error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	// Modified sort should have newest first
	if entries[0].Name != "new.md" {
		t.Errorf("newest file should be first, got %q", entries[0].Name)
	}
}

func TestLoadEntries_HeadingExtraction(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "test.md"), []byte("# My Heading\n\nSome content"), 0644)

	entries, err := LoadEntries(dir, "alpha")
	if err != nil {
		t.Fatalf("LoadEntries error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Heading != "My Heading" {
		t.Errorf("expected heading 'My Heading', got %q", entries[0].Heading)
	}
}

func TestLoadEntries_NoHeading(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "test.md"), []byte("No heading here"), 0644)

	entries, err := LoadEntries(dir, "alpha")
	if err != nil {
		t.Fatalf("LoadEntries error: %v", err)
	}
	if entries[0].Heading != "" {
		t.Errorf("expected empty heading, got %q", entries[0].Heading)
	}
}

func TestLoadEntries_InvalidDir(t *testing.T) {
	_, err := LoadEntries("/nonexistent/dir/path", "alpha")
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"a\nb\nc", []string{"a", "b", "c"}},
		{"single", []string{"single"}},
		{"", nil},
		{"a\n", []string{"a"}},
		{"\n", []string{""}},
	}

	for _, tt := range tests {
		got := splitLines(tt.input)
		if len(got) != len(tt.expected) {
			t.Errorf("splitLines(%q) len = %d, want %d", tt.input, len(got), len(tt.expected))
			continue
		}
		for i := range got {
			if got[i] != tt.expected[i] {
				t.Errorf("splitLines(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.expected[i])
			}
		}
	}
}

func TestTrimSpace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  ", "hello"},
		{"\thello\t", "hello"},
		{"hello", "hello"},
		{"  ", ""},
		{"", ""},
		{"\rhello\r", "hello"}, // trims CR from both ends
	}

	for _, tt := range tests {
		got := trimSpace(tt.input)
		if got != tt.expected {
			t.Errorf("trimSpace(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestSortEntries_Alpha(t *testing.T) {
	entries := []Entry{
		{Name: "charlie"},
		{Name: "alpha"},
		{Name: "bravo"},
	}
	sortEntries(entries, "alpha")
	if entries[0].Name != "alpha" || entries[1].Name != "bravo" || entries[2].Name != "charlie" {
		t.Errorf("alpha sort failed: %v", entries)
	}
}

func TestSortEntries_Modified(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Name: "old", Modified: now.Add(-2 * time.Hour)},
		{Name: "new", Modified: now},
		{Name: "mid", Modified: now.Add(-1 * time.Hour)},
	}
	sortEntries(entries, "modified")
	if entries[0].Name != "new" || entries[1].Name != "mid" || entries[2].Name != "old" {
		t.Errorf("modified sort failed: got %s, %s, %s", entries[0].Name, entries[1].Name, entries[2].Name)
	}
}

func TestSortEntries_SingleElement(t *testing.T) {
	entries := []Entry{{Name: "only"}}
	sortEntries(entries, "alpha") // should not panic
	if entries[0].Name != "only" {
		t.Error("single element sort failed")
	}
}

func TestSortEntries_Empty(t *testing.T) {
	var entries []Entry
	sortEntries(entries, "alpha") // should not panic
}
