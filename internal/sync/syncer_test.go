package sync_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/vaultpull/internal/config"
	internalsync "github.com/user/vaultpull/internal/sync"
	"github.com/user/vaultpull/internal/vault"
)

func mockServer(t *testing.T, secrets string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(secrets))
	}))
}

func TestRun_Success(t *testing.T) {
	svr := mockServer(t, `{"data":{"data":{"API_KEY":"abc123","DB_PASS":"secret"}}}`)
	defer svr.Close()

	cfg := &config.Config{
		Profiles: []config.Profile{
			{Name: "dev", Path: "secret/data/dev", Output: filepath.Join(t.TempDir(), ".env")},
		},
	}

	client, err := vault.New(svr.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	s := internalsync.New(cfg, client)
	if err := s.Run("dev"); err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	content, err := os.ReadFile(cfg.Profiles[0].Output)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if len(content) == 0 {
		t.Error("expected non-empty .env file")
	}
}

func TestRun_ProfileNotFound(t *testing.T) {
	cfg := &config.Config{Profiles: []config.Profile{}}
	client, _ := vault.New("http://localhost", "token")
	s := internalsync.New(cfg, client)

	if err := s.Run("missing"); err == nil {
		t.Error("expected error for missing profile, got nil")
	}
}
