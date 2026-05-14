// Package normalizer provides utilities for normalizing .env key-value pairs
// before comparison or export. It can trim whitespace, uppercase keys, and
// canonicalize values to reduce false-positive diffs.
package normalizer

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options controls which normalization steps are applied.
type Options struct {
	// UppercaseKeys converts all keys to uppercase before processing.
	UppercaseKeys bool
	// TrimValues strips leading and trailing whitespace from values.
	TrimValues bool
	// TrimKeys strips leading and trailing whitespace from keys.
	TrimKeys bool
	// CollapseEmptyValues treats values that are only whitespace as empty string.
	CollapseEmptyValues bool
}

// DefaultOptions returns a sensible default normalization configuration.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys:       false,
		TrimValues:          true,
		TrimKeys:            true,
		CollapseEmptyValues: true,
	}
}

// NormalizeMap applies normalization to a raw key-value map and returns a new map.
func NormalizeMap(m map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if opts.TrimKeys {
			k = strings.TrimSpace(k)
		}
		if opts.UppercaseKeys {
			k = strings.ToUpper(k)
		}
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		if opts.CollapseEmptyValues && strings.TrimSpace(v) == "" {
			v = ""
		}
		if k != "" {
			out[k] = v
		}
	}
	return out
}

// Apply normalizes the Key, BaseValue, and OtherValue fields of each diff.Result
// according to the provided Options and returns a new slice.
func Apply(results []diff.Result, opts Options) []diff.Result {
	out := make([]diff.Result, len(results))
	for i, r := range results {
		if opts.TrimKeys {
			r.Key = strings.TrimSpace(r.Key)
		}
		if opts.UppercaseKeys {
			r.Key = strings.ToUpper(r.Key)
		}
		if opts.TrimValues {
			r.BaseValue = strings.TrimSpace(r.BaseValue)
			r.OtherValue = strings.TrimSpace(r.OtherValue)
		}
		if opts.CollapseEmptyValues {
			if strings.TrimSpace(r.BaseValue) == "" {
				r.BaseValue = ""
			}
			if strings.TrimSpace(r.OtherValue) == "" {
				r.OtherValue = ""
			}
		}
		out[i] = r
	}
	return out
}
