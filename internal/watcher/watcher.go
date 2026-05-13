// Package watcher monitors .env files for changes and triggers a callback
// when modifications are detected, enabling live diff updates.
package watcher

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// Event represents a file change notification.
type Event struct {
	Path    string
	OldHash string
	NewHash string
}

// Handler is called when a watched file changes.
type Handler func(Event)

// Options configures the watcher behaviour.
type Options struct {
	// PollInterval is how often files are checked for changes.
	PollInterval time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		PollInterval: 2 * time.Second,
	}
}

// Watcher polls a set of file paths and invokes a handler on change.
type Watcher struct {
	paths   []string
	opts    Options
	hashes  map[string]string
	handler Handler
	stop    chan struct{}
}

// New creates a Watcher for the given paths.
func New(paths []string, handler Handler, opts Options) *Watcher {
	return &Watcher{
		paths:   paths,
		opts:    opts,
		hashes:  make(map[string]string),
		handler: handler,
		stop:    make(chan struct{}),
	}
}

// Start begins polling in a background goroutine.
func (w *Watcher) Start() {
	for _, p := range w.paths {
		w.hashes[p], _ = hashFile(p)
	}
	go w.loop()
}

// Stop signals the watcher to cease polling.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) loop() {
	ticker := time.NewTicker(w.opts.PollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			w.check()
		case <-w.stop:
			return
		}
	}
}

func (w *Watcher) check() {
	for _, p := range w.paths {
		newHash, err := hashFile(p)
		if err != nil {
			continue
		}
		if old, ok := w.hashes[p]; ok && old != newHash {
			w.handler(Event{Path: p, OldHash: old, NewHash: newHash})
		}
		w.hashes[p] = newHash
	}
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
