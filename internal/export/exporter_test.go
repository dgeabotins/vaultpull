package export

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []string{"dotenv", "json", "export", "JSON", "Dotenv"} {
		_, err := New(f)
		if err != nil {
			t.Errorf("expected no error for format %q, got %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := New("yaml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_Dotenv_ToFile(t *testing.T) {
	e, _ := New("dotenv")
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")

	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := e.Write(secrets, path); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "FOO=") || !strings.Contains(content, "BAZ=") {
		t.Errorf("unexpected dotenv output: %s", content)
	}
}

func TestWrite_JSON_ContainsBraces(t *testing.T) {
	e, _ := New("json")
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")

	secrets := map[string]string{"KEY": "value"}
	if err := e.Write(secrets, path); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.HasPrefix(content, "{") || !strings.Contains(content, "KEY") {
		t.Errorf("unexpected JSON output: %s", content)
	}
}

func TestWrite_Export_HasExportPrefix(t *testing.T) {
	e, _ := New("export")
	dir := t.TempDir()
	path := filepath.Join(dir, "out.sh")

	secrets := map[string]string{"MY_VAR": "hello"}
	if err := e.Write(secrets, path); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "export MY_VAR=") {
		t.Errorf("expected export prefix, got: %s", string(data))
	}
}

func TestWrite_FilePermissions(t *testing.T) {
	e, _ := New("dotenv")
	dir := t.TempDir()
	path := filepath.Join(dir, "secure.env")

	if err := e.Write(map[string]string{"A": "b"}, path); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	info, _ := os.Stat(path)
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
