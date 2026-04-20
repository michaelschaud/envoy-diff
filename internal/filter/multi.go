package filter

import (
	"github.com/yourorg/envoy-diff/internal/diff"
)

// ApplyMulti applies filtering across multiple prefixes and an optional onlyChanged flag.
// If prefixes is empty, no prefix filtering is applied.
// If onlyChanged is true, added and removed keys are excluded from the result.
func ApplyMulti(r diff.Result, prefixes []string, onlyChanged bool) diff.Result {
	if len(prefixes) == 0 {
		return Apply(r, "", onlyChanged)
	}

	merged := diff.Result{
		Added:     make(map[string]string),
		Removed:   make(map[string]string),
		Changed:   make(map[string][2]string),
		Unchanged: make(map[string]string),
	}

	for _, prefix := range prefixes {
		filtered := Apply(r, prefix, onlyChanged)

		for k, v := range filtered.Added {
			merged.Added[k] = v
		}
		for k, v := range filtered.Removed {
			merged.Removed[k] = v
		}
		for k, v := range filtered.Changed {
			merged.Changed[k] = v
		}
		for k, v := range filtered.Unchanged {
			merged.Unchanged[k] = v
		}
	}

	return merged
}
