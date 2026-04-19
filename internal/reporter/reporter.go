package reporter

import (
	"io"
	"time"

	"github.com/envoy-diff/internal/diff"
)

// Report holds a full diff report with metadata.
type Report struct {
	GeneratedAt time.Time        `json:"generated_at"`
	SourceFile  string           `json:"source_file"`
	TargetFile  string           `json:"target_file"`
	Result      diff.Result      `json:"result"`
	Stats       Stats            `json:"stats"`
}

// Stats summarises the counts from a diff result.
type Stats struct {
	Added   int `json:"added"`
	Removed int `json:"removed"`
	Changed int `json:"changed"`
	Same    int `json:"same"`
}

// Writer is implemented by anything that can render a Report.
type Writer interface {
	Write(w io.Writer, r Report) error
}

// New builds a Report from a diff result and file paths.
func New(sourceFile, targetFile string, result diff.Result) Report {
	stats := Stats{
		Added:   len(result.Added),
		Removed: len(result.Removed),
		Changed: len(result.Changed),
		Same:    len(result.Same),
	}
	return Report{
		GeneratedAt: time.Now().UTC(),
		SourceFile:  sourceFile,
		TargetFile:  targetFile,
		Result:      result,
		Stats:       stats,
	}
}

// HasDiff returns true when any differences exist.
func (r Report) HasDiff() bool {
	return r.Stats.Added > 0 || r.Stats.Removed > 0 || r.Stats.Changed > 0
}
