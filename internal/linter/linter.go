// Package linter checks .env file entries for common style and correctness issues.
package linter

import (
	"fmt"
	"strings"
	"unicode"
)

// Severity indicates how serious a lint issue is.
type Severity string

const (
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Issue represents a single lint finding for a key.
type Issue struct {
	Key      string
	Message  string
	Severity Severity
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s: %s", i.Severity, i.Key, i.Message)
}

// Lint inspects a parsed env map and returns any issues found.
func Lint(env map[string]string) []Issue {
	var issues []Issue

	for key, value := range env {
		if key == "" {
			continue
		}

		// Keys should be uppercase.
		if key != strings.ToUpper(key) {
			issues = append(issues, Issue{
				Key:      key,
				Message:  "key should be uppercase",
				Severity: SeverityWarning,
			})
		}

		// Keys should not contain spaces.
		if strings.ContainsAny(key, " \t") {
			issues = append(issues, Issue{
				Key:      key,
				Message:  "key contains whitespace",
				Severity: SeverityError,
			})
		}

		// Keys must start with a letter or underscore.
		if len(key) > 0 && !unicode.IsLetter(rune(key[0])) && key[0] != '_' {
			issues = append(issues, Issue{
				Key:      key,
				Message:  "key must start with a letter or underscore",
				Severity: SeverityError,
			})
		}

		// Values should not contain unquoted leading/trailing whitespace.
		if value != strings.TrimSpace(value) {
			issues = append(issues, Issue{
				Key:      key,
				Message:  "value has leading or trailing whitespace",
				Severity: SeverityWarning,
			})
		}

		// Warn on empty values.
		if strings.TrimSpace(value) == "" {
			issues = append(issues, Issue{
				Key:      key,
				Message:  "value is empty",
				Severity: SeverityWarning,
			})
		}
	}

	return issues
}
