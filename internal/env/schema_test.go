package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempSchema(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.schema")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestApplySchema_FillsDefaults(t *testing.T) {
	schema := []SchemaEntry{
		{Key: "APP_ENV", Required: false, Default: "production"},
		{Key: "PORT", Required: true},
	}
	data := map[string]string{"PORT": "8080"}
	out, result := ApplySchema(data, schema)
	if result.HasErrors() {
		t.Fatalf("unexpected errors: %v", result.Missing)
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected default 'production', got %q", out["APP_ENV"])
	}
	if _, ok := result.Defaults["APP_ENV"]; !ok {
		t.Error("expected APP_ENV in defaults map")
	}
}

func TestApplySchema_MissingRequired(t *testing.T) {
	schema := []SchemaEntry{
		{Key: "DATABASE_URL", Required: true},
	}
	_, result := ApplySchema(map[string]string{}, schema)
	if !result.HasErrors() {
		t.Fatal("expected errors for missing required key")
	}
	if len(result.Missing) != 1 || result.Missing[0] != "DATABASE_URL" {
		t.Errorf("unexpected missing: %v", result.Missing)
	}
}

func TestApplySchema_OptionalNoDefault(t *testing.T) {
	schema := []SchemaEntry{
		{Key: "DEBUG", Required: false},
	}
	_, result := ApplySchema(map[string]string{}, schema)
	if result.HasErrors() {
		t.Error("optional key without default should not cause error")
	}
}

func TestSummary_NoIssues(t *testing.T) {
	r := SchemaResult{Defaults: map[string]string{}}
	if r.Summary() != "schema ok" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestSummary_WithIssues(t *testing.T) {
	r := SchemaResult{
		Missing:  []string{"SECRET_KEY"},
		Defaults: map[string]string{"APP_ENV": "production"},
	}
	s := r.Summary()
	if s == "schema ok" {
		t.Error("expected non-ok summary")
	}
}

func TestLoadSchema_ParsesEntries(t *testing.T) {
	content := "# comment\nREQUIRED_KEY\n?OPTIONAL_KEY\nDEFAULT_KEY=mydefault\n"
	path := writeTempSchema(t, content)
	entries, err := LoadSchema(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if !entries[0].Required {
		t.Error("REQUIRED_KEY should be required")
	}
	if entries[1].Required {
		t.Error("OPTIONAL_KEY should not be required")
	}
	if entries[2].Default != "mydefault" {
		t.Errorf("expected default 'mydefault', got %q", entries[2].Default)
	}
}

func TestLoadSchema_FileNotFound(t *testing.T) {
	_, err := LoadSchema(filepath.Join(t.TempDir(), "missing.schema"))
	if err == nil {
		t.Error("expected error for missing file")
	}
}
