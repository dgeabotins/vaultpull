package health

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func mockVault(code int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	}))
}

func TestRun_VaultReachable(t *testing.T) {
	srv := mockVault(http.StatusOK)
	defer srv.Close()

	c := New(srv.URL)
	report := c.Run()

	if !report.Healthy {
		t.Fatal("expected healthy report")
	}
	if len(report.Statuses) == 0 {
		t.Fatal("expected at least one status")
	}
	if !report.Statuses[0].OK {
		t.Errorf("vault status not OK: %s", report.Statuses[0].Message)
	}
}

func TestRun_VaultUnreachable(t *testing.T) {
	c := New("http://127.0.0.1:19999")
	report := c.Run()

	if report.Healthy {
		t.Fatal("expected unhealthy report")
	}
	if report.Statuses[0].OK {
		t.Error("expected vault status to be not OK")
	}
}

func TestRun_VaultBadStatus(t *testing.T) {
	srv := mockVault(http.StatusInternalServerError)
	defer srv.Close()

	c := New(srv.URL)
	report := c.Run()

	if report.Healthy {
		t.Fatal("expected unhealthy report for 500")
	}
	if !strings.Contains(report.Statuses[0].Message, "500") {
		t.Errorf("expected message to contain status code, got: %s", report.Statuses[0].Message)
	}
}

func TestNew_DefaultTimeout(t *testing.T) {
	c := New("http://localhost:8200")
	if c.client == nil {
		t.Fatal("expected non-nil http client")
	}
}
