package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/vaultpull/internal/env"
)

var lintFile string
var lintWarnOnly bool

func init() {
	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint an .env file for key naming and value issues",
		RunE:  runLint,
	}
	lintCmd.Flags().StringVarP(&lintFile, "file", "f", ".env", "path to .env file")
	lintCmd.Flags().BoolVar(&lintWarnOnly, "warn-only", false, "exit 0 even when errors are found")
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	entries, err := env.LoadFile(lintFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", lintFile, err)
	}

	m := env.ToMap(entries)
	result := env.LintMap(m)

	if len(result.Issues) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✔ no lint issues found")
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), result.Summary())

	if result.HasErrors() && !lintWarnOnly {
		os.Exit(1)
	}
	return nil
}
