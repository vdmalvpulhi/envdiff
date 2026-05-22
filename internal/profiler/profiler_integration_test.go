package profiler_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envdiff/internal/diff"
	"github.com/your-org/envdiff/internal/parser"
	"github.com/your-org/envdiff/internal/profiler"
)

func writeEnv(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeEnv: %v", err)
	}
	return path
}

func TestProfile_Integration_ParseAndProfile(t *testing.T) {
	pathA := writeEnv(t, ".env.a", "APP=foo\nDB_URL=postgres://a\nSECRET=s3cr3t\n")
	pathB := writeEnv(t, ".env.b", "APP=foo\nDB_URL=postgres://b\n")

	mapA, err := parser.ParseFile(pathA)
	if err != nil {
		t.Fatalf("parse A: %v", err)
	}
	mapB, err := parser.ParseFile(pathB)
	if err != nil {
		t.Fatalf("parse B: %v", err)
	}

	results := diff.Compare(mapA, mapB)
	p := profiler.Compute(results, profiler.DefaultOptions())

	if p.Total != 3 {
		t.Errorf("Total: want 3, got %d", p.Total)
	}
	if p.Conflicts != 1 {
		t.Errorf("Conflicts: want 1, got %d", p.Conflicts)
	}
	if p.Missing != 1 {
		t.Errorf("Missing: want 1, got %d", p.Missing)
	}

	// Ensure the profile is JSON-serialisable.
	raw, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	var roundTrip profiler.Profile
	if err := json.Unmarshal(raw, &roundTrip); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if roundTrip.Health != p.Health {
		t.Errorf("round-trip Health mismatch: %f vs %f", roundTrip.Health, p.Health)
	}
}
