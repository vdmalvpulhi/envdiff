package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envdiff/internal/baseline"
	"github.com/yourorg/envdiff/internal/diff"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_ENV", Status: diff.StatusMatch, Values: []string{"production", "production"}},
		{Key: "DB_URL", Status: diff.StatusMissing, Values: []string{"postgres://localhost", ""}},
		{Key: "SECRET", Status: diff.StatusConflict, Values: []string{"abc", "xyz"}},
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	results := sampleResults()
	if err := baseline.Save(path, results); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(snap.Results) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(snap.Results))
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if snap.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Error("CreatedAt is in the future")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := baseline.Load("/nonexistent/path/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	err := baseline.Save("/nonexistent/dir/baseline.json", sampleResults())
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json{"), 0o644)
	_, err := baseline.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
