package watcher_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/watcher"
)

func TestHashFile_Deterministic(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	content := []byte("KEY=value\nOTHER=123\n")
	if err := os.WriteFile(p, content, 0o644); err != nil {
		t.Fatal(err)
	}

	h1, err := watcher.HashFile(p)
	if err != nil {
		t.Fatalf("HashFile: %v", err)
	}
	h2, err := watcher.HashFile(p)
	if err != nil {
		t.Fatalf("HashFile: %v", err)
	}
	if h1 != h2 {
		t.Errorf("expected identical hashes, got %s vs %s", h1, h2)
	}
}

func TestHashFile_ChangesOnEdit(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte("KEY=before\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before, _ := watcher.HashFile(p)

	if err := os.WriteFile(p, []byte("KEY=after\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	after, _ := watcher.HashFile(p)

	if before == after {
		t.Error("hash should differ after file content changes")
	}
}

func TestHashFile_NotFound(t *testing.T) {
	_, err := watcher.HashFile("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestHashBytes_MatchesFile(t *testing.T) {
	content := []byte("KEY=value\n")
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, content, 0o644); err != nil {
		t.Fatal(err)
	}
	fileHash, _ := watcher.HashFile(p)
	bytesHash := watcher.HashBytes(content)
	if fileHash != bytesHash {
		t.Errorf("HashFile=%s HashBytes=%s should match", fileHash, bytesHash)
	}
}
