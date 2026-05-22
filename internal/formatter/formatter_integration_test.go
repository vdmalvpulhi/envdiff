package formatter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/formatter"
	"github.com/user/envdiff/internal/parser"
	"os"
)

func writeEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envdiff-fmt-*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestFormatter_Integration_ParseAndFormat(t *testing.T) {
	file1 := writeEnvFile(t, "APP=myapp\nDB_HOST=localhost\nSECRET=abc\n")
	file2 := writeEnvFile(t, "APP=myapp\nDB_HOST=prod.db\nNEW_KEY=true\n")

	map1, err := parser.ParseFile(file1)
	if err != nil {
		t.Fatalf("parse file1: %v", err)
	}
	map2, err := parser.ParseFile(file2)
	if err != nil {
		t.Fatalf("parse file2: %v", err)
	}

	results := diff.Compare(map1, map2)

	opts := formatter.DefaultOptions()
	opts.ShowValues = true
	out := formatter.Format(results, opts)

	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST conflict in formatted output")
	}
	if !strings.Contains(out, "SECRET") {
		t.Error("expected SECRET missing key in formatted output")
	}
	if !strings.Contains(out, "NEW_KEY") {
		t.Error("expected NEW_KEY missing key in formatted output")
	}
	if !strings.Contains(out, "localhost") {
		t.Error("expected value 'localhost' shown in conflict")
	}
}

func TestFormatter_Integration_EmptyFiles(t *testing.T) {
	file1 := writeEnvFile(t, "")
	file2 := writeEnvFile(t, "")

	map1, _ := parser.ParseFile(file1)
	map2, _ := parser.ParseFile(file2)

	results := diff.Compare(map1, map2)
	out := formatter.Format(results, formatter.DefaultOptions())

	if out != "No differences found." {
		t.Errorf("expected no-diff message for empty files, got %q", out)
	}
}

func TestFormatter_Integration_IdenticalFiles(t *testing.T) {
	content := "APP=myapp\nDB_HOST=localhost\nSECRET=abc\n"
	file1 := writeEnvFile(t, content)
	file2 := writeEnvFile(t, content)

	map1, err := parser.ParseFile(file1)
	if err != nil {
		t.Fatalf("parse file1: %v", err)
	}
	map2, err := parser.ParseFile(file2)
	if err != nil {
		t.Fatalf("parse file2: %v", err)
	}

	results := diff.Compare(map1, map2)
	out := formatter.Format(results, formatter.DefaultOptions())

	if out != "No differences found." {
		t.Errorf("expected no-diff message for identical files, got %q", out)
	}
}
