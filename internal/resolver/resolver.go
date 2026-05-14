// Package resolver provides utilities for resolving the effective value
// of an environment variable across multiple env file sources, applying
// precedence rules and returning a merged view with source attribution.
package resolver

import (
	"fmt"

	"github.com/user/envdiff/internal/diff"
)

// Source represents a named env file and its parsed key-value pairs.
type Source struct {
	Name string
	Vars map[string]string
}

// Resolution holds the resolved value for a single key along with metadata
// about which source it came from and whether overrides occurred.
type Resolution struct {
	Key      string
	Value    string
	Source   string
	Overridden bool
	Conflict   bool
}

// Options controls resolver behaviour.
type Options struct {
	// FirstWins uses the first source that defines a key (default: last wins).
	FirstWins bool
}

// DefaultOptions returns sensible defaults (last-wins precedence).
func DefaultOptions() Options {
	return Options{FirstWins: false}
}

// Resolve applies precedence rules across the provided sources and returns
// one Resolution per unique key found across all sources.
func Resolve(sources []Source, opts Options) ([]Resolution, error) {
	if len(sources) == 0 {
		return nil, fmt.Errorf("resolver: at least one source is required")
	}

	type entry struct {
		value    string
		source   string
		seenFrom []string
	}

	index := map[string]*entry{}
	order := []string{}

	for _, src := range sources {
		for k, v := range src.Vars {
			if _, exists := index[k]; !exists {
				order = append(order, k)
				index[k] = &entry{value: v, source: src.Name, seenFrom: []string{src.Name}}
				continue
			}
			existing := index[k]
			existing.seenFrom = append(existing.seenFrom, src.Name)
			if !opts.FirstWins {
				existing.value = v
				existing.source = src.Name
			}
		}
	}

	results := make([]Resolution, 0, len(order))
	for _, k := range order {
		e := index[k]
		conflict := hasConflictingValues(k, sources)
		results = append(results, Resolution{
			Key:        k,
			Value:      e.value,
			Source:     e.source,
			Overridden: len(e.seenFrom) > 1,
			Conflict:   conflict,
		})
	}
	return results, nil
}

// ToDiffResults converts Resolutions back into the canonical diff.Result slice
// so downstream pipeline stages can consume resolver output unchanged.
func ToDiffResults(resolutions []Resolution) []diff.Result {
	out := make([]diff.Result, 0, len(resolutions))
	for _, r := range resolutions {
		status := diff.StatusMatch
		if r.Conflict {
			status = diff.StatusConflict
		}
		out = append(out, diff.Result{
			Key:    r.Key,
			Status: status,
		})
	}
	return out
}

func hasConflictingValues(key string, sources []Source) bool {
	var first string
	found := false
	for _, src := range sources {
		v, ok := src.Vars[key]
		if !ok {
			continue
		}
		if !found {
			first = v
			found = true
			continue
		}
		if v != first {
			return true
		}
	}
	return false
}
