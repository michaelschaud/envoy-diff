package filter

import "github.com/your-org/envoy-diff/internal/diff"

// BaselineConfig controls how a baseline comparison is applied.
// When a baseline result is provided, keys whose values are identical
// to the baseline are suppressed from the diff output.
type BaselineConfig struct {
	// Baseline holds the reference diff.Result to compare against.
	Baseline *diff.Result
	// IncludeNewKeys, when true, retains keys that are absent from the baseline.
	IncludeNewKeys bool
}

// ApplyBaseline removes entries from result that are unchanged relative to
// the provided baseline. This is useful when diffing against a known-good
// snapshot so that only regressions or new deviations are surfaced.
func ApplyBaseline(result diff.Result, cfg BaselineConfig) diff.Result {
	if cfg.Baseline == nil {
		return result
	}

	base := cfg.Baseline

	return diff.Result{
		Added:   filterBaselineMap(result.Added, base.Added, cfg.IncludeNewKeys),
		Removed: filterBaselineMap(result.Removed, base.Removed, cfg.IncludeNewKeys),
		Same:    result.Same,
		Changed: filterBaselineChanged(result.Changed, base.Changed, cfg.IncludeNewKeys),
	}
}

// filterBaselineMap drops keys whose value in current matches the baseline map.
func filterBaselineMap(current, base map[string]string, includeNew bool) map[string]string {
	out := make(map[string]string)
	for k, v := range current {
		baseVal, exists := base[k]
		if !exists {
			if includeNew {
				out[k] = v
			}
			continue
		}
		if v != baseVal {
			out[k] = v
		}
	}
	return out
}

// filterBaselineChanged drops changed entries whose old/new values match the baseline.
func filterBaselineChanged(current, base map[string]diff.ChangedValue, includeNew bool) map[string]diff.ChangedValue {
	out := make(map[string]diff.ChangedValue)
	for k, cv := range current {
		baseCV, exists := base[k]
		if !exists {
			if includeNew {
				out[k] = cv
			}
			continue
		}
		if cv.Old != baseCV.Old || cv.New != baseCV.New {
			out[k] = cv
		}
	}
	return out
}
