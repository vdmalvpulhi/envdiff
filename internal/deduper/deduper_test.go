package deduper_test

import (
	"testing"

	"github.com/user/envdiff/internal/deduper"
	"github.com/user/envdiff/internal/diff"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_HOST", Status: "match", Value1: "localhost", Value2: "localhost"},
		{Key: "APP_PORT", Status: "conflict", Value1: "8080", Value2: "9090"},
		{Key: "DB_URL", Status: "missing", Value1: "postgres://", Value2: ""},
	}
}

func TestApply_NoDuplicates_ReturnsAll(t *testing.T) {
	input := sampleResults()
	got := deduper.Apply(input, deduper.DefaultOptions())
	if len(got) != len(input) {
		t.Fatalf("expected %d results, got %d", len(input), len(got))
	}
}

func TestApply_RemovesDuplicateKey(t *testing.T) {
	input := []diff.Result{
		{Key: "APP_PORT", Status: "conflict", Value1: "8080", Value2: "9090"},
		{Key: "APP_PORT", Status: "match", Value1: "8080", Value2: "8080"},
	}
	got := deduper.Apply(input, deduper.DefaultOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	// First occurrence wins by default.
	if got[0].Status != "conflict" {
		t.Errorf("expected status 'conflict', got %q", got[0].Status)
	}
}

func TestApply_PreferStatus_WinsOverFirst(t *testing.T) {
	input := []diff.Result{
		{Key: "DB_URL", Status: "conflict", Value1: "a", Value2: "b"},
		{Key: "DB_URL", Status: "missing", Value1: "a", Value2: ""},
	}
	opts := deduper.Options{PreferStatus: "missing"}
	got := deduper.Apply(input, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].Status != "missing" {
		t.Errorf("expected status 'missing', got %q", got[0].Status)
	}
}

func TestApply_PreferStatus_NoMatch_KeepsFirst(t *testing.T) {
	input := []diff.Result{
		{Key: "KEY", Status: "conflict", Value1: "x", Value2: "y"},
		{Key: "KEY", Status: "conflict", Value1: "x", Value2: "z"},
	}
	opts := deduper.Options{PreferStatus: "missing"}
	got := deduper.Apply(input, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].Value2 != "y" {
		t.Errorf("expected first result to be kept, got Value2=%q", got[0].Value2)
	}
}

func TestApply_PreservesOrder(t *testing.T) {
	input := []diff.Result{
		{Key: "Z_KEY", Status: "match"},
		{Key: "A_KEY", Status: "missing"},
		{Key: "M_KEY", Status: "conflict"},
		{Key: "A_KEY", Status: "match"}, // duplicate
	}
	got := deduper.Apply(input, deduper.DefaultOptions())
	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d", len(got))
	}
	expectedOrder := []string{"Z_KEY", "A_KEY", "M_KEY"}
	for i, r := range got {
		if r.Key != expectedOrder[i] {
			t.Errorf("position %d: expected key %q, got %q", i, expectedOrder[i], r.Key)
		}
	}
}

func TestApply_EmptyInput(t *testing.T) {
	got := deduper.Apply(nil, deduper.DefaultOptions())
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	input := sampleResults()
	copy := make([]diff.Result, len(input))
	for i, r := range input {
		copy[i] = r
	}
	deduper.Apply(input, deduper.DefaultOptions())
	for i, r := range input {
		if r != copy[i] {
			t.Errorf("input mutated at index %d", i)
		}
	}
}
