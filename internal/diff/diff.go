package diff

import "fmt"

// Result holds the comparison between existing and incoming secrets.
type Result struct {
	Added   map[string]string
	Changed map[string]string
	Removed map[string]string
	Unchanged map[string]string
}

// Compare returns a Result describing differences between old and new secret maps.
func Compare(old, incoming map[string]string) Result {
	r := Result{
		Added:     make(map[string]string),
		Changed:   make(map[string]string),
		Removed:   make(map[string]string),
		Unchanged: make(map[string]string),
	}

	for k, v := range incoming {
		if oldVal, exists := old[k]; !exists {
			r.Added[k] = v
		} else if oldVal != v {
			r.Changed[k] = v
		} else {
			r.Unchanged[k] = v
		}
	}

	for k, v := range old {
		if _, exists := incoming[k]; !exists {
			r.Removed[k] = v
		}
	}

	return r
}

// Summary returns a human-readable summary of the diff result.
func Summary(r Result) string {
	return fmt.Sprintf("+%d added, ~%d changed, -%d removed, %d unchanged",
		len(r.Added), len(r.Changed), len(r.Removed), len(r.Unchanged))
}
