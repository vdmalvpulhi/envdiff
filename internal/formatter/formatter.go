// Package formatter provides utilities for formatting diff results
// into human-readable or structured representations beyond the default report.
package formatter

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Style controls the visual style of formatted output.
type Style string

const (
	StylePlain   Style = "plain"
	StyleColored Style = "colored"
	StyleCompact Style = "compact"
)

// Options configures formatter behaviour.
type Options struct {
	Style       Style
	ShowValues  bool
	Indent      string
}

// DefaultOptions returns sensible formatter defaults.
func DefaultOptions() Options {
	return Options{
		Style:      StylePlain,
		ShowValues: false,
		Indent:     "  ",
	}
}

// Format converts a slice of diff.Result into a formatted string.
func Format(results []diff.Result, opts Options) string {
	if len(results) == 0 {
		return "No differences found."
	}

	var sb strings.Builder
	for _, r := range results {
		line := formatLine(r, opts)
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}

func formatLine(r diff.Result, opts Options) string {
	prefix := statusPrefix(r.Status, opts.Style)
	if opts.ShowValues {
		return fmt.Sprintf("%s%s%s", opts.Indent, prefix, formatWithValues(r))
	}
	return fmt.Sprintf("%s%s%s", opts.Indent, prefix, r.Key)
}

func formatWithValues(r diff.Result) string {
	switch r.Status {
	case diff.StatusConflict:
		return fmt.Sprintf("%s (%q vs %q)", r.Key, r.Value1, r.Value2)
	case diff.StatusMissingInFirst:
		return fmt.Sprintf("%s (missing, second=%q)", r.Key, r.Value2)
	case diff.StatusMissingInSecond:
		return fmt.Sprintf("%s (first=%q, missing)", r.Key, r.Value1)
	default:
		return fmt.Sprintf("%s (%q)", r.Key, r.Value1)
	}
}

func statusPrefix(status diff.Status, style Style) string {
	switch status {
	case diff.StatusConflict:
		if style == StyleColored {
			return "\033[33m~ \033[0m"
		}
		return "~ "
	case diff.StatusMissingInFirst:
		if style == StyleColored {
			return "\033[32m+ \033[0m"
		}
		return "+ "
	case diff.StatusMissingInSecond:
		if style == StyleColored {
			return "\033[31m- \033[0m"
		}
		return "- "
	default:
		return "  "
	}
}
