package filter

import (
	"strings"

	"github.com/user/envoy-diff/internal/diff"
)

// HighlightConfig controls which keys get a highlight marker appended to their values.
type HighlightConfig struct {
	// Substrings is a list of key substrings to match (case-insensitive).
	Substrings []string
	// Marker is the string appended to matching values, e.g. " [!]".
	Marker string
}

// ApplyHighlight appends a marker to values whose key contains any of the
// configured substrings. This is useful for drawing attention to sensitive or
// important keys in rendered output.
func ApplyHighlight(result diff.Result, cfg HighlightConfig) diff.Result {
	if len(cfg.Substrings) == 0 || cfg.Marker == "" {
		return result
	}
	return diff.Result{
		Added:   highlightMap(result.Added, cfg),
		Removed: highlightMap(result.Removed, cfg),
		Same:    highlightMap(result.Same, cfg),
		Changed: highlightChanged(result.Changed, cfg),
	}
}

func highlightMap(m map[string]string, cfg HighlightConfig) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if keyMatchesSubstring(k, cfg.Substrings) {
			out[k] = v + cfg.Marker
		} else {
			out[k] = v
		}
	}
	return out
}

func highlightChanged(m map[string][2]string, cfg HighlightConfig) map[string][2]string {
	out := make(map[string][2]string, len(m))
	for k, pair := range m {
		if keyMatchesSubstring(k, cfg.Substrings) {
			out[k] = [2]string{pair[0] + cfg.Marker, pair[1] + cfg.Marker}
		} else {
			out[k] = pair
		}
	}
	return out
}

// keyMatchesSubstring returns true if key contains any substring (case-insensitive).
// NOTE: reuses the helper already defined in mask.go within this package.
func highlightKeyMatches(key string, substrings []string) bool {
	lower := strings.ToLower(key)
	for _, s := range substrings {
		if strings.Contains(lower, strings.ToLower(s)) {
			return true
		}
	}
	return false
}
