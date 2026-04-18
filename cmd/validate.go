package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/validate"
)

func init() {
	var requiredKeys []string
	var warnOnly bool

	validateCmd := &cobra.Command{
		Use:   "validate [env-file]",
		Short: "Validate a .env file against required keys",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate(args[0], requiredKeys, warnOnly)
		},
	}

	validateCmd.Flags().StringSliceVarP(&requiredKeys, "keys", "k", nil, "Required keys to check (comma-separated)")
	validateCmd.Flags().BoolVarP(&warnOnly, "warn", "w", false, "Warn instead of failing on missing keys")

	rootCmd.AddCommand(validateCmd)
}

func runValidate(envFile string, requiredKeys []string, warnOnly bool) error {
	results := validate.CheckEnvFile(envFile, requiredKeys)

	for _, r := range results {
		switch r.Status {
		case validate.StatusOK:
			fmt.Printf("  ✓ %s\n", r.Key)
		case validate.StatusMissing:
			fmt.Printf("  ✗ %s — missing\n", r.Key)
		case validate.StatusEmpty:
			fmt.Printf("  ⚠ %s — empty value\n", r.Key)
		}
	}

	if validate.HasFailures(results) {
		if warnOnly {
			fmt.Fprintln(os.Stderr, "warning: validation issues found")
			return nil
		}
		return fmt.Errorf("validation failed: one or more required keys are missing or empty")
	}

	fmt.Println("validation passed")
	return nil
}
