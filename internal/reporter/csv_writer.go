package reporter

import (
	"encoding/csv"
	"io"
	"sort"

	"github.com/yourorg/envoy-diff/internal/diff"
)

type CSVReportWriter struct{}

func NewCSVReportWriter() *CSVReportWriter {
	return &CSVReportWriter{}
}

func (c *CSVReportWriter) Write(w io.Writer, r *Report) error {
	cw := csv.NewWriter(w)

	header := []string{"category", "key", "staging_value", "production_value"}
	if err := cw.Write(header); err != nil {
		return err
	}

	writeSection := func(category string, keys []string, vals map[string]string, secondary map[string]string) error {
		sorted := make([]string, len(keys))
		copy(sorted, keys)
		sort.Strings(sorted)
		for _, k := range sorted {
			staging, production := "", ""
			if category == "added" {
				production = vals[k]
			} else if category == "removed" {
				staging = vals[k]
			} else {
				staging = vals[k]
				production = secondary[k]
			}
			if err := cw.Write([]string{category, k, staging, production}); err != nil {
				return err
			}
		}
		return nil
	}

	res := r.Result

	addedKeys := keysOf(res.Added)
	if err := writeSection("added", addedKeys, res.Added, nil); err != nil {
		return err
	}

	removedKeys := keysOf(res.Removed)
	if err := writeSection("removed", removedKeys, res.Removed, nil); err != nil {
		return err
	}

	changedKeys := changedKeysOf(res.Changed)
	stagingChanged := map[string]string{}
	prodChanged := map[string]string{}
	for _, k := range changedKeys {
		stagingChanged[k] = res.Changed[k].From
		prodChanged[k] = res.Changed[k].To
	}
	if err := writeSection("changed", changedKeys, stagingChanged, prodChanged); err != nil {
		return err
	}

	cw.Flush()
	return cw.Error()
}

func keysOf(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func changedKeysOf(m map[string]diff.ValuePair) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
