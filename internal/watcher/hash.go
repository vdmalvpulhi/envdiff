package watcher

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// HashFile computes the SHA-256 hex digest of the file at path.
// It is exported so callers can snapshot file state before starting a Watcher.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("watcher: open %s: %w", path, err)
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("watcher: hash %s: %w", path, err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// HashBytes returns the SHA-256 hex digest of an in-memory byte slice.
// Useful in tests to compute expected hashes without touching the filesystem.
func HashBytes(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
