package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/vaultpull/vaultpull/internal/env"
)

const maxValueDisplayLen = 40

var inspectFile string

func init() {
	inspectCmd := &cobra.Command{
		Use:   "inspect",
		Short: "Display parsed key-value pairs from a .env file",
		RunE:  runInspect,
	}
	inspectCmd.Flags().StringVarP(&inspectFile, "file", "f", ".env", "path to .env file")
	rootCmd.AddCommand(inspectCmd)
}

func runInspect(cmd *cobra.Command, _ []string) error {
	entries, err := env.LoadFile(inspectFile)
	if err != nil {
		return fmt.Errorf("inspect: %w", err)
	}
	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no entries found")
		return nil
	}
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tVALUE")
	fmt.Fprintln(w, "---\t-----")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\n", e.Key, truncateValue(e.Value))
	}
	return w.Flush()
}

// truncateValue shortens a value to maxValueDisplayLen characters, appending
// "..." if the value was truncated.
func truncateValue(val string) string {
	if len(val) > maxValueDisplayLen {
		return val[:maxValueDisplayLen] + "..."
	}
	return val
}

// InspectEnvFile is exported for testing.
func InspectEnvFile(path string) ([]env.Entry, error) {
	return env.LoadFile(path)
}

var _ = os.Stderr // keep os import
