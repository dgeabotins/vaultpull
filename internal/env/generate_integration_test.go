package env

import (
	"path/filepath"
	"testing"
)

func TestGenerate_IntegrationRoundTrip(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")

	// First pass: generate two keys into an empty file.
	initial := map[string]string{}
	opts := GenerateOptions{
		Keys:    []string{"DB_PASSWORD", "SESSION_SECRET"},
		Length:  40,
		Charset: "symbol",
	}
	out, result, err := Generate(initial, opts)
	if err != nil {
		t.Fatalf("generate pass 1: %v", err)
	}
	if len(result.Generated) != 2 {
		t.Fatalf("expected 2 generated, got %d", len(result.Generated))
	}
	if err := WriteFile(file, out, 0600); err != nil {
		t.Fatalf("write: %v", err)
	}

	// Second pass: reload and generate same keys without overwrite — should skip.
	loaded, err := LoadFile(file)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	origDB := loaded["DB_PASSWORD"]

	opts2 := GenerateOptions{
		Keys:      []string{"DB_PASSWORD", "SESSION_SECRET"},
		Length:    40,
		Charset:   "symbol",
		Overwrite: false,
	}
	out2, result2, err := Generate(loaded, opts2)
	if err != nil {
		t.Fatalf("generate pass 2: %v", err)
	}
	if len(result2.Skipped) != 2 {
		t.Errorf("expected 2 skipped on second pass, got %d", len(result2.Skipped))
	}
	if out2["DB_PASSWORD"] != origDB {
		t.Errorf("DB_PASSWORD should be unchanged on second pass")
	}

	// Third pass: overwrite both keys.
	opts3 := GenerateOptions{
		Keys:      []string{"DB_PASSWORD"},
		Length:    40,
		Charset:   "hex",
		Overwrite: true,
	}
	out3, result3, err := Generate(loaded, opts3)
	if err != nil {
		t.Fatalf("generate pass 3: %v", err)
	}
	if len(result3.Generated) != 1 {
		t.Errorf("expected 1 generated in pass 3")
	}
	if out3["DB_PASSWORD"] == origDB {
		t.Errorf("expected DB_PASSWORD to change on overwrite")
	}
}
