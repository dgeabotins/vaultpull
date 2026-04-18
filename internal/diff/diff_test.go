package diff

import (
	"strings"
	"testing"
)

func TestCompare_Added(t *testing.T) {
	old := map[string]string{}
	incoming := map[string]string{"FOO": "bar"}
	r := Compare(old, incoming)
	if _, ok := r.Added["FOO"]; !ok {
		t.Error("expected FOO to be in Added")
	}
}

func TestCompare_Changed(t *testing.T) {
	old := map[string]string{"FOO": "old"}
	incoming := map[string]string{"FOO": "new"}
	r := Compare(old, incoming)
	if v, ok := r.Changed["FOO"]; !ok || v != "new" {
		t.Errorf("expected FOO in Changed with value 'new', got %q", v)
	}
}

func TestCompare_Removed(t *testing.T) {
	old := map[string]string{"FOO": "bar", "BAZ": "qux"}
	incoming := map[string]string{"BAZ": "qux"}
	r := Compare(old, incoming)
	if _, ok := r.Removed["FOO"]; !ok {
		t.Error("expected FOO to be in Removed")
	}
}

func TestCompare_Unchanged(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"FOO": "bar"}
	r := Compare(old, incoming)
	if _, ok := r.Unchanged["FOO"]; !ok {
		t.Error("expected FOO to be in Unchanged")
	}
}

func TestSummary_Format(t *testing.T) {
	r := Result{
		Added:     map[string]string{"A": "1"},
		Changed:   map[string]string{"B": "2"},
		Removed:   map[string]string{},
		Unchanged: map[string]string{"C": "3", "D": "4"},
	}
	s := Summary(r)
	if !strings.Contains(s, "+1") || !strings.Contains(s, "~1") || !strings.Contains(s, "-0") {
		t.Errorf("unexpected summary: %s", s)
	}
}
