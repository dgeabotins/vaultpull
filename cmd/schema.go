package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/env"
)

var (
	schemaFile  string
	schemaEnv   string
	schemaApply bool
)

func init() {
	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Validate a .env file against a schema and optionally apply defaults",
		RunE:  runSchema,
	}
	schemaCmd.Flags().StringVarP(&schemaFile, "schema", "s", ".env.schema", "path to schema file")
	schemaCmd.Flags().StringVarP(&schemaEnv, "file", "f", ".env", "path to .env file")
	schemaCmd.Flags().BoolVar(&schemaApply, "apply", false, "write defaults back to the .env file")
	rootCmd.AddCommand(schemaCmd)
}

func runSchema(cmd *cobra.Command, args []string) error {
	schema, err := env.LoadSchema(schemaFile)
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}

	data, err := env.LoadFile(schemaEnv)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load env file: %w", err)
	}
	if data == nil {
		data = map[string]string{}
	}

	out, result := env.ApplySchema(data, schema)

	fmt.Fprintln(cmd.OutOrStdout(), result.Summary())

	if result.HasErrors() {
		return fmt.Errorf("schema validation failed: missing required keys: %v", result.Missing)
	}

	if schemaApply && len(result.Defaults) > 0 {
		if err := env.WriteFile(schemaEnv, out, 0600); err != nil {
			return fmt.Errorf("write env file: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "applied %d default(s) to %s\n", len(result.Defaults), schemaEnv)
	}

	return nil
}
