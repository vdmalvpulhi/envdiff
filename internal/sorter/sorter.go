package sorter

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// SortField defines the field to sort diff results by.
type SortField string

const (
	SortByKey    SortField = "key"
	SortByStatus SortField = "status"
)

// Order defines ascending or descending sort order.
type Order string

const (
	Ascending  Order = "asc"
	Descending Order = "desc"
)

// Options configures how results are sorted.
type Options struct {
	Field SortField
	Order Order
}

// DefaultOptions returns sensible sort defaults.
func DefaultOptions() Options {
	return Options{
		Field: SortByKey,
		Order: Ascending,
	}
}

// Apply sorts a slice of diff.Result according to the given options.
// If opts.Field is unrecognised it falls back to SortByKey.
func Apply(results []diff.Result, opts Options) []diff.Result {
	out := make([]diff.Result, len(results))
	copy(out, results)

	less := byKey
	if opts.Field == SortByStatus {
		less = byStatus
	}

	sort.SliceStable(out, func(i, j int) bool {
		if opts.Order == Descending {
			return less(out[j], out[i])
		}
		return less(out[i], out[j])
	})

	return out
}

func byKey(a, b diff.Result) bool {
	return a.Key < b.Key
}

// statusRank assigns a numeric rank to each status for stable ordering.
func statusRank(r diff.Result) int {
	switch r.Status {
	case diff.StatusMissing:
		return 0
	case diff.StatusConflict:
		return 1
	case diff.StatusOK:
		return 2
	default:
		return 3
	}
}

func byStatus(a, b diff.Result) bool {
	ra, rb := statusRank(a), statusRank(b)
	if ra != rb {
		return ra < rb
	}
	return a.Key < b.Key
}
