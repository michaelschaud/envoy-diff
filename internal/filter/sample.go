package filter

import (
	"math/rand"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// SampleConfig controls how sampling is applied.
type SampleConfig struct {
	// Rate is a value between 0.0 and 1.0 representing the fraction of keys to keep.
	Rate float64
	// Seed is used to make sampling deterministic when non-zero.
	Seed int64
}

// ApplySample randomly retains a fraction of keys in each category of the
// diff result. A rate of 1.0 keeps all keys; 0.0 keeps none.
// When Seed is non-zero the selection is deterministic.
func ApplySample(result diff.Result, cfg SampleConfig) diff.Result {
	if cfg.Rate <= 0 {
		return diff.Result{}
	}
	if cfg.Rate >= 1.0 {
		return result
	}

	r := rand.New(rand.NewSource(cfg.Seed)) //nolint:gosec

	return diff.Result{
		Added:   sampleMap(result.Added, cfg.Rate, r),
		Removed: sampleMap(result.Removed, cfg.Rate, r),
		Same:    sampleMap(result.Same, cfg.Rate, r),
		Changed: sampleChanged(result.Changed, cfg.Rate, r),
	}
}

func sampleMap(m map[string]string, rate float64, r *rand.Rand) map[string]string {
	if len(m) == 0 {
		return m
	}
	out := make(map[string]string)
	for k, v := range m {
		if r.Float64() < rate {
			out[k] = v
		}
	}
	return out
}

func sampleChanged(m map[string]diff.ChangedValue, rate float64, r *rand.Rand) map[string]diff.ChangedValue {
	if len(m) == 0 {
		return m
	}
	out := make(map[string]diff.ChangedValue)
	for k, v := range m {
		if r.Float64() < rate {
			out[k] = v
		}
	}
	return out
}
