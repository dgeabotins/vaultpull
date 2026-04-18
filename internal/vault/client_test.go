package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockVaultServer(t *testing.T, data map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{
				"data": data,
				"metadata": map[string]interface{}{},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func TestNew_InvalidAddress(t *testing.T) {
	_, err := New("://bad-address", "token")
	if err == nil {
		t.Fatal("expected error for invalid address, got nil")
	}
}

func TestNew_ValidAddress(t *testing.T) {
	c, err := New("http://127.0.0.1:8200", "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestGetSecrets_Success(t *testing.T) {
	srv := mockVaultServer(t, map[string]interface{}{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	})
	defer srv.Close()

	c, err := New(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	secrets, err := c.GetSecrets(context.Background(), "secret", "myapp/dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if secrets["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", secrets["DB_HOST"])
	}
	if secrets["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", secrets["DB_PORT"])
	}
}
