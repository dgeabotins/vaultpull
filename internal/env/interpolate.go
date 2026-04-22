package env

import (
	"fmt"
	"regexp"
	"strings"
)

// InterpolateResult holds the outcome of an interpolation pass.
type InterpolateResult struct {
	Resolved map[string]string
	Unresolved []string
}

// Summary returns a human-readable description of the interpolation result.
func (r InterpolateResult) Summary() string {
	if len(r.Unresolved) == 0 {
		return fmt.Sprintf("all %d keys resolved", len(r.Resolved))
	}
	return fmt.Sprintf("%d keys resolved, %d unresolved: %s",
		len(r.Resolved), len(r.Unresolved), strings.Join(r.Unresolved, ", "))
}

var refPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Interpolate expands ${KEY} references within values using other entries in
// the same map. It performs a single pass and collects any keys that could
// not be resolved.
func Interpolate(env map[string]string) InterpolateResult {
	resolved := make(map[string]string, len(env))
	var unresolved []string

	for k, v := range env {
		expanded, missing := expandValue(v, env)
		resolved[k] = expanded
		unresolved = append(unresolved, missing...)
	}

	// Deduplicate unresolved list.
	seen := map[string]struct{}{}
	filtered := unresolved[:0]
	for _, u := range unresolved {
		if _, ok := seen[u]; !ok {
			seen[u] = struct{}{}
			filtered = append(filtered, u)
		}
	}

	return InterpolateResult{Resolved: resolved, Unresolved: filtered}
}

func expandValue(v string, env map[string]string) (string, []string) {
	var missing []string
	result := refPattern.ReplaceAllStringFunc(v, func(match string) string {
		key := match[2 : len(match)-1] // strip ${ and }
		if val, ok := env[key]; ok {
			return val
		}
		missing = append(missing, key)
		return match // leave original placeholder
	})
	return result, missing
}
