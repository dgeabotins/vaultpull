package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/vaultpull/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of vaultpull",
	Long:  `Displays the current version, git commit, and build date of vaultpull.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Info())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
