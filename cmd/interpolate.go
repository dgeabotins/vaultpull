package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/vaultpull/internal/env"
)

func init() {
	var inputFile string
	var outputFile string
	var strict bool

	cmd := &cobra.Command{
		Use:   "interpolate",
		Short: "Expand ${KEY} references within a .env file",
		Long:  "Reads a .env file and resolves ${KEY} references using values defined in the same file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInterpolate(inputFile, outputFile, strict)
		},
	}

	cmd.Flags().StringVarP(&inputFile, "file", "f", ".env", "Input .env file")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (defaults to input file)")
	cmd.Flags().BoolVar(&strict, "strict", false, "Exit with error if any references are unresolved")

	rootCmd.AddCommand(cmd)
}

func runInterpolate(inputFile, outputFile string, strict bool) error {
	envMap, err := env.LoadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", inputFile, err)
	}

	result := env.Interpolate(envMap)

	if len(result.Unresolved) > 0 {
		fmt.Fprintf(os.Stderr, "warning: unresolved references: %v\n", result.Unresolved)
		if strict {
			return fmt.Errorf("strict mode: %d unresolved reference(s)", len(result.Unresolved))
		}
	}

	dest := inputFile
	if outputFile != "" {
		dest = outputFile
	}

	if err := env.WriteFile(dest, result.Resolved); err != nil {
		return fmt.Errorf("writing %s: %w", dest, err)
	}

	fmt.Println(result.Summary())
	return nil
}
