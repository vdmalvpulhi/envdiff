// Package deduper removes duplicate diff results, keeping the first
// occurrence of each key when multiple results share the same key.
package deduper

import "github.com/user/envdiff/internal/diff"

// Options controls deduplication behaviour.
type Options struct {
	// PreferStatus, if set, causes the result with the given status to be
	// kept over the first-seen result when a duplicate key is encountered.
	// If the preferred status is not found, the first-seen result is kept.
	PreferStatus string
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{}
}

// Apply removes duplicate results for the same key from results.
// By default the first occurrence is retained. When opts.PreferStatus is
// set, a result whose status matches the preference wins over any other.
func Apply(results []diff.Result, opts Options) []diff.Result {
	if len(results) == 0 {
		return results
	}

	// Track insertion order so output is stable.
	order := make([]string, 0, len(results))
	seen := make(map[string]diff.Result, len(results))

	for _, r := range results {
		existing, dup := seen[r.Key]
		if !dup {
			order = append(order, r.Key)
			seen[r.Key] = r
			continue
		}

		// Replace if the new result matches the preferred status and the
		// existing one does not.
		if opts.PreferStatus != "" &&
			r.Status == opts.PreferStatus &&
			existing.Status != opts.PreferStatus {
			seen[r.Key] = r
		}
	}

	out := make([]diff.Result, 0, len(order))
	for _, k := range order {
		out = append(out, seen[k])
	}
	return out
}
