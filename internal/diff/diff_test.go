package diff

import (
	"strings"
	"testing"
)

func TestCompare_Added(t *testing.T) {
	old := map[string]string{}
	new := map[string]string{"FOO": "bar"}
	results := Compare(old, new)
	if len(results) != 1 || results[0].Status != "added" {
		t.Fatalf("expected 1 added result, got %+v", results)
	}
}

func TestCompare_Changed(t *testing.T) {
	old := map[string]string{"FOO": "old"}
	new := map[string]string{"FOO": "new"}
	results := Compare(old, new)
	if len(results) != 1 || results[0].Status != "changed" {
		t.Fatalf("expected 1 changed result, got %+v", results)
	}
	if results[0].OldVal != "old" || results[0].NewVal != "new" {
		t.Errorf("unexpected values: %+v", results[0])
	}
}

func TestCompare_Removed(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new := map[string]string{}
	results := Compare(old, new)
	if len(results) != 1 || results[0].Status != "removed" {
		t.Fatalf("expected 1 removed result, got %+v", results)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new := map[string]string{"FOO": "bar"}
	results := Compare(old, new)
	if len(results) != 1 || results[0].Status != "unchanged" {
		t.Fatalf("expected 1 unchanged result, got %+v", results)
	}
}

func TestSummary_Format(t *testing.T) {
	results := []Result{
		{Status: "added"},
		{Status: "added"},
		{Status: "changed"},
		{Status: "removed"},
	}
	s := Summary(results)
	if !strings.Contains(s, "+2") || !strings.Contains(s, "~1") || !strings.Contains(s, "-1") {
		t.Errorf("unexpected summary: %s", s)
	}
}
