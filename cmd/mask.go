package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/mask"
)

func init() {
	var file string

	cmd := &cobra.Command{
		Use:   "mask",
		Short: "Display .env file with secret values masked",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMask(file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", ".env", "Path to .env file")
	rootCmd.AddCommand(cmd)
}

func runMask(file string) error {
	entries, err := env.LoadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}

	m := mask.New()
	secrets := make(map[string]string, len(entries))
	for _, e := range entries {
		secrets[e.Key] = e.Value
	}

	masked := m.MaskMap(secrets)
	for _, e := range entries {
		fmt.Printf("%s=%s\n", e.Key, masked[e.Key])
	}
	return nil
}
