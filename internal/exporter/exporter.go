// Package exporter provides functionality to export diff results
// to various file formats such as .env patch files or shell scripts.
package exporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the export output format.
type Format string

const (
	FormatEnv   Format = "env"
	FormatShell Format = "shell"
)

// Options controls export behaviour.
type Options struct {
	Format      Format
	OnlyMissing bool // export only missing keys (skip conflicts)
}

// DefaultOptions returns sensible export defaults.
func DefaultOptions() Options {
	return Options{
		Format:      FormatEnv,
		OnlyMissing: false,
	}
}

// Write serialises diff results to w using the provided options.
// For FormatEnv it emits KEY=VALUE lines.
// For FormatShell it emits `export KEY=VALUE` lines.
func Write(w io.Writer, results []diff.Result, opts Options) error {
	for _, r := range results {
		switch r.Status {
		case diff.Missing:
			if err := writeLine(w, opts.Format, r.Key, r.BaseValue); err != nil {
				return err
			}
		case diff.Conflict:
			if opts.OnlyMissing {
				continue
			}
			// Emit a commented conflict block so the user can resolve it.
			if _, err := fmt.Fprintf(w, "# CONFLICT: %s\n", r.Key); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, "# base:  %s=%s\n", r.Key, r.BaseValue); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, "# other: %s=%s\n", r.Key, r.OtherValue); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeLine(w io.Writer, f Format, key, value string) error {
	quoted := quoteIfNeeded(value)
	var line string
	switch f {
	case FormatShell:
		line = fmt.Sprintf("export %s=%s\n", key, quoted)
	default:
		line = fmt.Sprintf("%s=%s\n", key, quoted)
	}
	_, err := io.WriteString(w, line)
	return err
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t\n#") {
		return `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
	}
	return v
}
