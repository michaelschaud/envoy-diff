package formatter

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envoy-diff/internal/diff"
)

// MarkdownWriter writes diff results as a Markdown table.
type MarkdownWriter struct {
	w io.Writer
}

// NewMarkdownWriter creates a new MarkdownWriter.
func NewMarkdownWriter(w io.Writer) *MarkdownWriter {
	return &MarkdownWriter{w: w}
}

// Write outputs the diff result as Markdown.
func (m *MarkdownWriter) Write(result diff.Result) error {
	fmt.Fprintln(m.w, "# Environment Diff Report")
	fmt.Fprintln(m.w)

	sections := []struct {
		title string
		keys  []string
		render func(key string) string
	}{
		{
			title: "## Added (only in production)",
			keys:  sortedKeys(result.Added),
			render: func(k string) string {
				return fmt.Sprintf("| `%s` | — | `%s` |", k, result.Added[k])
			},
		},
		{
			title: "## Removed (only in staging)",
			keys:  sortedKeys(result.Removed),
			render: func(k string) string {
				return fmt.Sprintf("| `%s` | `%s` | — |", k, result.Removed[k])
			},
		},
		{
			title: "## Changed",
			keys: func() []string {
				keys := make([]string, 0, len(result.Changed))
				for k := range result.Changed {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				return keys
			}(),
			render: func(k string) string {
				c := result.Changed[k]
				return fmt.Sprintf("| `%s` | `%s` | `%s` |", k, c.Staging, c.Production)
			},
		},
	}

	for _, sec := range sections {
		if len(sec.keys) == 0 {
			continue
		}
		fmt.Fprintln(m.w, sec.title)
		fmt.Fprintln(m.w)
		fmt.Fprintln(m.w, "| Key | Staging | Production |")
		fmt.Fprintln(m.w, "|-----|---------|------------|")
		for _, k := range sec.keys {
			fmt.Fprintln(m.w, sec.render(k))
		}
		fmt.Fprintln(m.w)
	}

	return nil
}
