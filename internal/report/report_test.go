package report_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/report"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMissingInSecond},
		{Key: "API_KEY", Status: diff.StatusMissingInFirst},
		{Key: "PORT", ValueA: "8080", ValueB: "9090", Status: diff.StatusConflict},
	}
}

func TestWriter_TextFormat_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(report.FormatText, &buf)
	if err := w.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-differences message, got: %q", buf.String())
	}
}

func TestWriter_TextFormat_WithResults(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(report.FormatText, &buf)
	if err := w.Write(sampleResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "MISSING_IN_SECOND") {
		t.Error("expected MISSING_IN_SECOND label")
	}
	if !strings.Contains(out, "MISSING_IN_FIRST") {
		t.Error("expected MISSING_IN_FIRST label")
	}
	if !strings.Contains(out, "CONFLICT") {
		t.Error("expected CONFLICT label")
	}
	if !strings.Contains(out, "8080") || !strings.Contains(out, "9090") {
		t.Error("expected conflict values in output")
	}
}

func TestWriter_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(report.FormatJSON, &buf)
	if err := w.Write(sampleResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed []diff.Result
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(parsed) != 3 {
		t.Errorf("expected 3 results, got %d", len(parsed))
	}
}

func TestWriter_JSONFormat_Empty(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(report.FormatJSON, &buf)
	if err := w.Write([]diff.Result{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[]") {
		t.Errorf("expected empty JSON array, got: %q", buf.String())
	}
}

func TestSummary(t *testing.T) {
	got := report.Summary(sampleResults())
	if !strings.Contains(got, "missing in first") {
		t.Errorf("summary missing 'missing in first': %q", got)
	}
	if !strings.Contains(got, "missing in second") {
		t.Errorf("summary missing 'missing in second': %q", got)
	}
	if !strings.Contains(got, "conflict") {
		t.Errorf("summary missing 'conflict': %q", got)
	}
}

func TestSummary_NoDiff(t *testing.T) {
	got := report.Summary(nil)
	if got != "No differences found." {
		t.Errorf("unexpected summary: %q", got)
	}
}
