package whoami

import (
	"os"
	"strings"
	"testing"
)

func TestGather_ReturnsInfo(t *testing.T) {
	info, err := Gather("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Username == "" {
		t.Error("expected non-empty username")
	}
	if info.Hostname == "" {
		t.Error("expected non-empty hostname")
	}
	if info.Profile != "staging" {
		t.Errorf("expected profile 'staging', got %q", info.Profile)
	}
	if info.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestGather_VaultAddrFromEnv(t *testing.T) {
	os.Setenv("VAULT_ADDR", "https://vault.example.com")
	defer os.Unsetenv("VAULT_ADDR")

	info, err := Gather("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.VaultAddr != "https://vault.example.com" {
		t.Errorf("expected vault addr from env, got %q", info.VaultAddr)
	}
}

func TestGather_VaultAddrDefault(t *testing.T) {
	os.Unsetenv("VAULT_ADDR")
	info, err := Gather("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.VaultAddr != "(not set)" {
		t.Errorf("expected '(not set)', got %q", info.VaultAddr)
	}
}

func TestFormat_ContainsFields(t *testing.T) {
	info, _ := Gather("default")
	out := info.Format()
	for _, want := range []string{"User:", "Host:", "Vault Addr:", "Profile:", "Time:"} {
		if !strings.Contains(out, want) {
			t.Errorf("Format() missing field %q", want)
		}
	}
}
