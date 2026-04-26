package filter

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/diff"
)

func makeBaselineResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"NEW_KEY": "new_val", "SHARED_ADD": "same"},
		Removed: map[string]string{"OLD_KEY": "old_val", "SHARED_REM": "same"},
		Same:    map[string]string{"STABLE": "stable_val"},
		Changed: map[string]diff.ChangedValue{
			"MOD_KEY":    {Old: "a", New: "b"},
			"SHARED_MOD": {Old: "x", New: "y"},
		},
	}
}

func TestApplyBaseline_NilBaseline(t *testing.T) {
	result := makeBaselineResult()
	out := ApplyBaseline(result, BaselineConfig{Baseline: nil})
	if len(out.Added) != 2 || len(out.Removed) != 2 || len(out.Changed) != 2 {
		t.Errorf("expected result unchanged when baseline is nil, got added=%d removed=%d changed=%d",
			len(out.Added), len(out.Removed), len(out.Changed))
	}
}

func TestApplyBaseline_SuppressesMatchingKeys(t *testing.T) {
	result := makeBaselineResult()
	base := &diff.Result{
		Added:   map[string]string{"SHARED_ADD": "same"},
		Removed: map[string]string{"SHARED_REM": "same"},
		Changed: map[string]diff.ChangedValue{
			"SHARED_MOD": {Old: "x", New: "y"},
		},
	}
	out := ApplyBaseline(result, BaselineConfig{Baseline: base, IncludeNewKeys: false})

	if _, ok := out.Added["SHARED_ADD"]; ok {
		t.Error("expected SHARED_ADD to be suppressed by baseline")
	}
	if _, ok := out.Removed["SHARED_REM"]; ok {
		t.Error("expected SHARED_REM to be suppressed by baseline")
	}
	if _, ok := out.Changed["SHARED_MOD"]; ok {
		t.Error("expected SHARED_MOD to be suppressed by baseline")
	}
}

func TestApplyBaseline_IncludeNewKeys(t *testing.T) {
	result := makeBaselineResult()
	base := &diff.Result{
		Added:   map[string]string{"SHARED_ADD": "same"},
		Removed: map[string]string{},
		Changed: map[string]diff.ChangedValue{},
	}
	out := ApplyBaseline(result, BaselineConfig{Baseline: base, IncludeNewKeys: true})

	if _, ok := out.Added["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY to be retained when IncludeNewKeys is true")
	}
}

func TestApplyBaseline_PreservesDeviatingValues(t *testing.T) {
	result := makeBaselineResult()
	base := &diff.Result{
		Added:   map[string]string{"NEW_KEY": "different_val"},
		Removed: map[string]string{},
		Changed: map[string]diff.ChangedValue{
			"MOD_KEY": {Old: "a", New: "DIFFERENT"},
		},
	}
	out := ApplyBaseline(result, BaselineConfig{Baseline: base, IncludeNewKeys: true})

	if _, ok := out.Added["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY to be retained because value differs from baseline")
	}
	if _, ok := out.Changed["MOD_KEY"]; !ok {
		t.Error("expected MOD_KEY to be retained because changed value differs from baseline")
	}
}
