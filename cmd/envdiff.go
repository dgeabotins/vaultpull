package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

func init() {
	envDiffCmd := &cobra.Command{
		Use:   "envdiff <old-file> <new-file>",
		Short: "Compare two .env files and show differences",
		Args:  cobra.ExactArgs(2),
		RunE:  runEnvDiff,
	}
	rootCmd.AddCommand(envDiffCmd)
}

func runEnvDiff(cmd *cobra.Command, args []string) error {
	old, err := env.LoadFile(args[0])
	if err != nil {
		return fmt.Errorf("loading %s: %w", args[0], err)
	}
	next, err := env.LoadFile(args[1])
	if err != nil {
		return fmt.Errorf("loading %s: %w", args[1], err)
	}

	r := env.Diff(env.ToMap(old), env.ToMap(next))

	if !r.HasChanges() {
		fmt.Fprintln(cmd.OutOrStdout(), "No differences found.")
		return nil
	}

	w := cmd.OutOrStdout()
	for k, v := range r.Added {
		fmt.Fprintf(w, "+ %s=%s\n", k, v)
	}
	for k, v := range r.Removed {
		fmt.Fprintf(w, "- %s=%s\n", k, v)
	}
	for k, v := range r.Changed {
		fmt.Fprintf(w, "~ %s: %s -> %s\n", k, v[0], v[1])
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, r.Summary())

	os.Exit(1)
	return nil
}
