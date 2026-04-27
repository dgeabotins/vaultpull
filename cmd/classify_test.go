package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForClassify(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write: %v", err)
	}
	return p
}

func executeClassifyCmd(t *testing.T, args ...string) string {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"classify"}, args...))
	_ = rootCmd.Execute()
	return buf.String()
}

func TestClassifyCmd_ShowsTable(t *testing.T) {
	p := writeTempEnvForClassify(t, "PORT=8080\nDATABASE_URL=https://db.example.com\nAPP_NAME=myapp\n")
	out := executeClassifyCmd(t, p)
	if !strings.Contains(out, "integer") {
		t.Errorf("expected 'integer' in output, got: %s", out)
	}
	if !strings.Contains(out, "url") {
		t.Errorf("expected 'url' in output, got: %s", out)
	}
}

func TestClassifyCmd_Summary(t *testing.T) {
	p := writeTempEnvForClassify(t, "PORT=8080\nDEBUG=true\nHOST=localhost\n")
	out := executeClassifyCmd(t, "--summary", p)
	if !strings.Contains(out, "integer") {
		t.Errorf("expected 'integer' in summary, got: %s", out)
	}
	if !strings.Contains(out, "boolean") {
		t.Errorf("expected 'boolean' in summary, got: %s", out)
	}
}

func TestClassifyCmd_Filter(t *testing.T) {
	p := writeTempEnvForClassify(t, "PORT=8080\nDATABASE_URL=https://db.example.com\nDEBUG=true\n")
	out := executeClassifyCmd(t, "--filter", "url", p)
	if !strings.Contains(out, "DATABASE_URL") {
		t.Errorf("expected DATABASE_URL in output, got: %s", out)
	}
	if strings.Contains(out, "PORT") {
		t.Errorf("PORT should be filtered out, got: %s", out)
	}
}

func TestClassifyCmd_SecretMasked(t *testing.T) {
	p := writeTempEnvForClassify(t, "API_SECRET=supersensitive\n")
	out := executeClassifyCmd(t, p)
	if strings.Contains(out, "supersensitive") {
		t.Errorf("secret value should be masked, got: %s", out)
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected *** mask in output, got: %s", out)
	}
}

func TestClassifyCmd_MissingFile(t *testing.T) {
	buf := &bytes.Buffer{}
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	err := runClassify(cmd, []string{"/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
