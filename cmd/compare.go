package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/vaultpull/internal/env"
)

func init() {
	var showSame bool
	cmd := &cobra.Command{
		Use:   "compare <file-a> <file-b>",
		Short: "Compare two .env files and show differences",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompare(args[0], args[1], showSame)
		},
	}
	cmd.Flags().BoolVar(&showSame, "show-same", false, "also print unchanged keys")
	rootCmd.AddCommand(cmd)
}

func runCompare(fileA, fileB string, showSame bool) error {
	a, err := env.LoadFile(fileA)
	if err != nil {
		return fmt.Errorf("loading %s: %w", fileA, err)
	}
	b, err := env.LoadFile(fileB)
	if err != nil {
		return fmt.Errorf("loading %s: %w", fileB, err)
	}
	res := env.Compare(env.ToMap(a), env.ToMap(b))
	for k, v := range res.OnlyInA {
		fmt.Fprintf(os.Stdout, "- %s=%s\n", k, v)
	}
	for k, v := range res.OnlyInB {
		fmt.Fprintf(os.Stdout, "+ %s=%s\n", k, v)
	}
	for k, pair := range res.Changed {
		fmt.Fprintf(os.Stdout, "~ %s: %s -> %s\n", k, pair[0], pair[1])
	}
	if showSame {
		for k, v := range res.Same {
			fmt.Fprintf(os.Stdout, "  %s=%s\n", k, v)
		}
	}
	fmt.Fprintln(os.Stdout, res.Summary())
	return nil
}
