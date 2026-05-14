package renamer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/renamer"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMatch, Value1: "localhost", Value2: "localhost"},
		{Key: "DB_PORT", Status: diff.StatusConflict, Value1: "5432", Value2: "3306"},
		{Key: "API_KEY", Status: diff.StatusMissingInSecond, Value1: "abc", Value2: ""},
	}
}

func TestApply_NoRenames(t *testing.T) {
	opts := renamer.DefaultOptions()
	got := renamer.Apply(sampleResults(), opts)
	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d", len(got))
	}
	if got[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", got[0].Key)
	}
}

func TestApply_RenamesKey(t *testing.T) {
	opts := renamer.DefaultOptions()
	opts.Renames = map[string]string{"DB_HOST": "DATABASE_HOST"}
	got := renamer.Apply(sampleResults(), opts)
	if got[0].Key != "DATABASE_HOST" {
		t.Errorf("expected DATABASE_HOST, got %s", got[0].Key)
	}
	// Status and values should be preserved
	if got[0].Status != diff.StatusMatch {
		t.Errorf("expected StatusMatch, got %v", got[0].Status)
	}
}

func TestApply_CollisionDropsSecond(t *testing.T) {
	// Rename DB_PORT to DB_HOST so they collide
	opts := renamer.DefaultOptions()
	opts.Renames = map[string]string{"DB_PORT": "DB_HOST"}
	got := renamer.Apply(sampleResults(), opts)
	// DB_HOST already exists first, so renamed DB_PORT should be dropped
	if len(got) != 2 {
		t.Fatalf("expected 2 results after collision drop, got %d", len(got))
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	input := sampleResults()
	opts := renamer.DefaultOptions()
	opts.Renames = map[string]string{"DB_HOST": "DATABASE_HOST"}
	renamer.Apply(input, opts)
	if input[0].Key != "DB_HOST" {
		t.Errorf("input was mutated: got %s", input[0].Key)
	}
}

func TestBuildRenames_Valid(t *testing.T) {
	m, err := renamer.BuildRenames([]string{"OLD_KEY=NEW_KEY", "FOO=BAR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["OLD_KEY"] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY, got %s", m["OLD_KEY"])
	}
	if m["FOO"] != "BAR" {
		t.Errorf("expected BAR, got %s", m["FOO"])
	}
}

func TestBuildRenames_InvalidPair(t *testing.T) {
	_, err := renamer.BuildRenames([]string{"NODIVIDER"})
	if err == nil {
		t.Fatal("expected error for invalid pair")
	}
}

func TestBuildRenames_EmptyValue(t *testing.T) {
	_, err := renamer.BuildRenames([]string{"KEY="})
	if err == nil {
		t.Fatal("expected error for empty new key")
	}
}
