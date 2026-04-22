package env

import (
	"strings"
	"testing"
)

func TestInterpolate_ResolvesReference(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api",
	}
	res := Interpolate(env)
	if got := res.Resolved["API_URL"]; got != "https://example.com/api" {
		t.Errorf("expected expanded URL, got %q", got)
	}
	if len(res.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %v", res.Unresolved)
	}
}

func TestInterpolate_UnresolvedReference(t *testing.T) {
	env := map[string]string{
		"API_URL": "${MISSING_KEY}/api",
	}
	res := Interpolate(env)
	if len(res.Unresolved) != 1 || res.Unresolved[0] != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY unresolved, got %v", res.Unresolved)
	}
	if got := res.Resolved["API_URL"]; !strings.Contains(got, "${MISSING_KEY}") {
		t.Errorf("expected placeholder preserved, got %q", got)
	}
}

func TestInterpolate_MultipleReferences(t *testing.T) {
	env := map[string]string{
		"HOST":    "localhost",
		"PORT":    "5432",
		"DB_URL":  "postgres://${HOST}:${PORT}/mydb",
	}
	res := Interpolate(env)
	if got := res.Resolved["DB_URL"]; got != "postgres://localhost:5432/mydb" {
		t.Errorf("unexpected DB_URL: %q", got)
	}
}

func TestInterpolate_NoReferences(t *testing.T) {
	env := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	res := Interpolate(env)
	if res.Resolved["FOO"] != "bar" || res.Resolved["BAZ"] != "qux" {
		t.Errorf("plain values should be unchanged")
	}
	if len(res.Unresolved) != 0 {
		t.Errorf("expected no unresolved")
	}
}

func TestInterpolate_DeduplicatesUnresolved(t *testing.T) {
	env := map[string]string{
		"A": "${GHOST}",
		"B": "${GHOST}/extra",
	}
	res := Interpolate(env)
	if len(res.Unresolved) != 1 {
		t.Errorf("expected 1 unique unresolved key, got %d: %v", len(res.Unresolved), res.Unresolved)
	}
}

func TestInterpolateResult_Summary(t *testing.T) {
	r := InterpolateResult{Resolved: map[string]string{"A": "1", "B": "2"}, Unresolved: nil}
	if !strings.Contains(r.Summary(), "2 keys resolved") {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
	r2 := InterpolateResult{Resolved: map[string]string{"A": "1"}, Unresolved: []string{"MISSING"}}
	if !strings.Contains(r2.Summary(), "unresolved") {
		t.Errorf("expected unresolved in summary: %s", r2.Summary())
	}
}
