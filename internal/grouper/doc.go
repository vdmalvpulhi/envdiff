// Package grouper partitions a slice of diff.Result values into named groups
// based on key prefixes (e.g. "DB_", "AWS_", "APP_").
//
// Usage:
//
//	results := diff.Compare(a, b)
//	groups := grouper.ByPrefix(results, []string{"DB_", "AWS_"}, grouper.DefaultOptions())
//	for _, g := range groups {
//		fmt.Printf("[%s] — %d keys\n", g.Prefix, len(g.Results))
//	}
//
// Keys that do not match any supplied prefix are collected into an "OTHER"
// group (configurable via Options.OtherLabel). Empty groups are omitted from
// the output unless MinGroupSize is zero.
package grouper
