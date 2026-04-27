package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/vaultpull/internal/env"
)

func init() {
	classifyCmd := &cobra.Command{
		Use:   "classify [file]",
		Short: "Classify env var values by inferred type",
		Args:  cobra.ExactArgs(1),
		RunE:  runClassify,
	}
	classifyCmd.Flags().StringP("filter", "f", "", "Only show entries matching this category (url, secret, boolean, integer, float, path, json, empty, unknown)")
	classifyCmd.Flags().BoolP("summary", "s", false, "Print category counts only")
	rootCmd.AddCommand(classifyCmd)
}

func runClassify(cmd *cobra.Command, args []string) error {
	filter, _ := cmd.Flags().GetString("filter")
	summary, _ := cmd.Flags().GetBool("summary")

	m, err := env.LoadFile(args[0])
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	results := env.Classify(m)

	if summary {
		counts := make(map[env.Category]int)
		for _, r := range results {
			counts[r.Category]++
		}
		for _, cat := range []env.Category{
			env.CategoryURL, env.CategorySecret, env.CategoryBoolean,
			env.CategoryInteger, env.CategoryFloat, env.CategoryPath,
			env.CategoryJSON, env.CategoryEmpty, env.CategoryUnknown,
		} {
			if n := counts[cat]; n > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "%-12s %d\n", cat, n)
			}
		}
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tCATEGORY\tVALUE")
	for _, r := range results {
		if filter != "" && string(r.Category) != filter {
			continue
		}
		display := r.Value
		if r.Category == env.CategorySecret {
			display = "***"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", r.Key, r.Category, display)
	}
	w.Flush()

	_ = os.Stderr // satisfy import
	return nil
}
