package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Entry struct {
	Word string // first word, lowercased
	Line string // full original line
}

type Result struct {
	Entry
	Score          int
	MatchedIndices []int
}

func loadWords(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		word := strings.ToLower(parts[0])
		entries = append(entries, Entry{Word: word, Line: line})
	}
	return entries, scanner.Err()
}

func fuzzyMatch(query, word string) (int, []int, bool) {
	query = strings.ToLower(query)
	word = strings.ToLower(word)

	qi := 0
	var indices []int

	for wi := 0; wi < len(word) && qi < len(query); wi++ {
		if word[wi] == query[qi] {
			indices = append(indices, wi)
			qi++
		}
	}

	if qi < len(query) {
		return 0, nil, false
	}

	// Scoring
	score := len(indices) // +1 per matched char

	// First char bonus
	if len(indices) > 0 && indices[0] == 0 {
		score += 10
	}

	for i := 1; i < len(indices); i++ {
		gap := indices[i] - indices[i-1]
		// Consecutive bonus
		if gap == 1 {
			score += 5
		}
		// Tight clustering bonus
		if bonus := 3 - gap; bonus > 0 {
			score += bonus
		}
	}

	// Length penalty: prefer shorter, more precise matches
	score -= len(word) - len(query)

	return score, indices, true
}

func search(query string, entries []Entry) []Result {
	var results []Result
	for _, e := range entries {
		score, indices, ok := fuzzyMatch(query, e.Word)
		if ok {
			results = append(results, Result{Entry: e, Score: score, MatchedIndices: indices})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	if len(results) > 10 {
		results = results[:10]
	}
	return results
}

func highlight(word string, indices []int) string {
	idxSet := make(map[int]bool, len(indices))
	for _, i := range indices {
		idxSet[i] = true
	}

	const (
		boldYellow = "\033[1;33m"
		reset      = "\033[0m"
	)

	var b strings.Builder
	inHighlight := false
	for i, ch := range word {
		if idxSet[i] {
			if !inHighlight {
				b.WriteString(boldYellow)
				inHighlight = true
			}
		} else {
			if inHighlight {
				b.WriteString(reset)
				inHighlight = false
			}
		}
		b.WriteRune(ch)
	}
	if inHighlight {
		b.WriteString(reset)
	}
	return b.String()
}

func printResults(results []Result) {
	if len(results) == 0 {
		fmt.Println("  No matches found.")
		return
	}

	// Find the longest word for alignment
	maxWordLen := 0
	for _, r := range results {
		if len(r.Word) > maxWordLen {
			maxWordLen = len(r.Word)
		}
	}

	for i, r := range results {
		highlighted := highlight(r.Word, r.MatchedIndices)
		// Pad using the original word length (ANSI codes don't take visual space)
		padding := strings.Repeat(" ", maxWordLen-len(r.Word))
		// Truncate definition
		def := ""
		parts := strings.SplitN(r.Line, " ", 2)
		if len(parts) > 1 {
			def = parts[1]
			if len(def) > 70 {
				def = def[:70] + "..."
			}
		}
		fmt.Printf(" %2d. %s%s  %s\n", i+1, highlighted, padding, def)
	}
}

func main() {
	entries, err := loadWords("words.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading words: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded %d words. Type a query to fuzzy search (q to quit).\n\n", len(entries))

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("fuzzy> ")
		if !scanner.Scan() {
			break
		}
		query := strings.TrimSpace(scanner.Text())
		if query == "" || query == "q" || query == "quit" {
			break
		}
		results := search(query, entries)
		printResults(results)
		fmt.Println()
	}
}
