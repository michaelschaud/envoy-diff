package reporter

import (
	"fmt"
	"io"
)

// Writer is the interface for all report output formats.
type Writer interface {
	Write(out io.Writer, r *Report) error
}

// NewWriter returns a Writer for the given format name.
// Supported formats: "text", "json", "csv", "markdown".
func NewWriter(format string) (Writer, error) {
	switch format {
	case "text", "":
		return &TextReportWriter{}, nil
	case "json":
		return &JSONReportWriter{}, nil
	case "csv":
		return NewCSVReportWriter(), nil
	case "markdown":
		return NewMarkdownReportWriter(), nil
	default:
		return nil, fmt.Errorf("unsupported report format: %q", format)
	}
}
