package reporter

import (
	"fmt"
	"io"
	"sort"
	"time"
)

type MarkdownReportWriter struct{}

func NewMarkdownReportWriter() *MarkdownReportWriter {
	return &MarkdownReportWriter{}
}

func (w *MarkdownReportWriter) Write(out io.Writer, r *Report) error {
	fmt.Fprintf(out, "# Envoy Diff Report\n\n")
	fmt.Fprintf(out, "| Field | Value |\n")
	fmt.Fprintf(out, "|-------|-------|\n")
	fmt.Fprintf(out, "| Staging | `%s` |\n", r.Metadata.StagingFile)
	fmt.Fprintf(out, "| Production | `%s` |\n", r.Metadata.ProductionFile)
	fmt.Fprintf(out, "| Generated | %s |\n\n", r.Metadata.GeneratedAt.Format(time.RFC3339))

	fmt.Fprintf(out, "## Stats\n\n")
	fmt.Fprintf(out, "- Added: %d\n", r.Stats.Added)
	fmt.Fprintf(out, "- Removed: %d\n", r.Stats.Removed)
	fmt.Fprintf(out, "- Changed: %d\n", r.Stats.Changed)
	fmt.Fprintf(out, "- Identical: %d\n\n", r.Stats.Identical)

	if len(r.Result.Added) > 0 {
		fmt.Fprintf(out, "## Added\n\n| Key | Value |\n|-----|-------|\n")
		for _, k := range sortedKeys(r.Result.Added) {
			fmt.Fprintf(out, "| `%s` | `%s` |\n", k, r.Result.Added[k])
		}
		fmt.Fprintln(out)
	}

	if len(r.Result.Removed) > 0 {
		fmt.Fprintf(out, "## Removed\n\n| Key | Value |\n|-----|-------|\n")
		for _, k := range sortedKeys(r.Result.Removed) {
			fmt.Fprintf(out, "| `%s` | `%s` |\n", k, r.Result.Removed[k])
		}
		fmt.Fprintln(out)
	}

	if len(r.Result.Changed) > 0 {
		fmt.Fprintf(out, "## Changed\n\n| Key | Staging | Production |\n|-----|---------|------------|\n")
		keys := make([]string, 0, len(r.Result.Changed))
		for k := range r.Result.Changed {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			c := r.Result.Changed[k]
			fmt.Fprintf(out, "| `%s` | `%s` | `%s` |\n", k, c.Staging, c.Production)
		}
		fmt.Fprintln(out)
	}

	return nil
}
