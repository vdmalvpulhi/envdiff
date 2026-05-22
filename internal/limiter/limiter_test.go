package limiter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/limiter"
)

func makeResults(statuses ...diff.Status) []diff.Result {
	out := make([]diff.Result, len(statuses))
	for i, s := range statuses {
		out[i] = diff.Result{Key: fmt.Sprintf("KEY_%d", i), Status: s}
	}
	return out
}

import "fmt"

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "A", Status: diff.StatusMatch},
		{Key: "B", Status: diff.StatusMissing},
		{Key: "C", Status: diff.StatusConflict},
		{Key: "D", Status: diff.StatusMatch},
		{Key: "E", Status: diff.StatusMissing},
		{Key: "F", Status: diff.StatusConflict},
	}
}

func TestApply_NoLimit_ReturnsAll(t *testing.T) {
	results := sampleResults()
	opts := limiter.DefaultOptions()
	opts.Max = 0

	out, meta := limiter.Apply(results, opts)
	if len(out) != len(results) {
		t.Fatalf("expected %d, got %d", len(results), len(out))
	}
	if meta.Truncated {
		t.Error("expected Truncated=false")
	}
}

func TestApply_WithinLimit_ReturnsAll(t *testing.T) {
	results := sampleResults()
	opts := limiter.DefaultOptions()
	opts.Max = 100

	out, meta := limiter.Apply(results, opts)
	if len(out) != len(results) {
		t.Fatalf("expected %d, got %d", len(results), len(out))
	}
	if meta.Truncated {
		t.Error("expected Truncated=false")
	}
}

func TestApply_Truncates(t *testing.T) {
	results := sampleResults()
	opts := limiter.DefaultOptions()
	opts.Max = 3

	out, meta := limiter.Apply(results, opts)
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
	if !meta.Truncated {
		t.Error("expected Truncated=true")
	}
	if meta.Total != 6 || meta.Returned != 3 {
		t.Errorf("unexpected meta: %+v", meta)
	}
}

func TestApply_Balanced_DistributesAcrossStatuses(t *testing.T) {
	results := sampleResults() // 2 match, 2 missing, 2 conflict
	opts := limiter.Options{Max: 3, Balance: true}

	out, meta := limiter.Apply(results, opts)
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
	if !meta.Truncated {
		t.Error("expected Truncated=true")
	}
	// Each bucket should contribute at most 1 item (3 buckets, max=3)
	seen := map[diff.Status]int{}
	for _, r := range out {
		seen[r.Status]++
	}
	for s, count := range seen {
		if count > 1 {
			t.Errorf("status %s appears %d times, want <=1", s, count)
		}
	}
}

func TestMeta_String(t *testing.T) {
	m := limiter.Meta{Total: 20, Returned: 5, Truncated: true}
	s := m.String()
	if s == "" {
		t.Error("expected non-empty string")
	}

	m2 := limiter.Meta{Total: 5, Returned: 5, Truncated: false}
	if m2.String() == "" {
		t.Error("expected non-empty string")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := limiter.DefaultOptions()
	if opts.Max <= 0 {
		t.Error("expected positive default Max")
	}
	if opts.Balance {
		t.Error("expected Balance=false by default")
	}
}
