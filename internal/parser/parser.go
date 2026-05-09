package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a parsed .env file as a map of key-value pairs.
type EnvMap map[string]string

// ParseFile reads a .env file and returns an EnvMap.
// Lines starting with '#' are treated as comments and ignored.
// Empty lines are skipped. Keys without values are stored with an empty string.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: failed to open file %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("parser: %s:%d: %w", path, lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: error reading file %q: %w", path, err)
	}

	return env, nil
}

// parseLine splits a single line into a key and value.
func parseLine(line string) (string, string, error) {
	// Remove optional "export " prefix
	line = strings.TrimPrefix(line, "export ")

	parts := strings.SplitN(line, "=", 2)
	key := strings.TrimSpace(parts[0])

	if key == "" {
		return "", "", fmt.Errorf("empty key in line %q", line)
	}

	if len(parts) == 1 {
		return key, "", nil
	}

	value := strings.TrimSpace(parts[1])
	value = stripQuotes(value)

	return key, value, nil
}

// stripQuotes removes surrounding single or double quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
