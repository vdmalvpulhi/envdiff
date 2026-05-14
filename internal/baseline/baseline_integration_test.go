package baseline_test

import (
	"path/filepath"
	"testing"

	"github.com/yourorg/envdiff/internal/baseline"
	"github.com/yourorg/envdiff/internal/diff"
)

// TestBaseline_Integration_SaveCompareCycle exercises the full save → load → compare pipeline.
func TestBaseline_Integration_SaveCompareCycle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	initial := []diff.Result{
		{Key: "HOST", Status: diff.StatusMatch, Values: []string{"localhost", "localhost"}},
		{Key: "PORT", Status: diff.StatusConflict, Values: []string{"8080", "9090"}},
		{Key: "TOKEN", Status: diff.StatusMissing, Values: []string{"abc123", ""}},
	}

	if err := baseline.Save(path, initial); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	// Simulate: PORT conflict resolved, TOKEN added to second env, NEW_KEY appeared.
	updated := []diff.Result{
		{Key: "HOST", Status: diff.StatusMatch, Values: []string{"localhost", "localhost"}},
		{Key: "PORT", Status: diff.StatusMatch, Values: []string{"8080", "8080"}},
		{Key: "TOKEN", Status: diff.StatusMissing, Values: []string{"abc123", ""}},
		{Key: "NEW_KEY", Status: diff.StatusMissing, Values: []string{"val", ""}},
	}

	delta := baseline.Compare(snap, updated)

	if len(delta.New) != 1 || delta.New[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY as new, got %v", delta.New)
	}
	if len(delta.Resolved) != 0 {
		t.Errorf("expected no resolved keys, got %v", delta.Resolved)
	}
	if len(delta.Changed) != 1 || delta.Changed[0].Key != "PORT" {
		t.Errorf("expected PORT as changed, got %v", delta.Changed)
	}
}
