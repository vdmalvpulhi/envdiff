// Package limiter truncates diff results to a configurable maximum count,
// optionally preserving a representative sample across status types.
package limiter

import (
	"fmt"

	"github.com/user/envdiff/internal/diff"
)

// DefaultOptions returns a Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Max:     50,
		Balance: false,
	}
}

// Options controls how Apply truncates results.
type Options struct {
	// Max is the maximum number of results to return. Zero means no limit.
	Max int
	// Balance distributes the limit evenly across each distinct status type
	// before applying the overall cap, so no single status dominates output.
	Balance bool
}

// Apply returns a (possibly truncated) copy of results according to opts.
// When Balance is true the quota is spread evenly across status buckets;
// any remaining slots are filled from the leftover in original order.
func Apply(results []diff.Result, opts Options) ([]diff.Result, Meta) {
	if opts.Max <= 0 || len(results) <= opts.Max {
		return results, Meta{Total: len(results), Returned: len(results), Truncated: false}
	}

	var out []diff.Result

	if opts.Balance {
		out = balanced(results, opts.Max)
	} else {
		out = results[:opts.Max]
	}

	return out, Meta{
		Total:     len(results),
		Returned:  len(out),
		Truncated: len(out) < len(results),
	}
}

// Meta describes what Apply did to the result set.
type Meta struct {
	Total     int
	Returned  int
	Truncated bool
}

// String returns a human-readable summary of the meta.
func (m Meta) String() string {
	if !m.Truncated {
		return fmt.Sprintf("%d results", m.Total)
	}
	return fmt.Sprintf("%d of %d results (truncated)", m.Returned, m.Total)
}

// balanced picks up to perBucket items from each status bucket, then fills
// remaining capacity from leftovers in original order.
func balanced(results []diff.Result, max int) []diff.Result {
	buckets := map[string][]diff.Result{}
	order := []string{}

	for _, r := range results {
		k := string(r.Status)
		if _, seen := buckets[k]; !seen {
			order = append(order, k)
		}
		buckets[k] = append(buckets[k], r)
	}

	perBucket := max / len(order)
	if perBucket < 1 {
		perBucket = 1
	}

	out := make([]diff.Result, 0, max)
	for _, k := range order {
		slice := buckets[k]
		if len(slice) > perBucket {
			slice = slice[:perBucket]
		}
		out = append(out, slice...)
	}

	if len(out) > max {
		out = out[:max]
	}
	return out
}
