package reporter

import (
	"fmt"
	"io"
)

// Stats holds summary counts for a diff report.
type Stats struct {
	Added   int
	Removed int
	Changed int
	Total   int
}

// HasDiff returns true if any differences were found.
func (s Stats) HasDiff() bool {
	return s.Added > 0 || s.Removed > 0 || s.Changed > 0
}

// String returns a human-readable one-line summary.
func (s Stats) String() string {
	return fmt.Sprintf(
		"total=%d added=%d removed=%d changed=%d",
		s.Total, s.Added, s.Removed, s.Changed,
	)
}

// WriteStats writes a formatted stats block to w.
func WriteStats(w io.Writer, s Stats) error {
	_, err := fmt.Fprintf(w,
		"--- Stats ---\nTotal keys : %d\nAdded      : %d\nRemoved    : %d\nChanged    : %d\n",
		s.Total, s.Added, s.Removed, s.Changed,
	)
	return err
}

// StatsFromReport derives Stats from a Report.
func StatsFromReport(r Report) Stats {
	return Stats{
		Added:   len(r.Result.Added),
		Removed: len(r.Result.Removed),
		Changed: len(r.Result.Changed),
		Total:   r.Stats.Total,
	}
}
