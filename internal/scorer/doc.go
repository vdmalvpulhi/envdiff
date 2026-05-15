// Package scorer assigns numeric scores to diff results based on
// their status and key characteristics. Scores can be used to
// prioritize which variables need the most attention, or to filter
// out low-signal results below a minimum threshold.
//
// Usage:
//
//	opts := scorer.DefaultOptions()
//	opts.MinScore = 2
//	scored := scorer.Score(results, opts)
//	total := scorer.TotalScore(scored)
package scorer
