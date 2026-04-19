package formatter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/envoy-diff/internal/diff"
)

// Format controls the output format.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// TextWriter writes a human-readable diff to w.
func TextWriter(w io.Writer, result diff.Result) {
	if len(result.OnlyInA) > 0 {
		fmt.Fprintln(w, "--- only in staging ---")
		keys := sortedKeys(result.OnlyInA)
		for _, k := range keys {
			fmt.Fprintf(w, "  - %s=%s\n", k, result.OnlyInA[k])
		}
	}

	if len(result.OnlyInB) > 0 {
		fmt.Fprintln(w, "+++ only in production +++")
		keys := sortedKeys(result.OnlyInB)
		for _, k := range keys {
			fmt.Fprintf(w, "  + %s=%s\n", k, result.OnlyInB[k])
		}
	}

	if len(result.Changed) > 0 {
		fmt.Fprintln(w, "~~~ changed values ~~~")
		keys := make([]string, 0, len(result.Changed))
		for k := range result.Changed {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			pair := result.Changed[k]
			fmt.Fprintf(w, "  ~ %s: %q -> %q\n", k, pair[0], pair[1])
		}
	}

	if len(result.OnlyInA) == 0 && len(result.OnlyInB) == 0 && len(result.Changed) == 0 {
		fmt.Fprintln(w, "No differences found.")
	}
}

// Summary returns a one-line summary string.
func Summary(result diff.Result) string {
	parts := []string{}
	if n := len(result.OnlyInA); n > 0 {
		parts = append(parts, fmt.Sprintf("%d only-in-staging", n))
	}
	if n := len(result.OnlyInB); n > 0 {
		parts = append(parts, fmt.Sprintf("%d only-in-production", n))
	}
	if n := len(result.Changed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", n))
	}
	if len(parts) == 0 {
		return "identical"
	}
	return strings.Join(parts, ", ")
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
