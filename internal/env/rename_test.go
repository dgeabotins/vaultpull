package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvForRename(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRename_Success(t *testing.T) {
	p := writeTempEnvForRename(t, "FOO=bar\nBAZ=qux\n")
	res, err := Rename(p, "FOO", "NEW_FOO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Found {
		t.Fatal("expected key to be found")
	}
	entries, _ := LoadFile(p)
	m := ToMap(entries)
	if _, ok := m["FOO"]; ok {
		t.Error("old key should not exist")
	}
	if m["NEW_FOO"] != "bar" {
		t.Errorf("expected NEW_FOO=bar, got %q", m["NEW_FOO"])
	}
}

func TestRename_KeyNotFound(t *testing.T) {
	p := writeTempEnvForRename(t, "FOO=bar\n")
	res, err := Rename(p, "MISSING", "OTHER")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Found {
		t.Error("expected key not to be found")
	}
	if res.Summary() != `key "MISSING" not found` {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestRename_SameKey(t *testing.T) {
	p := writeTempEnvForRename(t, "FOO=bar\n")
	res, err := Rename(p, "FOO", "FOO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Found {
		t.Error("expected found=true for same key rename")
	}
}

func TestRename_FileNotFound(t *testing.T) {
	_, err := Rename("/nonexistent/.env", "FOO", "BAR")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRenameResult_Summary(t *testing.T) {
	r := RenameResult{OldKey: "A", NewKey: "B", Found: true}
	if r.Summary() != `renamed "A" -> "B"` {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}
