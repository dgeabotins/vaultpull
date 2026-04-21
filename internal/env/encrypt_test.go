package env

import (
	"strings"
	"testing"
)

func testKey() []byte {
	// 32-byte AES-256 key for tests
	return []byte("00000000000000000000000000000000")
}

func TestEncryptValues_EncryptsPlaintext(t *testing.T) {
	m := map[string]string{
		"DB_PASS": "secret",
		"API_KEY": "hunter2",
	}
	out, res, err := EncryptValues(m, testKey())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Encrypted != 2 {
		t.Errorf("expected 2 encrypted, got %d", res.Encrypted)
	}
	for k, v := range out {
		if !strings.HasPrefix(v, "enc:") {
			t.Errorf("key %s value not prefixed with enc:", k)
		}
	}
}

func TestEncryptValues_SkipsAlreadyEncrypted(t *testing.T) {
	m := map[string]string{
		"DB_PASS": "enc:alreadydone",
		"API_KEY": "plain",
	}
	_, res, err := EncryptValues(m, testKey())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
	if res.Encrypted != 1 {
		t.Errorf("expected 1 encrypted, got %d", res.Encrypted)
	}
}

func TestDecryptValues_RoundTrip(t *testing.T) {
	orig := map[string]string{
		"DB_PASS": "supersecret",
		"TOKEN":   "abc123",
	}
	encrypted, _, err := EncryptValues(orig, testKey())
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}
	decrypted, err := DecryptValues(encrypted, testKey())
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}
	for k, want := range orig {
		if got := decrypted[k]; got != want {
			t.Errorf("key %s: got %q, want %q", k, got, want)
		}
	}
}

func TestDecryptValues_LeavesPlaintextAlone(t *testing.T) {
	m := map[string]string{"HOST": "localhost"}
	out, err := DecryptValues(m, testKey())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", out["HOST"])
	}
}

func TestEncryptResult_Summary_NoEncryptions(t *testing.T) {
	r := EncryptResult{Encrypted: 0, Skipped: 2}
	if r.Summary() != "no values encrypted" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestEncryptResult_Summary_WithEncryptions(t *testing.T) {
	r := EncryptResult{Encrypted: 3, Skipped: 1}
	if !strings.Contains(r.Summary(), "3") {
		t.Errorf("summary should mention count: %s", r.Summary())
	}
}
