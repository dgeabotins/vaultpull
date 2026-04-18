package diff

import (
	"strings"
	"testing"
)

func TestCompare_Added(t *testing.T) {
	old := map[string]string{}
	new := map[string]string{"FOO": "bar"}
	changes := Compare(old, new)
	if len(changes) != 1 || changes[0].Type != Added {
		t.Fatalf("expected 1 Added change, got %+v", changes)
	}
}

func TestCompare_Changed(t *testing.T) {
	old := map[string]string{"FOO": "old"}
	new := map[string]string{"FOO": "new"}
	changes := Compare(old, new)
	if len(changes) != 1 || changes[0].Type != Changed {
		t.Fatalf("expected 1 Changed change, got %+v", changes)
	}
	if changes[0].OldValue != "old" || changes[0].NewValue != "new" {
		t.Errorf("unexpected values: %+v", changes[0])
	}
}

func TestCompare_Removed(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new := map[string]string{}
	changes := Compare(old, new)
	if len(changes) != 1 || changes[0].Type != Removed {
		t.Fatalf("expected 1 Removed change, got %+v", changes)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new := map[string]string{"FOO": "bar"}
	changes := Compare(old, new)
	if len(changes) != 1 || changes[0].Type != Unchanged {
		t.Fatalf("expected 1 Unchanged change, got %+v", changes)
	}
}

func TestSummary_Format(t *testing.T) {
	changes := []Change{
		{Type: Added},
		{Type: Added},
		{Type: Changed},
		{Type: Removed},
	}
	s := Summary(changes)
	if !strings.Contains(s, "+2") || !strings.Contains(s, "~1") || !strings.Contains(s, "-1") {
		t.Errorf("unexpected summary: %s", s)
	}
}
