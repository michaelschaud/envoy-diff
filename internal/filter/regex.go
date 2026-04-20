package filter

import (
	"regexp"
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// MatchesRegex reports whether key matches any of the provided regex patterns.
// Matching is case-insensitive.
func MatchesRegex(key string, patterns []string) bool {
	if len(patterns) == 0 {
		return false
	}
	lower := strings.ToLower(key)
	for _, p := range patterns {
		re, err := regexp.Compile("(?i)" + p)
		if err != nil {
			continue
		}
		if re.MatchString(lower) {
			return true
		}
	}
	return false
}

// ApplyRegexExclude removes keys from result that match any of the given regex patterns.
func ApplyRegexExclude(result diff.Result, patterns []string) diff.Result {
	if len(patterns) == 0 {
		return result
	}
	return diff.Result{
		Added:   excludeByRegex(result.Added, patterns),
		Removed: excludeByRegex(result.Removed, patterns),
		Changed: excludeChangedByRegex(result.Changed, patterns),
		Same:    excludeByRegex(result.Same, patterns),
	}
}

func excludeByRegex(m map[string]string, patterns []string) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		if !MatchesRegex(k, patterns) {
			out[k] = v
		}
	}
	return out
}

func excludeChangedByRegex(m map[string][2]string, patterns []string) map[string][2]string {
	out := make(map[string][2]string)
	for k, v := range m {
		if !MatchesRegex(k, patterns) {
			out[k] = v
		}
	}
	return out
}
