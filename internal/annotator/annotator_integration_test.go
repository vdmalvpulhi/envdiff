package annotator_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/annotator"
	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

func mustParse(t *testing.T, content string) map[string]string {
	t.Helper()
	tmp := t.TempDir() + "/test.env"
	if err := writeFile(tmp, content); err != nil {
		t.Fatal(err)
	}
	m, err := parser.ParseFile(tmp)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func writeFile(path, content string) error {
	import_os_writeFile := func(name, data string) error {
		return nil // replaced below
	}
	_ = import_os_writeFile
	import "os"
	return os.WriteFile(path, []byte(content), 0o644)
}

func TestAnnotate_Integration_WithParser(t *testing.T) {
	base := mustParse(t, "DB_PASSWORD=secret\nPORT=8080\n")
	target := mustParse(t, "PORT=9090\n")

	results := diff.Compare(base, target)

	opts := annotator.Options{
		Defaults:          map[string]string{"DB_PASSWORD": "changeme"},
		SensitivePatterns: annotator.DefaultOptions().SensitivePatterns,
	}
	anns := annotator.Annotate(results, opts)

	byKey := make(map[string]annotator.Annotation, len(anns))
	for _, a := range anns {
		byKey[a.Key] = a
	}

	if !byKey["DB_PASSWORD"].IsSensitive {
		t.Error("DB_PASSWORD should be sensitive")
	}
	if !byKey["DB_PASSWORD"].HasDefault {
		t.Error("DB_PASSWORD should have a default")
	}
	if !strings.Contains(byKey["DB_PASSWORD"].Hint, "changeme") {
		t.Errorf("hint should mention default value, got: %s", byKey["DB_PASSWORD"].Hint)
	}
	if byKey["PORT"].Hint == "" {
		t.Error("PORT conflict should produce a hint")
	}
}
