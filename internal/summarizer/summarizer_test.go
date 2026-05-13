package summarizer_test

import (
	"testing"

	"github.com/yourusername/envdiff/internal/diff"
	"github.com/yourusername/envdiff/internal/summarizer"
)

func makeResults(statuses ...diff.Status) []diff.Result {
	results := make([]diff.Result, len(statuses))
	for i, s := range statuses {
		results[i] = diff.Result{Key: "KEY", Status: s}
	}
	return results
}

func TestCompute_Empty(t *testing.T) {
	sum := summarizer.Compute(nil)
	if sum.Total != 0 || sum.Missing != 0 || sum.Conflicts != 0 || sum.Matching != 0 {
		t.Errorf("expected all zeros, got %+v", sum.Stats)
	}
	if sum.Status != summarizer.Healthy {
		t.Errorf("expected Healthy, got %q", sum.Status)
	}
}

func TestCompute_AllMatching(t *testing.T) {
	sum := summarizer.Compute(makeResults(diff.Match, diff.Match, diff.Match))
	if sum.Total != 3 || sum.Matching != 3 {
		t.Errorf("unexpected stats: %+v", sum.Stats)
	}
	if sum.Status != summarizer.Healthy {
		t.Errorf("expected Healthy, got %q", sum.Status)
	}
}

func TestCompute_WithMissing(t *testing.T) {
	sum := summarizer.Compute(makeResults(diff.Match, diff.Missing))
	if sum.Missing != 1 || sum.Matching != 1 {
		t.Errorf("unexpected stats: %+v", sum.Stats)
	}
	if sum.Status != summarizer.Warning {
		t.Errorf("expected Warning, got %q", sum.Status)
	}
}

func TestCompute_WithConflicts(t *testing.T) {
	sum := summarizer.Compute(makeResults(diff.Match, diff.Missing, diff.Conflict))
	if sum.Conflicts != 1 || sum.Missing != 1 || sum.Matching != 1 {
		t.Errorf("unexpected stats: %+v", sum.Stats)
	}
	if sum.Status != summarizer.Critical {
		t.Errorf("expected Critical, got %q", sum.Status)
	}
}

func TestCompute_ConflictTakesPriorityOverMissing(t *testing.T) {
	sum := summarizer.Compute(makeResults(diff.Conflict, diff.Missing))
	if sum.Status != summarizer.Critical {
		t.Errorf("expected Critical when both conflict and missing present, got %q", sum.Status)
	}
}

func TestCompute_TotalCount(t *testing.T) {
	input := makeResults(diff.Match, diff.Match, diff.Missing, diff.Conflict, diff.Conflict)
	sum := summarizer.Compute(input)
	if sum.Total != 5 {
		t.Errorf("expected Total=5, got %d", sum.Total)
	}
}
