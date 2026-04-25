package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

func writeTempEnvForGenerate(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func executeGenerateCmd(args []string) (string, error) {
	buf := new(strings.Builder)
	root := &cobra.Command{Use: "vaultpull"}
	var (
		outFile   string
		length    int
		charset   string
		overwrite bool
		dryRun    bool
	)
	cmd := &cobra.Command{
		Use:  "generate",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, a []string) error {
			return runGenerate(a, outFile, length, charset, overwrite, dryRun)
		},
	}
	cmd.Flags().StringVarP(&outFile, "file", "f", ".env", "")
	cmd.Flags().IntVarP(&length, "length", "l", 32, "")
	cmd.Flags().StringVarP(&charset, "charset", "c", "base64", "")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "")
	root.AddCommand(cmd)
	root.SetOut(buf)
	root.SetArgs(append([]string{"generate"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestGenerateCmd_DryRun(t *testing.T) {
	f := writeTempEnvForGenerate(t, "EXISTING=hello\n")
	_, err := executeGenerateCmd([]string{"--file", f, "--dry-run", "NEW_SECRET"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := env.LoadFile(f)
	if _, ok := data["NEW_SECRET"]; ok {
		t.Errorf("dry run should not write NEW_SECRET")
	}
}

func TestGenerateCmd_WritesNewKey(t *testing.T) {
	f := writeTempEnvForGenerate(t, "")
	_, err := executeGenerateCmd([]string{"--file", f, "--length", "24", "--charset", "alphanum", "MY_SECRET"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := env.LoadFile(f)
	if err != nil {
		t.Fatalf("load file: %v", err)
	}
	if len(data["MY_SECRET"]) != 24 {
		t.Errorf("expected length 24, got %d", len(data["MY_SECRET"]))
	}
}

func TestGenerateCmd_SkipsExisting(t *testing.T) {
	f := writeTempEnvForGenerate(t, "ALREADY=set\n")
	_, err := executeGenerateCmd([]string{"--file", f, "ALREADY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := env.LoadFile(f)
	if data["ALREADY"] != "set" {
		t.Errorf("expected ALREADY to remain 'set', got %q", data["ALREADY"])
	}
}
