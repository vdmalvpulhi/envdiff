package redactor_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/redactor"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", Status: diff.StatusConflict, Value1: "myapp", Value2: "otherapp"},
		{Key: "DB_PASSWORD", Status: diff.StatusConflict, Value1: "hunter2", Value2: "s3cr3t"},
		{Key: "API_KEY", Status: diff.StatusMissingInSecond, Value1: "abc123", Value2: ""},
		{Key: "AUTH_TOKEN", Status: diff.StatusMissingInFirst, Value1: "", Value2: "tok_xyz"},
		{Key: "PORT", Status: diff.StatusConflict, Value1: "8080", Value2: "9090"},
	}
}

func TestApply_NonSensitiveUnchanged(t *testing.T) {
	results := sampleResults()
	out := redactor.Apply(results, redactor.Options{})

	for _, r := range out {
		if r.Key == "APP_NAME" || r.Key == "PORT" {
			if r.Value1 == "***" || r.Value2 == "***" {
				t.Errorf("key %s should not be redacted", r.Key)
			}
		}
	}
}

func TestApply_SensitiveValuesAreMasked(t *testing.T) {
	out := redactor.Apply(sampleResults(), redactor.Options{})

	keys := map[string]diff.Result{}
	for _, r := range out {
		keys[r.Key] = r
	}

	if keys["DB_PASSWORD"].Value1 != "***" || keys["DB_PASSWORD"].Value2 != "***" {
		t.Error("DB_PASSWORD values should be masked")
	}
	if keys["API_KEY"].Value1 != "***" {
		t.Error("API_KEY Value1 should be masked")
	}
	if keys["AUTH_TOKEN"].Value2 != "***" {
		t.Error("AUTH_TOKEN Value2 should be masked")
	}
}

func TestApply_ExtraKeys(t *testing.T) {
	opts := redactor.Options{
		Patterns:  []string{},
		ExtraKeys: []string{"APP_NAME"},
	}
	out := redactor.Apply(sampleResults(), opts)
	for _, r := range out {
		if r.Key == "APP_NAME" {
			if r.Value1 != "***" || r.Value2 != "***" {
				t.Error("APP_NAME should be redacted via ExtraKeys")
			}
		}
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	orig := sampleResults()
	copy := make([]diff.Result, len(orig))
	for i, r := range orig {
		copy[i] = r
	}

	redactor.Apply(orig, redactor.Options{})

	for i, r := range orig {
		if r.Value1 != copy[i].Value1 || r.Value2 != copy[i].Value2 {
			t.Errorf("input mutated at index %d", i)
		}
	}
}

func TestApply_CustomPatterns(t *testing.T) {
	opts := redactor.Options{
		Patterns: []string{"port"},
	}
	out := redactor.Apply(sampleResults(), opts)
	for _, r := range out {
		if r.Key == "PORT" && r.Value1 != "***" {
			t.Error("PORT should be masked with custom pattern 'port'")
		}
	}
}
