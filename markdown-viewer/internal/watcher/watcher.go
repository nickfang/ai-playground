package watcher

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

// FileChangedMsg is sent when a watched file or directory changes.
type FileChangedMsg struct {
	Path string
}

// Watcher watches files and directories for changes.
type Watcher struct {
	watcher    *fsnotify.Watcher
	debounceMs int
}

// New creates a new file watcher.
func New() (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Watcher{
		watcher:    w,
		debounceMs: 100,
	}, nil
}

// Watch adds a path to watch.
func (w *Watcher) Watch(path string) error {
	// Remove all existing watches first
	for _, p := range w.watcher.WatchList() {
		_ = w.watcher.Remove(p)
	}
	return w.watcher.Add(path)
}

// WatchMultiple watches multiple paths.
func (w *Watcher) WatchMultiple(paths ...string) error {
	for _, p := range w.watcher.WatchList() {
		_ = w.watcher.Remove(p)
	}
	for _, p := range paths {
		if err := w.watcher.Add(p); err != nil {
			return err
		}
	}
	return nil
}

// Close closes the watcher.
func (w *Watcher) Close() error {
	return w.watcher.Close()
}

// ListenCmd returns a tea.Cmd that listens for file changes.
func (w *Watcher) ListenCmd() tea.Cmd {
	return func() tea.Msg {
		var lastEvent time.Time
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return nil
				}
				// Debounce: ignore events within 100ms of the last one
				now := time.Now()
				if now.Sub(lastEvent) < time.Duration(w.debounceMs)*time.Millisecond {
					continue
				}
				lastEvent = now

				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) ||
					event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					return FileChangedMsg{Path: event.Name}
				}

			case _, ok := <-w.watcher.Errors:
				if !ok {
					return nil
				}
				// Ignore errors, just continue watching
			}
		}
	}
}
