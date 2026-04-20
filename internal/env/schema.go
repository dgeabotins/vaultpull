package env

import (
	"fmt"
	"os"
	"strings"
)

// SchemaEntry describes a single expected env key.
type SchemaEntry struct {
	Key      string
	Required bool
	Default  string
}

// SchemaResult holds the outcome of a schema validation.
type SchemaResult struct {
	Missing  []string
	Defaults map[string]string // keys that were filled with defaults
}

// HasErrors returns true if any required keys are missing.
func (r SchemaResult) HasErrors() bool {
	return len(r.Missing) > 0
}

// Summary returns a human-readable summary of the schema result.
func (r SchemaResult) Summary() string {
	var sb strings.Builder
	if len(r.Missing) > 0 {
		sb.WriteString(fmt.Sprintf("missing required keys: %s\n", strings.Join(r.Missing, ", ")))
	}
	for k, v := range r.Defaults {
		sb.WriteString(fmt.Sprintf("default applied: %s=%s\n", k, v))
	}
	if sb.Len() == 0 {
		return "schema ok"
	}
	return strings.TrimRight(sb.String(), "\n")
}

// ApplySchema validates the provided map against a schema, filling defaults
// where applicable and collecting missing required keys.
func ApplySchema(data map[string]string, schema []SchemaEntry) (map[string]string, SchemaResult) {
	out := make(map[string]string, len(data))
	for k, v := range data {
		out[k] = v
	}

	result := SchemaResult{
		Defaults: make(map[string]string),
	}

	for _, entry := range schema {
		val, exists := out[entry.Key]
		if !exists || val == "" {
			if entry.Default != "" {
				out[entry.Key] = entry.Default
				result.Defaults[entry.Key] = entry.Default
			} else if entry.Required {
				result.Missing = append(result.Missing, entry.Key)
			}
		}
	}
	return out, result
}

// LoadSchema reads a simple schema file where each line is:
//   KEY          (required, no default)
//   KEY=default  (optional with default)
//   ?KEY         (optional, no default)
func LoadSchema(path string) ([]SchemaEntry, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read schema %s: %w", path, err)
	}
	var entries []SchemaEntry
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		optional := strings.HasPrefix(line, "?")
		if optional {
			line = line[1:]
		}
		if idx := strings.IndexByte(line, '='); idx >= 0 {
			entries = append(entries, SchemaEntry{
				Key:      strings.TrimSpace(line[:idx]),
				Required: false,
				Default:  strings.TrimSpace(line[idx+1:]),
			})
		} else {
			entries = append(entries, SchemaEntry{
				Key:      line,
				Required: !optional,
			})
		}
	}
	return entries, nil
}
