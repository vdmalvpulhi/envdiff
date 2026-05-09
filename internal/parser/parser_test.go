package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDEBUG=false\nSECRET=abc123\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := EnvMap{
		"APP_ENV": "production",
		"DEBUG":   "false",
		"SECRET":  "abc123",
	}

	for k, v := range expected {
		if env[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, env[k])
		}
	}
}

func TestParseFile_CommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nKEY=value\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(env) != 1 {
		t.Errorf("expected 1 key, got %d", len(env))
	}
	if env["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", env["KEY"])
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SINGLE='hello world'` + "\n" + `DOUBLE="foo bar"` + "\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if env["SINGLE"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", env["SINGLE"])
	}
	if env["DOUBLE"] != "foo bar" {
		t.Errorf("expected 'foo bar', got %q", env["DOUBLE"])
	}
}

func TestParseFile_ExportPrefix(t *testing.T) {
	path := writeTempEnv(t, "export PORT=8080\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", env["PORT"])
	}
}

func TestParseFile_MissingFile(t *testing.T) {
	_, err := ParseFile("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
