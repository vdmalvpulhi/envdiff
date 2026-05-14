// Package snapshot captures and restores the state of parsed env files,
// enabling point-in-time comparisons and audit trails.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// Snapshot holds a labelled capture of diff results at a point in time.
type Snapshot struct {
	Label     string           `json:"label"`
	CreatedAt time.Time        `json:"created_at"`
	Results   []diff.Result    `json:"results"`
}

// Take creates a new Snapshot with the given label and results.
func Take(label string, results []diff.Result) Snapshot {
	return Snapshot{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Results:   results,
	}
}

// Save writes the snapshot to a JSON file at path.
func Save(path string, s Snapshot) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at path.
func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return s, nil
}

// Diff compares two snapshots and returns keys whose status changed.
func Diff(before, after Snapshot) []Change {
	beforeMap := make(map[string]diff.Result, len(before.Results))
	for _, r := range before.Results {
		beforeMap[r.Key] = r
	}

	var changes []Change
	for _, r := range after.Results {
		if prev, ok := beforeMap[r.Key]; ok {
			if prev.Status != r.Status {
				changes = append(changes, Change{Key: r.Key, Before: prev.Status, After: r.Status})
			}
		} else {
			changes = append(changes, Change{Key: r.Key, Before: "", After: r.Status})
		}
	}
	return changes
}

// Change describes a status transition for a single key between two snapshots.
type Change struct {
	Key    string `json:"key"`
	Before string `json:"before"`
	After  string `json:"after"`
}
