package export

import (
	"fmt"
	"os"
	"strings"
)

// Format represents an export output format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatExport Format = "export"
)

// Exporter writes secrets to stdout or a file in a given format.
type Exporter struct {
	format Format
}

// New returns a new Exporter for the given format string.
func New(format string) (*Exporter, error) {
	f := Format(strings.ToLower(format))
	switch f {
	case FormatDotenv, FormatJSON, FormatExport:
		return &Exporter{format: f}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %q", format)
	}
}

// Write renders secrets in the configured format to the given file (or stdout if path is "-").
func (e *Exporter) Write(secrets map[string]string, path string) error {
	var out string
	switch e.format {
	case FormatDotenv:
		out = renderDotenv(secrets)
	case FormatJSON:
		out = renderJSON(secrets)
	case FormatExport:
		out = renderExport(secrets)
	}

	if path == "-" {
		_, err := fmt.Print(out)
		return err
	}
	return os.WriteFile(path, []byte(out), 0600)
}

func renderDotenv(secrets map[string]string) string {
	var sb strings.Builder
	for k, v := range secrets {
		fmt.Fprintf(&sb, "%s=%q\n", k, v)
	}
	return sb.String()
}

func renderJSON(secrets map[string]string) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	i := 0
	for k, v := range secrets {
		comma := ","
		if i == len(secrets)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", k, v, comma)
		i++
	}
	sb.WriteString("}\n")
	return sb.String()
}

func renderExport(secrets map[string]string) string {
	var sb strings.Builder
	for k, v := range secrets {
		fmt.Fprintf(&sb, "export %s=%q\n", k, v)
	}
	return sb.String()
}
