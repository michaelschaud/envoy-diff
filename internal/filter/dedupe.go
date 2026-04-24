package filter

import "github.com/your-org/envoy-diff/internal/diff"

// DedupeMode controls how duplicate values are handled.
type DedupeMode string

const (
	DedupeKeepFirst DedupeMode = "first"
	DedupeKeepLast  DedupeMode = "last"
	DedupeRemoveAll DedupeMode = "all"
)

// ApplyDedupe removes entries whose values appear more than once across
// the added, removed, or same key sets. Changed keys are never deduped.
func ApplyDedupe(result diff.Result, mode DedupeMode) diff.Result {
	if mode == "" {
		return result
	}

	result.Added = dedupeMap(result.Added, mode)
	result.Removed = dedupeMap(result.Removed, mode)
	result.Same = dedupeMap(result.Same, mode)
	return result
}

// dedupeMap returns a new map with duplicated values handled per mode.
func dedupeMap(m map[string]string, mode DedupeMode) map[string]string {
	if len(m) == 0 {
		return m
	}

	// count occurrences of each value
	counts := make(map[string]int, len(m))
	for _, v := range m {
		counts[v]++
	}

	// track insertion order for first/last semantics
	keys := SortedKeys(m)
	seen := make(map[string]string) // value -> key
	out := make(map[string]string, len(m))

	switch mode {
	case DedupeRemoveAll:
		for _, k := range keys {
			if counts[m[k]] == 1 {
				out[k] = m[k]
			}
		}
	case DedupeKeepLast:
		for _, k := range keys {
			v := m[k]
			if counts[v] > 1 {
				if prev, ok := seen[v]; ok {
					delete(out, prev)
				}
				seen[v] = k
			}
			out[k] = v
		}
	default: // DedupeKeepFirst
		for _, k := range keys {
			v := m[k]
			if counts[v] == 1 {
				out[k] = v
				continue
			}
			if _, already := seen[v]; !already {
				seen[v] = k
				out[k] = v
			}
		}
	}
	return out
}
