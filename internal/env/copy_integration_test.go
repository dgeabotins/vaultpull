package env_test

import (
	"os"
	"path/filepath"
	"testing"

	"vaultpull/internal/env"
)

func TestCopy_MergesIntoExistingFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.env")
	dst := filepath.Join(dir, "dst.env")

	_ = os.WriteFile(src, []byte("NEW=value\nFOO=updated\n"), 0600)
	_ = os.WriteFile(dst, []byte("FOO=original\nKEEP=me\n"), 0600)

	_, err := env.Copy(src, dst, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m, err := env.LoadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if m["FOO"] != "updated" {
		t.Errorf("expected FOO=updated, got %s", m["FOO"])
	}
	if m["KEEP"] != "me" {
		t.Errorf("expected KEEP=me to be preserved")
	}
	if m["NEW"] != "value" {
		t.Errorf("expected NEW=value to be added")
	}
}
