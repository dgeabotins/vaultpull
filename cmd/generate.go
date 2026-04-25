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
		outFile   string
		length    int
		charset   string
		overwrite bool
		dryRun    bool
	)

	cmd := &cobra.Command{
		Use:   "generate [flags] KEY [KEY...]",
		Short: "Generate random secret values for env keys",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(args, outFile, length, charset, overwrite, dryRun)
		},
	}

	cmd.Flags().StringVarP(&outFile, "file", "f", ".env", "env file to read/write")
	cmd.Flags().IntVarP(&length, "length", "l", 32, "length of generated value")
	cmd.Flags().StringVarP(&charset, "charset", "c", "base64", "charset: alpha|alphanum|hex|base64|symbol")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview without writing")

	rootCmd.AddCommand(cmd)
}

func runGenerate(keys []string, outFile string, length int, charset string, overwrite, dryRun bool) error {
	existing := map[string]string{}
	if data, err := env.LoadFile(outFile); err == nil {
		existing = data
	}

	opts := env.GenerateOptions{
		Keys:      keys,
		Length:    length,
		Charset:   charset,
		Overwrite: overwrite,
		DryRun:    dryRun,
	}

	out, result, err := env.Generate(existing, opts)
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	for _, k := range result.Generated {
		action := "generated"
		if dryRun {
			action = "would generate"
		}
		fmt.Printf("  %s  %s\n", action, k)
	}
	for _, k := range result.Skipped {
		fmt.Printf("  skipped   %s (already set)\n", k)
	}

	if dryRun {
		fmt.Println("\n(dry run — no changes written)")
		return nil
	}

	if err := env.WriteFile(outFile, out, 0600); err != nil {
		return fmt.Errorf("write %s: %w", outFile, err)
	}

	fmt.Printf("\n%s — written to %s\n", result.Summary(), outFile)
	_ = strings.Join // suppress unused import if needed
	_ = os.Stderr
	return nil
}
