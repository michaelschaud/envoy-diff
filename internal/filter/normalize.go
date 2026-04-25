package filter

import (
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// NormalizeMode controls how values are normalized.
type NormalizeMode string

const (
	NormalizeTrimSpace NormalizeMode = "trim"
	NormalizeTrimQuotes NormalizeMode = "unquote"
	NormalizeBoth NormalizeMode = "both"
)

// ApplyNormalize applies value normalization to all entries in the result.
// Supported modes: "trim" (trim whitespace), "unquote" (strip surrounding quotes),
// "both" (trim then unquote).
func ApplyNormalize(result diff.Result, mode NormalizeMode) diff.Result {
	if mode == "" {
		return result
	}
	normFn := resolveNormalize(mode)
	return diff.Result{
		Added:   normalizeMap(result.Added, normFn),
		Removed: normalizeMap(result.Removed, normFn),
		Same:    normalizeMap(result.Same, normFn),
		Changed: normalizeChanged(result.Changed, normFn),
	}
}

func resolveNormalize(mode NormalizeMode) func(string) string {
	switch mode {
	case NormalizeTrimSpace:
		return strings.TrimSpace
	case NormalizeTrimQuotes:
		return stripSurroundingQuotes
	case NormalizeBoth:
		return func(s string) string {
			return stripSurroundingQuotes(strings.TrimSpace(s))
		}
	default:
		return func(s string) string { return s }
	}
}

func normalizeMap(m map[string]string, fn func(string) string) map[string]string {
	if len(m) == 0 {
		return m
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = fn(v)
	}
	return out
}

func normalizeChanged(m map[string][2]string, fn func(string) string) map[string][2]string {
	if len(m) == 0 {
		return m
	}
	out := make(map[string][2]string, len(m))
	for k, pair := range m {
		out[k] = [2]string{fn(pair[0]), fn(pair[1])}
	}
	return out
}

func stripSurroundingQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
