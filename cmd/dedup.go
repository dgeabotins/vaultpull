package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

var dedupWrite bool

func init() {
	dedupCmd := &cobra.Command{
		Use:   "dedup [file]",
		Short: "Remove duplicate keys from an env file",
		Args:  cobra.ExactArgs(1),
		RunE:  runDedup,
	}
	dedupCmd.Flags().BoolVarP(&dedupWrite, "write", "w", false, "write deduplicated output back to file")
	rootCmd.AddCommand(dedupCmd)
}

func runDedup(cmd *cobra.Command, args []string) error {
	path := args[0]
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	// trim trailing empty element from split
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	out, result := env.Dedup(lines)

	fmt.Fprintln(cmd.OutOrStdout(), result.Summary())

	if result.Removed == 0 {
		return nil
	}

	if dedupWrite {
		content := strings.Join(out, "\n") + "\n"
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "file updated")
	}
	return nil
}
