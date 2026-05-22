package profiler_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/diff"
	"github.com/your-org/envdiff/internal/profiler"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", Value1: "myapp", Value2: "myapp", Status: diff.StatusMatch},
		{Key: "DEBUG", Value1: "true", Value2: "false", Status: diff.StatusConflict},
		{Key: "SECRET", Value1: "abc", Value2: "", Status: diff.StatusMissing},
		{Key: "PORT", Value1: "8080", Value2: "8080", Status: diff.StatusMatch},
	}
}

func TestCompute_Counts(t *testing.T) {
	p := profiler.Compute(sampleResults(), profiler.DefaultOptions())
	if p.Total != 4 {
		t.Errorf("Total: want 4, got %d", p.Total)
	}
	if p.Matching != 2 {
		t.Errorf("Matching: want 2, got %d", p.Matching)
	}
	if p.Missing != 1 {
		t.Errorf("Missing: want 1, got %d", p.Missing)
	}
	if p.Conflicts != 1 {
		t.Errorf("Conflicts: want 1, got %d", p.Conflicts)
	}
}

func TestCompute_Health(t *testing.T) {
	p := profiler.Compute(sampleResults(), profiler.DefaultOptions())
	// 2 matching out of 4 → 50%
	if p.Health != 50.0 {
		t.Errorf("Health: want 50.0, got %f", p.Health)
	}
}

func TestCompute_HealthAllMatch(t *testing.T) {
	results := []diff.Result{
		{Key: "A", Value1: "1", Value2: "1", Status: diff.StatusMatch},
	}
	p := profiler.Compute(results, profiler.DefaultOptions())
	if p.Health != 100.0 {
		t.Errorf("Health: want 100.0, got %f", p.Health)
	}
}

func TestCompute_Empty(t *testing.T) {
	p := profiler.Compute(nil, profiler.DefaultOptions())
	if p.Total != 0 || p.Health != 0 || p.AvgLen != 0 {
		t.Errorf("expected zero profile, got %+v", p)
	}
}

func TestCompute_ByStatus(t *testing.T) {
	p := profiler.Compute(sampleResults(), profiler.DefaultOptions())
	if p.ByStatus["match"] != 2 {
		t.Errorf("ByStatus[match]: want 2, got %d", p.ByStatus["match"])
	}
	if p.ByStatus["conflict"] != 1 {
		t.Errorf("ByStatus[conflict]: want 1, got %d", p.ByStatus["conflict"])
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := profiler.DefaultOptions()
	if !opts.IncludeMatching {
		t.Error("DefaultOptions.IncludeMatching should be true")
	}
}
