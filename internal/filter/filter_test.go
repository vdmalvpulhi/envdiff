package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/filter"
)

var sampleResults = []diff.Result{
	{Key: "DB_HOST", Status: diff.MissingInSecond, ValueA: "localhost"},
	{Key: "DB_PORT", Status: diff.MissingInFirst, ValueB: "5432"},
	{Key: "APP_SECRET", Status: diff.Conflict, ValueA: "abc", ValueB: "xyz"},
	{Key: "APP_ENV", Status: diff.Conflict, ValueA: "dev", ValueB: "prod"},
	{Key: "CACHE_URL", Status: diff.MissingInSecond, ValueA: "redis://localhost"},
}

func TestApply_NoFilter(t *testing.T) {
	out := filter.Apply(sampleResults, filter.Options{})
	if len(out) != len(sampleResults) {
		t.Errorf("expected %d results, got %d", len(sampleResults), len(out))
	}
}

func TestApply_OnlyMissing(t *testing.T) {
	out := filter.Apply(sampleResults, filter.Options{OnlyMissing: true})
	for _, r := range out {
		if r.Status == diff.Conflict {
			t.Errorf("expected no conflicts, got key %q with status %v", r.Key, r.Status)
		}
	}
	if len(out) != 3 {
		t.Errorf("expected 3 missing results, got %d", len(out))
	}
}

func TestApply_OnlyConflicts(t *testing.T) {
	out := filter.Apply(sampleResults, filter.Options{OnlyConflicts: true})
	for _, r := range out {
		if r.Status != diff.Conflict {
			t.Errorf("expected only conflicts, got key %q with status %v", r.Key, r.Status)
		}
	}
	if len(out) != 2 {
		t.Errorf("expected 2 conflict results, got %d", len(out))
	}
}

func TestApply_Prefix(t *testing.T) {
	out := filter.Apply(sampleResults, filter.Options{Prefix: "APP_"})
	if len(out) != 2 {
		t.Errorf("expected 2 results with APP_ prefix, got %d", len(out))
	}
	for _, r := range out {
		if r.Key != "APP_SECRET" && r.Key != "APP_ENV" {
			t.Errorf("unexpected key %q", r.Key)
		}
	}
}

func TestApply_Keys(t *testing.T) {
	out := filter.Apply(sampleResults, filter.Options{Keys: []string{"DB_HOST", "APP_ENV"}})
	if len(out) != 2 {
		t.Errorf("expected 2 results, got %d", len(out))
	}
}

func TestApply_PrefixAndOnlyConflicts(t *testing.T) {
	out := filter.Apply(sampleResults, filter.Options{Prefix: "APP_", OnlyConflicts: true})
	if len(out) != 2 {
		t.Errorf("expected 2 results, got %d", len(out))
	}
}

func TestApply_EmptyResults(t *testing.T) {
	out := filter.Apply([]diff.Result{}, filter.Options{OnlyMissing: true})
	if len(out) != 0 {
		t.Errorf("expected 0 results, got %d", len(out))
	}
}
