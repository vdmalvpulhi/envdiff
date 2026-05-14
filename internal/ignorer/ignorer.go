// Package ignorer provides functionality to exclude specific keys from diff results
// based on patterns, exact matches, or glob expressions.
package ignorer

import (
	"path"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options configures the ignorer behaviour.
type Options struct {
	// Keys is a list of exact key names to ignore.
	Keys []string
	// Patterns is a list of glob patterns (e.g. "AWS_*") to ignore.
	Patterns []string
	// Prefixes is a list of key prefixes to ignore (e.g. "INTERNAL_").
	Prefixes []string
}

// DefaultOptions returns an Options with no ignore rules applied.
func DefaultOptions() Options {
	return Options{}
}

// Apply filters out diff.Result entries whose keys match any ignore rule
// defined in opts. The original slice is not mutated.
func Apply(results []diff.Result, opts Options) []diff.Result {
	if len(opts.Keys) == 0 && len(opts.Patterns) == 0 && len(opts.Prefixes) == 0 {
		return results
	}

	keySet := buildKeySet(opts.Keys)

	out := make([]diff.Result, 0, len(results))
	for _, r := range results {
		if shouldIgnore(r.Key, keySet, opts.Patterns, opts.Prefixes) {
			continue
		}
		out = append(out, r)
	}
	return out
}

func shouldIgnore(key string, keySet map[string]struct{}, patterns, prefixes []string) bool {
	if _, ok := keySet[key]; ok {
		return true
	}
	for _, p := range patterns {
		if matched, _ := path.Match(p, key); matched {
			return true
		}
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

func buildKeySet(keys []string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}
