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
	if profile == "" {
		return fmt.Errorf("profile name must not be empty")
	}

	info, err := whoami.Gather(profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error gathering info for profile %q: %v\n", profile, err)
		return err
	}
	fmt.Println(info.Format())
	return nil
}
