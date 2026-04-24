package env

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const encryptedPrefix = "enc:"

// EncryptResult holds the outcome of an encryption pass over an env map.
type EncryptResult struct {
	Encrypted int
	Skipped   int
}

func (r EncryptResult) Summary() string {
	if r.Encrypted == 0 {
		return "no values encrypted"
	}
	return fmt.Sprintf("%d value(s) encrypted, %d skipped", r.Encrypted, r.Skipped)
}

// EncryptValues encrypts all plaintext values in m using AES-GCM with the
// provided 32-byte key. Values already prefixed with "enc:" are skipped.
func EncryptValues(m map[string]string, key []byte) (map[string]string, EncryptResult, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, EncryptResult{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, EncryptResult{}, err
	}

	out := make(map[string]string, len(m))
	var res EncryptResult
	for k, v := range m {
		if strings.HasPrefix(v, encryptedPrefix) {
			out[k] = v
			res.Skipped++
			continue
		}
		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, EncryptResult{}, err
		}
		ciphertext := gcm.Seal(nonce, nonce, []byte(v), nil)
		out[k] = encryptedPrefix + base64.StdEncoding.EncodeToString(ciphertext)
		res.Encrypted++
	}
	return out, res, nil
}

// DecryptValues decrypts all values prefixed with "enc:" using AES-GCM.
func DecryptValues(m map[string]string, key []byte) (map[string]string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string, len(m))
	for k, v := range m {
		if !strings.HasPrefix(v, encryptedPrefix) {
			out[k] = v
			continue
		}
		raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(v, encryptedPrefix))
		if err != nil {
			return nil, fmt.Errorf("base64 decode failed for key %q: %w", k, err)
		}
		if len(raw) < gcm.NonceSize() {
			return nil, errors.New("ciphertext too short for key: " + k)
		}
		nonce, ct := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
		plain, err := gcm.Open(nil, nonce, ct, nil)
		if err != nil {
			return nil, fmt.Errorf("decryption failed for key %q: %w", k, err)
		}
		out[k] = string(plain)
	}
	return out, nil
}
