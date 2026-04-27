package filter

import (
	"strconv"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// ClampConfig holds the numeric bounds for clamping values.
type ClampConfig struct {
	Min *float64
	Max *float64
}

// ApplyClamp restricts numeric env var values to the [Min, Max] range.
// Non-numeric values are left unchanged. Nil bounds are treated as unbounded.
func ApplyClamp(result diff.Result, cfg ClampConfig) diff.Result {
	if cfg.Min == nil && cfg.Max == nil {
		return result
	}
	return diff.Result{
		Added:   clampMap(result.Added, cfg),
		Removed: clampMap(result.Removed, cfg),
		Same:    clampMap(result.Same, cfg),
		Changed: clampChangedMap(result.Changed, cfg),
	}
}

func clampMap(m map[string]string, cfg ClampConfig) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = clampValue(v, cfg)
	}
	return out
}

func clampChangedMap(m map[string][2]string, cfg ClampConfig) map[string][2]string {
	out := make(map[string][2]string, len(m))
	for k, pair := range m {
		out[k] = [2]string{
			clampValue(pair[0], cfg),
			clampValue(pair[1], cfg),
		}
	}
	return out
}

func clampValue(v string, cfg ClampConfig) string {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return v
	}
	if cfg.Min != nil && f < *cfg.Min {
		f = *cfg.Min
	}
	if cfg.Max != nil && f > *cfg.Max {
		f = *cfg.Max
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}
