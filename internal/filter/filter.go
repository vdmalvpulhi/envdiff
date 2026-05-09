package filter

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options holds filtering criteria for diff results.
type Options struct {
	// OnlyMissing filters to show only missing variables (not conflicts).
	OnlyMissing bool
	// OnlyConflicts filters to show only conflicting variables (not missing).
	OnlyConflicts bool
	// Prefix filters results to only variables with the given prefix.
	Prefix string
	// Keys filters results to only the specified variable names.
	Keys []string
}

// Apply filters a slice of diff.Result according to the provided Options.
// If no filter criteria are set, the original slice is returned unchanged.
func Apply(results []diff.Result, opts Options) []diff.Result {
	if opts.isEmpty() {
		return results
	}

	keySet := buildKeySet(opts.Keys)
	filtered := make([]diff.Result, 0, len(results))

	for _, r := range results {
		if opts.OnlyMissing && r.Status == diff.Conflict {
			continue
		}
		if opts.OnlyConflicts && r.Status != diff.Conflict {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(r.Key, opts.Prefix) {
			continue
		}
		if len(keySet) > 0 {
			if _, ok := keySet[r.Key]; !ok {
				continue
			}
		}
		filtered = append(filtered, r)
	}

	return filtered
}

func (o Options) isEmpty() bool {
	return !o.OnlyMissing && !o.OnlyConflicts && o.Prefix == "" && len(o.Keys) == 0
}

func buildKeySet(keys []string) map[string]struct{} {
	if len(keys) == 0 {
		return nil
	}
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}
