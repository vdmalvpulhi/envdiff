package auditor_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/auditor"
	"github.com/user/envdiff/internal/diff"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.Match, Values: []string{"localhost", "localhost"}},
		{Key: "DB_PASS", Status: diff.Missing, Values: []string{"", "secret"}},
		{Key: "API_URL", Status: diff.Conflict, Values: []string{"http://a", "http://b"}},
	}
}

func TestRecord_SetsTimestampAndActor(t *testing.T) {
	before := time.Now().UTC()
	entry := auditor.Record([]string{"a.env", "b.env"}, sampleResults(), auditor.Options{Actor: "alice"})
	if entry.Actor != "alice" {
		t.Errorf("expected actor alice, got %s", entry.Actor)
	}
	if entry.Timestamp.Before(before) {
		t.Error("timestamp should be at or after test start")
	}
	if len(entry.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(entry.Files))
	}
}

func TestRecord_DefaultActorFallback(t *testing.T) {
	t.Setenv("USER", "")
	opts := auditor.DefaultOptions()
	entry := auditor.Record(nil, nil, opts)
	if entry.Actor != "unknown" {
		t.Errorf("expected unknown actor, got %s", entry.Actor)
	}
}

func TestRecord_SummaryCorrect(t *testing.T) {
	entry := auditor.Record(nil, sampleResults(), auditor.Options{Actor: "bot"})
	s := entry.Summary
	if s.Total != 3 {
		t.Errorf("expected total 3, got %d", s.Total)
	}
	if s.Missing != 1 {
		t.Errorf("expected 1 missing, got %d", s.Missing)
	}
	if s.Conflicts != 1 {
		t.Errorf("expected 1 conflict, got %d", s.Conflicts)
	}
	if s.Matching != 1 {
		t.Errorf("expected 1 matching, got %d", s.Matching)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.json")

	entry := auditor.Record([]string{"x.env"}, sampleResults(), auditor.Options{Actor: "ci"})
	if err := auditor.Save(path, entry); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entries, err := auditor.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Actor != "ci" {
		t.Errorf("actor mismatch: %s", entries[0].Actor)
	}
	if entries[0].Summary.Total != 3 {
		t.Errorf("summary total mismatch: %d", entries[0].Summary.Total)
	}
}

func TestSave_AppendsEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.json")

	for _, actor := range []string{"alice", "bob"} {
		e := auditor.Record(nil, sampleResults(), auditor.Options{Actor: actor})
		if err := auditor.Save(path, e); err != nil {
			t.Fatalf("Save(%s): %v", actor, err)
		}
	}

	entries, err := auditor.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	_, err := auditor.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := auditor.Load("/nonexistent/audit.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRecord_ResultsPreserved(t *testing.T) {
	results := sampleResults()
	entry := auditor.Record(nil, results, auditor.Options{Actor: "x"})
	data, _ := json.Marshal(entry.Results)
	var back []diff.Result
	_ = json.Unmarshal(data, &back)
	if len(back) != len(results) {
		t.Errorf("results length mismatch after round-trip")
	}
}
