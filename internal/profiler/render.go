package profiler

import (
	"fmt"
	"io"
	"strings"
)

// Render writes a human-readable summary of p to w.
// It is intentionally lightweight — for full structured output use the
// report or exporter packages.
func Render(w io.Writer, p Profile) error {
	lines := []string{
		fmt.Sprintf("Total keys : %d", p.Total),
		fmt.Sprintf("Matching   : %d", p.Matching),
		fmt.Sprintf("Missing    : %d", p.Missing),
		fmt.Sprintf("Conflicts  : %d", p.Conflicts),
		fmt.Sprintf("Health     : %.2f%%", p.Health),
		fmt.Sprintf("Avg val len: %.2f chars", p.AvgLen),
	}
	_, err := fmt.Fprintln(w, strings.Join(lines, "\n"))
	return err
}
