package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/loader"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestLoadFile_Basic(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, ".env", "FOO=bar\nBAZ=qux\n")

	ef, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Path != path {
		t.Errorf("expected path %s, got %s", path, ef.Path)
	}
	if ef.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", ef.Vars["FOO"])
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := loader.LoadFile("/nonexistent/.env.missing")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadGlob_MatchesFiles(t *testing.T) {
	dir := t.TempDir()
	writeTempEnv(t, dir, ".env.production", "APP=prod\n")
	writeTempEnv(t, dir, ".env.staging", "APP=staging\n")

	files, err := loader.LoadGlob(filepath.Join(dir, ".env.*"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

func TestLoadGlob_NoMatches(t *testing.T) {
	dir := t.TempDir()
	_, err := loader.LoadGlob(filepath.Join(dir, ".env.*"))
	if err == nil {
		t.Fatal("expected error for no matches, got nil")
	}
}

func TestLoadPair_Success(t *testing.T) {
	dir := t.TempDir()
	pathA := writeTempEnv(t, dir, ".env.a", "KEY=alpha\n")
	pathB := writeTempEnv(t, dir, ".env.b", "KEY=beta\n")

	a, b, err := loader.LoadPair(pathA, pathB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Vars["KEY"] != "alpha" {
		t.Errorf("expected KEY=alpha in a, got %s", a.Vars["KEY"])
	}
	if b.Vars["KEY"] != "beta" {
		t.Errorf("expected KEY=beta in b, got %s", b.Vars["KEY"])
	}
}

func TestLoadPair_SecondMissing(t *testing.T) {
	dir := t.TempDir()
	pathA := writeTempEnv(t, dir, ".env.a", "KEY=alpha\n")

	_, _, err := loader.LoadPair(pathA, "/no/such/file")
	if err == nil {
		t.Fatal("expected error when second file is missing")
	}
}
