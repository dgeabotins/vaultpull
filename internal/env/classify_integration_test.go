package env_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultpull/internal/env"
)

func TestClassify_IntegrationRoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")

	content := strings.Join([]string{
		"DATABASE_URL=https://db.example.com",
		"PORT=5432",
		"RATIO=1.5",
		"DEBUG=true",
		"CONFIG_PATH=/etc/app.yaml",
		"SETTINGS={\"timeout\":30}",
		"EMPTY_VAR=",
		"API_TOKEN=abc123secret",
		"APP_NAME=myapp",
	}, "\n")

	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write: %v", err)
	}

	m, err := env.LoadFile(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	results := env.Classify(m)
	if len(results) != 9 {
		t.Fatalf("expected 9 results, got %d", len(results))
	}

	byKey := make(map[string]env.Category)
	for _, r := range results {
		byKey[r.Key] = r.Category
	}

	expect := map[string]env.Category{
		"DATABASE_URL": env.CategoryURL,
		"PORT":         env.CategoryInteger,
		"RATIO":        env.CategoryFloat,
		"DEBUG":        env.CategoryBoolean,
		"CONFIG_PATH":  env.CategoryPath,
		"SETTINGS":     env.CategoryJSON,
		"EMPTY_VAR":    env.CategoryEmpty,
		"API_TOKEN":    env.CategorySecret,
		"APP_NAME":     env.CategoryUnknown,
	}

	for k, want := range expect {
		if got := byKey[k]; got != want {
			t.Errorf("key %s: expected %s, got %s", k, want, got)
		}
	}
}
