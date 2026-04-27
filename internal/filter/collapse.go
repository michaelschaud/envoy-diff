package filter

import (
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// CollapseConfig controls how values are collapsed.
type CollapseConfig struct {
	// Delimiter is the string used to join repeated values (default: ",").
	Delimiter string
	// Keys is the list of key substrings to target. Empty means all keys.
	Keys []string
}

// ApplyCollapse merges duplicate values within each key's value by splitting
// on the delimiter, deduplicating, and rejoining. This is useful when env
// values are comma-separated lists that may contain repeated entries.
func ApplyCollapse(result diff.Result, cfg CollapseConfig) diff.Result {
	if cfg.Delimiter == "" {
		return result
	}

	return diff.Result{
		Added:   collapseMap(result.Added, cfg),
		Removed: collapseMap(result.Removed, cfg),
		Same:    collapseMap(result.Same, cfg),
		Changed: collapseChangedMap(result.Changed, cfg),
	}
}

func collapseMap(m map[string]string, cfg CollapseConfig) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if collapseKeyMatches(k, cfg.Keys) {
			out[k] = collapseValue(v, cfg.Delimiter)
		} else {
			out[k] = v
		}
	}
	return out
}

func collapseChangedMap(m map[string]diff.ChangedValue, cfg CollapseConfig) map[string]diff.ChangedValue {
	out := make(map[string]diff.ChangedValue, len(m))
	for k, cv := range m {
		if collapseKeyMatches(k, cfg.Keys) {
			out[k] = diff.ChangedValue{
				Old: collapseValue(cv.Old, cfg.Delimiter),
				New: collapseValue(cv.New, cfg.Delimiter),
			}
		} else {
			out[k] = cv
		}
	}
	return out
}

func collapseValue(v, delimiter string) string {
	parts := strings.Split(v, delimiter)
	seen := make(map[string]struct{}, len(parts))
	uniq := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if _, ok := seen[p]; !ok {
			seen[p] = struct{}{}
			uniq = append(uniq, p)
		}
	}
	return strings.Join(uniq, delimiter)
}

func collapseKeyMatches(key string, keys []string) bool {
	if len(keys) == 0 {
		return true
	}
	lower := strings.ToLower(key)
	for _, sub := range keys {
		if strings.Contains(lower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}
