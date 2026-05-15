package differ_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
)

func TestCompare_SameValues(t *testing.T) {
	before := map[string]string{"FOO": "bar", "BAZ": "qux"}
	after := map[string]string{"FOO": "bar", "BAZ": "qux"}

	results := differ.Compare(before, after, differ.DefaultOptions())
	for _, r := range results {
		if !r.Same {
			t.Errorf("expected key %q to be Same, got diff", r.Key)
		}
	}
}

func TestCompare_ChangedValue(t *testing.T) {
	before := map[string]string{"FOO": "old"}
	after := map[string]string{"FOO": "new"}

	results := differ.Compare(before, after, differ.DefaultOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Same {
		t.Error("expected Same=false for changed value")
	}
	if r.Before != "old" || r.After != "new" {
		t.Errorf("unexpected Before/After: %q / %q", r.Before, r.After)
	}
}

func TestCompare_AddedKey(t *testing.T) {
	before := map[string]string{}
	after := map[string]string{"NEW_KEY": "value"}

	results := differ.Compare(before, after, differ.DefaultOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Before != "" || results[0].After != "value" {
		t.Errorf("unexpected values: %+v", results[0])
	}
}

func TestCompare_RemovedKey(t *testing.T) {
	before := map[string]string{"OLD_KEY": "value"}
	after := map[string]string{}

	results := differ.Compare(before, after, differ.DefaultOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Before != "value" || results[0].After != "" {
		t.Errorf("unexpected values: %+v", results[0])
	}
}

func TestCompare_MaskSensitive(t *testing.T) {
	opts := differ.DefaultOptions()
	opts.MaskSensitive = true

	before := map[string]string{"DB_PASSWORD": "secret123", "APP_NAME": "myapp"}
	after := map[string]string{"DB_PASSWORD": "newpass", "APP_NAME": "myapp"}

	results := differ.Compare(before, after, opts)
	for _, r := range results {
		if r.Key == "DB_PASSWORD" {
			if r.Before != "***" || r.After != "***" {
				t.Errorf("expected masked values, got Before=%q After=%q", r.Before, r.After)
			}
		}
		if r.Key == "APP_NAME" && (r.Before == "***" || r.After == "***") {
			t.Error("non-sensitive key should not be masked")
		}
	}
}

func TestFormat_Added(t *testing.T) {
	d := differ.LineDiff{Key: "FOO", Before: "", After: "bar", Same: false}
	out := differ.Format(d)
	if !strings.HasPrefix(out, "+") {
		t.Errorf("expected '+' prefix for added key, got %q", out)
	}
}

func TestFormat_Removed(t *testing.T) {
	d := differ.LineDiff{Key: "FOO", Before: "bar", After: "", Same: false}
	out := differ.Format(d)
	if !strings.HasPrefix(out, "-") {
		t.Errorf("expected '-' prefix for removed key, got %q", out)
	}
}

func TestFormat_Changed(t *testing.T) {
	d := differ.LineDiff{Key: "FOO", Before: "old", After: "new", Same: false}
	out := differ.Format(d)
	if !strings.HasPrefix(out, "~") {
		t.Errorf("expected '~' prefix for changed key, got %q", out)
	}
}
