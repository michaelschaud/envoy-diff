package filter

import (
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// FlattenConfig controls how nested key structures are flattened.
type FlattenConfig struct {
	// Delimiter is the separator used to join nested key segments (e.g. "__" or ".").
	Delimiter string
	// Depth is the maximum nesting depth to flatten. 0 means unlimited.
	Depth int
}

// ApplyFlatten collapses keys that share a common prefix+delimiter into a
// single canonical key by stripping the prefix segment up to Depth levels.
// For example, with delimiter "__" and depth 1:
//   APP__HOST -> HOST
//   APP__DB__PORT -> DB__PORT
func ApplyFlatten(result diff.Result, cfg FlattenConfig) diff.Result {
	if cfg.Delimiter == "" {
		return result
	}
	return diff.Result{
		Added:   flattenMap(result.Added, cfg),
		Removed: flattenMap(result.Removed, cfg),
		Same:    flattenMap(result.Same, cfg),
		Changed: flattenChangedMap(result.Changed, cfg),
	}
}

func flattenMap(m map[string]string, cfg FlattenConfig) map[string]string {
	if len(m) == 0 {
		return m
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[flattenKey(k, cfg)] = v
	}
	return out
}

func flattenChangedMap(m map[string][2]string, cfg FlattenConfig) map[string][2]string {
	if len(m) == 0 {
		return m
	}
	out := make(map[string][2]string, len(m))
	for k, v := range m {
		out[flattenKey(k, cfg)] = v
	}
	return out
}

// flattenKey strips leading prefix segments separated by cfg.Delimiter up to
// cfg.Depth times. A Depth of 0 removes all leading segments.
func flattenKey(key string, cfg FlattenConfig) string {
	parts := strings.Split(key, cfg.Delimiter)
	if len(parts) <= 1 {
		return key
	}
	strips := len(parts) - 1
	if cfg.Depth > 0 && cfg.Depth < strips {
		strips = cfg.Depth
	}
	return strings.Join(parts[strips:], cfg.Delimiter)
}
