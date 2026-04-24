package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

var (
	sanitizeTrim      bool
	sanitizeControl   bool
	sanitizeNewlines  bool
	sanitizeWrite     bool
	sanitizeDryRun    bool
)

func init() {
	sanitizeCmd := &cobra.Command{
		Use:   "sanitize <file>",
		Short: "Clean env values by stripping whitespace or control characters",
		Args:  cobra.ExactArgs(1),
		RunE:  runSanitize,
	}

	sanitizeCmd.Flags().BoolVar(&sanitizeTrim, "trim", true, "trim leading/trailing whitespace from values")
	sanitizeCmd.Flags().BoolVar(&sanitizeControl, "strip-control", false, "strip ASCII control characters from values")
	sanitizeCmd.Flags().BoolVar(&sanitizeNewlines, "normalize-newlines", false, "normalize CRLF and CR to LF")
	sanitizeCmd.Flags().BoolVar(&sanitizeWrite, "write", false, "write sanitized output back to the file")
	sanitizeCmd.Flags().BoolVar(&sanitizeDryRun, "dry-run", false, "show what would change without writing")

	rootCmd.AddCommand(sanitizeCmd)
}

func runSanitize(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	data, err := env.LoadFile(filePath)
	if err != nil {
		return fmt.Errorf("loading %s: %w", filePath, err)
	}

	opts := env.SanitizeOptions{
		TrimWhitespace:    sanitizeTrim,
		StripControlChars: sanitizeControl,
		NormalizeNewlines: sanitizeNewlines,
	}

	result := env.Sanitize(data, opts)

	if len(result.Changes) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no changes needed")
		return nil
	}

	for _, c := range result.Changes {
		fmt.Fprintf(cmd.OutOrStdout(), "  ~ %s: %q -> %q\n", c.Key, c.Before, c.After)
	}

	if sanitizeDryRun {
		fmt.Fprintln(cmd.OutOrStdout(), "(dry-run: no file written)")
		return nil
	}

	if sanitizeWrite {
		if err := env.WriteFile(filePath, result.Sanitized, os.FileMode(0600)); err != nil {
			return fmt.Errorf("writing %s: %w", filePath, err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "wrote sanitized file: %s\n", filePath)
	}

	return nil
}
