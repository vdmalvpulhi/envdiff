// Package grouper organises diff results into named groups based on key prefix.
// Keys sharing a common prefix (e.g. "DB_", "AWS_") are collected together,
// making large comparisons easier to navigate.
package grouper

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Group holds a labelled collection of diff results.
type Group struct {
	Prefix  string
	Results []diff.Result
}

// Options controls grouper behaviour.
type Options struct {
	// MinGroupSize skips groups with fewer results than this value (0 = include all).
	MinGroupSize int
	// OtherLabel is the group name for keys that match no prefix.
	OtherLabel string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MinGroupSize: 0,
		OtherLabel:   "OTHER",
	}
}

// ByPrefix partitions results into groups using the provided prefixes.
// Each prefix is matched case-insensitively against the start of the key.
// Results that match no prefix are placed in the "other" group.
func ByPrefix(results []diff.Result, prefixes []string, opts Options) []Group {
	index := make(map[string]*Group, len(prefixes)+1)
	order := make([]string, 0, len(prefixes)+1)

	normalised := make([]string, len(prefixes))
	for i, p := range prefixes {
		normalised[i] = strings.ToUpper(p)
		if _, exists := index[normalised[i]]; !exists {
			index[normalised[i]] = &Group{Prefix: prefixes[i]}
			order = append(order, normalised[i])
		}
	}

	otherKey := strings.ToUpper(opts.OtherLabel)
	if _, exists := index[otherKey]; !exists {
		index[otherKey] = &Group{Prefix: opts.OtherLabel}
		order = append(order, otherKey)
	}

	for _, r := range results {
		upper := strings.ToUpper(r.Key)
		matched := false
		for _, p := range normalised {
			if strings.HasPrefix(upper, p) {
				index[p].Results = append(index[p].Results, r)
				matched = true
				break
			}
		}
		if !matched {
			index[otherKey].Results = append(index[otherKey].Results, r)
		}
	}

	out := make([]Group, 0, len(order))
	for _, key := range order {
		g := index[key]
		if opts.MinGroupSize > 0 && len(g.Results) < opts.MinGroupSize {
			continue
		}
		if len(g.Results) == 0 {
			continue
		}
		sort.Slice(g.Results, func(i, j int) bool {
			return g.Results[i].Key < g.Results[j].Key
		})
		out = append(out, *g)
	}
	return out
}
