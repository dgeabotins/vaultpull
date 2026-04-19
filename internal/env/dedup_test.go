package env

import (
	"strings"
	"testing"
)

func TestDedup_NoDuplicates(t *testing.T) {
	lines := []string{"FOO=bar", "BAZ=qux"}
	out, result := Dedup(lines)
	if len(result.Duplicates) != 0 {
		t.Fatalf("expected no duplicates, got %d", len(result.Duplicates))
	}
	if result.Removed != 0 {
		t.Fatalf("expected 0 removed, got %d", result.Removed)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
}

func TestDedup_RemovesDuplicates(t *testing.T) {
	lines := []string{"FOO=first", "BAR=keep", "FOO=second", "FOO=third"}
	out, result := Dedup(lines)
	if result.Removed != 2 {
		t.Fatalf("expected 2 removed, got %d", result.Removed)
	}
	if len(result.Duplicates) != 1 {
		t.Fatalf("expected 1 duplicate entry, got %d", len(result.Duplicates))
	}
	for _, line := range out {
		if strings.Contains(line, "FOO=") && !strings.Contains(line, "first") {
			t.Errorf("unexpected line in output: %s", line)
		}
	}
}

func TestDedup_PreservesComments(t *testing.T) {
	lines := []string{"# comment", "FOO=bar", "# another", "FOO=baz"}
	out, result := Dedup(lines)
	if result.Removed != 1 {
		t.Fatalf("expected 1 removed, got %d", result.Removed)
	}
	comments := 0
	for _, l := range out {
		if strings.HasPrefix(strings.TrimSpace(l), "#") {
			comments++
		}
	}
	if comments != 2 {
		t.Errorf("expected 2 comments preserved, got %d", comments)
	}
}

func TestDedupResult_Summary_NoIssues(t *testing.T) {
	r := DedupResult{}
	if r.Summary() != "no duplicates found" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestDedupResult_Summary_WithDupes(t *testing.T) {
	r := DedupResult{
		Duplicates: []DuplicateEntry{{Key: "FOO", Lines: []int{1, 3}}},
		Removed:    1,
	}
	s := r.Summary()
	if !strings.Contains(s, "FOO") {
		t.Errorf("expected FOO in summary, got: %s", s)
	}
	if !strings.Contains(s, "1 duplicate(s) removed") {
		t.Errorf("expected removal count in summary, got: %s", s)
	}
}
