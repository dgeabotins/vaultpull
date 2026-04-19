package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from an env file.
type Entry struct {
	Key   string
	Value string
}

// LoadFile reads a .env file and returns all key-value entries.
// Comments and blank lines are skipped.
func LoadFile(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("env: open %s: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, `"`)
		if key == "" {
			continue
		}
		entries = append(entries, Entry{Key: key, Value: val})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("env: scan %s: %w", path, err)
	}
	return entries, nil
}

// ToMap converts a slice of Entry into a map for quick lookup.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
