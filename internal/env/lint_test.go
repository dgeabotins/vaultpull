package env

import (
	"strings"
	"testing"
)

func TestLintMap_NoIssues(t *testing.T) {
	m := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "abc123",
	}
	r := LintMap(m)
	if len(r.Issues) != 0 {
		t.Fatalf("expected no issues, got %d: %s", len(r.Issues), r.Summary())
	}
	if r.HasErrors() {
		t.Error("expected HasErrors=false")
	}
}

func TestLintMap_InvalidKey(t *testing.T) {
	m := map[string]string{"bad-key": "value"}
	r := LintMap(m)
	if len(r.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(r.Issues))
	}
	if r.Issues[0].Severity != LintError {
		t.Errorf("expected error severity, got %s", r.Issues[0].Severity)
	}
	if !r.HasErrors() {
		t.Error("expected HasErrors=true")
	}
}

func TestLintMap_EmptyValue(t *testing.T) {
	m := map[string]string{"MY_VAR": ""}
	r := LintMap(m)
	if len(r.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(r.Issues))
	}
	if r.Issues[0].Severity != LintWarning {
		t.Errorf("expected warning severity, got %s", r.Issues[0].Severity)
	}
	if r.HasErrors() {
		t.Error("expected HasErrors=false for warning-only result")
	}
}

func TestLintMap_WhitespaceValue(t *testing.T) {
	m := map[string]string{"MY_VAR": " padded "}
	r := LintMap(m)
	if len(r.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(r.Issues))
	}
	if r.Issues[0].Severity != LintWarning {
		t.Errorf("expected warning, got %s", r.Issues[0].Severity)
	}
}

func TestLintResult_Summary_NoIssues(t *testing.T) {
	r := LintResult{}
	if r.Summary() != "no lint issues found" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestLintResult_Summary_WithIssues(t *testing.T) {
	r := LintResult{
		Issues: []LintIssue{
			{Key: "bad-key", Message: "invalid key", Severity: LintError},
		},
	}
	s := r.Summary()
	if !strings.Contains(s, "bad-key") {
		t.Errorf("expected key in summary, got: %s", s)
	}
	if !strings.Contains(s, "error") {
		t.Errorf("expected severity in summary, got: %s", s)
	}
}
