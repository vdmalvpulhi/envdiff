package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Writer writes a diff result to an output stream.
type Writer struct {
	format Format
	out    io.Writer
}

// NewWriter creates a new report Writer.
func NewWriter(format Format, out io.Writer) *Writer {
	return &Writer{format: format, out: out}
}

// Write renders the diff results to the configured output.
func (w *Writer) Write(results []diff.Result) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(results)
	default:
		return w.writeText(results)
	}
}

func (w *Writer) writeText(results []diff.Result) error {
	if len(results) == 0 {
		fmt.Fprintln(w.out, "✓ No differences found.")
		return nil
	}
	for _, r := range results {
		var line string
		switch r.Status {
		case diff.StatusMissingInFirst:
			line = fmt.Sprintf("[MISSING_IN_FIRST]  %s", r.Key)
		case diff.StatusMissingInSecond:
			line = fmt.Sprintf("[MISSING_IN_SECOND] %s", r.Key)
		case diff.StatusConflict:
			line = fmt.Sprintf("[CONFLICT]          %s  (%q vs %q)", r.Key, r.ValueA, r.ValueB)
		}
		fmt.Fprintln(w.out, line)
	}
	return nil
}

func (w *Writer) writeJSON(results []diff.Result) error {
	if results == nil {
		results = []diff.Result{}
	}
	enc := json.NewEncoder(w.out)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

// Summary returns a human-readable summary string.
func Summary(results []diff.Result) string {
	if len(results) == 0 {
		return "No differences found."
	}
	var counts [3]int
	for _, r := range results {
		switch r.Status {
		case diff.StatusMissingInFirst:
			counts[0]++
		case diff.StatusMissingInSecond:
			counts[1]++
		case diff.StatusConflict:
			counts[2]++
		}
	}
	parts := []string{}
	if counts[0] > 0 {
		parts = append(parts, fmt.Sprintf("%d missing in first", counts[0]))
	}
	if counts[1] > 0 {
		parts = append(parts, fmt.Sprintf("%d missing in second", counts[1]))
	}
	if counts[2] > 0 {
		parts = append(parts, fmt.Sprintf("%d conflict(s)", counts[2]))
	}
	return strings.Join(parts, ", ")
}
