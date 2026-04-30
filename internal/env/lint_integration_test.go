package env_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultpull/vaultpull/internal/env"
)

func TestLint_IntegrationRoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")

	content := "DATABASE_URL=postgres://localhost/db\nAPI_KEY=secret\nBAD KEY=oops\nEMPTY=\n"
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	entries, err := env.LoadFile(p)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	m := env.ToMap(entries)
	result := env.LintMap(m)

	// BAD KEY → invalid key (error)
	// EMPTY   → empty value (warning)
	var errCount, warnCount int
	for _, issue := range result.Issues {
		switch issue.Severity {
		case env.LintError:
			errCount++
		case env.LintWarning:
			warnCount++
		}
	}

	if errCount < 1 {
		t.Errorf("expected at least 1 error, got %d", errCount)
	}
	if warnCount < 1 {
		t.Errorf("expected at least 1 warning, got %d", warnCount)
	}
	if !result.HasErrors() {
		t.Error("expected HasErrors=true")
	}
	if result.Summary() == "no lint issues found" {
		t.Error("expected non-empty summary")
	}
}
