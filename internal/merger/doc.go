// Package merger provides utilities for combining more than two parsed env
// maps into a single unified view.
//
// Unlike the diff package — which compares exactly two files — merger accepts
// an arbitrary number of named env maps and produces:
//
//   - A merged map where the first-seen value wins for each key.
//   - A conflicts map that records every key whose value differs across sources,
//     along with all observed values and their origin file names.
//
// The ToDiffResults helper converts merger output into the canonical
// []diff.Result slice so results can flow directly into the existing
// filter, sorter, and report pipelines without modification.
//
// Typical usage:
//
//	res, err := merger.Merge(names, maps)
//	if err != nil { ... }
//	results := merger.ToDiffResults(res)
package merger
