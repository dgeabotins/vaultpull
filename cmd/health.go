package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/health"
)

func init() {
	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "Check connectivity to Vault and local configuration",
		RunE:  runHealth,
	}
	rootCmd.AddCommand(healthCmd)
}

func runHealth(cmd *cobra.Command, args []string) error {
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		vaultAddr = "http://127.0.0.1:8200"
	}

	checker := health.New(vaultAddr)
	report := checker.Run()

	for _, s := range report.Statuses {
		icon := "✓"
		if !s.OK {
			icon = "✗"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s %s: %s\n", icon, s.Name, s.Message)
	}

	if !report.Healthy {
		fmt.Fprintln(cmd.OutOrStdout(), "\nHealth check failed.")
		return fmt.Errorf("one or more checks failed")
	}

	fmt.Fprintln(cmd.OutOrStdout(), "\nAll checks passed.")
	return nil
}
