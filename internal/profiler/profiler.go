// Package profiler analyses a set of diff results and produces a
// statistical profile: key counts by status, value-length distribution,
// and a rough "health" percentage for the compared environment pair.
package profiler

import (
	"math"

	"github.com/your-org/envdiff/internal/diff"
)

// Profile holds aggregate statistics computed from a slice of diff results.
type Profile struct {
	Total     int            `json:"total"`
	Matching  int            `json:"matching"`
	Missing   int            `json:"missing"`
	Conflicts int            `json:"conflicts"`
	Health    float64        `json:"health_pct"`  // 0-100
	AvgLen    float64        `json:"avg_value_len"`
	ByStatus  map[string]int `json:"by_status"`
}

// DefaultOptions returns a zero-value Options (reserved for future flags).
type Options struct {
	// IncludeMatching controls whether matching keys contribute to AvgLen.
	IncludeMatching bool
}

func DefaultOptions() Options {
	return Options{IncludeMatching: true}
}

// Compute derives a Profile from results using the provided options.
func Compute(results []diff.Result, opts Options) Profile {
	p := Profile{
		ByStatus: make(map[string]int),
	}

	var totalLen int
	var lenCount int

	for _, r := range results {
		p.Total++
		switch r.Status {
		case diff.StatusMatch:
			p.Matching++
		case diff.StatusMissing:
			p.Missing++
		case diff.StatusConflict:
			p.Conflicts++
		}
		p.ByStatus[string(r.Status)]++

		if opts.IncludeMatching || r.Status != diff.StatusMatch {
			if r.Value1 != "" {
				totalLen += len(r.Value1)
				lenCount++
			}
			if r.Value2 != "" {
				totalLen += len(r.Value2)
				lenCount++
			}
		}
	}

	if p.Total > 0 {
		p.Health = math.Round(float64(p.Matching)/float64(p.Total)*10000) / 100
	}
	if lenCount > 0 {
		p.AvgLen = math.Round(float64(totalLen)/float64(lenCount)*100) / 100
	}

	return p
}
