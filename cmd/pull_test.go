package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, addr string) string {
	t.Helper()
	content := `profiles:
  default:
    address: ` + addr + `
    token: test-token
    secrets:
      - secret/data/app
`
	f, err := os.CreateTemp(t.TempDir(), "vaultpull-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestPullCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"data":{"KEY":"value"}}}`))
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	cfgFile = writeTempConfig(t, server.URL)
	profile = "default"
	output = filepath.Join(tmpDir, ".env")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"pull", "--config", cfgFile, "--profile", profile, "--output", output})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Error("expected output file to be created")
	}
}
