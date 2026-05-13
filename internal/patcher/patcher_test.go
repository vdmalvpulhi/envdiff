package patcher_test

import (
	"os"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/patcher"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.Missing, Values: []string{"localhost"}},
		{Key: "API_KEY", Status: diff.Missing, Values: []string{"secret"}},
		{Key: "PORT", Status: diff.Conflict, Values: []string{"8080", "9090"}},
	}
}

func TestPatch_DryRun_OnlyMissing(t *testing.T) {
	opts := patcher.DefaultOptions()
	opts.DryRun = true

	out, err := patcher.Patch("", sampleResults(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "API_KEY=secret") {
		t.Errorf("expected API_KEY in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST in output, got:\n%s", out)
	}
	if strings.Contains(out, "PORT") {
		t.Errorf("conflict key PORT should not appear in patch output")
	}
}

func TestPatch_DryRun_SectionHeader(t *testing.T) {
	opts := patcher.DefaultOptions()
	opts.DryRun = true

	out, err := patcher.Patch("", sampleResults(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "# Added by envdiff") {
		t.Errorf("expected section header comment, got:\n%s", out)
	}
}

func TestPatch_DryRun_NoHeader(t *testing.T) {
	opts := patcher.DefaultOptions()
	opts.DryRun = true
	opts.AddSectionHeader = false

	out, err := patcher.Patch("", sampleResults(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "#") {
		t.Errorf("expected no comment when AddSectionHeader=false, got:\n%s", out)
	}
}

func TestPatch_DryRun_NoMissing(t *testing.T) {
	results := []diff.Result{
		{Key: "PORT", Status: diff.Conflict, Values: []string{"8080", "9090"}},
	}
	opts := patcher.DefaultOptions()
	opts.DryRun = true

	out, err := patcher.Patch("", results, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty output when no missing keys, got: %q", out)
	}
}

func TestPatch_WritesToFile(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = tmp.WriteString("EXISTING=yes\n")
	tmp.Close()

	opts := patcher.DefaultOptions()
	_, err = patcher.Patch(tmp.Name(), sampleResults(), opts)
	if err != nil {
		t.Fatalf("patch write failed: %v", err)
	}

	data, _ := os.ReadFile(tmp.Name())
	content := string(data)

	if !strings.Contains(content, "EXISTING=yes") {
		t.Error("original content should be preserved")
	}
	if !strings.Contains(content, "API_KEY=secret") {
		t.Error("missing key API_KEY should have been appended")
	}
}
