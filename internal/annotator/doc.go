// Package annotator enriches diff results with human-readable hints and
// metadata such as sensitivity classification and default-value availability.
//
// The primary entry point is [Annotate], which accepts a slice of diff results
// and an [Options] value, and returns a slice of [Annotation] values.
//
// Usage:
//
//	rs, _ := diff.Compare(a, b)
//	opts := annotator.DefaultOptions()
//	opts.Defaults = map[string]string{"PORT": "8080"}
//	annotations := annotator.Annotate(rs, opts)
//	for _, ann := range annotations {
//		fmt.Printf("%s: %s\n", ann.Key, ann.Hint)
//	}
//
// # Sensitive Detection
//
// Sensitive detection is based on case-insensitive substring matching against
// a configurable list of patterns (e.g. SECRET, TOKEN, PASSWORD). The default
// pattern list is available via [DefaultSensitivePatterns].
//
// # Default Values
//
// Defaults are supplied via a simple key→value map on [Options.Defaults] and
// appear in the generated hint text when a key is added or removed.
package annotator
