package normalizer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/normalizer"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "  db_host  ", Status: diff.StatusMatch, BaseValue: "  localhost  ", OtherValue: "  localhost  "},
		{Key: "api_key", Status: diff.StatusConflict, BaseValue: "abc123", OtherValue: "  xyz789  "},
		{Key: "EMPTY_VAL", Status: diff.StatusMatch, BaseValue: "   ", OtherValue: ""},
	}
}

func TestApply_TrimKeysAndValues(t *testing.T) {
	opts := normalizer.DefaultOptions()
	results := sampleResults()
	out := normalizer.Apply(results, opts)

	if out[0].Key != "db_host" {
		t.Errorf("expected trimmed key 'db_host', got %q", out[0].Key)
	}
	if out[0].BaseValue != "localhost" {
		t.Errorf("expected trimmed base value 'localhost', got %q", out[0].BaseValue)
	}
	if out[1].OtherValue != "xyz789" {
		t.Errorf("expected trimmed other value 'xyz789', got %q", out[1].OtherValue)
	}
}

func TestApply_CollapseEmptyValues(t *testing.T) {
	opts := normalizer.DefaultOptions()
	out := normalizer.Apply(sampleResults(), opts)

	if out[2].BaseValue != "" {
		t.Errorf("expected collapsed empty base value, got %q", out[2].BaseValue)
	}
}

func TestApply_UppercaseKeys(t *testing.T) {
	opts := normalizer.DefaultOptions()
	opts.UppercaseKeys = true
	out := normalizer.Apply(sampleResults(), opts)

	if out[0].Key != "DB_HOST" {
		t.Errorf("expected uppercase key 'DB_HOST', got %q", out[0].Key)
	}
	if out[1].Key != "API_KEY" {
		t.Errorf("expected uppercase key 'API_KEY', got %q", out[1].Key)
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	original := sampleResults()
	copy := sampleResults()
	normalizer.Apply(original, normalizer.DefaultOptions())

	for i := range original {
		if original[i].Key != copy[i].Key {
			t.Errorf("input mutated at index %d: key changed from %q to %q", i, copy[i].Key, original[i].Key)
		}
	}
}

func TestNormalizeMap_Basic(t *testing.T) {
	m := map[string]string{
		"  key1  ": "  value1  ",
		"key2":     "   ",
	}
	out := normalizer.NormalizeMap(m, normalizer.DefaultOptions())

	if v, ok := out["key1"]; !ok || v != "value1" {
		t.Errorf("expected key1=value1, got %q", v)
	}
	if v, ok := out["key2"]; !ok || v != "" {
		t.Errorf("expected key2 collapsed to empty, got %q", v)
	}
}

func TestNormalizeMap_UppercaseKeys(t *testing.T) {
	m := map[string]string{"db_pass": "secret"}
	opts := normalizer.DefaultOptions()
	opts.UppercaseKeys = true
	out := normalizer.NormalizeMap(m, opts)

	if _, ok := out["DB_PASS"]; !ok {
		t.Error("expected uppercased key 'DB_PASS' in output")
	}
	if _, ok := out["db_pass"]; ok {
		t.Error("expected original lowercase key to be absent")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := normalizer.DefaultOptions()
	if opts.UppercaseKeys {
		t.Error("expected UppercaseKeys to be false by default")
	}
	if !opts.TrimValues {
		t.Error("expected TrimValues to be true by default")
	}
	if !opts.TrimKeys {
		t.Error("expected TrimKeys to be true by default")
	}
	if !opts.CollapseEmptyValues {
		t.Error("expected CollapseEmptyValues to be true by default")
	}
}
