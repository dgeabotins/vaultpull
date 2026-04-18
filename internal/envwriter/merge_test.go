package envwriter

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	return p
}

func TestMerge_OverlaysNewSecrets(t *testing.T) {
	path := writeTempEnv(t, "EXISTING=old\nKEEP=me\n")

	merged, err := Merge(path, map[string]string{"EXISTING": "new", "ADDED": "yes"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if merged["EXISTING"] != "new" {
		t.Errorf("expected EXISTING=new, got %q", merged["EXISTING"])
	}
	if merged["KEEP"] != "me" {
		t.Errorf("expected KEEP=me, got %q", merged["KEEP"])
	}
	if merged["ADDED"] != "yes" {
		t.Errorf("expected ADDED=yes, got %q", merged["ADDED"])
	}
}

func TestMerge_NonExistentFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")

	merged, err := Merge(path, map[string]string{"KEY": "value"})
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}
	if merged["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", merged["KEY"])
	}
}

func TestReadEnvFile_IgnoresComments(t *testing.T) {
	path := writeTempEnv(t, "# comment\nVALID=yes\n\nANOTHER=1\n")

	result, err := readEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}
