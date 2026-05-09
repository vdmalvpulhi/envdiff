package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/sorter"
)

var sample = []diff.Result{
	{Key: "PORT", Status: diff.StatusOK, ValueA: "8080", ValueB: "8080"},
	{Key: "API_KEY", Status: diff.StatusConflict, ValueA: "abc", ValueB: "xyz"},
	{Key: "DB_URL", Status: diff.StatusMissing, ValueA: "postgres://", ValueB: ""},
	{Key: "APP_ENV", Status: diff.StatusMissing, ValueA: "", ValueB: "production"},
}

func TestApply_SortByKey_Ascending(t *testing.T) {
	opts := sorter.Options{Field: sorter.SortByKey, Order: sorter.Ascending}
	out := sorter.Apply(sample, opts)

	expected := []string{"API_KEY", "APP_ENV", "DB_URL", "PORT"}
	for i, r := range out {
		if r.Key != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, r.Key, expected[i])
		}
	}
}

func TestApply_SortByKey_Descending(t *testing.T) {
	opts := sorter.Options{Field: sorter.SortByKey, Order: sorter.Descending}
	out := sorter.Apply(sample, opts)

	expected := []string{"PORT", "DB_URL", "APP_ENV", "API_KEY"}
	for i, r := range out {
		if r.Key != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, r.Key, expected[i])
		}
	}
}

func TestApply_SortByStatus_Ascending(t *testing.T) {
	opts := sorter.Options{Field: sorter.SortByStatus, Order: sorter.Ascending}
	out := sorter.Apply(sample, opts)

	// missing first, then conflict, then ok; ties broken by key
	if out[0].Status != diff.StatusMissing {
		t.Errorf("expected first result to be Missing, got %q", out[0].Status)
	}
	if out[2].Status != diff.StatusConflict {
		t.Errorf("expected third result to be Conflict, got %q", out[2].Status)
	}
	if out[3].Status != diff.StatusOK {
		t.Errorf("expected last result to be OK, got %q", out[3].Status)
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	input := []diff.Result{
		{Key: "Z", Status: diff.StatusOK},
		{Key: "A", Status: diff.StatusOK},
	}
	originalFirst := input[0].Key

	_ = sorter.Apply(input, sorter.DefaultOptions())

	if input[0].Key != originalFirst {
		t.Error("Apply mutated the input slice")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := sorter.DefaultOptions()
	if opts.Field != sorter.SortByKey {
		t.Errorf("expected default field SortByKey, got %q", opts.Field)
	}
	if opts.Order != sorter.Ascending {
		t.Errorf("expected default order Ascending, got %q", opts.Order)
	}
}
