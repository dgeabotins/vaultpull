package env

import (
	"fmt"
	"sort"
	"strings"
)

// PinResult holds the outcome of a pin operation.
type PinResult struct {
	Pinned  []string
	Skipped []string
}

// PinOptions controls how keys are pinned.
type PinOptions struct {
	// Keys is the explicit list of keys to pin. Empty means all keys.
	Keys []string
	// Prefix restricts pinning to keys with this prefix.
	Prefix string
	// DryRun reports what would change without writing.
	DryRun bool
}

// Pin locks the values of selected keys by wrapping them in double-quotes
// and appending a "# pinned" comment marker. This signals downstream
// tooling (merge, promote, etc.) to treat those values as immutable.
func Pin(env map[string]string, opts PinOptions) (map[string]string, PinResult, error) {
	if env == nil {
		return nil, PinResult{}, fmt.Errorf("pin: env map must not be nil")
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	out := make(map[string]string, len(env))
	var result PinResult

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]

		selected := (len(keySet) == 0 && (opts.Prefix == "" || strings.HasPrefix(k, opts.Prefix))) ||
			keySet[k]

		if opts.Prefix != "" && len(keySet) == 0 && !strings.HasPrefix(k, opts.Prefix) {
			selected = false
		}

		if selected && !strings.HasSuffix(strings.TrimSpace(v), "# pinned") {
			if !opts.DryRun {
				out[k] = fmt.Sprintf("%s # pinned", v)
			} else {
				out[k] = v
			}
			result.Pinned = append(result.Pinned, k)
		} else {
			out[k] = v
			result.Skipped = append(result.Skipped, k)
		}
	}

	return out, result, nil
}

// IsPinned reports whether the given value carries a pin marker.
func IsPinned(value string) bool {
	return strings.HasSuffix(strings.TrimSpace(value), "# pinned")
}

// UnpinValue strips the pin marker from a value.
func UnpinValue(value string) string {
	v := strings.TrimSuffix(strings.TrimSpace(value), "# pinned")
	return strings.TrimSpace(v)
}
