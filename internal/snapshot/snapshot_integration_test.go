package snapshot_test

import (
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/snapshot"
)

// TestSnapshot_Integration_FullCycle exercises Take → Save → Load → Diff
// using two successive snapshots to verify the end-to-end pipeline.
func TestSnapshot_Integration_FullCycle(t *testing.T) {
	dir := t.TempDir()

	// First snapshot: DB_HOST is missing in second env.
	v1Results := []diff.Result{
		{Key: "APP_ENV", Status: "match", Value1: "staging", Value2: "staging"},
		{Key: "DB_HOST", Status: "missing_in_second", Value1: "db.internal", Value2: ""},
	}
	v1 := snapshot.Take("v1", v1Results)
	v1Path := filepath.Join(dir, "v1.json")
	if err := snapshot.Save(v1Path, v1); err != nil {
		t.Fatalf("save v1: %v", err)
	}

	// Second snapshot: DB_HOST is now present but conflicts; SECRET_KEY is new.
	v2Results := []diff.Result{
		{Key: "APP_ENV", Status: "match", Value1: "staging", Value2: "staging"},
		{Key: "DB_HOST", Status: "conflict", Value1: "db.internal", Value2: "db.external"},
		{Key: "SECRET_KEY", Status: "missing_in_first", Value1: "", Value2: "s3cr3t"},
	}
	v2 := snapshot.Take("v2", v2Results)
	v2Path := filepath.Join(dir, "v2.json")
	if err := snapshot.Save(v2Path, v2); err != nil {
		t.Fatalf("save v2: %v", err)
	}

	// Reload both from disk to ensure persistence.
	loaded1, err := snapshot.Load(v1Path)
	if err != nil {
		t.Fatalf("load v1: %v", err)
	}
	loaded2, err := snapshot.Load(v2Path)
	if err != nil {
		t.Fatalf("load v2: %v", err)
	}

	changes := snapshot.Diff(loaded1, loaded2)

	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d: %+v", len(changes), changes)
	}

	changeMap := make(map[string]snapshot.Change)
	for _, c := range changes {
		changeMap[c.Key] = c
	}

	if c, ok := changeMap["DB_HOST"]; !ok {
		t.Error("expected change for DB_HOST")
	} else if c.Before != "missing_in_second" || c.After != "conflict" {
		t.Errorf("DB_HOST: got before=%q after=%q", c.Before, c.After)
	}

	if c, ok := changeMap["SECRET_KEY"]; !ok {
		t.Error("expected change for SECRET_KEY")
	} else if c.Before != "" || c.After != "missing_in_first" {
		t.Errorf("SECRET_KEY: got before=%q after=%q", c.Before, c.After)
	}
}
