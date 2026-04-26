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
		file     string
		incoming string
		keys     []string
		suffix   string
		dryRun   bool
	)

	cmd := &cobra.Command{
		Use:   "rotate",
		Short: "Rotate secrets: archive current values and apply new ones",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRotate(file, incoming, keys, suffix, dryRun)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", ".env", "target .env file")
	cmd.Flags().StringVarP(&incoming, "incoming", "i", "", "source .env file with new values (required)")
	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "keys to rotate (default: all keys in incoming)")
	cmd.Flags().StringVar(&suffix, "suffix", "_PREVIOUS", "suffix appended to archived key names")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print changes without writing")
	_ = cmd.MarkFlagRequired("incoming")

	rootCmd.AddCommand(cmd)
}

func runRotate(file, incomingFile string, keys []string, suffix string, dryRun bool) error {
	current, err := env.LoadFile(file)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading %s: %w", file, err)
	}
	if current == nil {
		current = map[string]string{}
	}

	incomingMap, err := env.LoadFile(incomingFile)
	if err != nil {
		return fmt.Errorf("reading incoming file %s: %w", incomingFile, err)
	}

	result, err := env.Rotate(current, incomingMap, env.RotateOptions{
		Keys:   keys,
		Suffix: strings.ToUpper(suffix),
		DryRun: dryRun,
	})
	if err != nil {
		return err
	}

	fmt.Println(result.Summary())
	for _, k := range result.Rotated {
		fmt.Printf("  ~ %s\n", k)
	}

	if dryRun {
		return nil
	}

	return env.WriteFile(file, current, 0o600)
}
