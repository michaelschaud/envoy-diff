package reporter

import (
	"encoding/json"
	"io"
)

// JSONReportWriter renders a Report as indented JSON.
type JSONReportWriter struct {
	Indent bool
}

// Write serialises the report to w as JSON.
func (jw JSONReportWriter) Write(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	if jw.Indent {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(r)
}
