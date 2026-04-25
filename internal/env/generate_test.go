package env

import (
	"strings"
	"testing"
)

func TestGenerate_NewKeys(t *testing.T) {
	env := map[string]string{}
	opts := GenerateOptions{
		Keys:   []string{"SECRET_KEY", "API_TOKEN"},
		Length: 32,
	}
	out, result, err := Generate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Generated) != 2 {
		t.Errorf("expected 2 generated, got %d", len(result.Generated))
	}
	if len(out["SECRET_KEY"]) != 32 {
		t.Errorf("expected length 32 for SECRET_KEY, got %d", len(out["SECRET_KEY"]))
	}
}

func TestGenerate_SkipsExistingWithoutOverwrite(t *testing.T) {
	env := map[string]string{"EXISTING": "old_value"}
	opts := GenerateOptions{
		Keys:    []string{"EXISTING"},
		Length:  16,
		Overwrite: false,
	}
	out, result, err := Generate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if out["EXISTING"] != "old_value" {
		t.Errorf("expected original value preserved")
	}
}

func TestGenerate_OverwriteFlag(t *testing.T) {
	env := map[string]string{"KEY": "original"}
	opts := GenerateOptions{
		Keys:      []string{"KEY"},
		Length:    24,
		Overwrite: true,
	}
	out, result, err := Generate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Generated) != 1 {
		t.Errorf("expected 1 generated")
	}
	if out["KEY"] == "original" {
		t.Errorf("expected value to be overwritten")
	}
}

func TestGenerate_DryRun(t *testing.T) {
	env := map[string]string{}
	opts := GenerateOptions{
		Keys:   []string{"DRY_KEY"},
		Length: 16,
		DryRun: true,
	}
	out, result, err := Generate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Generated) != 1 {
		t.Errorf("expected 1 in generated list")
	}
	if _, ok := out["DRY_KEY"]; ok {
		t.Errorf("dry run should not write key")
	}
}

func TestGenerate_Charsets(t *testing.T) {
	charsets := []string{"alpha", "alphanum", "hex", "base64", "symbol"}
	for _, cs := range charsets {
		opts := GenerateOptions{
			Keys:    []string{"K"},
			Length:  20,
			Charset: cs,
		}
		out, _, err := Generate(map[string]string{}, opts)
		if err != nil {
			t.Errorf("charset %s: unexpected error: %v", cs, err)
		}
		if len(out["K"]) != 20 {
			t.Errorf("charset %s: expected length 20, got %d", cs, len(out["K"]))
		}
	}
}

func TestGenerateResult_Summary(t *testing.T) {
	r := GenerateResult{
		Generated: []string{"A", "B"},
		Skipped:   []string{"C"},
	}
	s := r.Summary()
	if !strings.Contains(s, "2 generated") {
		t.Errorf("expected '2 generated' in summary, got: %s", s)
	}
	if !strings.Contains(s, "1 skipped") {
		t.Errorf("expected '1 skipped' in summary, got: %s", s)
	}
}
