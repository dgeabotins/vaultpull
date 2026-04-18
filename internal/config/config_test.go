package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/vaultpull/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultpull-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

const validConfig = `
default_profile: dev
profiles:
  - name: dev
    vault_addr: http://127.0.0.1:8200
    vault_path: secret/data/myapp/dev
    output_file: .env.dev
  - name: prod
    vault_addr: https://vault.example.com
    vault_path: secret/data/myapp/prod
    output_file: .env.prod
    mapping:
      DB_PASSWORD: database_password
`

func TestLoad_Valid(t *testing.T) {
	path := writeTemp(t, validConfig)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultProfile != "dev" {
		t.Errorf("expected default_profile=dev, got %q", cfg.DefaultProfile)
	}
	if len(cfg.Profiles) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(cfg.Profiles))
	}
}

func TestGetProfile_Found(t *testing.T) {
	path := writeTemp(t, validConfig)
	cfg, _ := config.Load(path)

	p, err := cfg.GetProfile("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.OutputFile != ".env.prod" {
		t.Errorf("expected .env.prod, got %q", p.OutputFile)
	}
}

func TestGetProfile_Default(t *testing.T) {
	path := writeTemp(t, validConfig)
	cfg, _ := config.Load(path)

	p, err := cfg.GetProfile("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "dev" {
		t.Errorf("expected dev profile, got %q", p.Name)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	path := writeTemp(t, validConfig)
	cfg, _ := config.Load(path)

	_, err := cfg.GetProfile("staging")
	if err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := writeTemp(t, ": bad: yaml: [")
	_, err := config.Load(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoad_DuplicateProfile(t *testing.T) {
	dup := `
profiles:
  - name: dev
    vault_addr: http://127.0.0.1:8200
    vault_path: secret/data/app
    output_file: .env
  - name: dev
    vault_addr: http://127.0.0.1:8200
    vault_path: secret/data/app2
    output_file: .env2
`
	path := writeTemp(t, dup)
	_, err := config.Load(path)
	if err == nil {
		t.Error("expected error for duplicate profile names")
	}
}
