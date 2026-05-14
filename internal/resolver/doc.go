// Package resolver implements multi-source precedence resolution for
// environment variables.
//
// When working with several .env files (e.g. .env.local, .env.staging,
// .env.production), it is often necessary to determine the effective value
// of a variable after applying a defined precedence order. The resolver
// package handles this by accepting an ordered slice of named Sources and
// returning a Resolution for every unique key found across all sources.
//
// Precedence modes:
//
//	- Last-wins (default): the value from the final source that defines the
//	  key is used. Mirrors how most shell and Docker tooling behaves.
//	- First-wins: the value from the earliest source is preserved, subsequent
//	  definitions are recorded but ignored. Useful for local overrides.
//
// Conflict detection:
//
//	A Resolution is flagged as Conflict=true whenever the same key is defined
//	with different values across two or more sources, regardless of which
//	value ultimately wins. This allows callers to surface discrepancies even
//	when a clear winner exists.
//
// Integration:
//
//	ToDiffResults converts []Resolution to []diff.Result so the output can
//	flow into the existing filter, sorter, reporter, and exporter pipeline
//	without modification.
package resolver
