package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/export"
	"vaultpull/internal/vault"
)

var (
	exportFormat  string
	exportOutput  string
	exportProfile string
)

func init() {
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export Vault secrets in a chosen format (dotenv, json, export)",
		RunE:  runExport,
	}
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "dotenv", "Output format: dotenv, json, export")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "-", "Output file path (use '-' for stdout)")
	exportCmd.Flags().StringVarP(&exportProfile, "profile", "p", "default", "Config profile to use")
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Root().PersistentFlags().GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	profile, err := cfg.GetProfile(exportProfile)
	if err != nil {
		return fmt.Errorf("profile %q not found: %w", exportProfile, err)
	}

	client, err := vault.New(profile.Address, profile.Token)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := client.GetSecrets(profile.Path)
	if err != nil {
		return fmt.Errorf("fetching secrets: %w", err)
	}

	if len(secrets) == 0 {
		fmt.Fprintf(os.Stderr, "warning: no secrets found at path %q\n", profile.Path)
	}

	ex, err := export.New(exportFormat)
	if err != nil {
		return err
	}

	if err := ex.Write(secrets, exportOutput); err != nil {
		fmt.Fprintf(os.Stderr, "export error: %v\n", err)
		return err
	}
	return nil
}
