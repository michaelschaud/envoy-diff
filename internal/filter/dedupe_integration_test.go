package filter

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/diff"
)

func TestApplyDedupe_WithApplyMulti(t *testing.T) {
	results := []diff.Result{
		{
			Added:   map[string]string{"A": "shared", "B": "shared", "C": "unique"},
			Removed: map[string]string{"X": "dup", "Y": "dup"},
			Same:    map[string]string{},
			Changed: map[string]diff.ChangedValue{},
		},
		{
			Added:   map[string]string{"D": "only"},
			Removed: map[string]string{"Z": "solo"},
			Same:    map[string]string{},
			Changed: map[string]diff.ChangedValue{},
		},
	}

	merged := ApplyMulti(results)
	out := ApplyDedupe(merged, DedupeRemoveAll)

	// A and B share "shared" — both removed
	if _, ok := out.Added["A"]; ok {
		t.Error("expected A deduped")
	}
	if _, ok := out.Added["B"]; ok {
		t.Error("expected B deduped")
	}
	// C is unique
	if _, ok := out.Added["C"]; !ok {
		t.Error("expected C retained")
	}
	// D from second result
	if _, ok := out.Added["D"]; !ok {
		t.Error("expected D retained")
	}
	// X and Y share "dup"
	if _, ok := out.Removed["X"]; ok {
		t.Error("expected X deduped")
	}
	if _, ok := out.Removed["Y"]; ok {
		t.Error("expected Y deduped")
	}
}
