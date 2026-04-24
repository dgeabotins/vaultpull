package env

import (
	"strings"
	"testing"
)

func TestSanitize_TrimWhitespace(t *testing.T) {
	input := map[string]string{
		"KEY": "  hello world  ",
		"CLEAN": "already",
	}
	res := Sanitize(input, SanitizeOptions{TrimWhitespace: true})
	if res.Sanitized["KEY"] != "hello world" {
		t.Errorf("expected trimmed value, got %q", res.Sanitized["KEY"])
	}
	if res.Sanitized["CLEAN"] != "already" {
		t.Errorf("expected unchanged value, got %q", res.Sanitized["CLEAN"])
	}
	if len(res.Changes) != 1 || res.Changes[0].Key != "KEY" {
		t.Errorf("expected 1 change for KEY, got %v", res.Changes)
	}
}

func TestSanitize_StripControlChars(t *testing.T) {
	input := map[string]string{
		"KEY": "hello\x01world\x7F",
	}
	res := Sanitize(input, SanitizeOptions{StripControlChars: true})
	if res.Sanitized["KEY"] != "helloworld" {
		t.Errorf("expected control chars stripped, got %q", res.Sanitized["KEY"])
	}
	if len(res.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(res.Changes))
	}
}

func TestSanitize_NormalizeNewlines(t *testing.T) {
	input := map[string]string{
		"KEY": "line1\r\nline2\rline3",
	}
	res := Sanitize(input, SanitizeOptions{NormalizeNewlines: true})
	expected := "line1\nline2\nline3"
	if res.Sanitized["KEY"] != expected {
		t.Errorf("expected %q, got %q", expected, res.Sanitized["KEY"])
	}
}

func TestSanitize_NoChanges(t *testing.T) {
	input := map[string]string{
		"KEY": "clean_value",
	}
	res := Sanitize(input, SanitizeOptions{TrimWhitespace: true, StripControlChars: true})
	if len(res.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(res.Changes))
	}
}

func TestSanitize_Summary_NoChanges(t *testing.T) {
	res := SanitizeResult{Sanitized: map[string]string{}}
	if res.Summary() != "no changes made" {
		t.Errorf("unexpected summary: %q", res.Summary())
	}
}

func TestSanitize_Summary_WithChanges(t *testing.T) {
	res := SanitizeResult{
		Sanitized: map[string]string{},
		Changes: []SanitizeChange{
			{Key: "FOO", Before: " x ", After: "x"},
			{Key: "BAR", Before: "y\x01", After: "y"},
		},
	}
	summary := res.Summary()
	if !strings.Contains(summary, "FOO") || !strings.Contains(summary, "BAR") {
		t.Errorf("summary missing keys: %q", summary)
	}
}
