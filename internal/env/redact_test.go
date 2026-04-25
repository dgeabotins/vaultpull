package env

import (
	"testing"
)

func TestRedact_ExplicitKeys(t *testing.T) {
	values := map[string]string{
		"DB_PASSWORD": "secret",
		"APP_NAME":    "myapp",
	}
	res, err := Redact(values, RedactOptions{Keys: []string{"DB_PASSWORD"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Values["DB_PASSWORD"] != "***" {
		t.Errorf("expected DB_PASSWORD to be redacted")
	}
	if res.Values["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be unchanged")
	}
	if len(res.Redacted) != 1 || res.Redacted[0] != "DB_PASSWORD" {
		t.Errorf("expected one redacted key, got %v", res.Redacted)
	}
}

func TestRedact_PatternMatch(t *testing.T) {
	values := map[string]string{
		"API_SECRET": "topsecret",
		"API_KEY":    "key123",
		"APP_PORT":   "8080",
	}
	res, err := Redact(values, RedactOptions{Patterns: []string{"(?i)(secret|key)"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Values["API_SECRET"] != "***" {
		t.Errorf("expected API_SECRET redacted")
	}
	if res.Values["API_KEY"] != "***" {
		t.Errorf("expected API_KEY redacted")
	}
	if res.Values["APP_PORT"] != "8080" {
		t.Errorf("expected APP_PORT unchanged")
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	values := map[string]string{"TOKEN": "abc"}
	res, err := Redact(values, RedactOptions{Keys: []string{"TOKEN"}, Placeholder: "<hidden>"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Values["TOKEN"] != "<hidden>" {
		t.Errorf("expected custom placeholder, got %s", res.Values["TOKEN"])
	}
}

func TestRedact_NothingRedacted(t *testing.T) {
	values := map[string]string{"HOST": "localhost"}
	res, err := Redact(values, RedactOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Redacted) != 0 {
		t.Errorf("expected no redactions")
	}
	if res.Summary() != "no keys redacted" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestRedact_InvalidPattern(t *testing.T) {
	values := map[string]string{"X": "y"}
	_, err := Redact(values, RedactOptions{Patterns: []string{"[invalid"}})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestRedactResult_Summary(t *testing.T) {
	res := RedactResult{Redacted: []string{"A", "B"}}
	s := res.Summary()
	if s != "A, B redacted" {
		t.Errorf("unexpected summary: %s", s)
	}
}
