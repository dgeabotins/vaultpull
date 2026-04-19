package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"vaultpull/internal/env"
)

var copyKeys string

func init() {
	copyCmd := &cobra.Command{
		Use:   "copy <src> <dst>",
		Short: "Copy keys from one .env file into another",
		Args:  cobra.ExactArgs(2),
		RunE:  runCopy,
	}
	copyCmd.Flags().StringVarP(&copyKeys, "keys", "k", "", "comma-separated list of keys to copy (default: all)")
	rootCmd.AddCommand(copyCmd)
}

func runCopy(cmd *cobra.Command, args []string) error {
	src, dst := args[0], args[1]

	var keys []string
	if copyKeys != "" {
		for _, k := range strings.Split(copyKeys, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				keys = append(keys, k)
			}
		}
	}

	res, err := env.Copy(src, dst, keys)
	if err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), res.Summary())
	return nil
}
