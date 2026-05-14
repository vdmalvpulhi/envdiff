// Package baseline provides functionality for saving and comparing
// .env diff results against a stored baseline snapshot.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/envdiff/internal/diff"
)

// Snapshot represents a saved baseline of diff results.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Results   []diff.Result     `json:"results"`
}

// Save writes a snapshot of the given results to the specified file path.
func Save(path string, results []diff.Result) error {
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Results:   results,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal snapshot: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write file %q: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the specified file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: read file %q: %w", path, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal snapshot: %w", err)
	}
	return &snap, nil
}

// Delta represents the difference between a current result set and a baseline.
type Delta struct {
	New      []diff.Result // present in current, absent from baseline
	Resolved []diff.Result // present in baseline, absent from current
	Changed  []diff.Result // present in both but status differs
}

// Compare returns a Delta describing what changed relative to the baseline snapshot.
func Compare(snap *Snapshot, current []diff.Result) Delta {
	baseMap := make(map[string]diff.Result, len(snap.Results))
	for _, r := range snap.Results {
		baseMap[r.Key] = r
	}
	currMap := make(map[string]diff.Result, len(current))
	for _, r := range current {
		currMap[r.Key] = r
	}

	var d Delta
	for _, r := range current {
		if b, ok := baseMap[r.Key]; !ok {
			d.New = append(d.New, r)
		} else if b.Status != r.Status {
			d.Changed = append(d.Changed, r)
		}
	}
	for _, r := range snap.Results {
		if _, ok := currMap[r.Key]; !ok {
			d.Resolved = append(d.Resolved, r)
		}
	}
	return d
}
