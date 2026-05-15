package formatter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/formatter"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", Status: diff.StatusMatch, Value1: "myapp", Value2: "myapp"},
		{Key: "DB_HOST", Status: diff.StatusConflict, Value1: "localhost", Value2: "prod.db"},
		{Key: "SECRET_KEY", Status: diff.StatusMissingInSecond, Value1: "abc123", Value2: ""},
		{Key: "NEW_FLAG", Status: diff.StatusMissingInFirst, Value1: "", Value2: "true"},
	}
}

func TestFormat_EmptyResults(t *testing.T) {
	out := formatter.Format(nil, formatter.DefaultOptions())
	if out != "No differences found." {
		t.Errorf("expected no-diff message, got %q", out)
	}
}

func TestFormat_PlainNoValues(t *testing.T) {
	opts := formatter.DefaultOptions()
	out := formatter.Format(sampleResults(), opts)

	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "~ ") {
		t.Error("expected conflict prefix '~ '")
	}
	if !strings.Contains(out, "- ") {
		t.Error("expected missing-in-second prefix '- '")
	}
	if !strings.Contains(out, "+ ") {
		t.Error("expected missing-in-first prefix '+ '")
	}
}

func TestFormat_ShowValues_Conflict(t *testing.T) {
	opts := formatter.DefaultOptions()
	opts.ShowValues = true
	out := formatter.Format(sampleResults(), opts)

	if !strings.Contains(out, "localhost") {
		t.Error("expected value1 'localhost' in conflict line")
	}
	if !strings.Contains(out, "prod.db") {
		t.Error("expected value2 'prod.db' in conflict line")
	}
}

func TestFormat_ShowValues_Missing(t *testing.T) {
	opts := formatter.DefaultOptions()
	opts.ShowValues = true
	out := formatter.Format(sampleResults(), opts)

	if !strings.Contains(out, "missing") {
		t.Error("expected 'missing' annotation in output")
	}
}

func TestFormat_ColoredStyle_ContainsEscapes(t *testing.T) {
	opts := formatter.Options{
		Style:      formatter.StyleColored,
		ShowValues: false,
		Indent:     "",
	}
	out := formatter.Format(sampleResults(), opts)
	if !strings.Contains(out, "\033[") {
		t.Error("expected ANSI escape codes in colored output")
	}
}

func TestFormat_CustomIndent(t *testing.T) {
	opts := formatter.DefaultOptions()
	opts.Indent = "    "
	out := formatter.Format(sampleResults(), opts)
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if len(line) > 0 && !strings.HasPrefix(line, "    ") {
			t.Errorf("line missing custom indent: %q", line)
		}
	}
}

func TestFormat_MatchStatus_NoSpecialPrefix(t *testing.T) {
	results := []diff.Result{
		{Key: "STABLE", Status: diff.StatusMatch, Value1: "v1", Value2: "v1"},
	}
	opts := formatter.DefaultOptions()
	out := formatter.Format(results, opts)
	if strings.Contains(out, "~ ") || strings.Contains(out, "+ ") || strings.Contains(out, "- ") {
		t.Error("match status should not have a change prefix")
	}
}
