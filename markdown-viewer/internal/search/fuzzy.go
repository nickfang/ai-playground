package search

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sahilm/fuzzy"
)

// Result represents a single search match.
type Result struct {
	FilePath    string
	FileName    string
	LineNum     int    // 0 for filename matches
	LineContent string // the matching line content
	MatchType   string // "file" or "content"
}

// Search performs a fuzzy search across all .md files in dir (recursive).
func Search(dir string, query string) []Result {
	if query == "" {
		return nil
	}

	var allFiles []string
	var allLines []lineEntry

	// Walk directory tree
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		// Skip hidden directories
		if info.IsDir() && info.Name()[0] == '.' && path != dir {
			return filepath.SkipDir
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			relPath, _ := filepath.Rel(dir, path)
			allFiles = append(allFiles, relPath)

			// Read file lines for content search
			data, err := os.ReadFile(path)
			if err == nil {
				lines := strings.Split(string(data), "\n")
				for i, line := range lines {
					trimmed := strings.TrimSpace(line)
					if trimmed != "" {
						allLines = append(allLines, lineEntry{
							filePath: relPath,
							fullPath: path,
							lineNum:  i + 1,
							content:  trimmed,
						})
					}
				}
			}
		}
		return nil
	})

	var results []Result

	// Fuzzy match filenames
	fileMatches := fuzzy.Find(query, allFiles)
	for _, m := range fileMatches {
		if len(results) >= 50 {
			break
		}
		results = append(results, Result{
			FilePath:  filepath.Join(dir, allFiles[m.Index]),
			FileName:  allFiles[m.Index],
			MatchType: "file",
		})
	}

	// Fuzzy match content
	contentStrings := make([]string, len(allLines))
	for i, l := range allLines {
		contentStrings[i] = l.content
	}

	contentMatches := fuzzy.Find(query, contentStrings)
	for _, m := range contentMatches {
		if len(results) >= 100 {
			break
		}
		entry := allLines[m.Index]
		results = append(results, Result{
			FilePath:    entry.fullPath,
			FileName:    entry.filePath,
			LineNum:     entry.lineNum,
			LineContent: entry.content,
			MatchType:   "content",
		})
	}

	return results
}

type lineEntry struct {
	filePath string
	fullPath string
	lineNum  int
	content  string
}
