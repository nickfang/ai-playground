package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew_Watcher(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer w.Close()

	if w.watcher == nil {
		t.Error("internal watcher should not be nil")
	}
}

func TestWatch(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	dir := t.TempDir()
	err = w.Watch(dir)
	if err != nil {
		t.Fatalf("Watch error: %v", err)
	}

	// Verify it's in the watch list
	list := w.watcher.WatchList()
	if len(list) != 1 {
		t.Errorf("expected 1 watch, got %d", len(list))
	}
}

func TestWatch_ReplacesExisting(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	dir1 := t.TempDir()
	dir2 := t.TempDir()

	w.Watch(dir1)
	w.Watch(dir2)

	list := w.watcher.WatchList()
	if len(list) != 1 {
		t.Errorf("Watch should replace existing; expected 1 watch, got %d", len(list))
	}
}

func TestWatchMultiple(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	dir1 := t.TempDir()
	dir2 := t.TempDir()

	err = w.WatchMultiple(dir1, dir2)
	if err != nil {
		t.Fatalf("WatchMultiple error: %v", err)
	}

	list := w.watcher.WatchList()
	if len(list) != 2 {
		t.Errorf("expected 2 watches, got %d", len(list))
	}
}

func TestWatchMultiple_ReplacesExisting(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	dir1 := t.TempDir()
	dir2 := t.TempDir()
	dir3 := t.TempDir()

	w.Watch(dir1)
	w.WatchMultiple(dir2, dir3)

	list := w.watcher.WatchList()
	if len(list) != 2 {
		t.Errorf("WatchMultiple should replace existing; expected 2 watches, got %d", len(list))
	}
}

func TestClose(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		t.Errorf("Close error: %v", err)
	}
}

func TestFileChangedMsg(t *testing.T) {
	msg := FileChangedMsg{Path: "/test/file.md"}
	if msg.Path != "/test/file.md" {
		t.Errorf("expected path '/test/file.md', got %q", msg.Path)
	}
}

func TestListenCmd_DetectsWrite(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	dir := t.TempDir()
	testFile := filepath.Join(dir, "test.md")
	os.WriteFile(testFile, []byte("initial"), 0644)

	err = w.Watch(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Start listening in a goroutine
	done := make(chan interface{}, 1)
	go func() {
		cmd := w.ListenCmd()
		msg := cmd()
		done <- msg
	}()

	// Give watcher time to fully register
	time.Sleep(200 * time.Millisecond)

	// Trigger a write event — use Create + Write pattern for reliability
	tmpFile := testFile + ".tmp"
	os.WriteFile(tmpFile, []byte("modified"), 0644)
	os.Rename(tmpFile, testFile)

	// Also try direct write
	time.Sleep(50 * time.Millisecond)
	f, _ := os.OpenFile(testFile, os.O_WRONLY|os.O_TRUNC, 0644)
	if f != nil {
		f.Write([]byte("modified again"))
		f.Sync()
		f.Close()
	}

	select {
	case msg := <-done:
		if msg == nil {
			t.Error("expected FileChangedMsg, got nil")
		}
		if fcm, ok := msg.(FileChangedMsg); ok {
			if fcm.Path == "" {
				t.Error("FileChangedMsg Path should not be empty")
			}
		} else {
			t.Errorf("expected FileChangedMsg, got %T", msg)
		}
	case <-time.After(5 * time.Second):
		t.Error("timed out waiting for file change event")
	}
}
