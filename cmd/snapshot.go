package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/env"
)

var snapshotOutput string

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot <file>",
		Short: "Capture a point-in-time snapshot of an env file",
		Args:  cobra.ExactArgs(1),
		RunE:  runSnapshot,
	}

	snapshotCmd.Flags().StringVarP(&snapshotOutput, "output", "o", "",
		"Destination path for the snapshot JSON (default: <file>.snap-<timestamp>.json)")

	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshot(cmd *cobra.Command, args []string) error {
	src := args[0]

	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", src)
	}

	dest := snapshotOutput
	if dest == "" {
		base := filepath.Base(src)
		timestamp := time.Now().UTC().Format("20060102T150405Z")
		dest = fmt.Sprintf("%s.snap-%s.json", base, timestamp)
	}

	res, err := env.TakeSnapshot(src, dest)
	if err != nil {
		return fmt.Errorf("snapshot failed: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Snapshot saved: %s\n", res.Path)
	fmt.Fprintf(cmd.OutOrStdout(), "Captured %d entries from %s\n",
		len(res.Snapshot.Entries), res.Snapshot.File)

	return nil
}
