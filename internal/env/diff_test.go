package env

import (
	"strings"
	"testing"
)

func TestDiff_Added(t *testing.T) {
	old := map[string]string{"A": "1"}
	next := map[string]string{"A": "1", "B": "2"}
	r := Diff(old, next)
	if _, ok := r.Added["B"]; !ok {
		t.Error("expected B to be added")
	}
	if len(r.Removed) != 0 || len(r.Changed) != 0 {
		t.Error("unexpected changes")
	}
}

func TestDiff_Removed(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1"}
	r := Diff(old, next)
	if _, ok := r.Removed["B"]; !ok {
		t.Error("expected B to be removed")
	}
}

func TestDiff_Changed(t *testing.T) {
	old := map[string]string{"A": "old"}
	next := map[string]string{"A": "new"}
	r := Diff(old, next)
	v, ok := r.Changed["A"]
	if !ok {
		t.Fatal("expected A to be changed")
	}
	if v[0] != "old" || v[1] != "new" {
		t.Errorf("unexpected change values: %v", v)
	}
}

func TestDiff_Same(t *testing.T) {
	old := map[string]string{"A": "1"}
	next := map[string]string{"A": "1"}
	r := Diff(old, next)
	if r.HasChanges() {
		t.Error("expected no changes")
	}
	if _, ok := r.Same["A"]; !ok {
		t.Error("expected A in Same")
	}
}

func TestDiff_Summary(t *testing.T) {
	old := map[string]string{"A": "1", "B": "old"}
	next := map[string]string{"B": "new", "C": "3"}
	r := Diff(old, next)
	s := r.Summary()
	if !strings.Contains(s, "+1") || !strings.Contains(s, "-1") || !strings.Contains(s, "~1") {
		t.Errorf("unexpected summary: %s", s)
	}
}
