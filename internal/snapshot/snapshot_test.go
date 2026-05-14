package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/snapshot"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_ENV", Status: "match", Value1: "production", Value2: "production"},
		{Key: "DB_HOST", Status: "missing_in_second", Value1: "localhost", Value2: ""},
		{Key: "API_KEY", Status: "conflict", Value1: "abc", Value2: "xyz"},
	}
}

func TestTake_SetsFields(t *testing.T) {
	before := time.Now()
	s := snapshot.Take("test-label", sampleResults())
	if s.Label != "test-label" {
		t.Errorf("expected label 'test-label', got %q", s.Label)
	}
	if len(s.Results) != 3 {
		t.Errorf("expected 3 results, got %d", len(s.Results))
	}
	if s.CreatedAt.Before(before) {
		t.Error("CreatedAt should be >= start of test")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	original := snapshot.Take("round-trip", sampleResults())
	if err := snapshot.Save(path, original); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Label != original.Label {
		t.Errorf("label mismatch: got %q, want %q", loaded.Label, original.Label)
	}
	if len(loaded.Results) != len(original.Results) {
		t.Errorf("results length mismatch: got %d, want %d", len(loaded.Results), len(original.Results))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	s := snapshot.Take("x", nil)
	err := snapshot.Save("/nonexistent/dir/snap.json", s)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestDiff_DetectsChanges(t *testing.T) {
	before := snapshot.Take("before", []diff.Result{
		{Key: "APP_ENV", Status: "match"},
		{Key: "DB_HOST", Status: "missing_in_second"},
	})
	after := snapshot.Take("after", []diff.Result{
		{Key: "APP_ENV", Status: "match"},
		{Key: "DB_HOST", Status: "conflict"},
		{Key: "NEW_KEY", Status: "missing_in_first"},
	})

	changes := snapshot.Diff(before, after)
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}

	changeMap := make(map[string]snapshot.Change)
	for _, c := range changes {
		changeMap[c.Key] = c
	}

	if c, ok := changeMap["DB_HOST"]; !ok || c.Before != "missing_in_second" || c.After != "conflict" {
		t.Errorf("unexpected DB_HOST change: %+v", changeMap["DB_HOST"])
	}
	if c, ok := changeMap["NEW_KEY"]; !ok || c.Before != "" || c.After != "missing_in_first" {
		t.Errorf("unexpected NEW_KEY change: %+v", c)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	results := sampleResults()
	before := snapshot.Take("a", results)
	after := snapshot.Take("b", results)
	changes := snapshot.Diff(before, after)
	if len(changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(changes))
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o644)
	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
