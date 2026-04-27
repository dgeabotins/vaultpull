package env

import (
	"strings"
	"testing"
)

func TestTypeCheck_AllPass(t *testing.T) {
	m := map[string]string{
		"PORT":    "8080",
		"RATIO":   "0.75",
		"ENABLED": "true",
		"API_URL": "https://example.com",
		"NAME":    "myapp",
	}
	hints := map[string]TypeHint{
		"PORT":    TypeInt,
		"RATIO":   TypeFloat,
		"ENABLED": TypeBool,
		"API_URL": TypeURL,
		"NAME":    TypeNonempty,
	}
	res := TypeCheck(m, hints)
	if res.HasIssues() {
		t.Fatalf("expected no issues, got: %s", res.Summary())
	}
}

func TestTypeCheck_InvalidInt(t *testing.T) {
	m := map[string]string{"PORT": "not-a-number"}
	hints := map[string]TypeHint{"PORT": TypeInt}
	res := TypeCheck(m, hints)
	if !res.HasIssues() {
		t.Fatal("expected issue for invalid int")
	}
	if res.Issues[0].Key != "PORT" {
		t.Errorf("expected key PORT, got %s", res.Issues[0].Key)
	}
}

func TestTypeCheck_InvalidBool(t *testing.T) {
	m := map[string]string{"DEBUG": "maybe"}
	hints := map[string]TypeHint{"DEBUG": TypeBool}
	res := TypeCheck(m, hints)
	if !res.HasIssues() {
		t.Fatal("expected issue for invalid bool")
	}
}

func TestTypeCheck_InvalidURL(t *testing.T) {
	m := map[string]string{"WEBHOOK": "ftp://old-school.example"}
	hints := map[string]TypeHint{"WEBHOOK": TypeURL}
	res := TypeCheck(m, hints)
	if !res.HasIssues() {
		t.Fatal("expected issue for non-http URL")
	}
}

func TestTypeCheck_EmptyValueNonempty(t *testing.T) {
	m := map[string]string{"SECRET": "   "}
	hints := map[string]TypeHint{"SECRET": TypeNonempty}
	res := TypeCheck(m, hints)
	if !res.HasIssues() {
		t.Fatal("expected issue for empty value")
	}
}

func TestTypeCheck_MissingKeySkipped(t *testing.T) {
	m := map[string]string{}
	hints := map[string]TypeHint{"MISSING": TypeInt}
	res := TypeCheck(m, hints)
	if res.HasIssues() {
		t.Fatal("missing key should not produce a type error")
	}
}

func TestTypeCheckResult_Summary_NoIssues(t *testing.T) {
	res := TypeCheckResult{}
	if !strings.Contains(res.Summary(), "pass") {
		t.Errorf("expected pass message, got: %s", res.Summary())
	}
}

func TestTypeCheckResult_Summary_WithIssues(t *testing.T) {
	res := TypeCheckResult{
		Issues: []TypeIssue{
			{Key: "PORT", Value: "abc", Expected: TypeInt, Reason: "not a valid integer"},
		},
	}
	if !strings.Contains(res.Summary(), "PORT") {
		t.Errorf("expected PORT in summary, got: %s", res.Summary())
	}
}
