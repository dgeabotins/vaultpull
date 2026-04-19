package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMergeIntoFile_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, err := MergeIntoFile(path, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added != 2 {
		t.Errorf("expected 2 added, got %d", res.Added)
	}
	if res.Updated != 0 || res.Unchanged != 0 {
		t.Errorf("expected 0 updated/unchanged")
	}
}

func TestMergeIntoFile_UpdatesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	// Write initial file
	initial := map[string]string{"FOO": "old", "KEEP": "same"}
	if err := WriteFile(path, initial); err != nil {
		t.Fatal(err)
	}

	secrets := map[string]string{"FOO": "new", "KEEP": "same", "EXTRA": "val"}
	res, err := MergeIntoFile(path, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added != 1 {
		t.Errorf("expected 1 added, got %d", res.Added)
	}
	if res.Updated != 1 {
		t.Errorf("expected 1 updated, got %d", res.Updated)
	}
	if res.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", res.Unchanged)
	}

	loaded, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %s", loaded["FOO"])
	}
}

func TestMergeResult_Summary(t *testing.T) {
	r := MergeResult{Added: 2, Updated: 1, Unchanged: 3}
	s := r.Summary()
	for _, want := range []string{"2 added", "1 updated", "3 unchanged"} {
		if !containsStr(s, want) {
			t.Errorf("summary missing %q: %s", want, s)
		}
	}
}

func TestMergeResult_Summary_NoChanges(t *testing.T) {
	r := MergeResult{}
	if r.Summary() != "no changes" {
		t.Errorf("expected 'no changes', got %s", r.Summary())
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStrHelper(s, sub))
}

func containsStrHelper(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestMergeIntoFile_BadPath(t *testing.T) {
	_, err := MergeIntoFile("/nonexistent/dir/.env", map[string]string{"K": "v"})
	if err == nil {
		t.Error("expected error for bad path")
	}
	_ = os.Remove("/nonexistent/dir/.env")
}
