package resolver_test

import (
	"testing"

	"github.com/user/envdiff/internal/resolver"
)

func sources(pairs ...resolver.Source) []resolver.Source { return pairs }

func src(name string, vars map[string]string) resolver.Source {
	return resolver.Source{Name: name, Vars: vars}
}

func TestResolve_SingleSource(t *testing.T) {
	res, err := resolver.Resolve(sources(src("a.env", map[string]string{"FOO": "bar"})), resolver.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 || res[0].Key != "FOO" || res[0].Value != "bar" {
		t.Errorf("unexpected result: %+v", res)
	}
}

func TestResolve_LastWins(t *testing.T) {
	res, err := resolver.Resolve(sources(
		src("a.env", map[string]string{"FOO": "first"}),
		src("b.env", map[string]string{"FOO": "second"}),
	), resolver.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res[0].Value != "second" {
		t.Errorf("expected last-wins value 'second', got %q", res[0].Value)
	}
	if !res[0].Overridden {
		t.Error("expected Overridden=true")
	}
}

func TestResolve_FirstWins(t *testing.T) {
	opts := resolver.Options{FirstWins: true}
	res, err := resolver.Resolve(sources(
		src("a.env", map[string]string{"FOO": "first"}),
		src("b.env", map[string]string{"FOO": "second"}),
	), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res[0].Value != "first" {
		t.Errorf("expected first-wins value 'first', got %q", res[0].Value)
	}
}

func TestResolve_ConflictDetected(t *testing.T) {
	res, err := resolver.Resolve(sources(
		src("a.env", map[string]string{"DB_URL": "postgres://a"}),
		src("b.env", map[string]string{"DB_URL": "postgres://b"}),
	), resolver.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res[0].Conflict {
		t.Error("expected Conflict=true when values differ")
	}
}

func TestResolve_NoConflictWhenValuesMatch(t *testing.T) {
	res, err := resolver.Resolve(sources(
		src("a.env", map[string]string{"PORT": "8080"}),
		src("b.env", map[string]string{"PORT": "8080"}),
	), resolver.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res[0].Conflict {
		t.Error("expected Conflict=false when values are identical")
	}
}

func TestResolve_EmptySourcesError(t *testing.T) {
	_, err := resolver.Resolve(nil, resolver.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty sources")
	}
}

func TestToDiffResults_ConflictMapped(t *testing.T) {
	res, _ := resolver.Resolve(sources(
		src("a.env", map[string]string{"X": "1"}),
		src("b.env", map[string]string{"X": "2"}),
	), resolver.DefaultOptions())
	dr := resolver.ToDiffResults(res)
	if len(dr) != 1 {
		t.Fatalf("expected 1 diff result, got %d", len(dr))
	}
	if dr[0].Status != "conflict" {
		t.Errorf("expected status 'conflict', got %q", dr[0].Status)
	}
}
