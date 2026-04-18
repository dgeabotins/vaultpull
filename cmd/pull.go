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
		return fmt.Errorf("syncing secrets: %w", err)
	}

	fmt.Printf("Secrets written to %s\n", output)
	return nil
}

// confirmOverwrite checks if the output file already exists and prompts the
// user for confirmation before overwriting it. Returns true if it is safe to
// proceed.
func confirmOverwrite(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("checking output file: %w", err)
	}

	fmt.Printf("File %s already exists. Overwrite? [y/N]: ", path)
	var answer string
	fmt.Fscan(os.Stdin, &answer)
	return answer == "y" || answer == "Y", nil
}
