package filter

import (
	"strings"

	"github.com/your-org/envoy-diff/internal/diff"
)

// MaskPatterns holds substrings that, when found in a key name, cause
// the corresponding value(s) to be replaced with a redaction marker.
const maskMarker = "***REDACTED***"

// ApplyMask redacts values whose key contains any of the given substrings
// (case-insensitive). It operates on a copy of the result so the original
// is not mutated.
func ApplyMask(result diff.Result, substrings []string) diff.Result {
	if len(substrings) == 0 {
		return result
	}

	return diff.Result{
		Added:   maskMap(result.Added, substrings),
		Removed: maskMap(result.Removed, substrings),
		Same:    maskMap(result.Same, substrings),
		Changed: maskChanged(result.Changed, substrings),
	}
}

func maskMap(m map[string]string, substrings []string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if keyMatchesSubstring(k, substrings) {
			out[k] = maskMarker
		} else {
			out[k] = v
		}
	}
	return out
}

func maskChanged(m map[string][2]string, substrings []string) map[string][2]string {
	out := make(map[string][2]string, len(m))
	for k, pair := range m {
		if keyMatchesSubstring(k, substrings) {
			out[k] = [2]string{maskMarker, maskMarker}
		} else {
			out[k] = pair
		}
	}
	return out
}

func keyMatchesSubstring(key string, substrings []string) bool {
	lower := strings.ToLower(key)
	for _, s := range substrings {
		if strings.Contains(lower, strings.ToLower(s)) {
			return true
		}
	}
	return false
}
