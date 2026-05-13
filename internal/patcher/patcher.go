// Package patcher writes missing keys from a diff result back into an .env file.
// It appends missing variables to the end of the target file, preserving
// existing content and optionally adding a section header comment.
package patcher

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options controls patcher behaviour.
type Options struct {
	// DryRun prints the patch to w instead of writing to the file.
	DryRun bool
	// AddSectionHeader prepends a comment block before the appended keys.
	AddSectionHeader bool
	// SectionComment is the text used when AddSectionHeader is true.
	// Defaults to "# Added by envdiff".
	SectionComment string
}

// DefaultOptions returns sensible defaults for Options.
func DefaultOptions() Options {
	return Options{
		AddSectionHeader: true,
		SectionComment:   "# Added by envdiff",
	}
}

// Patch appends keys that are missing in targetPath (status diff.Missing)
// using the values found in the results slice. When opts.DryRun is true the
// lines that would be written are returned as a string instead.
func Patch(targetPath string, results []diff.Result, opts Options) (string, error) {
	if opts.SectionComment == "" {
		opts.SectionComment = "# Added by envdiff"
	}

	var missing []diff.Result
	for _, r := range results {
		if r.Status == diff.Missing {
			missing = append(missing, r)
		}
	}
	if len(missing) == 0 {
		return "", nil
	}

	sort.Slice(missing, func(i, j int) bool {
		return missing[i].Key < missing[j].Key
	})

	var sb strings.Builder
	if opts.AddSectionHeader {
		sb.WriteString("\n")
		sb.WriteString(opts.SectionComment + "\n")
	}
	for _, r := range missing {
		val := r.Values[0]
		sb.WriteString(fmt.Sprintf("%s=%s\n", r.Key, val))
	}
	patch := sb.String()

	if opts.DryRun {
		return patch, nil
	}

	f, err := os.OpenFile(targetPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", fmt.Errorf("patcher: open %s: %w", targetPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(patch); err != nil {
		return "", fmt.Errorf("patcher: write %s: %w", targetPath, err)
	}
	return "", nil
}
