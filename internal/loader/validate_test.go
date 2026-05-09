package loader_test

import (
	"testing"

	"github.com/user/envdiff/internal/loader"
)

func makeEnvFile(path string, vars map[string]string) *loader.EnvFile {
	return &loader.EnvFile{Path: path, Vars: vars}
}

func TestValidate_NoIssues(t *testing.T) {
	ef := makeEnvFile(".env", map[string]string{"FOO": "bar", "BAZ": "qux"})
	err := loader.Validate(ef, loader.ValidationOptions{})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_DisallowEmpty(t *testing.T) {
	ef := makeEnvFile(".env", map[string]string{"FOO": "", "BAR": "ok"})
	err := loader.Validate(ef, loader.ValidationOptions{DisallowEmpty: true})
	if err == nil {
		t.Fatal("expected validation error for empty value, got nil")
	}
	ve, ok := err.(*loader.ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(ve.Issues))
	}
}

func TestValidate_RequiredKeys_AllPresent(t *testing.T) {
	ef := makeEnvFile(".env", map[string]string{"DB_URL": "postgres://...", "SECRET": "abc"})
	err := loader.Validate(ef, loader.ValidationOptions{
		RequiredKeys: []string{"DB_URL", "SECRET"},
	})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_RequiredKeys_Missing(t *testing.T) {
	ef := makeEnvFile(".env", map[string]string{"FOO": "bar"})
	err := loader.Validate(ef, loader.ValidationOptions{
		RequiredKeys: []string{"FOO", "MISSING_KEY"},
	})
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
	ve := err.(*loader.ValidationError)
	if len(ve.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d: %v", len(ve.Issues), ve.Issues)
	}
}

func TestValidate_MultipleIssues(t *testing.T) {
	ef := makeEnvFile(".env", map[string]string{"FOO": ""})
	err := loader.Validate(ef, loader.ValidationOptions{
		DisallowEmpty: true,
		RequiredKeys:  []string{"MUST_EXIST"},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	ve := err.(*loader.ValidationError)
	if len(ve.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d: %v", len(ve.Issues), ve.Issues)
	}
}

func TestValidationError_Message(t *testing.T) {
	ve := &loader.ValidationError{
		Path:   ".env.prod",
		Issues: []string{"key \"X\" has empty value"},
	}
	msg := ve.Error()
	if msg == "" {
		t.Error("expected non-empty error message")
	}
}
