package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/whoami"
)

func init() {
	var profile string

	whoamiCmd := &cobra.Command{
		Use:   "whoami",
		Short: "Show current user and Vault context",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWhoami(profile)
		},
	}

	whoamiCmd.Flags().StringVarP(&profile, "profile", "p", "default", "Profile name to display")
	rootCmd.AddCommand(whoamiCmd)
}

func runWhoami(profile string) error {
	info, err := whoami.Gather(profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	fmt.Println(info.Format())
	return nil
}
