package env

import (
	"strings"
)

// TrimResult holds the outcome of a trim operation.
type TrimResult struct {
	Trimmed []string
	Total   int
}

// Summary returns a human-readable summary of the trim result.
func (r TrimResult) Summary() string {
	if len(r.Trimmed) == 0 {
		return "no values needed trimming"
	}
	return strings.Join(
		[]string{
			strconv.Itoa(len(r.Trimmed)) + " value(s) trimmed",
		},
		", ",
	)
}

// Trim removes leading/trailing whitespace from all values in the map.
// It returns the modified map and a TrimResult describing what changed.
func Trim(m map[string]string) (map[string]string, TrimResult) {
	out := make(map[string]string, len(m))
	result := TrimResult{Total: len(m)}
	for k, v := range m {
		trimmed := strings.TrimSpace(v)
		if trimmed != v {
			result.Trimmed = append(result.Trimmed, k)
		}
		out[k] = trimmed
	}
	return out, result
}
