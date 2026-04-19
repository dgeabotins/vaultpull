package env

import (
	"testing"
)

func TestTrim_NoWhitespace(t *testing.T) {
	m := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}
	out, res := Trim(m)
	if len(res.Trimmed) != 0 {
		t.Errorf("expected 0 trimmed, got %d", len(res.Trimmed))
	}
	if out["KEY1"] != "value1" {
		t.Errorf("unexpected value: %s", out["KEY1"])
	}
}

func TestTrim_TrimsLeadingTrailing(t *testing.T) {
	m := map[string]string{
		"KEY1": "  hello  ",
		"KEY2": "world",
	}
	out, res := Trim(m)
	if out["KEY1"] != "hello" {
		t.Errorf("expected 'hello', got '%s'", out["KEY1"])
	}
	if out["KEY2"] != "world" {
		t.Errorf("expected 'world', got '%s'", out["KEY2"])
	}
	if len(res.Trimmed) != 1 || res.Trimmed[0] != "KEY1" {
		t.Errorf("expected KEY1 in trimmed list, got %v", res.Trimmed)
	}
}

func TestTrim_MultipleKeys(t *testing.T) {
	m := map[string]string{
		"A": "\t tab \t",
		"B": "\n newline \n",
		"C": "clean",
	}
	out, res := Trim(m)
	if out["A"] != "tab" {
		t.Errorf("expected 'tab', got '%s'", out["A"])
	}
	if out["B"] != "newline" {
		t.Errorf("expected 'newline', got '%s'", out["B"])
	}
	if res.Total != 3 {
		t.Errorf("expected total 3, got %d", res.Total)
	}
	if len(res.Trimmed) != 2 {
		t.Errorf("expected 2 trimmed, got %d", len(res.Trimmed))
	}
}

func TestTrim_EmptyMap(t *testing.T) {
	out, res := Trim(map[string]string{})
	if len(out) != 0 {
		t.Error("expected empty output map")
	}
	if res.Total != 0 {
		t.Errorf("expected total 0, got %d", res.Total)
	}
}
