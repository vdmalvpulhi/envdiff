package loader

import (
	"fmt"
	"strings"
)

// ValidationError holds all issues found during validation of an EnvFile.
type ValidationError struct {
	Path   string
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", e.Path, strings.Join(e.Issues, "; "))
}

// ValidationOptions controls which checks are applied during Validate.
type ValidationOptions struct {
	// DisallowEmpty rejects keys that have an empty value.
	DisallowEmpty bool
	// RequiredKeys lists keys that must be present in the file.
	RequiredKeys []string
}

// Validate checks an EnvFile against the provided options and returns a
// ValidationError if any issues are found, or nil if the file is valid.
func Validate(ef *EnvFile, opts ValidationOptions) error {
	var issues []string

	if opts.DisallowEmpty {
		for k, v := range ef.Vars {
			if strings.TrimSpace(v) == "" {
				issues = append(issues, fmt.Sprintf("key %q has empty value", k))
			}
		}
	}

	for _, req := range opts.RequiredKeys {
		if _, ok := ef.Vars[req]; !ok {
			issues = append(issues, fmt.Sprintf("required key %q is missing", req))
		}
	}

	if len(issues) == 0 {
		return nil
	}

	return &ValidationError{Path: ef.Path, Issues: issues}
}
