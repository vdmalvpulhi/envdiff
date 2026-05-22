// Package auditor records and replays diff events for audit trail purposes.
// It tracks when comparisons were performed, what changed, and by whom.
package auditor

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Actor     string            `json:"actor"`
	Files     []string          `json:"files"`
	Results   []diff.Result     `json:"results"`
	Summary   Summary           `json:"summary"`
}

// Summary holds aggregate counts for an audit entry.
type Summary struct {
	Total     int `json:"total"`
	Missing   int `json:"missing"`
	Conflicts int `json:"conflicts"`
	Matching  int `json:"matching"`
}

// DefaultOptions returns sensible defaults for the auditor.
func DefaultOptions() Options {
	return Options{
		Actor: os.Getenv("USER"),
	}
}

// Options configures audit behaviour.
type Options struct {
	Actor string
}

// Record creates an audit Entry from the given results and file paths.
func Record(files []string, results []diff.Result, opts Options) Entry {
	if opts.Actor == "" {
		opts.Actor = "unknown"
	}
	s := buildSummary(results)
	return Entry{
		Timestamp: time.Now().UTC(),
		Actor:     opts.Actor,
		Files:     files,
		Results:   results,
		Summary:   s,
	}
}

// Save writes an Entry as JSON to the given path, appending to existing entries.
func Save(path string, entry Entry) error {
	entries, _ := Load(path) // ignore error — file may not exist yet
	entries = append(entries, entry)
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("auditor: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("auditor: write %s: %w", path, err)
	}
	return nil
}

// Load reads all audit entries from a JSON file.
func Load(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("auditor: read %s: %w", path, err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("auditor: unmarshal: %w", err)
	}
	return entries, nil
}

func buildSummary(results []diff.Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case diff.Missing:
			s.Missing++
		case diff.Conflict:
			s.Conflicts++
		case diff.Match:
			s.Matching++
		}
	}
	return s
}
