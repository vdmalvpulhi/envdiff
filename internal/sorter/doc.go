// Package sorter provides utilities for ordering []diff.Result slices.
//
// Results can be sorted by key name or by status (missing → conflict → ok),
// in ascending or descending order. The Apply function never modifies the
// input slice; it always returns a new copy.
//
// Example usage:
//
//	import (
//		"github.com/user/envdiff/internal/sorter"
//	)
//
//	opts := sorter.Options{
//		Field: sorter.SortByStatus,
//		Order: sorter.Ascending,
//	}
//	sorted := sorter.Apply(results, opts)
package sorter
