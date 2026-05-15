// Package differ provides line-level diffing of individual env variable values,
// helping surface exactly what changed between two versions of the same key.
package differ

import (
	"fmt"
	"strings"
)

// LineDiff represents a single value comparison for a key.
type LineDiff struct {
	Key    string
	Before string
	After  string
	Same   bool
}

// Options controls how value diffs are rendered.
type Options struct {
	// MaskSensitive replaces values of sensitive keys with "***".
	MaskSensitive bool
	// SensitiveKeys is a set of key substrings considered sensitive.
	SensitiveKeys []string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaskSensitive: false,
		SensitiveKeys: []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PASS"},
	}
}

// Compare computes a LineDiff for each key present in either map.
func Compare(before, after map[string]string, opts Options) []LineDiff {
	seen := make(map[string]struct{})
	var results []LineDiff

	for k := range before {
		seen[k] = struct{}{}
	}
	for k := range after {
		seen[k] = struct{}{}
	}

	for k := range seen {
		bv := before[k]
		av := after[k]
		if opts.MaskSensitive && isSensitive(k, opts.SensitiveKeys) {
			if bv != "" {
				bv = "***"
			}
			if av != "" {
				av = "***"
			}
		}
		results = append(results, LineDiff{
			Key:    k,
			Before: bv,
			After:  av,
			Same:   before[k] == after[k],
		})
	}
	return results
}

// Format returns a human-readable representation of a LineDiff.
func Format(d LineDiff) string {
	if d.Same {
		return fmt.Sprintf("  %s = %s", d.Key, d.Before)
	}
	if d.Before == "" {
		return fmt.Sprintf("+ %s = %s", d.Key, d.After)
	}
	if d.After == "" {
		return fmt.Sprintf("- %s = %s", d.Key, d.Before)
	}
	return fmt.Sprintf("~ %s: %q -> %q", d.Key, d.Before, d.After)
}

func isSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}
