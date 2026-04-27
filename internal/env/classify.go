package env

import (
	"regexp"
	"strings"
)

// Category represents the inferred type/category of an env var value.
type Category string

const (
	CategoryURL      Category = "url"
	CategorySecret   Category = "secret"
	CategoryBoolean  Category = "boolean"
	CategoryInteger  Category = "integer"
	CategoryFloat    Category = "float"
	CategoryPath     Category = "path"
	CategoryJSON     Category = "json"
	CategoryEmpty    Category = "empty"
	CategoryUnknown  Category = "unknown"
)

var (
	urlRe     = regexp.MustCompile(`(?i)^https?://`)
	intRe     = regexp.MustCompile(`^-?\d+$`)
	floatRe   = regexp.MustCompile(`^-?\d+\.\d+$`)
	pathRe    = regexp.MustCompile(`^[/~.]`)
	secretKey = regexp.MustCompile(`(?i)(password|secret|token|key|api_key|auth|credential|private)`)
)

// ClassifyResult holds the category and key that was classified.
type ClassifyResult struct {
	Key      string
	Value    string
	Category Category
}

// Classify inspects each key/value pair and infers a Category.
func Classify(m map[string]string) []ClassifyResult {
	results := make([]ClassifyResult, 0, len(m))
	for _, k := range SortKeys(m, SortOptions{}) {
		v := m[k]
		results = append(results, ClassifyResult{
			Key:      k,
			Value:    v,
			Category: classifyValue(k, v),
		})
	}
	return results
}

func classifyValue(key, value string) Category {
	if value == "" {
		return CategoryEmpty
	}
	if urlRe.MatchString(value) {
		return CategoryURL
	}
	if strings.HasPrefix(value, "{") || strings.HasPrefix(value, "[") {
		return CategoryJSON
	}
	if intRe.MatchString(value) {
		return CategoryInteger
	}
	if floatRe.MatchString(value) {
		return CategoryFloat
	}
	v := strings.ToLower(value)
	if v == "true" || v == "false" || v == "yes" || v == "no" || v == "1" || v == "0" {
		return CategoryBoolean
	}
	if pathRe.MatchString(value) {
		return CategoryPath
	}
	if secretKey.MatchString(key) {
		return CategorySecret
	}
	return CategoryUnknown
}
