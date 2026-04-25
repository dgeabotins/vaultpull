package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
)

func init() {
	var keys []string
	var patterns []string
	var placeholder string
	var dryRun bool
	var output string

	cmd := &cobra.Command{
		Use:   "redact <file>",
		Short: "Replace sensitive values in an env file with a placeholder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRedact(args[0], keys, patterns, placeholder, output, dryRun)
		},
	}

	cmd.Flags().StringSliceVarP(&keys, "key", "k", nil, "explicit key names to redact")
	cmd.Flags().StringSliceVarP(&patterns, "pattern", "p", nil, "regex patterns to match key names")
	cmd.Flags().StringVar(&placeholder, "placeholder", "***", "replacement value for redacted keys")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print result without writing")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write result to this file (default: overwrite input)")

	rootCmd.AddCommand(cmd)
}

func runRedact(file string, keys, patterns []string, placeholder, output string, dryRun bool) error {
	values, err := env.LoadFile(file)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	res, err := env.Redact(values, env.RedactOptions{
		Keys:        keys,
		Patterns:    patterns,
		Placeholder: placeholder,
	})
	if err != nil {
		return fmt.Errorf("redact: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Summary: %s\n", res.Summary())

	if dryRun {
		for _, k := range sortedKeys(res.Values) {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, res.Values[k])
		}
		return nil
	}

	dest := output
	if dest == "" {
		dest = file
	}
	return env.WriteFile(dest, res.Values)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sorted := make([]string, len(keys))
	copy(sorted, keys)
	// simple insertion sort for small maps
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && strings.ToLower(sorted[j]) < strings.ToLower(sorted[j-1]); j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}
	return sorted
}
