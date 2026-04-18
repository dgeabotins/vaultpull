package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	profile string
	output  string
)

var rootCmd = &cobra.Command{
	Use:   "vaultpull",
	Short: "Sync HashiCorp Vault secrets into local .env files",
	Long: `vaultpull fetches secrets from HashiCorp Vault and writes
them into local .env files, with support for multiple profiles.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "vaultpull.yaml", "config file path")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "profile to use")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", ".env", "output .env file path")
}
