package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/vaultpull/internal/rotate"
)

var (
	backupDir  string
	maxBackups int
)

var backupCmd = &cobra.Command{
	Use:   "backup <envfile>",
	Short: "Backup an existing .env file to the backup directory",
	Args:  cobra.ExactArgs(1),
	RunE:  runBackup,
}

func init() {
	backupCmd.Flags().StringVar(&backupDir, "backup-dir", ".vaultpull/backups", "directory to store backups")
	backupCmd.Flags().IntVar(&maxBackups, "max-backups", 5, "maximum number of backups to retain per file")
	rootCmd.AddCommand(backupCmd)
}

func runBackup(cmd *cobra.Command, args []string) error {
	src := args[0]
	if !filepath.IsAbs(src) {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}
		src = filepath.Join(cwd, src)
	}

	r := rotate.New(backupDir, maxBackups)
	dest, err := r.Backup(src)
	if err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}
	if dest == "" {
		fmt.Fprintf(cmd.OutOrStdout(), "nothing to back up: %s does not exist\n", src)
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "backed up to %s\n", dest)
	return nil
}
