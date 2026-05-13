// Package watcher provides file-change detection for .env files using
// periodic SHA-256 hash comparison.
//
// Usage:
//
//	opts := watcher.DefaultOptions()
//	opts.PollInterval = 1 * time.Second
//
//	w := watcher.New([]string{".env", ".env.production"}, func(e watcher.Event) {
//		fmt.Printf("changed: %s\n", e.Path)
//	}, opts)
//	w.Start()
//	defer w.Stop()
//
// The watcher runs in a background goroutine and can be stopped cleanly via
// Stop. It does not use OS filesystem notifications, so it works uniformly
// across platforms including containers and network mounts.
package watcher
