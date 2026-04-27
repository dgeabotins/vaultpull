package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

func init() {
	typecheckCmd := &cobra.Command{
		Use:   "typecheck [file]",
		Short: "Validate env var values against expected types",
		Args:  cobra.ExactArgs(1),
		RunE:  runTypecheck,
	}
	typecheckCmd.Flags().StringSliceP("hint", "H", nil,
		`type hints in KEY=TYPE format (types: string, int, float, bool, url, nonempty)`)
	typecheckCmd.Flags().BoolP("strict", "s", false, "exit non-zero if any issues are found")
	rootCmd.AddCommand(typecheckCmd)
}

func runTypecheck(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	hintStrs, _ := cmd.Flags().GetStringSlice("hint")
	strict, _ := cmd.Flags().GetBool("strict")

	m, err := env.LoadFile(filePath)
	if err != nil {
		return fmt.Errorf("loading %s: %w", filePath, err)
	}

	hints, err := parseTypeHints(hintStrs)
	if err != nil {
		return err
	}

	if len(hints) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no hints provided — nothing to check")
		return nil
	}

	result := env.TypeCheck(m, hints)
	if !result.HasIssues() {
		fmt.Fprintln(cmd.OutOrStdout(), "✓ all values pass type checks")
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), "type check issues found:")
	fmt.Fprintln(cmd.OutOrStdout(), result.Summary())

	if strict {
		os.Exit(1)
	}
	return nil
}

func parseTypeHints(raw []string) (map[string]env.TypeHint, error) {
	hints := make(map[string]env.TypeHint, len(raw))
	valid := map[string]env.TypeHint{
		"string":   env.TypeString,
		"int":      env.TypeInt,
		"float":    env.TypeFloat,
		"bool":     env.TypeBool,
		"url":      env.TypeURL,
		"nonempty": env.TypeNonempty,
	}
	for _, h := range raw {
		parts := strings.SplitN(h, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid hint %q: expected KEY=TYPE", h)
		}
		typeStr := strings.ToLower(strings.TrimSpace(parts[1]))
		t, ok := valid[typeStr]
		if !ok {
			return nil, fmt.Errorf("unknown type %q for key %s", typeStr, parts[0])
		}
		hints[strings.TrimSpace(parts[0])] = t
	}
	return hints, nil
}
