package merger_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/merger"
)

func TestMerge_SingleMap(t *testing.T) {
	m := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, err := merger.Merge([]string{"a.env"}, []map[string]string{m})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Merged) != 2 {
		t.Errorf("expected 2 merged keys, got %d", len(res.Merged))
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
}

func TestMerge_NoConflicts(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"BAR": "2"}
	res, err := merger.Merge([]string{"a.env", "b.env"}, []map[string]string{a, b})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Merged) != 2 {
		t.Errorf("expected 2 merged keys, got %d", len(res.Merged))
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts")
	}
}

func TestMerge_WithConflict(t *testing.T) {
	a := map[string]string{"FOO": "one", "SHARED": "alpha"}
	b := map[string]string{"FOO": "two", "SHARED": "alpha"}
	res, err := merger.Merge([]string{"a.env", "b.env"}, []map[string]string{a, b})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Conflicts["FOO"]; !ok {
		t.Error("expected FOO to be a conflict")
	}
	if _, ok := res.Conflicts["SHARED"]; ok {
		t.Error("SHARED has same value, should not be a conflict")
	}
}

func TestMerge_MismatchedLengths(t *testing.T) {
	_, err := merger.Merge([]string{"a.env"}, []map[string]string{{"A": "1"}, {"B": "2"}})
	if err == nil {
		t.Error("expected error for mismatched lengths")
	}
}

func TestMerge_FirstSeenWins(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}
	res, err := merger.Merge([]string{"a.env", "b.env"}, []map[string]string{a, b})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"].Value != "first" {
		t.Errorf("expected first-seen value 'first', got %q", res.Merged["KEY"].Value)
	}
	if res.Merged["KEY"].Source != "a.env" {
		t.Errorf("expected source 'a.env', got %q", res.Merged["KEY"].Source)
	}
}

func TestToDiffResults(t *testing.T) {
	a := map[string]string{"FOO": "one"}
	b := map[string]string{"FOO": "two"}
	res, _ := merger.Merge([]string{"a.env", "b.env"}, []map[string]string{a, b})
	diffRes := merger.ToDiffResults(res)
	if len(diffRes) != 1 {
		t.Fatalf("expected 1 diff result, got %d", len(diffRes))
	}
	if diffRes[0].Key != "FOO" {
		t.Errorf("expected key FOO, got %q", diffRes[0].Key)
	}
}
