package filter

import (
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// ApplyValueFilter removes keys whose staging or production value matches any of the given substrings.
// This is useful for hiding keys that contain sensitive placeholder values like "REDACTED" or "TODO".
func ApplyValueFilter(result diff.Result, substrings []string) diff.Result {
	if len(substrings) == 0 {
		return result
	}

	return diff.Result{
		Added:   filterByValue(result.Added, substrings),
		Removed: filterByValue(result.Removed, substrings),
		Changed: filterChangedByValue(result.Changed, substrings),
		Same:    filterByValue(result.Same, substrings),
	}
}

// filterByValue removes keys whose value contains any of the given substrings (case-insensitive).
func filterByValue(m map[string]string, substrings []string) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		if !valueMatchesAny(v, substrings) {
			out[k] = v
		}
	}
	return out
}

// filterChangedByValue removes changed keys where either the old or new value matches a substring.
func filterChangedByValue(m map[string][2]string, substrings []string) map[string][2]string {
	out := make(map[string][2]string)
	for k, pair := range m {
		if !valueMatchesAny(pair[0], substrings) && !valueMatchesAny(pair[1], substrings) {
			out[k] = pair
		}
	}
	return out
}

// valueMatchesAny returns true if the value contains any of the given substrings (case-insensitive).
func valueMatchesAny(value string, substrings []string) bool {
	lower := strings.ToLower(value)
	for _, sub := range substrings {
		if strings.Contains(lower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}
