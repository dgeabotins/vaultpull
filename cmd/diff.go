package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/diff"
	"github.com/yourusername/vaultpull/internal/envwriter"
	"github.com/yourusername/vaultpull/internal/vault"
)

func init() {
	diffCmd := &cobra.Command{
		Use:   "diff",
		Short: "Show differences between local .env and Vault secrets",
		RunE:  runDiff,
	}
	diffCmd.Flags().StringP("profile", "p", "default", "Profile to use")
	diffCmd.Flags().StringP("config", "c", "vaultpull.yaml", "Config file path")
	RootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, _ []string) error {
	profileName, _ := cmd.Flags().GetString("profile")
	configPath, _ := cmd.Flags().GetString("config")

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	profile, err := cfg.GetProfile(profileName)
	if err != nil {
		return fmt.Errorf("profile %q not found", profileName)
	}

	client, err := vault.New(profile.Address, profile.Token)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	incoming, err := client.GetSecrets(profile.SecretPath)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	existing, _ := envwriter.ReadEnvFilePublic(profile.OutputFile)

	r := diff.Compare(existing, incoming)
	fmt.Fprintln(os.Stdout, diff.Summary(r))
	for k := range r.Added {
		fmt.Fprintf(os.Stdout, "  + %s\n", k)
	}
	for k := range r.Changed {
		fmt.Fprintf(os.Stdout, "  ~ %s\n", k)
	}
	for k := range r.Removed {
		fmt.Fprintf(os.Stdout, "  - %s\n", k)
	}
	return nil
}
