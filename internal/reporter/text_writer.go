package reporter

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// TextReportWriter renders a Report as human-readable text.
type TextReportWriter struct{}

// Write outputs the report to w in plain-text format.
func (tw TextReportWriter) Write(w io.Writer, r Report) error {
	tab := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tab, "Source : %s\n", r.SourceFile)
	fmt.Fprintf(tab, "Target : %s\n", r.TargetFile)
	fmt.Fprintf(tab, "Generated: %s\n\n", r.GeneratedAt.Format("2006-01-02 15:04:05 UTC"))

	fmt.Fprintf(tab, "Stats: +%d added  -%d removed  ~%d changed  =%d same\n\n",
		r.Stats.Added, r.Stats.Removed, r.Stats.Changed, r.Stats.Same)

	if len(r.Result.Added) > 0 {
		fmt.Fprintln(tab, "[ADDED]")
		for _, k := range sortedKeys(r.Result.Added) {
			fmt.Fprintf(tab, "  + %s\t= %s\n", k, r.Result.Added[k])
		}
		fmt.Fprintln(tab)
	}

	if len(r.Result.Removed) > 0 {
		fmt.Fprintln(tab, "[REMOVED]")
		for _, k := range sortedKeys(r.Result.Removed) {
			fmt.Fprintf(tab, "  - %s\t= %s\n", k, r.Result.Removed[k])
		}
		fmt.Fprintln(tab)
	}

	if len(r.Result.Changed) > 0 {
		fmt.Fprintln(tab, "[CHANGED]")
		for _, k := range sortedChangedKeys(r.Result.Changed) {
			vals := r.Result.Changed[k]
			fmt.Fprintf(tab, "  ~ %s\t%s -> %s\n", k, vals[0], vals[1])
		}
		fmt.Fprintln(tab)
	}

	return tab
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedChangedKeys(m map[string][2]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
