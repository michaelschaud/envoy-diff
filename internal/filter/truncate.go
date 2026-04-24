package filter

import (
	"fmt"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// ApplyTruncate shortens values longer than maxLen to maxLen characters,
// appending a suffix (e.g. "...") to indicate truncation.
// If maxLen <= 0 the result is returned unchanged.
func ApplyTruncate(result diff.Result, maxLen int, suffix string) diff.Result {
	if maxLen <= 0 {
		return result
	}
	if suffix == "" {
		suffix = "..."
	}
	return diff.Result{
		Added:   truncateMap(result.Added, maxLen, suffix),
		Removed: truncateMap(result.Removed, maxLen, suffix),
		Same:    truncateMap(result.Same, maxLen, suffix),
		Changed: truncateChanged(result.Changed, maxLen, suffix),
	}
}

func truncateMap(m map[string]string, maxLen int, suffix string) map[string]string {
	if len(m) == 0 {
		return m
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = truncateValue(v, maxLen, suffix)
	}
	return out
}

func truncateChanged(m map[string]diff.ChangedValue, maxLen int, suffix string) map[string]diff.ChangedValue {
	if len(m) == 0 {
		return m
	}
	out := make(map[string]diff.ChangedValue, len(m))
	for k, cv := range m {
		out[k] = diff.ChangedValue{
			Old: truncateValue(cv.Old, maxLen, suffix),
			New: truncateValue(cv.New, maxLen, suffix),
		}
	}
	return out
}

func truncateValue(v string, maxLen int, suffix string) string {
	if len(v) <= maxLen {
		return v
	}
	if maxLen <= len(suffix) {
		return fmt.Sprintf("%.*s", maxLen, suffix)
	}
	return v[:maxLen-len(suffix)] + suffix
}
