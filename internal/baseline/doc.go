// Package baseline provides save/load/compare functionality for envdiff
// baseline snapshots.
//
// A baseline is a point-in-time JSON snapshot of diff results. It can be
// persisted to disk and later reloaded to compute a Delta — describing
// which keys are newly introduced, which have been resolved, and which
// have changed status since the snapshot was taken.
//
// Typical usage:
//
//	// Save a baseline
//	if err := baseline.Save(".envdiff-baseline.json", results); err != nil {
//		log.Fatal(err)
//	}
//
//	// Load and compare later
//	snap, err := baseline.Load(".envdiff-baseline.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	delta := baseline.Compare(snap, currentResults)
//	fmt.Printf("New: %d, Resolved: %d, Changed: %d\n",
//		len(delta.New), len(delta.Resolved), len(delta.Changed))
package baseline
