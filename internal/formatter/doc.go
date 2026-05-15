// Package formatter renders diff.Result slices into human-readable strings.
//
// It supports multiple visual styles (plain, colored, compact) and can
// optionally include the actual values of each key to aid debugging.
//
// Basic usage:
//
//	opts := formatter.DefaultOptions()
//	opts.ShowValues = true
//	output := formatter.Format(results, opts)
//	fmt.Println(output)
//
// Colored output uses ANSI escape codes and is suitable for terminal display.
// Use StylePlain when writing to files or piping to other tools.
package formatter
