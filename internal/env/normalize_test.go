package env

import (
	"strings"
	"testing"
)

func TestNormalize_UppercaseKeys(t *testing.T) {
	input := map[string]string{"db_host": "localhost", "api_key": "abc"}
	res := Normalize(input, NormalizeOptions{UppercaseKeys: true})

	if _, ok := res.Output["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in output")
	}
	if _, ok := res.Output["API_KEY"]; !ok {
		t.Error("expected API_KEY in output")
	}
	if len(res.Changed) != 2 {
		t.Errorf("expected 2 changed, got %d", len(res.Changed))
	}
}

func TestNormalize_SnakeCaseKeys(t *testing.T) {
	input := map[string]string{"db-host": "localhost", "my key": "val"}
	res := Normalize(input, NormalizeOptions{SnakeCaseKeys: true})

	if _, ok := res.Output["db_host"]; !ok {
		t.Error("expected db_host in output")
	}
	if _, ok := res.Output["my_key"]; !ok {
		t.Error("expected my_key in output")
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	input := map[string]string{"KEY": "  value  "}
	res := Normalize(input, NormalizeOptions{TrimValues: true})

	if got := res.Output["KEY"]; got != "value" {
		t.Errorf("expected 'value', got %q", got)
	}
	if len(res.Changed) != 1 {
		t.Errorf("expected 1 changed, got %d", len(res.Changed))
	}
}

func TestNormalize_StripQuotes(t *testing.T) {
	input := map[string]string{"A": `"hello"`, "B": "'world'", "C": "plain"}
	res := Normalize(input, NormalizeOptions{StripQuotes: true})

	if res.Output["A"] != "hello" {
		t.Errorf("expected hello, got %q", res.Output["A"])
	}
	if res.Output["B"] != "world" {
		t.Errorf("expected world, got %q", res.Output["B"])
	}
	if res.Output["C"] != "plain" {
		t.Errorf("expected plain unchanged, got %q", res.Output["C"])
	}
}

func TestNormalize_NoChanges(t *testing.T) {
	input := map[string]string{"KEY": "value"}
	res := Normalize(input, NormalizeOptions{})

	if len(res.Changed) != 0 {
		t.Errorf("expected no changes, got %v", res.Changed)
	}
}

func TestNormalize_Summary_NoChanges(t *testing.T) {
	res := NormalizeResult{Output: map[string]string{}, Changed: []string{}}
	if res.Summary() != "normalize: no changes" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestNormalize_Summary_WithChanges(t *testing.T) {
	res := NormalizeResult{Output: map[string]string{}, Changed: []string{"FOO", "BAR"}}
	if !strings.Contains(res.Summary(), "FOO") || !strings.Contains(res.Summary(), "BAR") {
		t.Errorf("summary missing keys: %s", res.Summary())
	}
}
