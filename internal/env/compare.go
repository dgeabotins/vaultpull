package env

import "fmt"

// CompareResult holds the result of comparing two env files.
type CompareResult struct {
	OnlyInA  map[string]string
	OnlyInB  map[string]string
	Changed  map[string][2]string // key -> [valueA, valueB]
	Same     map[string]string
}

// Compare returns a CompareResult between two env maps.
func Compare(a, b map[string]string) CompareResult {
	res := CompareResult{
		OnlyInA: make(map[string]string),
		OnlyInB: make(map[string]string),
		Changed: make(map[string][2]string),
		Same:    make(map[string]string),
	}
	for k, va := range a {
		if vb, ok := b[k]; !ok {
			res.OnlyInA[k] = va
		} else if va != vb {
			res.Changed[k] = [2]string{va, vb}
		} else {
			res.Same[k] = va
		}
	}
	for k, vb := range b {
		if _, ok := a[k]; !ok {
			res.OnlyInB[k] = vb
		}
	}
	return res
}

// Summary returns a human-readable summary of the comparison.
func (r CompareResult) Summary() string {
	return fmt.Sprintf("only-in-a=%d only-in-b=%d changed=%d same=%d",
		len(r.OnlyInA), len(r.OnlyInB), len(r.Changed), len(r.Same))
}
