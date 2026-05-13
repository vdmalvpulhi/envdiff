// Package annotator attaches human-readable annotations to diff results,
// providing contextual hints such as whether a key looks like a secret,
// whether it has a default value, or whether it appears in multiple groups.
package annotator

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Annotation holds metadata attached to a single diff result.
type Annotation struct {
	Key        string
	Hint       string
	IsSensitive bool
	HasDefault  bool
	DefaultVal  string
}

// Options controls annotation behaviour.
type Options struct {
	// Defaults maps key names to their default values.
	Defaults map[string]string
	// SensitivePatterns is a list of substrings that mark a key as sensitive.
	SensitivePatterns []string
}

// DefaultOptions returns sensible defaults for annotation.
func DefaultOptions() Options {
	return Options{
		SensitivePatterns: []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PASS", "PRIVATE"},
	}
}

// Annotate produces an Annotation for every result in rs.
func Annotate(rs []diff.Result, opts Options) []Annotation {
	if opts.SensitivePatterns == nil {
		opts = DefaultOptions()
	}

	out := make([]Annotation, 0, len(rs))
	for _, r := range rs {
		a := Annotation{Key: r.Key}
		a.IsSensitive = isSensitive(r.Key, opts.SensitivePatterns)

		if opts.Defaults != nil {
			if dv, ok := opts.Defaults[r.Key]; ok {
				a.HasDefault = true
				a.DefaultVal = dv
			}
		}

		a.Hint = buildHint(r, a)
		out = append(out, a)
	}
	return out
}

func isSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

func buildHint(r diff.Result, a Annotation) string {
	switch r.Status {
	case diff.StatusMissing:
		if a.HasDefault {
			return "missing; default available: " + a.DefaultVal
		}
		if a.IsSensitive {
			return "missing sensitive key — add via secret manager"
		}
		return "missing in target environment"
	case diff.StatusConflict:
		if a.IsSensitive {
			return "sensitive key has conflicting values"
		}
		return "value differs between environments"
	default:
		return ""
	}
}
