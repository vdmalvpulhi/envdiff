package diff_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
)

func TestCompare_NoChanges(t *testing.T) {
	first := map[string]string{"FOO": "bar", "BAZ": "qux"}
	second := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := diff.Compare(first, second)

	if result.HasDifferences() {
		t.Errorf("expected no differences, got %+v", result)
	}
}

func TestCompare_MissingInSecond(t *testing.T) {
	first := map[string]string{"FOO": "bar", "ONLY_FIRST": "val"}
	second := map[string]string{"FOO": "bar"}

	result := diff.Compare(first, second)

	if len(result.MissingInSecond) != 1 || result.MissingInSecond[0] != "ONLY_FIRST" {
		t.Errorf("expected ONLY_FIRST missing in second, got %v", result.MissingInSecond)
	}
	if len(result.MissingInFirst) != 0 {
		t.Errorf("expected no keys missing in first, got %v", result.MissingInFirst)
	}
}

func TestCompare_MissingInFirst(t *testing.T) {
	first := map[string]string{"FOO": "bar"}
	second := map[string]string{"FOO": "bar", "ONLY_SECOND": "val"}

	result := diff.Compare(first, second)

	if len(result.MissingInFirst) != 1 || result.MissingInFirst[0] != "ONLY_SECOND" {
		t.Errorf("expected ONLY_SECOND missing in first, got %v", result.MissingInFirst)
	}
}

func TestCompare_Conflicts(t *testing.T) {
	first := map[string]string{"FOO": "original", "BAR": "same"}
	second := map[string]string{"FOO": "changed", "BAR": "same"}

	result := diff.Compare(first, second)

	if len(result.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(result.Conflicts))
	}
	c := result.Conflicts[0]
	if c.Key != "FOO" || c.FirstValue != "original" || c.SecondValue != "changed" {
		t.Errorf("unexpected conflict: %+v", c)
	}
}

func TestCompare_Mixed(t *testing.T) {
	first := map[string]string{"A": "1", "B": "2", "C": "3"}
	second := map[string]string{"A": "1", "B": "changed", "D": "4"}

	result := diff.Compare(first, second)

	if !result.HasDifferences() {
		t.Fatal("expected differences")
	}
	if len(result.MissingInSecond) != 1 {
		t.Errorf("expected 1 missing in second (C), got %v", result.MissingInSecond)
	}
	if len(result.MissingInFirst) != 1 {
		t.Errorf("expected 1 missing in first (D), got %v", result.MissingInFirst)
	}
	if len(result.Conflicts) != 1 {
		t.Errorf("expected 1 conflict (B), got %v", result.Conflicts)
	}
}
