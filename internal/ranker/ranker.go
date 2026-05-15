// Package ranker orders diff results by their computed score,
// allowing callers to surface the highest-priority variables first.
package ranker

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/scorer"
)

// Options controls ranking behaviour.
type Options struct {
	// Descending sorts highest score first when true (default).
	Descending bool
	// MinScore excludes results with a score below this value.
	MinScore int
	// ScorerOptions are forwarded to the scorer when computing scores.
	ScorerOptions scorer.Options
}

// DefaultOptions returns sensible defaults: descending order, no minimum.
func DefaultOptions() Options {
	return Options{
		Descending:    true,
		MinScore:      0,
		ScorerOptions: scorer.DefaultOptions(),
	}
}

// Ranked pairs a DiffResult with its computed score.
type Ranked struct {
	Result diff.Result
	Score  int
}

// Apply scores each result, filters by MinScore, and returns them
// ordered according to opts.Descending.
func Apply(results []diff.Result, opts Options) []Ranked {
	scored := scorer.Score(results, opts.ScorerOptions)

	ranked := make([]Ranked, 0, len(scored))
	for _, s := range scored {
		if s.Score >= opts.MinScore {
			ranked = append(ranked, Ranked{
				Result: s.Result,
				Score:  s.Score,
			})
		}
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		if opts.Descending {
			return ranked[i].Score > ranked[j].Score
		}
		return ranked[i].Score < ranked[j].Score
	})

	return ranked
}
