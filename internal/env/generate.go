package env

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
)

// GenerateOptions controls secret generation behaviour.
type GenerateOptions struct {
	Length     int
	Charset    string // "alpha", "alphanum", "hex", "base64", "symbol"
	Keys       []string
	Overwrite  bool
	DryRun     bool
}

// GenerateResult holds the outcome of a Generate call.
type GenerateResult struct {
	Generated []string
	Skipped   []string
}

func (r GenerateResult) Summary() string {
	parts := []string{}
	if len(r.Generated) > 0 {
		parts = append(parts, fmt.Sprintf("%d generated", len(r.Generated)))
	}
	if len(r.Skipped) > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped (already set)", len(r.Skipped)))
	}
	if len(parts) == 0 {
		return "nothing to generate"
	}
	return strings.Join(parts, ", ")
}

// Generate fills the given env map with random values for each key.
func Generate(env map[string]string, opts GenerateOptions) (map[string]string, GenerateResult, error) {
	if opts.Length <= 0 {
		opts.Length = 32
	}
	if opts.Charset == "" {
		opts.Charset = "base64"
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var result GenerateResult
	for _, key := range opts.Keys {
		if _, exists := out[key]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, key)
			continue
		}
		val, err := randomValue(opts.Length, opts.Charset)
		if err != nil {
			return nil, result, fmt.Errorf("generate %s: %w", key, err)
		}
		if !opts.DryRun {
			out[key] = val
		}
		result.Generated = append(result.Generated, key)
	}
	return out, result, nil
}

func randomValue(length int, charset string) (string, error) {
	switch charset {
	case "hex":
		b := make([]byte, (length+1)/2)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		return fmt.Sprintf("%x", b)[:length], nil
	case "base64":
		b := make([]byte, length)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		return base64.URLEncoding.EncodeToString(b)[:length], nil
	case "alpha":
		return randomFromAlphabet(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	case "alphanum":
		return randomFromAlphabet(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	case "symbol":
		return randomFromAlphabet(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*")
	default:
		return "", fmt.Errorf("unknown charset %q", charset)
	}
}

func randomFromAlphabet(length int, alphabet string) (string, error) {
	var sb strings.Builder
	sb.Grow(length)
	max := big.NewInt(int64(len(alphabet)))
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		sb.WriteByte(alphabet[n.Int64()])
	}
	return sb.String(), nil
}
