// Package scorer assigns a numeric severity score to diff results,
// allowing callers to prioritize or threshold based on criticality.
package scorer

import "github.com/user/envdiff/internal/diff"

// Severity levels assigned to each result kind.
const (
	SeverityMissing  = 2
	SeverityConflict = 3
	SeverityExtra    = 1
)

// Result wraps a diff.Result with an attached Score.
type Result struct {
	diff.Result
	Score int
}

// Options controls scoring behaviour.
type Options struct {
	// MinScore filters out results below this threshold (inclusive).
	// Zero means include everything.
	MinScore int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{MinScore: 0}
}

// Score attaches a numeric severity score to each diff result and
// optionally filters results below opts.MinScore.
func Score(results []diff.Result, opts Options) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		s := scoreOne(r)
		if s < opts.MinScore {
			continue
		}
		out = append(out, Result{Result: r, Score: s})
	}
	return out
}

// TotalScore sums the scores of all results.
func TotalScore(results []Result) int {
	total := 0
	for _, r := range results {
		total += r.Score
	}
	return total
}

func scoreOne(r diff.Result) int {
	switch r.Status {
	case diff.StatusMissing:
		return SeverityMissing
	case diff.StatusConflict:
		return SeverityConflict
	case diff.StatusExtra:
		return SeverityExtra
	default:
		return 0
	}
}
