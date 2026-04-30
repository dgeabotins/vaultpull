package env

import (
	"fmt"
	"regexp"
	"strings"
)

// LintSeverity represents the severity level of a lint issue.
type LintSeverity string

const (
	LintError   LintSeverity = "error"
	LintWarning LintSeverity = "warning"
)

// LintIssue describes a single linting problem found in an env map.
type LintIssue struct {
	Key      string
	Message  string
	Severity LintSeverity
}

// LintResult holds all issues found during linting.
type LintResult struct {
	Issues []LintIssue
}

// HasErrors returns true if any issue has error severity.
func (r LintResult) HasErrors() bool {
	for _, i := range r.Issues {
		if i.Severity == LintError {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary of lint results.
func (r LintResult) Summary() string {
	if len(r.Issues) == 0 {
		return "no lint issues found"
	}
	var sb strings.Builder
	for _, i := range r.Issues {
		fmt.Fprintf(&sb, "[%s] %s: %s\n", i.Severity, i.Key, i.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}

var validKeyRe = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// LintMap runs lint checks against an env map and returns a LintResult.
func LintMap(m map[string]string) LintResult {
	var issues []LintIssue
	for k, v := range m {
		if !validKeyRe.MatchString(k) {
			issues = append(issues, LintIssue{
				Key:      k,
				Message:  "key should be uppercase with underscores only (A-Z, 0-9, _)",
				Severity: LintError,
			})
		}
		if v == "" {
			issues = append(issues, LintIssue{
				Key:      k,
				Message:  "value is empty",
				Severity: LintWarning,
			})
		}
		if strings.HasPrefix(v, " ") || strings.HasSuffix(v, " ") {
			issues = append(issues, LintIssue{
				Key:      k,
				Message:  "value has leading or trailing whitespace",
				Severity: LintWarning,
			})
		}
	}
	return LintResult{Issues: issues}
}
