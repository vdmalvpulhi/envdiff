package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/report"
)

func main() {
	format := flag.String("format", "text", "Output format: text or json")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [flags] <file1> <file2>\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	fileA, fileB := args[0], args[1]

	envA, err := parser.ParseFile(fileA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", fileA, err)
		os.Exit(1)
	}

	envB, err := parser.ParseFile(fileB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", fileB, err)
		os.Exit(1)
	}

	results := diff.Compare(envA, envB)

	fmt := report.Format(*format)
	writer := report.NewWriter(fmt, os.Stdout)
	if err := writer.Write(results); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}

	if len(results) > 0 {
		fmt.Fprintf(os.Stderr, "Summary: %s\n", report.Summary(results))
		os.Exit(2)
	}
}
