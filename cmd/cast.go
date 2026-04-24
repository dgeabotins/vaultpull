package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

func init() {
	var file string
	var hints []string
	var strictBool bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "cast",
		Short: "Coerce env values to typed representations",
		Long:  "Cast reads an .env file and normalises values according to type hints (int, float, bool, string).",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCast(file, hints, strictBool, dryRun)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", ".env", "path to .env file")
	cmd.Flags().StringArrayVar(&hints, "hint", nil, "type hint in KEY=TYPE format (e.g. PORT=int)")
	cmd.Flags().BoolVar(&strictBool, "strict-bool", true, "normalise yes/no/on/off/1/0 to true/false")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print changes without writing")

	rootCmd.AddCommand(cmd)
}

func runCast(file string, hints []string, strictBool, dryRun bool) error {
	envMap, err := env.LoadFile(file)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	typeHints := make(map[string]string, len(hints))
	for _, h := range hints {
		parts := strings.SplitN(h, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid hint %q: expected KEY=TYPE", h)
		}
		typeHints[parts[0]] = parts[1]
	}

	result := env.Cast(envMap, env.CastOptions{
		TypeHints:  typeHints,
		StrictBool: strictBool,
	})

	for _, e := range result.Errors {
		fmt.Fprintf(os.Stderr, "warn: %s\n", e)
	}

	if result.Changed == 0 {
		fmt.Println("no values changed")
		return nil
	}

	fmt.Printf("%d value(s) changed\n", result.Changed)

	if dryRun {
		for k, v := range result.Casted {
			if envMap[k] != v {
				fmt.Printf("  %s: %q -> %q\n", k, envMap[k], v)
			}
		}
		return nil
	}

	return env.WriteFile(file, result.Casted)
}
