package env

import (
	"fmt"
	"strconv"
	"strings"
)

// CastResult holds the outcome of a type-cast operation on an env map.
type CastResult struct {
	Casted  map[string]string
	Errors  []string
	Changed int
}

// CastOptions controls how values are coerced.
type CastOptions struct {
	// TypeHints maps key names to desired types: "int", "bool", "float", "string"
	TypeHints map[string]string
	// StrictBool normalises truthy strings (yes/no/1/0/on/off) to true/false
	StrictBool bool
}

// Cast coerces env values according to TypeHints and returns a CastResult.
func Cast(env map[string]string, opts CastOptions) CastResult {
	out := make(map[string]string, len(env))
	result := CastResult{Casted: out}

	for k, v := range env {
		hint, ok := opts.TypeHints[k]
		if !ok {
			out[k] = v
			continue
		}

		coerced, err := coerce(v, hint, opts.StrictBool)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", k, err))
			out[k] = v
			continue
		}
		if coerced != v {
			result.Changed++
		}
		out[k] = coerced
	}
	return result
}

func coerce(v, hint string, strictBool bool) (string, error) {
	switch strings.ToLower(hint) {
	case "int":
		n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to int", v)
		}
		return strconv.FormatInt(n, 10), nil
	case "float":
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to float", v)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case "bool":
		return normaliseBool(v, strictBool)
	case "string":
		return v, nil
	default:
		return "", fmt.Errorf("unknown type hint %q", hint)
	}
}

func normaliseBool(v string, strict bool) (string, error) {
	lower := strings.ToLower(strings.TrimSpace(v))
	if strict {
		switch lower {
		case "true", "1", "yes", "on":
			return "true", nil
		case "false", "0", "no", "off", "":
			return "false", nil
		}
		return "", fmt.Errorf("cannot cast %q to bool", v)
	}
	b, err := strconv.ParseBool(lower)
	if err != nil {
		return "", fmt.Errorf("cannot cast %q to bool", v)
	}
	return strconv.FormatBool(b), nil
}
