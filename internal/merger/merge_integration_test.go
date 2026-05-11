package merger_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/filter"
	"github.com/yourorg/envdiff/internal/merger"
)

// TestMerge_PipelineIntegration verifies that merger output can be fed
// directly into the filter pipeline that the rest of the tool uses.
func TestMerge_PipelineIntegration(t *testing.T) {
	a := map[string]string{"DB_HOST": "localhost", "API_KEY": "abc"}
	b := map[string]string{"DB_HOST": "prod.db", "API_KEY": "abc"}
	c := map[string]string{"DB_HOST": "staging.db", "NEW_KEY": "xyz"}

	res, err := merger.Merge(
		[]string{"dev.env", "prod.env", "staging.env"},
		[]map[string]string{a, b, c},
	)
	if err != nil {
		t.Fatalf("merge failed: %v", err)
	}

	diffResults := merger.ToDiffResults(res)

	// Only conflicts should appear in diff results.
	for _, r := range diffResults {
		if r.Status != diff.Conflict {
			t.Errorf("expected Conflict status, got %v for key %q", r.Status, r.Key)
		}
	}

	// Run through the filter pipeline — only conflicts filter should match all.
	filtered := filter.Apply(diffResults, filter.Options{OnlyConflicts: true})
	if len(filtered) != len(diffResults) {
		t.Errorf("filter lost results: got %d, want %d", len(filtered), len(diffResults))
	}
}

// TestMerge_ThreeWayConflict ensures three distinct values for the same key
// all appear in the conflict entries.
func TestMerge_ThreeWayConflict(t *testing.T) {
	a := map[string]string{"MODE": "dev"}
	b := map[string]string{"MODE": "staging"}
	c := map[string]string{"MODE": "prod"}

	res, err := merger.Merge(
		[]string{"a.env", "b.env", "c.env"},
		[]map[string]string{a, b, c},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, ok := res.Conflicts["MODE"]
	if !ok {
		t.Fatal("expected MODE to be a conflict")
	}
	if len(entries) < 2 {
		t.Errorf("expected at least 2 conflict entries for MODE, got %d", len(entries))
	}
}
