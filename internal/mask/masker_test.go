package mask

import (
	"strings"
	"testing"
)

func TestMask_ShortValue(t *testing.T) {
	m := New()
	result := m.Mask("abc")
	if result != strings.Repeat("*", 8) {
		t.Errorf("expected full mask, got %q", result)
	}
}

func TestMask_LongValue(t *testing.T) {
	m := New()
	result := m.Mask("supersecretvalue")
	if !strings.HasSuffix(result, "alue") {
		t.Errorf("expected last 4 chars visible, got %q", result)
	}
	if !strings.HasPrefix(result, "********") {
		t.Errorf("expected mask prefix, got %q", result)
	}
}

func TestMask_ExactBoundary(t *testing.T) {
	m := New()
	result := m.Mask("1234")
	if result != strings.Repeat("*", 8) {
		t.Errorf("expected full mask at boundary, got %q", result)
	}
}

func TestMaskMap_MasksAllValues(t *testing.T) {
	m := New()
	secrets := map[string]string{
		"DB_PASS": "hunter2",
		"API_KEY": "abcdefghijklmnop",
	}
	masked := m.MaskMap(secrets)
	for k, v := range masked {
		if v == secrets[k] {
			t.Errorf("key %s was not masked", k)
		}
		if !strings.HasPrefix(v, "*") {
			t.Errorf("key %s mask missing prefix, got %q", k, v)
		}
	}
}

func TestIsSafe_Masked(t *testing.T) {
	if !IsSafe("********alue") {
		t.Error("expected masked value to be safe")
	}
}

func TestIsSafe_Plain(t *testing.T) {
	if IsSafe("plaintext") {
		t.Error("expected plain value to not be safe")
	}
}
