package linter_test

import (
	"testing"

	"github.com/user/envdiff/internal/linter"
)

func TestLint_NoIssues(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/mydb",
		"PORT":         "8080",
		"_INTERNAL":    "value",
	}
	issues := linter.Lint(env)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := map[string]string{"db_host": "localhost"}
	issues := linter.Lint(env)
	if !hasIssue(issues, "db_host", linter.SeverityWarning) {
		t.Errorf("expected warning for lowercase key, got %v", issues)
	}
}

func TestLint_KeyWithWhitespace(t *testing.T) {
	env := map[string]string{"MY KEY": "value"}
	issues := linter.Lint(env)
	if !hasIssue(issues, "MY KEY", linter.SeverityError) {
		t.Errorf("expected error for key with whitespace, got %v", issues)
	}
}

func TestLint_KeyStartsWithDigit(t *testing.T) {
	env := map[string]string{"1BAD_KEY": "value"}
	issues := linter.Lint(env)
	if !hasIssue(issues, "1BAD_KEY", linter.SeverityError) {
		t.Errorf("expected error for key starting with digit, got %v", issues)
	}
}

func TestLint_ValueWithLeadingSpace(t *testing.T) {
	env := map[string]string{"MY_VAR": " value"}
	issues := linter.Lint(env)
	if !hasIssue(issues, "MY_VAR", linter.SeverityWarning) {
		t.Errorf("expected warning for leading whitespace in value, got %v", issues)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	env := map[string]string{"EMPTY_VAR": ""}
	issues := linter.Lint(env)
	if !hasIssue(issues, "EMPTY_VAR", linter.SeverityWarning) {
		t.Errorf("expected warning for empty value, got %v", issues)
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	env := map[string]string{
		"good_key": " bad value ",
	}
	issues := linter.Lint(env)
	// expect at least warning for lowercase key and whitespace value
	if len(issues) < 2 {
		t.Errorf("expected at least 2 issues, got %d: %v", len(issues), issues)
	}
}

func TestIssue_String(t *testing.T) {
	i := linter.Issue{Key: "FOO", Message: "something wrong", Severity: linter.SeverityError}
	s := i.String()
	if s == "" {
		t.Error("expected non-empty string from Issue.String()")
	}
}

// hasIssue checks whether issues contains an entry matching key and severity.
func hasIssue(issues []linter.Issue, key string, sev linter.Severity) bool {
	for _, i := range issues {
		if i.Key == key && i.Severity == sev {
			return true
		}
	}
	return false
}
