package validate

import (
	"fmt"
	"os"
	"strings"
)

// Result holds the outcome of a validation check.
type Result struct {
	Key     string
	Missing bool
	Empty   bool
}

// String returns a human-readable description of the result.
func (r Result) String() string {
	if r.Missing {
		return fmt.Sprintf("MISSING  %s", r.Key)
	}
	if r.Empty {
		return fmt.Sprintf("EMPTY    %s", r.Key)
	}
	return fmt.Sprintf("OK       %s", r.Key)
}

// CheckEnvFile reads an .env file and validates that all expected keys are
// present and non-empty. Returns one Result per expected key.
func CheckEnvFile(path string, expected []string) ([]Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	present := parseKeys(string(data))

	results := make([]Result, 0, len(expected))
	for _, key := range expected {
		r := Result{Key: key}
		val, ok := present[key]
		if !ok {
			r.Missing = true
		} else if strings.TrimSpace(val) == "" {
			r.Empty = true
		}
		results = append(results, r)
	}
	return results, nil
}

// HasFailures returns true if any result is missing or empty.
func HasFailures(results []Result) bool {
	for _, r := range results {
		if r.Missing || r.Empty {
			return true
		}
	}
	return false
}

func parseKeys(content string) map[string]string {
	m := make(map[string]string)
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		val := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
		m[strings.TrimSpace(parts[0])] = val
	}
	return m
}
