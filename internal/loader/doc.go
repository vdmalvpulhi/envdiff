// Package loader handles discovery and loading of .env files from the
// filesystem. It supports loading individual files by path, matching
// multiple files via glob patterns, and a convenience helper for loading
// exactly two files for side-by-side comparison.
//
// Typical usage:
//
//	a, b, err := loader.LoadPair(".env.staging", ".env.production")
//	if err != nil {
//		log.Fatal(err)
//	}
//	results := diff.Compare(a.Vars, b.Vars)
//
// For multi-environment workflows, use LoadGlob to collect all matching
// files and iterate over them:
//
//	files, err := loader.LoadGlob(".env.*")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, f := range files {
//		fmt.Println(f.Path, f.Vars)
//	}
//
// # File Format
//
// Each .env file is expected to contain key=value pairs, one per line.
// Lines beginning with '#' are treated as comments and ignored. Blank
// lines are also ignored. Values may optionally be quoted with single
// or double quotes.
//
// Files that do not exist or cannot be parsed return a descriptive error
// that includes the offending file path for easier debugging.
package loader
