package env

import (
	"testing"
)

func TestCast_IntHint(t *testing.T) {
	env := map[string]string{"PORT": "8080", "NAME": "app"}
	result := Cast(env, CastOptions{
		TypeHints: map[string]string{"PORT": "int"},
	})
	if result.Casted["PORT"] != "8080" {
		t.Errorf("expected 8080, got %s", result.Casted["PORT"])
	}
	if result.Changed != 0 {
		t.Errorf("expected 0 changes, got %d", result.Changed)
	}
}

func TestCast_IntHint_Invalid(t *testing.T) {
	env := map[string]string{"PORT": "not-a-number"}
	result := Cast(env, CastOptions{
		TypeHints: map[string]string{"PORT": "int"},
	})
	if len(result.Errors) == 0 {
		t.Error("expected an error for invalid int")
	}
	if result.Casted["PORT"] != "not-a-number" {
		t.Error("original value should be preserved on error")
	}
}

func TestCast_BoolStrictNormalises(t *testing.T) {
	cases := map[string]string{
		"yes": "true",
		"no":  "false",
		"1":   "true",
		"0":   "false",
		"on":  "true",
		"off": "false",
	}
	for input, want := range cases {
		env := map[string]string{"FLAG": input}
		result := Cast(env, CastOptions{
			TypeHints:  map[string]string{"FLAG": "bool"},
			StrictBool: true,
		})
		if result.Casted["FLAG"] != want {
			t.Errorf("input %q: expected %q, got %q", input, want, result.Casted["FLAG"])
		}
	}
}

func TestCast_FloatHint(t *testing.T) {
	env := map[string]string{"RATIO": "3.14"}
	result := Cast(env, CastOptions{
		TypeHints: map[string]string{"RATIO": "float"},
	})
	if result.Casted["RATIO"] != "3.14" {
		t.Errorf("expected 3.14, got %s", result.Casted["RATIO"])
	}
	if len(result.Errors) != 0 {
		t.Errorf("unexpected errors: %v", result.Errors)
	}
}

func TestCast_UnknownHint_ReturnsError(t *testing.T) {
	env := map[string]string{"X": "value"}
	result := Cast(env, CastOptions{
		TypeHints: map[string]string{"X": "datetime"},
	})
	if len(result.Errors) == 0 {
		t.Error("expected error for unknown type hint")
	}
}

func TestCast_NoHints_PassThrough(t *testing.T) {
	env := map[string]string{"A": "1", "B": "hello"}
	result := Cast(env, CastOptions{})
	if result.Changed != 0 {
		t.Errorf("expected 0 changes, got %d", result.Changed)
	}
	if len(result.Errors) != 0 {
		t.Errorf("unexpected errors: %v", result.Errors)
	}
}
