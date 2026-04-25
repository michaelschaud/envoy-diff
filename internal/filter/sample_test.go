package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeSampleResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"A": "1", "B": "2", "C": "3", "D": "4"},
		Removed: map[string]string{"W": "9", "X": "8", "Y": "7", "Z": "6"},
		Same:    map[string]string{"P": "p", "Q": "q"},
		Changed: map[string]diff.ChangedValue{
			"M": {Old: "old", New: "new"},
			"N": {Old: "a", New: "b"},
		},
	}
}

func TestApplySample_ZeroRateReturnsEmpty(t *testing.T) {
	result := makeSampleResult()
	out := ApplySample(result, SampleConfig{Rate: 0})
	if len(out.Added)+len(out.Removed)+len(out.Same)+len(out.Changed) != 0 {
		t.Errorf("expected empty result for rate=0, got %+v", out)
	}
}

func TestApplySample_FullRateReturnsAll(t *testing.T) {
	result := makeSampleResult()
	out := ApplySample(result, SampleConfig{Rate: 1.0})
	if len(out.Added) != len(result.Added) {
		t.Errorf("expected all added keys, got %d", len(out.Added))
	}
	if len(out.Removed) != len(result.Removed) {
		t.Errorf("expected all removed keys, got %d", len(out.Removed))
	}
	if len(out.Same) != len(result.Same) {
		t.Errorf("expected all same keys, got %d", len(out.Same))
	}
	if len(out.Changed) != len(result.Changed) {
		t.Errorf("expected all changed keys, got %d", len(out.Changed))
	}
}

func TestApplySample_DeterministicWithSeed(t *testing.T) {
	result := makeSampleResult()
	cfg := SampleConfig{Rate: 0.5, Seed: 42}
	out1 := ApplySample(result, cfg)
	out2 := ApplySample(result, cfg)
	if len(out1.Added) != len(out2.Added) {
		t.Errorf("expected deterministic sampling: run1=%d run2=%d", len(out1.Added), len(out2.Added))
	}
	for k := range out1.Added {
		if _, ok := out2.Added[k]; !ok {
			t.Errorf("key %q present in first run but not second", k)
		}
	}
}

func TestApplySample_SubsetOfOriginal(t *testing.T) {
	result := makeSampleResult()
	out := ApplySample(result, SampleConfig{Rate: 0.5, Seed: 7})
	for k := range out.Added {
		if _, ok := result.Added[k]; !ok {
			t.Errorf("sampled key %q not in original added set", k)
		}
	}
	for k := range out.Changed {
		if _, ok := result.Changed[k]; !ok {
			t.Errorf("sampled changed key %q not in original changed set", k)
		}
	}
}
