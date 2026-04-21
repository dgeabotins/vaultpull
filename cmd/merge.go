package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

var (
	mergeStrategy string
	mergePrefix   string
	mergeOutput   string
)

func init() {
	mergeCmd := &cobra.Command{
		Use:   "merge <src> <dst>",
		Short: "Merge two .env files using a chosen conflict strategy",
		Args:  cobra.ExactArgs(2),
		RunE:  runMerge,
	}

	mergeCmd.Flags().StringVarP(&mergeStrategy, "strategy", "s", "overwrite",
		"Conflict strategy: overwrite | keep | error")
	mergeCmd.Flags().StringVar(&mergePrefix, "prefix", "",
		"Only merge keys with this prefix")
	mergeCmd.Flags().StringVarP(&mergeOutput, "output", "o", "",
		"Write result to file instead of stdout")

	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
	srcPath, dstPath := args[0], args[1]

	srcMap, err := env.LoadFile(srcPath)
	if err != nil {
		return fmt.Errorf("loading src: %w", err)
	}
	dstMap, err := env.LoadFile(dstPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("loading dst: %w", err)
	}
	if dstMap == nil {
		dstMap = make(map[string]string)
	}

	var strat env.MergeStrategy
	switch mergeStrategy {
	case "overwrite":
		strat = env.StrategyOverwrite
	case "keep":
		strat = env.StrategyKeepExisting
	case "error":
		strat = env.StrategyError
	default:
		return fmt.Errorf("unknown strategy %q: use overwrite, keep, or error", mergeStrategy)
	}

	res, err := env.MergeMap(dstMap, srcMap, env.MergeOptions{Strategy: strat, Prefix: mergePrefix})
	if err != nil {
		return err
	}

	out := mergeOutput
	if out == "" {
		out = dstPath
	}
	if err := env.WriteFile(out, dstMap); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "merged: +%d updated:%d skipped:%d → %s\n",
		res.Added, res.Updated, res.Skipped, out)
	return nil
}
