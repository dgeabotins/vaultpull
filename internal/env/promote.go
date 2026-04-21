package env

import "fmt"

// PromoteOptions controls how keys are promoted between environments.
type PromoteOptions struct {
	Keys    []string // if empty, promote all keys
	DryRun  bool
	Force   bool // overwrite existing keys in destination
}

// PromoteResult summarises what happened during promotion.
type PromoteResult struct {
	Promoted  []string
	Skipped   []string
	Overwrite []string
}

func (r PromoteResult) Summary() string {
	if len(r.Promoted) == 0 && len(r.Skipped) == 0 {
		return "nothing to promote"
	}
	return fmt.Sprintf("promoted: %d, skipped: %d, overwritten: %d",
		len(r.Promoted), len(r.Skipped), len(r.Overwrite))
}

// Promote copies keys from src map into dst map according to opts.
// It returns the merged destination map and a result summary.
func Promote(src, dst map[string]string, opts PromoteOptions) (map[string]string, PromoteResult) {
	keys := opts.Keys
	if len(keys) == 0 {
		for k := range src {
			keys = append(keys, k)
		}
	}

	result := PromoteResult{}
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	for _, k := range keys {
		val, ok := src[k]
		if !ok {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		if _, exists := out[k]; exists && !opts.Force {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		if _, exists := out[k]; exists && opts.Force {
			result.Overwrite = append(result.Overwrite, k)
		}
		if !opts.DryRun {
			out[k] = val
		}
		result.Promoted = append(result.Promoted, k)
	}
	return out, result
}
