package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/vaultpull/internal/template"
)

func init() {
	var templatePath string
	var outputPath string

	renderCmd := &cobra.Command{
		Use:   "render",
		Short: "Render a template file using secrets as context",
		Long:  "Render a Go template file substituting secret values and write the result to an output file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRender(templatePath, outputPath)
		},
	}

	renderCmd.Flags().StringVarP(&templatePath, "template", "t", "", "Path to the template file (required)")
	renderCmd.Flags().StringVarP(&outputPath, "output", "o", ".env", "Path to write rendered output")
	_ = renderCmd.MarkFlagRequired("template")

	rootCmd.AddCommand(renderCmd)
}

func runRender(templatePath, outputPath string) error {
	vars := make(map[string]string)
	for _, e := range os.Environ() {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				vars[e[:i]] = e[i+1:]
				break
			}
		}
	}

	r := template.New(outputPath)
	if err := r.Render(templatePath, vars); err != nil {
		return fmt.Errorf("render failed: %w", err)
	}

	fmt.Fprintf(os.Stdout, "rendered template %s -> %s\n", templatePath, outputPath)
	return nil
}
