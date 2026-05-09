// Package loader provides functionality for discovering and loading
// .env files from the filesystem based on glob patterns or explicit paths.
package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envdiff/internal/parser"
)

// EnvFile represents a loaded environment file with its path and parsed variables.
type EnvFile struct {
	Path string
	Vars map[string]string
}

// LoadFile loads and parses a single .env file from the given path.
func LoadFile(path string) (*EnvFile, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	vars, err := parser.ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	return &EnvFile{Path: path, Vars: vars}, nil
}

// LoadGlob loads all .env files matching the given glob pattern.
// Returns a slice of EnvFile entries sorted by path.
func LoadGlob(pattern string) ([]*EnvFile, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid glob pattern %q: %w", pattern, err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no files matched pattern: %s", pattern)
	}

	var files []*EnvFile
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil || info.IsDir() {
			continue
		}

		ef, err := LoadFile(match)
		if err != nil {
			return nil, err
		}
		files = append(files, ef)
	}

	return files, nil
}

// LoadPair is a convenience function that loads exactly two files for comparison.
func LoadPair(pathA, pathB string) (*EnvFile, *EnvFile, error) {
	a, err := LoadFile(pathA)
	if err != nil {
		return nil, nil, err
	}

	b, err := LoadFile(pathB)
	if err != nil {
		return nil, nil, err
	}

	return a, b, nil
}
