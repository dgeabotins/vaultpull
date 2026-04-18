package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/sync"
	"github.com/user/vaultpull/internal/vault"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull secrets from Vault and write to .env file",
	RunE:  runPull,
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	p, err := cfg.GetProfile(profile)
	if err != nil {
		return fmt.Errorf("profile %q not found: %w", profile, err)
	}

	client, err := vault.New(p.Address, p.Token)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	syncer := sync.New(cfg, client, output)
	if err := syncer.Run(profile); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}

	fmt.Printf("Secrets written to %s\n", output)
	return nil
}
