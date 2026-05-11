// Package merger combines multiple parsed env maps into a single merged map,
// tracking the origin file for each key and flagging conflicts.
package merger

import (
	"fmt"

	"github.com/yourorg/envdiff/internal/diff"
)

// Entry holds a value and its source file.
type Entry struct {
	Value  string
	Source string
}

// Result is the output of merging N env files.
type Result struct {
	// Merged holds the winning value per key (first-seen wins by default).
	Merged map[string]Entry
	// Conflicts holds keys that appeared with differing values across sources.
	Conflicts map[string][]Entry
}

// Merge combines multiple named env maps. The name slice must align with maps.
// First-seen value wins; subsequent differing values are recorded as conflicts.
func Merge(names []string, maps []map[string]string) (*Result, error) {
	if len(names) != len(maps) {
		return nil, fmt.Errorf("merger: names length %d does not match maps length %d", len(names), len(maps))
	}

	result := &Result{
		Merged:    make(map[string]Entry),
		Conflicts: make(map[string][]Entry),
	}

	for i, m := range maps {
		source := names[i]
		for k, v := range m {
			existing, seen := result.Merged[k]
			if !seen {
				result.Merged[k] = Entry{Value: v, Source: source}
				continue
			}
			if existing.Value != v {
				// Record conflict; initialise slice with existing entry on first conflict.
				if _, hasConflict := result.Conflicts[k]; !hasConflict {
					result.Conflicts[k] = []Entry{existing}
				}
				result.Conflicts[k] = append(result.Conflicts[k], Entry{Value: v, Source: source})
			}
		}
	}

	return result, nil
}

// ToDiffResults converts a merger Result into the canonical diff.Result slice
// so it can be fed into existing report/filter/sorter pipelines.
func ToDiffResults(r *Result) []diff.Result {
	out := []diff.Result{}
	for k, entries := range r.Conflicts {
		out = append(out, diff.Result{
			Key:    k,
			Status: diff.Conflict,
			ValueA: entries[0].Value,
			ValueB: entries[len(entries)-1].Value,
		})
	}
	return out
}
