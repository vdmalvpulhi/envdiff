// Package profiler computes aggregate statistics over a slice of diff.Result
// values, providing a high-level "health" view of how well two .env files
// align.
//
// Usage:
//
//	results := diff.Compare(a, b)
//	profile := profiler.Compute(results, profiler.DefaultOptions())
//	fmt.Printf("Health: %.1f%%  Missing: %d  Conflicts: %d\n",
//		profile.Health, profile.Missing, profile.Conflicts)
//
// The Profile struct is JSON-serialisable, making it easy to embed in
// structured report pipelines.
package profiler
