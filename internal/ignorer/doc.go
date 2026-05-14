// Package ignorer provides key-based filtering to suppress specific entries
// from diff results before they are reported or exported.
//
// Rules can be specified as:
//   - Exact key names (e.g. "DEBUG")
//   - Glob patterns   (e.g. "AWS_*", "*_SECRET")
//   - Key prefixes    (e.g. "INTERNAL_", "CI_")
//
// Example usage:
//
//	opts := ignorer.Options{
//	    Keys:     []string{"LEGACY_KEY"},
//	    Patterns: []string{"AWS_*"},
//	    Prefixes: []string{"INTERNAL_"},
//	}
//	filtered := ignorer.Apply(results, opts)
//
// Apply never mutates the input slice; it returns a new slice containing
// only the results that did not match any ignore rule.
package ignorer
