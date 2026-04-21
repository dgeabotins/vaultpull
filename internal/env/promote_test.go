package env

import (
	"strings"
	"testing"
)

func TestPromote_AllKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{}
	out, res := Promote(src, dst, PromoteOptions{})
	if len(res.Promoted) != 2 {
		t.Fatalf("expected 2 promoted, got %d", len(res.Promoted))
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Error("unexpected values in output")
	}
}

func TestPromote_SelectedKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}
	out, res := Promote(src, dst, PromoteOptions{Keys: []string{"A", "C"}})
	if len(res.Promoted) != 2 {
		t.Fatalf("expected 2 promoted, got %d", len(res.Promoted))
	}
	if _, ok := out["B"]; ok {
		t.Error("B should not be in output")
	}
}

func TestPromote_SkipsExistingWithoutForce(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, res := Promote(src, dst, PromoteOptions{})
	if out["A"] != "old" {
		t.Error("expected original value preserved")
	}
	if len(res.Skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(res.Skipped))
	}
}

func TestPromote_ForceOverwrites(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, res := Promote(src, dst, PromoteOptions{Force: true})
	if out["A"] != "new" {
		t.Error("expected new value after force")
	}
	if len(res.Overwrite) != 1 {
		t.Fatalf("expected 1 overwrite, got %d", len(res.Overwrite))
	}
}

func TestPromote_DryRun(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{}
	out, res := Promote(src, dst, PromoteOptions{DryRun: true})
	if _, ok := out["A"]; ok {
		t.Error("dry run should not modify output")
	}
	if len(res.Promoted) != 1 {
		t.Fatalf("expected 1 in promoted list, got %d", len(res.Promoted))
	}
}

func TestPromoteResult_Summary(t *testing.T) {
	res := PromoteResult{
		Promoted:  []string{"A", "B"},
		Skipped:   []string{"C"},
		Overwrite: []string{},
	}
	s := res.Summary()
	if !strings.Contains(s, "promoted: 2") {
		t.Errorf("unexpected summary: %s", s)
	}
}
