package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

func init() {
	var srcFile, dstFile string
	var keys []string
	var force, dryRun bool

	cmd := &cobra.Command{
		Use:   "promote",
		Short: "Promote env vars from one file to another",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPromote(srcFile, dstFile, keys, force, dryRun)
		},
	}

	cmd.Flags().StringVar(&srcFile, "src", "", "source .env file (required)")
	cmd.Flags().StringVar(&dstFile, "dst", "", "destination .env file (required)")
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "specific keys to promote (default: all)")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing keys in destination")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be promoted without writing")
	_ = cmd.MarkFlagRequired("src")
	_ = cmd.MarkFlagRequired("dst")

	rootCmd.AddCommand(cmd)
}

func runPromote(srcFile, dstFile string, keys []string, force, dryRun bool) error {
	srcMap, err := env.LoadFile(srcFile)
	if err != nil {
		return fmt.Errorf("loading src: %w", err)
	}

	dstMap, err := env.LoadFile(dstFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("loading dst: %w", err)
	}
	if dstMap == nil {
		dstMap = map[string]string{}
	}

	opts := env.PromoteOptions{Keys: keys, Force: force, DryRun: dryRun}
	merged, result := env.Promote(srcMap, dstMap, opts)

	fmt.Println(result.Summary())
	for _, k := range result.Promoted {
		fmt.Printf("  + %s\n", k)
	}
	for _, k := range result.Skipped {
		fmt.Printf("  ~ %s (skipped)\n", k)
	}

	if dryRun {
		fmt.Println("[dry-run] no changes written")
		return nil
	}

	return env.WriteFile(dstFile, merged)
}
