package env

import (
	"fmt"
	"strconv"
	"strings"
)

// TypeHint represents an expected type for an env var value.
type TypeHint string

const (
	TypeString  TypeHint = "string"
	TypeInt     TypeHint = "int"
	TypeFloat   TypeHint = "float"
	TypeBool    TypeHint = "bool"
	TypeURL     TypeHint = "url"
	TypeNonempty TypeHint = "nonempty"
)

// TypeIssue describes a single type-check failure.
type TypeIssue struct {
	Key      string
	Value    string
	Expected TypeHint
	Reason   string
}

// TypeCheckResult holds the outcome of a type-check run.
type TypeCheckResult struct {
	Issues []TypeIssue
}

func (r TypeCheckResult) HasIssues() bool { return len(r.Issues) > 0 }

func (r TypeCheckResult) Summary() string {
	if !r.HasIssues() {
		return "all values pass type checks"
	}
	var sb strings.Builder
	for _, iss := range r.Issues {
		fmt.Fprintf(&sb, "  [%s] %q expected %s: %s\n", iss.Key, iss.Value, iss.Expected, iss.Reason)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// TypeCheck validates the values in m against the provided hints map (key -> TypeHint).
func TypeCheck(m map[string]string, hints map[string]TypeHint) TypeCheckResult {
	var issues []TypeIssue
	for key, hint := range hints {
		val, ok := m[key]
		if !ok {
			continue // missing keys are not a type error
		}
		if reason := checkValue(val, hint); reason != "" {
			issues = append(issues, TypeIssue{Key: key, Value: val, Expected: hint, Reason: reason})
		}
	}
	return TypeCheckResult{Issues: issues}
}

func checkValue(val string, hint TypeHint) string {
	switch hint {
	case TypeInt:
		if _, err := strconv.ParseInt(val, 10, 64); err != nil {
			return "not a valid integer"
		}
	case TypeFloat:
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return "not a valid float"
		}
	case TypeBool:
		norm := strings.ToLower(strings.TrimSpace(val))
		valid := map[string]bool{"true": true, "false": true, "1": true, "0": true, "yes": true, "no": true}
		if !valid[norm] {
			return "not a valid boolean (true/false/1/0/yes/no)"
		}
	case TypeURL:
		if !strings.HasPrefix(val, "http://") && !strings.HasPrefix(val, "https://") {
			return "not a valid URL (must start with http:// or https://)"
		}
	case TypeNonempty:
		if strings.TrimSpace(val) == "" {
			return "value must not be empty"
		}
	}
	return ""
}
