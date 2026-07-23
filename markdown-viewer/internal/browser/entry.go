package browser

import (
	"fmt"
	"os"
	"time"
)

// Entry represents a file or directory in the file browser.
type Entry struct {
	Name     string
	IsDir    bool
	Heading  string // first # heading (for .md files)
	Modified time.Time
	Path     string // full path
}

// RelativeTime returns a human-friendly relative time string.
func (e Entry) RelativeTime() string {
	d := time.Since(e.Modified)

	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		return fmt.Sprintf("%dm ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		return fmt.Sprintf("%dh ago", h)
	case d < 30*24*time.Hour:
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	default:
		months := int(d.Hours() / 24 / 30)
		if months < 1 {
			months = 1
		}
		return fmt.Sprintf("%dmo ago", months)
	}
}

// LoadEntries reads a directory and returns sorted entries.
// Directories come first, then .md files. Sorted by the given mode.
func LoadEntries(dir string, sortMode string) ([]Entry, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var dirs []Entry
	var files []Entry

	for _, de := range dirEntries {
		info, err := de.Info()
		if err != nil {
			continue
		}

		if de.IsDir() {
			// Skip hidden directories
			if de.Name()[0] == '.' {
				continue
			}
			dirs = append(dirs, Entry{
				Name:     de.Name(),
				IsDir:    true,
				Modified: info.ModTime(),
				Path:     dir + "/" + de.Name(),
			})
		} else if len(de.Name()) > 3 && de.Name()[len(de.Name())-3:] == ".md" {
			heading := readFirstHeading(dir + "/" + de.Name())
			files = append(files, Entry{
				Name:     de.Name(),
				IsDir:    false,
				Heading:  heading,
				Modified: info.ModTime(),
				Path:     dir + "/" + de.Name(),
			})
		}
	}

	// Sort directories and files
	sortEntries(dirs, sortMode)
	sortEntries(files, sortMode)

	return append(dirs, files...), nil
}

func readFirstHeading(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	content := string(data)
	lines := splitLines(content)
	for _, line := range lines {
		trimmed := trimSpace(line)
		if len(trimmed) > 2 && trimmed[0] == '#' && trimmed[1] == ' ' {
			return trimmed[2:]
		}
	}
	return ""
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func sortEntries(entries []Entry, mode string) {
	if len(entries) <= 1 {
		return
	}

	// Simple insertion sort (small lists)
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0; j-- {
			swap := false
			if mode == "modified" {
				swap = entries[j].Modified.After(entries[j-1].Modified)
			} else {
				swap = entries[j].Name < entries[j-1].Name
			}
			if swap {
				entries[j], entries[j-1] = entries[j-1], entries[j]
			} else {
				break
			}
		}
	}
}
