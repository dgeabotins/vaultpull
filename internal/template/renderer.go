package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// Renderer renders a template file using env vars as context.
type Renderer struct {
	outputPath string
}

// New creates a new Renderer that writes to outputPath.
func New(outputPath string) *Renderer {
	return &Renderer{outputPath: outputPath}
}

// Render executes the template at templatePath with the provided vars and
// writes the result to the configured output path.
func (r *Renderer) Render(templatePath string, vars map[string]string) error {
	tmplBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("template: read %s: %w", templatePath, err)
	}

	tmpl, err := template.New("env").Option("missingkey=error").Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("template: parse: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return fmt.Errorf("template: execute: %w", err)
	}

	if err := os.WriteFile(r.outputPath, buf.Bytes(), 0600); err != nil {
		return fmt.Errorf("template: write %s: %w", r.outputPath, err)
	}
	return nil
}
