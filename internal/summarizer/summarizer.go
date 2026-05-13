// Package summarizer produces a concise statistical summary of diff results,
// including counts of missing, conflicting, and matching variables across
// compared environment files.
package summarizer

import "github.com/yourusername/envdiff/internal/diff"

// Stats holds aggregated counts derived from a slice of diff results.
type Stats struct {
	Total     int `json:"total"`
	Missing   int `json:"missing"`
	Conflicts int `json:"conflicts"`
	Matching  int `json:"matching"`
}

// HealthStatus represents an overall assessment of the diff.
type HealthStatus string

const (
	Healthy  HealthStatus = "healthy"
	Warning  HealthStatus = "warning"
	Critical HealthStatus = "critical"
)

// Summary combines Stats with a derived health status.
type Summary struct {
	Stats
	Status HealthStatus `json:"status"`
}

// Compute calculates a Summary from the provided diff results.
func Compute(results []diff.Result) Summary {
	s := Stats{Total: len(results)}

	for _, r := range results {
		switch r.Status {
		case diff.Missing:
			s.Missing++
		case diff.Conflict:
			s.Conflicts++
		case diff.Match:
			s.Matching++
		}
	}

	return Summary{
		Stats:  s,
		Status: deriveStatus(s),
	}
}

// deriveStatus maps stats to a HealthStatus value.
func deriveStatus(s Stats) HealthStatus {
	switch {
	case s.Conflicts > 0:
		return Critical
	case s.Missing > 0:
		return Warning
	default:
		return Healthy
	}
}
