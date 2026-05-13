// Package templater generates .env template files from diff results,
// emitting keys with empty values for variables that are missing or conflicting.
package templater

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options controls which result types are included in the template output.
type Options struct {
	// IncludeMissing includes keys that are missing in one of the environments.
	IncludeMissing bool
	// IncludeConflicts includes keys whose values differ between environments.
	IncludeConflicts bool
	// CommentValues when true writes the observed values as comments above each key.
	CommentValues bool
}

// DefaultOptions returns sensible defaults: include all result types with comments.
func DefaultOptions() Options {
	return Options{
		IncludeMissing:   true,
		IncludeConflicts: true,
		CommentValues:    true,
	}
}

// Write renders a .env template to w based on the provided diff results and options.
// Keys are emitted in alphabetical order with empty values so the template can be
// filled in by the operator.
func Write(w io.Writer, results []diff.Result, opts Options) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "# No differences found — template is empty.")
		return err
	}

	filtered := make([]diff.Result, 0, len(results))
	for _, r := range results {
		switch r.Status {
		case diff.Missing:
			if opts.IncludeMissing {
				filtered = append(filtered, r)
			}
		case diff.Conflict:
			if opts.IncludeConflicts {
				filtered = append(filtered, r)
			}
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Key < filtered[j].Key
	})

	for _, r := range filtered {
		if opts.CommentValues {
			comment := buildComment(r)
			if _, err := fmt.Fprintln(w, comment); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "%s=\n", r.Key); err != nil {
			return err
		}
	}
	return nil
}

func buildComment(r diff.Result) string {
	switch r.Status {
	case diff.Missing:
		val := r.Value1
		if val == "" {
			val = r.Value2
		}
		return fmt.Sprintf("# [missing] known value: %s", strings.TrimSpace(val))
	case diff.Conflict:
		return fmt.Sprintf("# [conflict] file1=%s | file2=%s", r.Value1, r.Value2)
	default:
		return "# [unknown]"
	}
}
