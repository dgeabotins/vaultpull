package template

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempTemplate(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRender_ProducesOutput(t *testing.T) {
	tmplPath := writeTempTemplate(t, "DB_HOST={{ index . \"DB_HOST\" }}\nDB_PORT={{ index . \"DB_PORT\" }}\n")
	outPath := filepath.Join(t.TempDir(), ".env")

	r := New(outPath)
	err := r.Render(tmplPath, map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(outPath)
	got := string(data)
	if got != "DB_HOST=localhost\nDB_PORT=5432\n" {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestRender_MissingKey(t *testing.T) {
	tmplPath := writeTempTemplate(t, "VAL={{ index . \"MISSING\" }}\n")
	outPath := filepath.Join(t.TempDir(), ".env")

	r := New(outPath)
	err := r.Render(tmplPath, map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRender_FilePermissions(t *testing.T) {
	tmplPath := writeTempTemplate(t, "KEY=value\n")
	outPath := filepath.Join(t.TempDir(), ".env")

	r := New(outPath)
	if err := r.Render(tmplPath, map[string]string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestRender_MissingTemplate(t *testing.T) {
	r := New("/tmp/out.env")
	err := r.Render("/nonexistent/template.tmpl", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing template file")
	}
}
