package env_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
)

func TestPromote_IntegrationRoundTrip(t *testing.T) {
	dir := t.TempDir()
	srcFile := filepath.Join(dir, "src.env")
	dstFile := filepath.Join(dir, "dst.env")

	_ = os.WriteFile(srcFile, []byte("NEW_KEY=value1\nSHARED=from_src\n"), 0o600)
	_ = os.WriteFile(dstFile, []byte("EXISTING=keep\nSHARED=from_dst\n"), 0o600)

	srcMap, err := env.LoadFile(srcFile)
	if err != nil {
		t.Fatalf("load src: %v", err)
	}
	dstMap, err := env.LoadFile(dstFile)
	if err != nil {
		t.Fatalf("load dst: %v", err)
	}

	merged, result := env.Promote(srcMap, dstMap, env.PromoteOptions{Force: false})

	if merged["EXISTING"] != "keep" {
		t.Error("EXISTING should be preserved")
	}
	if merged["SHARED"] != "from_dst" {
		t.Error("SHARED should not be overwritten without force")
	}
	if merged["NEW_KEY"] != "value1" {
		t.Error("NEW_KEY should be promoted")
	}

	var foundShared bool
	for _, k := range result.Skipped {
		if k == "SHARED" {
			foundShared = true
		}
	}
	if !foundShared {
		t.Error("SHARED should appear in skipped list")
	}

	if err := env.WriteFile(dstFile, merged); err != nil {
		t.Fatalf("write dst: %v", err)
	}
	data, _ := os.ReadFile(dstFile)
	if len(data) == 0 {
		t.Error("dst file should not be empty after write")
	}
}
