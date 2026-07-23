package search

import (
	"os"
	"path/filepath"
	"testing"
)

func setupSearchDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "deploy.md"), []byte("# Deploy Checklist\n\n- [ ] Run migrations\n- [x] Build artifacts"), 0644)
	os.WriteFile(filepath.Join(dir, "notes.md"), []byte("# Meeting Notes\n\nDiscussed deployment pipeline"), 0644)
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "subdir", "nested.md"), []byte("# Nested File\n\nSome nested content about deploy"), 0644)
	os.Mkdir(filepath.Join(dir, ".hidden"), 0755)
	os.WriteFile(filepath.Join(dir, ".hidden", "secret.md"), []byte("# Secret"), 0644)

	return dir
}

func TestSearch_EmptyQuery(t *testing.T) {
	dir := setupSearchDir(t)
	results := Search(dir, "")
	if len(results) != 0 {
		t.Errorf("empty query should return 0 results, got %d", len(results))
	}
}

func TestSearch_FilenameMatch(t *testing.T) {
	dir := setupSearchDir(t)
	results := Search(dir, "deploy")

	hasFileMatch := false
	for _, r := range results {
		if r.MatchType == "file" && r.FileName == "deploy.md" {
			hasFileMatch = true
			break
		}
	}
	if !hasFileMatch {
		t.Error("should find deploy.md as a filename match")
	}
}

func TestSearch_ContentMatch(t *testing.T) {
	dir := setupSearchDir(t)
	results := Search(dir, "migration")

	hasContentMatch := false
	for _, r := range results {
		if r.MatchType == "content" {
			hasContentMatch = true
			break
		}
	}
	if !hasContentMatch {
		t.Error("should find content matches for 'migration'")
	}
}

func TestSearch_RecursiveSearch(t *testing.T) {
	dir := setupSearchDir(t)
	results := Search(dir, "nested")

	hasNested := false
	for _, r := range results {
		if filepath.Base(r.FileName) == "nested.md" || r.FileName == filepath.Join("subdir", "nested.md") {
			hasNested = true
			break
		}
	}
	if !hasNested {
		t.Error("should find nested.md in recursive search")
	}
}

func TestSearch_SkipsHiddenDirs(t *testing.T) {
	dir := setupSearchDir(t)
	results := Search(dir, "secret")

	for _, r := range results {
		if filepath.Base(r.FilePath) == "secret.md" {
			t.Error("should not search in hidden directories")
		}
	}
}

func TestSearch_ContentMatchHasLineNum(t *testing.T) {
	dir := setupSearchDir(t)
	results := Search(dir, "artifacts")

	for _, r := range results {
		if r.MatchType == "content" && r.LineNum > 0 {
			return // found valid content match with line number
		}
	}
	t.Error("content matches should have LineNum > 0")
}

func TestSearch_NonexistentDir(t *testing.T) {
	results := Search("/nonexistent/path", "test")
	if len(results) != 0 {
		t.Error("nonexistent dir should return 0 results")
	}
}

func TestSearch_NoMatches(t *testing.T) {
	dir := setupSearchDir(t)
	results := Search(dir, "zzzzxyzzy")
	if len(results) != 0 {
		t.Errorf("nonsense query should return 0 results, got %d", len(results))
	}
}

func TestSearch_ResultFields(t *testing.T) {
	dir := setupSearchDir(t)
	results := Search(dir, "deploy")

	if len(results) == 0 {
		t.Fatal("expected results for 'deploy'")
	}

	for _, r := range results {
		if r.FilePath == "" {
			t.Error("FilePath should not be empty")
		}
		if r.FileName == "" {
			t.Error("FileName should not be empty")
		}
		if r.MatchType != "file" && r.MatchType != "content" {
			t.Errorf("MatchType should be 'file' or 'content', got %q", r.MatchType)
		}
		if r.MatchType == "content" && r.LineContent == "" {
			t.Error("content match should have LineContent")
		}
	}
}
