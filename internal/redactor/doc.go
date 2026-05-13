// Package redactor provides value-masking for sensitive environment variables
// before they are written to any output (text, JSON, .env export, etc.).
//
// Keys are identified as sensitive by matching their name (case-insensitively)
// against a configurable list of substrings. A default list covering common
// patterns such as "password", "secret", "token", and "api_key" is provided
// via DefaultSensitivePatterns.
//
// # Usage
//
//	results := diff.Compare(base, override)
//	safe := redactor.Apply(results, redactor.Options{
//		ExtraKeys: []string{"MY_CUSTOM_KEY"},
//	})
//	// pass safe to report.Writer, exporter.Write, etc.
//
// Apply never mutates the input slice; it always returns a fresh copy.
package redactor
