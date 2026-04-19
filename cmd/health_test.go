package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/health"
)

func TestHealthCmd_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	t.Setenv("VAULT_ADDR", srv.URL)

	buf := &bytes.Buffer{}
	cmd := &cobra.Command{
		Use:  "health",
		RunE: runHealth,
	}
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "All checks passed") {
		t.Errorf("expected success message, got: %s", buf.String())
	}
}

func TestHealthCmd_Failure(t *testing.T) {
	checker := health.New("http://127.0.0.1:19999")
	report := checker.Run()

	if report.Healthy {
		t.Fatal("expected unhealthy")
	}
	if len(report.Statuses) == 0 {
		t.Fatal("expected statuses")
	}
	if report.Statuses[0].OK {
		t.Error("expected vault check to fail")
	}
}
