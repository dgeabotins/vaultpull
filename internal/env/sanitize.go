package env

import (
	"regexp"
	"strings"
)

// SanitizeResult holds the outcome of a sanitize operation.
type SanitizeResult struct {
	Sanitized map[string]string
	Changes   []SanitizeChange
}

// SanitizeChange describes a single value that was sanitized.
type SanitizeChange struct {
	Key    string
	Before string
	After  string
}

// SanitizeOptions controls which sanitization rules are applied.
type SanitizeOptions struct {
	StripControlChars bool
	TrimWhitespace    bool
	NormalizeNewlines bool
}

var controlCharRe = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)

// Sanitize cleans env values according to the given options and returns a
// new map along with a list of changes that were made.
func Sanitize(input map[string]string, opts SanitizeOptions) SanitizeResult {
	result := SanitizeResult{
		Sanitized: make(map[string]string, len(input)),
	}

	for k, v := range input {
		original := v

		if opts.TrimWhitespace {
			v = strings.TrimSpace(v)
		}
		if opts.NormalizeNewlines {
			v = strings.ReplaceAll(v, "\r\n", "\n")
			v = strings.ReplaceAll(v, "\r", "\n")
		}
		if opts.StripControlChars {
			v = controlCharRe.ReplaceAllString(v, "")
		}

		result.Sanitized[k] = v
		if v != original {
			result.Changes = append(result.Changes, SanitizeChange{
				Key:    k,
				Before: original,
				After:  v,
			})
		}
	}

	return result
}

// Summary returns a human-readable summary of the sanitize result.
func (r SanitizeResult) Summary() string {
	if len(r.Changes) == 0 {
		return "no changes made"
	}
	var sb strings.Builder
	for _, c := range r.Changes {
		sb.WriteString("sanitized " + c.Key + "\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
