package lint

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var validKeyRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

type Issue struct {
	Line    int
	Message string
}

type Result struct {
	Issues []Issue
}

func (r *Result) HasIssues() bool {
	return len(r.Issues) > 0
}

func (r *Result) Summary() string {
	if !r.HasIssues() {
		return "no issues found"
	}
	var sb strings.Builder
	for _, iss := range r.Issues {
		fmt.Fprintf(&sb, "line %d: %s\n", iss.Line, iss.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}

func Check(path string) (*Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	result := &Result{}
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			result.Issues = append(result.Issues, Issue{Line: lineNum, Message: "missing '=' separator"})
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if !validKeyRe.MatchString(key) {
			result.Issues = append(result.Issues, Issue{Line: lineNum, Message: fmt.Sprintf("key %q does not match [A-Z][A-Z0-9_]*", key)})
		}
		if val == "" {
			result.Issues = append(result.Issues, Issue{Line: lineNum, Message: fmt.Sprintf("key %q has empty value", key)})
		}
	}
	return result, scanner.Err()
}
