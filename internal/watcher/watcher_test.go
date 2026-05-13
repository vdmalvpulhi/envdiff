package watcher_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/user/envdiff/internal/watcher"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestWatcher_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	p := writeTempEnv(t, dir, ".env", "KEY=original\n")

	var mu sync.Mutex
	var events []watcher.Event

	opts := watcher.DefaultOptions()
	opts.PollInterval = 50 * time.Millisecond

	w := watcher.New([]string{p}, func(e watcher.Event) {
		mu.Lock()
		events = append(events, e)
		mu.Unlock()
	}, opts)
	w.Start()
	defer w.Stop()

	time.Sleep(100 * time.Millisecond)
	if err := os.WriteFile(p, []byte("KEY=changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(events) == 0 {
		t.Fatal("expected at least one change event")
	}
	if events[0].Path != p {
		t.Errorf("expected path %s, got %s", p, events[0].Path)
	}
	if events[0].OldHash == events[0].NewHash {
		t.Error("old and new hash should differ")
	}
}

func TestWatcher_NoEventWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	p := writeTempEnv(t, dir, ".env", "KEY=stable\n")

	var mu sync.Mutex
	var events []watcher.Event

	opts := watcher.DefaultOptions()
	opts.PollInterval = 50 * time.Millisecond

	w := watcher.New([]string{p}, func(e watcher.Event) {
		mu.Lock()
		events = append(events, e)
		mu.Unlock()
	}, opts)
	w.Start()
	defer w.Stop()

	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestWatcher_Stop(t *testing.T) {
	dir := t.TempDir()
	p := writeTempEnv(t, dir, ".env", "KEY=val\n")

	opts := watcher.DefaultOptions()
	opts.PollInterval = 50 * time.Millisecond

	w := watcher.New([]string{p}, func(_ watcher.Event) {}, opts)
	w.Start()
	w.Stop()
	// stopping twice must not panic
}
