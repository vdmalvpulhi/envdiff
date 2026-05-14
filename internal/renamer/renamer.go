// Package renamer provides utilities for renaming keys across diff results,
// useful when migrating variable names between environments.
package renamer

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options controls renaming behaviour.
type Options struct {
	// Map of oldKey -> newKey renames to apply.
	Renames map[string]string
	// CaseSensitive controls whether key matching is case-sensitive.
	CaseSensitive bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Renames:       make(map[string]string),
		CaseSensitive: true,
	}
}

// Apply returns a new slice of diff.Result with keys renamed according to opts.
// Results whose keys do not appear in the rename map are returned unchanged.
// If two entries would collide after renaming, the first one wins.
func Apply(results []diff.Result, opts Options) []diff.Result {
	seen := make(map[string]struct{}, len(results))
	out := make([]diff.Result, 0, len(results))

	for _, r := range results {
		newKey := rename(r.Key, opts)
		if _, exists := seen[newKey]; exists {
			// skip collision — first entry wins
			continue
		}
		seen[newKey] = struct{}{}
		r.Key = newKey
		out = append(out, r)
	}
	return out
}

// BuildRenames parses a slice of "OLD=NEW" strings into a rename map.
func BuildRenames(pairs []string) (map[string]string, error) {
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, &InvalidPairError{Pair: p}
		}
		m[parts[0]] = parts[1]
	}
	return m, nil
}

func rename(key string, opts Options) string {
	lookup := key
	if !opts.CaseSensitive {
		lookup = strings.ToUpper(key)
	}
	if v, ok := opts.Renames[lookup]; ok {
		return v
	}
	return key
}

// InvalidPairError is returned when a rename pair cannot be parsed.
type InvalidPairError struct {
	Pair string
}

func (e *InvalidPairError) Error() string {
	return "renamer: invalid pair (expected OLD=NEW): " + e.Pair
}
