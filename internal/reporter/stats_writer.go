package reporter

import (
	"encoding/json"
	"fmt"
	"io"
)

// StatsWriter writes only the Stats section of a Report.
type StatsWriter struct {
	w      io.Writer
	format string
}

// NewStatsWriter creates a StatsWriter for the given format ("text" or "json").
func NewStatsWriter(w io.Writer, format string) *StatsWriter {
	return &StatsWriter{w: w, format: format}
}

// Write outputs the stats derived from r in the configured format.
func (sw *StatsWriter) Write(r Report) error {
	s := StatsFromReport(r)
	switch sw.format {
	case "json":
		return sw.writeJSON(r, s)
	default:
		return sw.writeText(r, s)
	}
}

func (sw *StatsWriter) writeText(r Report, s Stats) error {
	_, err := fmt.Fprintf(sw.w,
		"Source : %s\nTarget : %s\n%s",
		r.Metadata.SourceFile,
		r.Metadata.TargetFile,
		s.String(),
	)
	return err
}

func (sw *StatsWriter) writeJSON(r Report, s Stats) error {
	payload := struct {
		Source  string `json:"source"`
		Target  string `json:"target"`
		Added   int    `json:"added"`
		Removed int    `json:"removed"`
		Changed int    `json:"changed"`
		Total   int    `json:"total"`
		HasDiff bool   `json:"has_diff"`
	}{
		Source:  r.Metadata.SourceFile,
		Target:  r.Metadata.TargetFile,
		Added:   s.Added,
		Removed: s.Removed,
		Changed: s.Changed,
		Total:   s.Total,
		HasDiff: s.HasDiff(),
	}
	enc := json.NewEncoder(sw.w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
