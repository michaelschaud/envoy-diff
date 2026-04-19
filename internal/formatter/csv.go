package formatter

import (
	"encoding/csv"
	"io"
	"sort"

	"github.com/example/envoy-diff/internal/diff"
)

// CSVWriter writes diff results as CSV rows.
type CSVWriter struct {
	w io.Writer
}

// NewCSVWriter creates a new CSVWriter.
func NewCSVWriter(w io.Writer) *CSVWriter {
	return &CSVWriter{w: w}
}

// Write outputs the diff result in CSV format with columns:
// category, key, staging_value, production_value
func (c *CSVWriter) Write(result diff.Result) error {
	cw := csv.NewWriter(c.w)

	if err := cw.Write([]string{"category", "key", "staging_value", "production_value"}); err != nil {
		return err
	}

	write := func(category string, keys []string, staging, production map[string]string) error {
		sorted := make([]string, len(keys))
		copy(sorted, keys)
		sort.Strings(sorted)
		for _, k := range sorted {
			row := []string{category, k, staging[k], production[k]}
			if err := cw.Write(row); err != nil {
				return err
			}
		}
		return nil
	}

	if err := write("added", result.Added, map[string]string{}, result.Production); err != nil {
		return err
	}
	if err := write("removed", result.Removed, result.Staging, map[string]string{}); err != nil {
		return err
	}
	if err := write("changed", result.Changed, result.Staging, result.Production); err != nil {
		return err
	}
	if err := write("identical", result.Identical, result.Staging, result.Production); err != nil {
		return err
	}

	cw.Flush()
	return cw.Error()
}
