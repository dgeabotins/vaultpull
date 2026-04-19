package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvForCopy(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(f.Name(), []byte(content), 0600)
	return f.Name()
}

func TestCopy_AllKeys(t *testing.T) {
	src := writeTempEnvForCopy(t, "FOO=bar\nBAZ=qux\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	res, err := Copy(src, dst, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Keys))
	}
	m, _ := LoadFile(dst)
	if m["FOO"] != "bar" {
		t.Errorf("expected FOO=bar")
	}
}

func TestCopy_SelectedKeys(t *testing.T) {
	src := writeTempEnvForCopy(t, "FOO=bar\nBAZ=qux\nSECRET=x\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	res, err := Copy(src, dst, []string{"FOO", "BAZ"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Keys) != 2 {
		t.Errorf("expected 2 keys copied")
	}
	m, _ := LoadFile(dst)
	if _, ok := m["SECRET"]; ok {
		t.Errorf("SECRET should not have been copied")
	}
}

func TestCopy_MissingKey(t *testing.T) {
	src := writeTempEnvForCopy(t, "FOO=bar\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	_, err := Copy(src, dst, []string{"MISSING"})
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestCopy_Summary(t *testing.T) {
	r := CopyResult{Source: "a.env", Dest: "b.env", Keys: []string{"X", "Y"}}
	s := r.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestCopy_SourceNotFound(t *testing.T) {
	_, err := Copy("/no/such/file.env", "/tmp/dst.env", nil)
	if err == nil {
		t.Error("expected error for missing source")
	}
}
