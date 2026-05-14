package ignorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/ignorer"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", Status: diff.StatusMatch},
		{Key: "AWS_ACCESS_KEY", Status: diff.StatusConflict, ValueA: "old", ValueB: "new"},
		{Key: "AWS_SECRET", Status: diff.StatusMissingInB},
		{Key: "INTERNAL_FLAG", Status: diff.StatusMissingInA},
		{Key: "DATABASE_URL", Status: diff.StatusConflict, ValueA: "a", ValueB: "b"},
		{Key: "DEBUG", Status: diff.StatusMatch},
	}
}

func TestApply_NoRules_ReturnsAll(t *testing.T) {
	results := sampleResults()
	out := ignorer.Apply(results, ignorer.DefaultOptions())
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestApply_ExactKey(t *testing.T) {
	opts := ignorer.Options{Keys: []string{"DEBUG", "APP_NAME"}}
	out := ignorer.Apply(sampleResults(), opts)
	for _, r := range out {
		if r.Key == "DEBUG" || r.Key == "APP_NAME" {
			t.Errorf("key %q should have been ignored", r.Key)
		}
	}
	if len(out) != 4 {
		t.Fatalf("expected 4 results, got %d", len(out))
	}
}

func TestApply_GlobPattern(t *testing.T) {
	opts := ignorer.Options{Patterns: []string{"AWS_*"}}
	out := ignorer.Apply(sampleResults(), opts)
	for _, r := range out {
		if r.Key == "AWS_ACCESS_KEY" || r.Key == "AWS_SECRET" {
			t.Errorf("key %q should have been ignored by glob", r.Key)
		}
	}
	if len(out) != 4 {
		t.Fatalf("expected 4 results, got %d", len(out))
	}
}

func TestApply_Prefix(t *testing.T) {
	opts := ignorer.Options{Prefixes: []string{"INTERNAL_"}}
	out := ignorer.Apply(sampleResults(), opts)
	for _, r := range out {
		if r.Key == "INTERNAL_FLAG" {
			t.Errorf("key %q should have been ignored by prefix", r.Key)
		}
	}
	if len(out) != 5 {
		t.Fatalf("expected 5 results, got %d", len(out))
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	results := sampleResults()
	orig := make([]diff.Result, len(results))
	copy(orig, results)

	opts := ignorer.Options{Keys: []string{"DEBUG"}}
	ignorer.Apply(results, opts)

	if len(results) != len(orig) {
		t.Fatal("input slice was mutated")
	}
}

func TestApply_CombinedRules(t *testing.T) {
	opts := ignorer.Options{
		Keys:     []string{"APP_NAME"},
		Patterns: []string{"AWS_*"},
		Prefixes: []string{"INTERNAL_"},
	}
	out := ignorer.Apply(sampleResults(), opts)
	// Remaining: DATABASE_URL, DEBUG
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}
