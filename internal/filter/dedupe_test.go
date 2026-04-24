package filter

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/diff"
)

func makeDedupeResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"A": "alpha", "B": "alpha", "C": "gamma"},
		Removed: map[string]string{"X": "delta", "Y": "delta", "Z": "unique"},
		Same:    map[string]string{"P": "same", "Q": "same", "R": "only"},
		Changed: map[string]diff.ChangedValue{"K": {Old: "a", New: "b"}},
	}
}

func TestApplyDedupe_NoMode(t *testing.T) {
	r := makeDedupeResult()
	out := ApplyDedupe(r, "")
	if len(out.Added) != 3 {
		t.Errorf("expected 3 added, got %d", len(out.Added))
	}
}

func TestApplyDedupe_KeepFirst(t *testing.T) {
	r := makeDedupeResult()
	out := ApplyDedupe(r, DedupeKeepFirst)
	// A and B share "alpha"; sorted first = A kept
	if _, ok := out.Added["A"]; !ok {
		t.Error("expected A to be kept (first)")
	}
	if _, ok := out.Added["B"]; ok {
		t.Error("expected B to be removed (duplicate, not first)")
	}
	if _, ok := out.Added["C"]; !ok {
		t.Error("expected C to be kept (unique value)")
	}
	// Changed keys must be untouched
	if len(out.Changed) != 1 {
		t.Error("changed keys must not be deduped")
	}
}

func TestApplyDedupe_KeepLast(t *testing.T) {
	r := makeDedupeResult()
	out := ApplyDedupe(r, DedupeKeepLast)
	// sorted order A < B; last = B
	if _, ok := out.Added["B"]; !ok {
		t.Error("expected B to be kept (last)")
	}
	if _, ok := out.Added["A"]; ok {
		t.Error("expected A to be removed (not last)")
	}
}

func TestApplyDedupe_RemoveAll(t *testing.T) {
	r := makeDedupeResult()
	out := ApplyDedupe(r, DedupeRemoveAll)
	if _, ok := out.Added["A"]; ok {
		t.Error("expected A removed")
	}
	if _, ok := out.Added["B"]; ok {
		t.Error("expected B removed")
	}
	if _, ok := out.Added["C"]; !ok {
		t.Error("expected C kept (unique value)")
	}
	if _, ok := out.Removed["Z"]; !ok {
		t.Error("expected Z kept in removed")
	}
}

func TestApplyDedupe_EmptyMaps(t *testing.T) {
	r := diff.Result{}
	out := ApplyDedupe(r, DedupeRemoveAll)
	if len(out.Added) != 0 || len(out.Removed) != 0 {
		t.Error("expected empty result for empty input")
	}
}
