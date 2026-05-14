package baseline_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/baseline"
	"github.com/yourorg/envdiff/internal/diff"
)

func TestCompare_NoChanges(t *testing.T) {
	results := sampleResults()
	snap := &baseline.Snapshot{Results: results}
	d := baseline.Compare(snap, results)
	if len(d.New) != 0 || len(d.Resolved) != 0 || len(d.Changed) != 0 {
		t.Errorf("expected empty delta, got new=%d resolved=%d changed=%d",
			len(d.New), len(d.Resolved), len(d.Changed))
	}
}

func TestCompare_NewKey(t *testing.T) {
	base := sampleResults()
	snap := &baseline.Snapshot{Results: base}
	current := append(base, diff.Result{
		Key: "NEW_KEY", Status: diff.StatusMissing, Values: []string{"val", ""},
	})
	d := baseline.Compare(snap, current)
	if len(d.New) != 1 || d.New[0].Key != "NEW_KEY" {
		t.Errorf("expected 1 new key, got %v", d.New)
	}
}

func TestCompare_ResolvedKey(t *testing.T) {
	base := sampleResults()
	snap := &baseline.Snapshot{Results: base}
	// Remove last entry from current
	current := base[:len(base)-1]
	d := baseline.Compare(snap, current)
	if len(d.Resolved) != 1 || d.Resolved[0].Key != "SECRET" {
		t.Errorf("expected 1 resolved key, got %v", d.Resolved)
	}
}

func TestCompare_ChangedStatus(t *testing.T) {
	base := sampleResults()
	snap := &baseline.Snapshot{Results: base}
	current := []diff.Result{
		{Key: "APP_ENV", Status: diff.StatusMatch, Values: []string{"production", "production"}},
		{Key: "DB_URL", Status: diff.StatusMatch, Values: []string{"postgres://localhost", "postgres://localhost"}}, // was missing
		{Key: "SECRET", Status: diff.StatusConflict, Values: []string{"abc", "xyz"}},
	}
	d := baseline.Compare(snap, current)
	if len(d.Changed) != 1 || d.Changed[0].Key != "DB_URL" {
		t.Errorf("expected 1 changed key (DB_URL), got %v", d.Changed)
	}
}
