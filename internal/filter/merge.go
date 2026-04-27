package filter

import "github.com/yourorg/envoy-diff/internal/diff"

// MergeStrategy controls how overlapping keys are resolved when merging.
type MergeStrategy string

const (
	MergeStrategyLeft  MergeStrategy = "left"
	MergeStrategyRight MergeStrategy = "right"
	MergeStrategyUnion MergeStrategy = "union"
)

// MergeConfig holds settings for ApplyMerge.
type MergeConfig struct {
	Overlay  map[string]string
	Strategy MergeStrategy
}

// ApplyMerge merges an overlay map into the diff result according to the
// chosen strategy, then re-runs the comparison so added/removed/changed
// categories reflect the merged state.
func ApplyMerge(result diff.Result, cfg MergeConfig) diff.Result {
	if len(cfg.Overlay) == 0 {
		return result
	}

	strategy := cfg.Strategy
	if strategy == "" {
		strategy = MergeStrategyRight
	}

	newSame := mergeMap(result.Same, cfg.Overlay, strategy)
	newAdded := mergeMap(result.Added, cfg.Overlay, strategy)
	newRemoved := mergeMap(result.Removed, cfg.Overlay, strategy)

	newChanged := make(map[string][2]string)
	for k, pair := range result.Changed {
		if ov, ok := cfg.Overlay[k]; ok {
			switch strategy {
			case MergeStrategyLeft:
				newChanged[k] = pair
			case MergeStrategyRight:
				newChanged[k] = [2]string{pair[0], ov}
			case MergeStrategyUnion:
				newChanged[k] = [2]string{pair[0], ov}
			}
		} else {
			newChanged[k] = pair
		}
	}

	return diff.Result{
		Added:   newAdded,
		Removed: newRemoved,
		Changed: newChanged,
		Same:    newSame,
	}
}

func mergeMap(base map[string]string, overlay map[string]string, strategy MergeStrategy) map[string]string {
	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}
	if strategy == MergeStrategyLeft {
		return out
	}
	for k, v := range overlay {
		if _, exists := out[k]; exists || strategy == MergeStrategyUnion {
			out[k] = v
		}
	}
	return out
}
