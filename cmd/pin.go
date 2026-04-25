package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

func init() {
	var (
		file   string
		keys   []string
		prefix string
		dryRun bool
	)

	cmd := &cobra.Command{
		Use:   "pin",
		Short: "Mark env keys as pinned to prevent accidental overwrite",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPin(file, keys, prefix, dryRun)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", ".env", "env file to pin")
	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "specific keys to pin (comma-separated)")
	cmd.Flags().StringVar(&prefix, "prefix", "", "pin all keys with this prefix")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be pinned without writing")

	rootCmd.AddCommand(cmd)
}

func runPin(file string, keys []string, prefix string, dryRun bool) error {
	current, err := env.LoadFile(file)
	if err != nil {
		return fmt.Errorf("pin: failed to load %s: %w", file, err)
	}

	opts := env.PinOptions{
		Keys:   keys,
		Prefix: prefix,
		DryRun: dryRun,
	}

	updated, result, err := env.Pin(current, opts)
	if err != nil {
		return fmt.Errorf("pin: %w", err)
	}

	if len(result.Pinned) == 0 {
		fmt.Fprintln(os.Stdout, "No keys pinned.")
		return nil
	}

	fmt.Fprintf(os.Stdout, "Pinned: %s\n", strings.Join(result.Pinned, ", "))
	if len(result.Skipped) > 0 {
		fmt.Fprintf(os.Stdout, "Skipped (already pinned): %s\n", strings.Join(result.Skipped, ", "))
	}

	if dryRun {
		fmt.Fprintln(os.Stdout, "(dry-run: no changes written)")
		return nil
	}

	if err := env.WriteFile(file, updated, 0o600); err != nil {
		return fmt.Errorf("pin: failed to write %s: %w", file, err)
	}
	return nil
}
