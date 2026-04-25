package env_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"vaultpull/internal/env"
)

func TestRedact_RoundTripWriteRead(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")

	initial := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "key-abc",
		"APP_ENV":     "production",
	}
	if err := env.WriteFile(file, initial); err != nil {
		t.Fatalf("write: %v", err)
	}

	loaded, err := env.LoadFile(file)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	res, err := env.Redact(loaded, env.RedactOptions{
		Patterns: []string{"(?i)(password|key)"},
	})
	if err != nil {
		t.Fatalf("redact: %v", err)
	}

	if err := env.WriteFile(file, res.Values); err != nil {
		t.Fatalf("write redacted: %v", err)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	content := string(data)

	if strings.Contains(content, "supersecret") {
		t.Error("DB_PASSWORD value should be redacted")
	}
	if strings.Contains(content, "key-abc") {
		t.Error("API_KEY value should be redacted")
	}
	if !strings.Contains(content, "production") {
		t.Error("APP_ENV should remain unchanged")
	}
	if len(res.Redacted) != 2 {
		t.Errorf("expected 2 redacted keys, got %d: %v", len(res.Redacted), res.Redacted)
	}
}
