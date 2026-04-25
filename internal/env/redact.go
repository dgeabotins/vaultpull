package env

import (
	"regexp"
	"strings"
)

// RedactOptions controls which keys are redacted and how.
type RedactOptions struct {
	// Keys is an explicit list of key names to redact.
	Keys []string
	// Patterns is a list of regex patterns matched against key names.
	Patterns []string
	// Placeholder replaces the secret value; defaults to "***".
	Placeholder string
}

// RedactResult holds the output map and metadata.
type RedactResult struct {
	Values   map[string]string
	Redacted []string
}

// Summary returns a human-readable summary of the redaction.
func (r RedactResult) Summary() string {
	if len(r.Redacted) == 0 {
		return "no keys redacted"
	}
	return strings.Join(r.Redacted, ", ") + " redacted"
}

// Redact returns a copy of values with sensitive keys replaced by a placeholder.
func Redact(values map[string]string, opts RedactOptions) (RedactResult, error) {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "***"
	}

	compiled := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return RedactResult{}, err
		}
		compiled = append(compiled, re)
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	out := make(map[string]string, len(values))
	var redacted []string

	for k, v := range values {
		if shouldRedact(k, keySet, compiled) {
			out[k] = placeholder
			redacted = append(redacted, k)
		} else {
			out[k] = v
		}
	}

	return RedactResult{Values: out, Redacted: redacted}, nil
}

func shouldRedact(key string, keySet map[string]struct{}, patterns []*regexp.Regexp) bool {
	if _, ok := keySet[key]; ok {
		return true
	}
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
