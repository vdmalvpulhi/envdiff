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
package loader
