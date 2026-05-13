// Package redactor masks sensitive values in diff results before output.
package redactor

import (
	"strings"

	"github.com/yourorg/envdiff/internal/diff"
)

const maskedValue = "***"

// DefaultSensitivePatterns contains common substrings that indicate a key holds
// a sensitive value (case-insensitive match against the key name).
var DefaultSensitivePatterns = []string{
	"password",
	"passwd",
	"secret",
	"token",
	"api_key",
	"apikey",
	"private_key",
	"auth",
	"credential",
}

// Options controls redaction behaviour.
type Options struct {
	// Patterns is the list of case-insensitive substrings matched against key
	// names. Defaults to DefaultSensitivePatterns when nil.
	Patterns []string
	// ExtraKeys lists exact key names (case-insensitive) that should always be
	// redacted regardless of pattern matching.
	ExtraKeys []string
}

// Apply returns a copy of results with sensitive values replaced by "***".
// Original results are never mutated.
func Apply(results []diff.Result, opts Options) []diff.Result {
	patterns := opts.Patterns
	if patterns == nil {
		patterns = DefaultSensitivePatterns
	}

	extraSet := make(map[string]struct{}, len(opts.ExtraKeys))
	for _, k := range opts.ExtraKeys {
		extraSet[strings.ToLower(k)] = struct{}{}
	}

	out := make([]diff.Result, len(results))
	for i, r := range results {
		if isSensitive(r.Key, patterns, extraSet) {
			r = maskResult(r)
		}
		out[i] = r
	}
	return out
}

func isSensitive(key string, patterns []string, extraKeys map[string]struct{}) bool {
	lower := strings.ToLower(key)
	if _, ok := extraKeys[lower]; ok {
		return true
	}
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

func maskResult(r diff.Result) diff.Result {
	if r.Value1 != "" {
		r.Value1 = maskedValue
	}
	if r.Value2 != "" {
		r.Value2 = maskedValue
	}
	return r
}
