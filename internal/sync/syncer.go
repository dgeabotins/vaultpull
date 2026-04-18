package sync

import (
	"fmt"

	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/envwriter"
	"github.com/user/vaultpull/internal/vault"
)

// Syncer orchestrates fetching secrets from Vault and writing them to a .env file.
type Syncer struct {
	cfg    *config.Config
	client *vault.Client
}

// New creates a new Syncer for the given config and vault client.
func New(cfg *config.Config, client *vault.Client) *Syncer {
	return &Syncer{cfg: cfg, client: client}
}

// Run executes the sync for the given profile name.
// It fetches secrets from Vault and merges them into the target .env file.
func (s *Syncer) Run(profileName string) error {
	profile, err := s.cfg.GetProfile(profileName)
	if err != nil {
		return fmt.Errorf("profile %q not found: %w", profileName, err)
	}

	secrets, err := s.client.GetSecrets(profile.Path)
	if err != nil {
		return fmt.Errorf("failed to fetch secrets from vault path %q: %w", profile.Path, err)
	}

	output := profile.Output
	if output == "" {
		output = ".env"
	}

	merged, err := envwriter.Merge(output, secrets)
	if err != nil {
		return fmt.Errorf("failed to merge secrets into %q: %w", output, err)
	}

	w := envwriter.New(output)
	if err := w.Write(merged); err != nil {
		return fmt.Errorf("failed to write %q: %w", output, err)
	}

	return nil
}
