package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	vc *vaultapi.Client
}

// New creates a new Vault client using the provided address and token.
func New(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	vc, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	vc.SetToken(token)

	return &Client{vc: vc}, nil
}

// GetSecrets reads a KV v2 secret at the given path and returns key/value pairs.
func (c *Client) GetSecrets(ctx context.Context, mount, path string) (map[string]string, error) {
	secret, err := c.vc.KVv2(mount).Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("reading secret %s/%s: %w", mount, path, err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found at %s/%s", mount, path)
	}

	result := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("value for key %q is not a string", k)
		}
		result[k] = str
	}

	return result, nil
}
