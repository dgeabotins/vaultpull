package env

import (
	"strings"
	"unicode"
)

// NormalizeOptions controls how keys and values are normalized.
type NormalizeOptions struct {
	// UppercaseKeys converts all keys to UPPER_CASE.
	UppercaseKeys bool
	// SnakeCaseKeys replaces hyphens and spaces in keys with underscores.
	SnakeCaseKeys bool
	// TrimValues strips surrounding whitespace from values.
	TrimValues bool
	// StripQuotes removes surrounding single or double quotes from values.
	StripQuotes bool
}

// NormalizeResult holds the outcome of a Normalize call.
type NormalizeResult struct {
	Output  map[string]string
	Changed []string // keys whose key or value was modified
}

// Summary returns a human-readable summary of the normalization.
func (r NormalizeResult) Summary() string {
	if len(r.Changed) == 0 {
		return "normalize: no changes"
	}
	return "normalize: " + strings.Join(r.Changed, ", ") + " modified"
}

// Normalize applies the given options to env, returning a new map and a result
// describing what changed. The original map is never mutated.
func Normalize(env map[string]string, opts NormalizeOptions) NormalizeResult {
	out := make(map[string]string, len(env))
	changed := []string{}

	for k, v := range env {
		newKey := k
		newVal := v

		if opts.SnakeCaseKeys {
			newKey = toSnakeCase(newKey)
		}
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		}
		if opts.TrimValues {
			newVal = strings.TrimSpace(newVal)
		}
		if opts.StripQuotes {
			newVal = stripSurroundingQuotes(newVal)
		}

		if newKey != k || newVal != v {
			changed = append(changed, newKey)
		}
		out[newKey] = newVal
	}

	return NormalizeResult{Output: out, Changed: changed}
}

func toSnakeCase(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '-' || unicode.IsSpace(r) {
			b.WriteRune('_')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func stripSurroundingQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
