package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/internal/env"
)

var (
	sortReverse         bool
	sortCaseInsensitive bool
	sortPrefixFirst     []string
)

func init() {
	sortCmd := &cobra.Command{
		Use:   "sort <file>",
		Short: "Print env file keys in sorted order",
		Args:  cobra.ExactArgs(1),
		RunE:  runSort,
	}

	sortCmd.Flags().BoolVarP(&sortReverse, "reverse", "r", false, "Sort in descending order")
	sortCmd.Flags().BoolVarP(&sortCaseInsensitive, "ignore-case", "i", false, "Case-insensitive sort")
	sortCmd.Flags().StringArrayVarP(&sortPrefixFirst, "prefix-first", "p", nil, "Prefixes to sort before others")

	rootCmd.AddCommand(sortCmd)
}

func runSort(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	m, err := env.LoadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}

	opts := env.SortOptions{
		Reverse:         sortReverse,
		CaseInsensitive: sortCaseInsensitive,
		PrefixFirst:     sortPrefixFirst,
	}

	sorted := env.SortMap(m, opts)
	for _, k := range sorted {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, m[k])
	}
	return nil
}
