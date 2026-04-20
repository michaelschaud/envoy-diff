package filter

import (
	"path"
	"strings"
)

// MatchesGlob reports whether key matches any of the provided glob patterns.
// Patterns are matched case-insensitively against the key.
func MatchesGlob(key string, patterns []string) bool {
	if len(patterns) == 0 {
		return false
	}
	lower := strings.ToLower(key)
	for _, p := range patterns {
		matched, err := path.Match(strings.ToLower(p), lower)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// ApplyGlobExclude returns a copy of result with any keys matching one of the
// provided glob patterns removed from all categories (added, removed, changed,
// unchanged).
func ApplyGlobExclude(result DiffResult, patterns []string) DiffResult {
	if len(patterns) == 0 {
		return result
	}
	return DiffResult{
		Added:     excludeKeys(result.Added, patterns),
		Removed:   excludeKeys(result.Removed, patterns),
		Changed:   excludeChanged(result.Changed, patterns),
		Unchanged: excludeKeys(result.Unchanged, patterns),
	}
}

func excludeKeys(m map[string]string, patterns []string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if !MatchesGlob(k, patterns) {
			out[k] = v
		}
	}
	return out
}

func excludeChanged(m map[string][2]string, patterns []string) map[string][2]string {
	out := make(map[string][2]string, len(m))
	for k, v := range m {
		if !MatchesGlob(k, patterns) {
			out[k] = v
		}
	}
	return out
}
