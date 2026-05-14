// Package snapshot provides utilities for capturing, persisting, and comparing
// point-in-time states of environment variable diff results.
//
// # Overview
//
// A Snapshot wraps a slice of diff.Result values together with a human-readable
// label and a UTC timestamp. Snapshots can be saved to and loaded from JSON
// files, making them suitable for audit trails, CI artefacts, and regression
// detection.
//
// # Usage
//
//	s := snapshot.Take("pre-deploy", results)
//	if err := snapshot.Save("./snapshots/pre-deploy.json", s); err != nil { ... }
//
//	loaded, err := snapshot.Load("./snapshots/pre-deploy.json")
//	changes := snapshot.Diff(loaded, newSnapshot)
//
// # Change Detection
//
// Diff compares two snapshots key-by-key and returns a slice of Change values
// describing any status transitions (e.g. "missing_in_second" → "conflict") as
// well as keys that are entirely new in the later snapshot.
package snapshot
