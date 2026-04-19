package formatter

import (
	"fmt"
	"io"

	"github.com/example/envoy-diff/internal/diff"
)

// Writer is the interface implemented by all output formatters.
type Writer interface {
	Write(result diff.Result) error
}

// NewWriter returns a Writer for the given format string.
// Supported formats: "text", "json", "csv".
func NewWriter(format string, w io.Writer) (Writer, error) {
	switch format {
	case "text":
		return NewTextWriter(w), nil
	case "json":
		return NewJSONWriter(w), nil
	case "csv":
		return NewCSVWriter(w), nil
	default:
		return nil, fmt.Errorf("unsupported format %q: choose text, json, or csv", format)
	}
}
