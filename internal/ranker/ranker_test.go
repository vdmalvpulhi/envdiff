package ranker_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/ranker"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMatch, ValueA: "localhost", ValueB: "localhost"},
		{Key: "SECRET_KEY", Status: diff.StatusConflict, ValueA: "abc", ValueB: "xyz"},
		{Key: "PORT", Status: diff.StatusMissingInB, ValueA: "8080", ValueB: ""},
		{Key: "API_URL", Status: diff.StatusMissingInA, ValueA: "", ValueB: "https://api.example.com"},
	}
}

func TestApply_DescendingOrder(t *testing.T) {
	opts := ranker.DefaultOptions()
	ranked := ranker.Apply(sampleResults(), opts)

	if len(ranked) == 0 {
		t.Fatal("expected ranked results, got none")
	}
	for i := 1; i < len(ranked); i++ {
		if ranked[i].Score > ranked[i-1].Score {
			t.Errorf("position %d has higher score than %d: %d > %d",
				i, i-1, ranked[i].Score, ranked[i-1].Score)
		}
	}
}

func TestApply_AscendingOrder(t *testing.T) {
	opts := ranker.DefaultOptions()
	opts.Descending = false
	ranked := ranker.Apply(sampleResults(), opts)

	for i := 1; i < len(ranked); i++ {
		if ranked[i].Score < ranked[i-1].Score {
			t.Errorf("position %d has lower score than %d: %d < %d",
				i, i-1, ranked[i].Score, ranked[i-1].Score)
		}
	}
}

func TestApply_MinScoreFilters(t *testing.T) {
	opts := ranker.DefaultOptions()
	opts.MinScore = 100 // unreachably high
	ranked := ranker.Apply(sampleResults(), opts)

	if len(ranked) != 0 {
		t.Errorf("expected 0 results with MinScore=100, got %d", len(ranked))
	}
}

func TestApply_EmptyInput(t *testing.T) {
	ranked := ranker.Apply([]diff.Result{}, ranker.DefaultOptions())
	if len(ranked) != 0 {
		t.Errorf("expected empty output, got %d", len(ranked))
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	input := sampleResults()
	copy := make([]diff.Result, len(input))
	for i, r := range input {
		copy[i] = r
	}

	ranker.Apply(input, ranker.DefaultOptions())

	for i := range input {
		if input[i].Key != copy[i].Key {
			t.Errorf("input mutated at index %d", i)
		}
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := ranker.DefaultOptions()
	if !opts.Descending {
		t.Error("expected Descending=true by default")
	}
	if opts.MinScore != 0 {
		t.Errorf("expected MinScore=0, got %d", opts.MinScore)
	}
}
