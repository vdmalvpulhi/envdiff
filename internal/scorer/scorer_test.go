package scorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/scorer"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMissing, FileA: "localhost", FileB: ""},
		{Key: "API_KEY", Status: diff.StatusConflict, FileA: "aaa", FileB: "bbb"},
		{Key: "LOG_LEVEL", Status: diff.StatusExtra, FileA: "", FileB: "debug"},
		{Key: "PORT", Status: diff.StatusOK, FileA: "8080", FileB: "8080"},
	}
}

func TestScore_AllIncluded(t *testing.T) {
	results := sampleResults()
	scored := scorer.Score(results, scorer.DefaultOptions())

	// StatusOK gets score 0 and MinScore is 0, so all should be included.
	if len(scored) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(scored))
	}
}

func TestScore_CorrectValues(t *testing.T) {
	results := sampleResults()
	scored := scorer.Score(results, scorer.DefaultOptions())

	expected := map[string]int{
		"DB_HOST":   scorer.SeverityMissing,
		"API_KEY":   scorer.SeverityConflict,
		"LOG_LEVEL": scorer.SeverityExtra,
		"PORT":      0,
	}
	for _, r := range scored {
		want, ok := expected[r.Key]
		if !ok {
			t.Errorf("unexpected key %q", r.Key)
			continue
		}
		if r.Score != want {
			t.Errorf("key %q: got score %d, want %d", r.Key, r.Score, want)
		}
	}
}

func TestScore_MinScoreFilters(t *testing.T) {
	results := sampleResults()
	opts := scorer.Options{MinScore: 2}
	scored := scorer.Score(results, opts)

	// Only Missing (2) and Conflict (3) should pass.
	if len(scored) != 2 {
		t.Fatalf("expected 2 results with MinScore=2, got %d", len(scored))
	}
	for _, r := range scored {
		if r.Score < 2 {
			t.Errorf("result %q has score %d below MinScore", r.Key, r.Score)
		}
	}
}

func TestTotalScore(t *testing.T) {
	results := sampleResults()
	scored := scorer.Score(results, scorer.DefaultOptions())
	total := scorer.TotalScore(scored)

	// Missing=2 + Conflict=3 + Extra=1 + OK=0 = 6
	const want = 6
	if total != want {
		t.Errorf("TotalScore: got %d, want %d", total, want)
	}
}

func TestScore_EmptyInput(t *testing.T) {
	scored := scorer.Score(nil, scorer.DefaultOptions())
	if len(scored) != 0 {
		t.Errorf("expected empty slice, got %d results", len(scored))
	}
	if total := scorer.TotalScore(scored); total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
}
