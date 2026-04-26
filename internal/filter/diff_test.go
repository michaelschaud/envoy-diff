package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeDiffResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"NEW_KEY": "newval"},
		Removed: map[string]string{"OLD_KEY": "oldval"},
		Changed: map[string]diff.ChangedValue{
			"MOD_KEY": {Old: "before", New: "after"},
		},
		Same: map[string]string{"STABLE": "same"},
	}
}

func TestParseDiffModes_Valid(t *testing.T) {
	modes := ParseDiffModes("added,removed")
	if len(modes) != 2 {
		t.Fatalf("expected 2 modes, got %d", len(modes))
	}
}

func TestParseDiffModes_Empty(t *testing.T) {
	modes := ParseDiffModes("")
	if len(modes) != 0 {
		t.Fatalf("expected 0 modes, got %d", len(modes))
	}
}

func TestParseDiffModes_UnknownIgnored(t *testing.T) {
	modes := ParseDiffModes("added,unknown,changed")
	if len(modes) != 2 {
		t.Fatalf("expected 2 modes, got %d", len(modes))
	}
}

func TestApplyDiffModeFilter_NoModes(t *testing.T) {
	r := makeDiffResult()
	out := ApplyDiffModeFilter(r, nil)
	if len(out.Added) != 1 || len(out.Removed) != 1 || len(out.Changed) != 1 || len(out.Same) != 1 {
		t.Error("expected result unchanged when no modes specified")
	}
}

func TestApplyDiffModeFilter_OnlyAdded(t *testing.T) {
	r := makeDiffResult()
	out := ApplyDiffModeFilter(r, []DiffMode{DiffModeAdded})
	if len(out.Added) != 1 {
		t.Errorf("expected 1 added key, got %d", len(out.Added))
	}
	if len(out.Removed) != 0 || len(out.Changed) != 0 || len(out.Same) != 0 {
		t.Error("expected removed/changed/same to be empty")
	}
}

func TestApplyDiffModeFilter_AddedAndChanged(t *testing.T) {
	r := makeDiffResult()
	out := ApplyDiffModeFilter(r, []DiffMode{DiffModeAdded, DiffModeChanged})
	if len(out.Added) != 1 {
		t.Errorf("expected 1 added, got %d", len(out.Added))
	}
	if len(out.Changed) != 1 {
		t.Errorf("expected 1 changed, got %d", len(out.Changed))
	}
	if len(out.Removed) != 0 || len(out.Same) != 0 {
		t.Error("expected removed/same to be empty")
	}
}

func TestApplyDiffModeFilter_SameOnly(t *testing.T) {
	r := makeDiffResult()
	out := ApplyDiffModeFilter(r, []DiffMode{DiffModeSame})
	if len(out.Same) != 1 {
		t.Errorf("expected 1 same key, got %d", len(out.Same))
	}
	if len(out.Added) != 0 || len(out.Removed) != 0 || len(out.Changed) != 0 {
		t.Error("expected added/removed/changed to be empty")
	}
}
