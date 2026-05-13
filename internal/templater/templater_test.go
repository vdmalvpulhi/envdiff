package templater_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/templater"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.Conflict, Value1: "localhost", Value2: "db.prod.example.com"},
		{Key: "SECRET_KEY", Status: diff.Missing, Value1: "abc123", Value2: ""},
		{Key: "API_URL", Status: diff.Missing, Value1: "", Value2: "https://api.example.com"},
	}
}

func TestWrite_EmptyResults(t *testing.T) {
	var buf strings.Builder
	err := templater.Write(&buf, nil, templater.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "empty") {
		t.Errorf("expected empty template message, got: %s", buf.String())
	}
}

func TestWrite_AllKeys_Sorted(t *testing.T) {
	var buf strings.Builder
	opts := templater.DefaultOptions()
	opts.CommentValues = false

	err := templater.Write(&buf, sampleResults(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := nonEmptyLines(buf.String())
	if len(lines) != 3 {
		t.Fatalf("expected 3 key lines, got %d: %v", len(lines), lines)
	}
	if lines[0] != "API_URL=" || lines[1] != "DB_HOST=" || lines[2] != "SECRET_KEY=" {
		t.Errorf("keys not sorted correctly: %v", lines)
	}
}

func TestWrite_WithComments(t *testing.T) {
	var buf strings.Builder
	err := templater.Write(&buf, sampleResults(), templater.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[conflict]") {
		t.Error("expected conflict comment in output")
	}
	if !strings.Contains(out, "[missing]") {
		t.Error("expected missing comment in output")
	}
	if !strings.Contains(out, "localhost") {
		t.Error("expected Value1 in conflict comment")
	}
}

func TestWrite_OnlyMissing(t *testing.T) {
	var buf strings.Builder
	opts := templater.Options{IncludeMissing: true, IncludeConflicts: false, CommentValues: false}

	err := templater.Write(&buf, sampleResults(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "DB_HOST") {
		t.Error("DB_HOST (conflict) should not appear when IncludeConflicts=false")
	}
	lines := nonEmptyLines(out)
	if len(lines) != 2 {
		t.Errorf("expected 2 missing keys, got %d", len(lines))
	}
}

func TestWrite_OnlyConflicts(t *testing.T) {
	var buf strings.Builder
	opts := templater.Options{IncludeMissing: false, IncludeConflicts: true, CommentValues: false}

	err := templater.Write(&buf, sampleResults(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := nonEmptyLines(buf.String())
	if len(lines) != 1 || lines[0] != "DB_HOST=" {
		t.Errorf("expected only DB_HOST=, got %v", lines)
	}
}

func nonEmptyLines(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		l = strings.TrimSpace(l)
		if l != "" && !strings.HasPrefix(l, "#") {
			out = append(out, l)
		}
	}
	return out
}
