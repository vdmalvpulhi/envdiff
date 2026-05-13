// Package annotator enriches diff results with human-readable hints and
// metadata such as sensitivity classification and default-value availability.
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
// Sensitive detection is based on substring matching against a configurable
// list of patterns (e.g. SECRET, TOKEN, PASSWORD). Defaults are supplied via
// a simple key→value map and appear in the generated hint text.
package annotator
