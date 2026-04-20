package filter

import (
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// ApplyKeyFilter removes entries whose keys contain any of the given substrings.
// It operates on added, removed, and changed sets independently.
// If substrings is empty, the result is returned unchanged.
func ApplyKeyFilter(result diff.Result, substrings []string) diff.Result {
	if len(substrings) == 0 {
		return result
	}

	return diff.Result{
		Added:   filterKeysBySubstring(result.Added, substrings),
		Removed: filterKeysBySubstring(result.Removed, substrings),
		Changed: filterChangedKeysBySubstring(result.Changed, substrings),
		Same:    filterKeysBySubstring(result.Same, substrings),
	}
}

// KeyMatchesAny returns true if the key contains any of the given substrings
// (case-insensitive).
func KeyMatchesAny(key string, substrings []string) bool {
	lower := strings.ToLower(key)
	for _, sub := range substrings {
		if strings.Contains(lower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}

func filterKeysBySubstring(m map[string]string, substrings []string) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		if !KeyMatchesAny(k, substrings) {
			out[k] = v
		}
	}
	return out
}

func filterChangedKeysBySubstring(m map[string][2]string, substrings []string) map[string][2]string {
	out := make(map[string][2]string)
	for k, v := range m {
		if !KeyMatchesAny(k, substrings) {
			out[k] = v
		}
	}
	return out
}
