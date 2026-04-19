package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

func init() {
	trimCmd := &cobra.Command{
		Use:   "trim <file>",
		Short: "Remove leading/trailing whitespace from .env values",
		Args:  cobra.ExactArgs(1),
		RunE:  runTrim,
	}
	trimCmd.Flags().BoolP("dry-run", "n", false, "show what would be trimmed without writing")
	rootCmd.AddCommand(trimCmd)
}

func runTrim(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	m, err := env.LoadFile(filePath)
	if err != nil {
		return fmt.Errorf("loading file: %w", err)
	}

	out, result := env.Trim(m)

	if len(result.Trimmed) == 0 {
		fmt.Fprintln(os.Stdout, "no values needed trimming")
		return nil
	}

	for _, k := range result.Trimmed {
		fmt.Fprintf(os.Stdout, "trimmed: %s\n", k)
	}

	if dryRun {
		fmt.Fprintf(os.Stdout, "dry-run: %d value(s) would be trimmed\n", len(result.Trimmed))
		return nil
	}

	if err := env.WriteFile(filePath, out, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	fmt.Fprintf(os.Stdout, "%d value(s) trimmed and written to %s\n", len(result.Trimmed), filePath)
	return nil
}
