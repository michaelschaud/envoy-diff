package filter

import (
	"strconv"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// ThresholdConfig controls filtering based on numeric value thresholds.
type ThresholdConfig struct {
	Min *float64
	Max *float64
}

// ApplyThreshold removes entries whose values fall outside the configured
// numeric range. Non-numeric values are always kept.
func ApplyThreshold(r diff.Result, cfg ThresholdConfig) diff.Result {
	if cfg.Min == nil && cfg.Max == nil {
		return r
	}
	return diff.Result{
		Added:   thresholdMap(r.Added, cfg),
		Removed: thresholdMap(r.Removed, cfg),
		Changed: thresholdChanged(r.Changed, cfg),
		Same:    thresholdMap(r.Same, cfg),
	}
}

func thresholdMap(m map[string]string, cfg ThresholdConfig) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if withinThreshold(v, cfg) {
			out[k] = v
		}
	}
	return out
}

func thresholdChanged(m map[string][2]string, cfg ThresholdConfig) map[string][2]string {
	out := make(map[string][2]string, len(m))
	for k, pair := range m {
		// Keep if either old or new value is within threshold
		if withinThreshold(pair[0], cfg) || withinThreshold(pair[1], cfg) {
			out[k] = pair
		}
	}
	return out
}

func withinThreshold(v string, cfg ThresholdConfig) bool {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		// Non-numeric: keep unconditionally
		return true
	}
	if cfg.Min != nil && f < *cfg.Min {
		return false
	}
	if cfg.Max != nil && f > *cfg.Max {
		return false
	}
	return true
}
