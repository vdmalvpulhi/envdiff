package exporter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/exporter"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.Missing, BaseValue: "localhost"},
		{Key: "API_KEY", Status: diff.Conflict, BaseValue: "abc123", OtherValue: "xyz789"},
		{Key: "PORT", Status: diff.Missing, BaseValue: "8080"},
	}
}

func TestWrite_EnvFormat_Default(t *testing.T) {
	var buf strings.Builder
	err := exporter.Write(&buf, sampleResults(), exporter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT=8080 in output, got:\n%s", out)
	}
	if !strings.Contains(out, "# CONFLICT: API_KEY") {
		t.Errorf("expected conflict comment for API_KEY, got:\n%s", out)
	}
}

func TestWrite_ShellFormat(t *testing.T) {
	var buf strings.Builder
	opts := exporter.Options{Format: exporter.FormatShell}
	err := exporter.Write(&buf, sampleResults(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export DB_HOST=localhost") {
		t.Errorf("expected export DB_HOST=localhost, got:\n%s", out)
	}
	if !strings.Contains(out, "export PORT=8080") {
		t.Errorf("expected export PORT=8080, got:\n%s", out)
	}
}

func TestWrite_OnlyMissing_SkipsConflicts(t *testing.T) {
	var buf strings.Builder
	opts := exporter.Options{Format: exporter.FormatEnv, OnlyMissing: true}
	err := exporter.Write(&buf, sampleResults(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "CONFLICT") {
		t.Errorf("expected no conflict lines when OnlyMissing=true, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected missing key DB_HOST in output")
	}
}

func TestWrite_QuotedValues(t *testing.T) {
	results := []diff.Result{
		{Key: "GREETING", Status: diff.Missing, BaseValue: "hello world"},
	}
	var buf strings.Builder
	err := exporter.Write(&buf, results, exporter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `GREETING="hello world"`) {
		t.Errorf("expected quoted value, got:\n%s", out)
	}
}

func TestWrite_EmptyResults(t *testing.T) {
	var buf strings.Builder
	err := exporter.Write(&buf, []diff.Result{}, exporter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for empty results, got: %q", buf.String())
	}
}

func TestWrite_ShellFormat_QuotedValues(t *testing.T) {
	results := []diff.Result{
		{Key: "GREETING", Status: diff.Missing, BaseValue: "hello world"},
	}
	var buf strings.Builder
	opts := exporter.Options{Format: exporter.FormatShell}
	err := exporter.Write(&buf, results, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `export GREETING="hello world"`) {
		t.Errorf("expected quoted shell export, got:\n%s", out)
	}
}
