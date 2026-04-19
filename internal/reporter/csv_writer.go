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

	writeSection := func(category string, keys []string, vals map[string]string, secondary map[string]string) {
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
			_ = cw.Write([]string{category, k, staging, production})
		}
	}

	res := r.Result

	addedKeys := keysOf(res.Added)
	writeSection("added", addedKeys, res.Added, nil)

	removedKeys := keysOf(res.Removed)
	writeSection("removed", removedKeys, res.Removed, nil)

	changedKeys := changedKeysOf(res.Changed)
	stagingChanged := map[string]string{}
	prodChanged := map[string]string{}
	for _, k := range changedKeys {
		stagingChanged[k] = res.Changed[k].From
		prodChanged[k] = res.Changed[k].To
	}
	writeSection("changed", changedKeys, stagingChanged, prodChanged)

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
