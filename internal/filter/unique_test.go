package filter

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/diff"
)

func makeUniqueResult() diff.Result {
	return diff.Result{
		Added: map[string]string{
			"A": "alpha",
			"B": "alpha", // duplicate value
			"C": "gamma",
		},
		Removed: map[string]string{
			"X": "x-val",
			"Y": "x-val", // duplicate value
		},
		Same: map[string]string{
			"S1": "shared",
			"S2": "unique",
		},
		Changed: map[string][2]string{
			"K1": {"old1", "new1"},
			"K2": {"old1", "new2"}, // old value duplicated
			"K3": {"old3", "new3"},
		},
	}
}

func TestApplyUnique_NoModes(t *testing.T) {
	res := ApplyUnique(makeUniqueResult(), nil)

	if _, ok := res.Added["A"]; ok {
		t.Error("expected A to be removed (duplicate value 'alpha')")
	}
	if _, ok := res.Added["B"]; ok {
		t.Error("expected B to be removed (duplicate value 'alpha')")
	}
	if v, ok := res.Added["C"]; !ok || v != "gamma" {
		t.Error("expected C to be retained")
	}
}

func TestApplyUnique_RemovedCategory(t *testing.T) {
	res := ApplyUnique(makeUniqueResult(), []UniqueMode{UniqueModeRemoved})

	if len(res.Removed) != 0 {
		t.Errorf("expected all removed entries deduplicated, got %d", len(res.Removed))
	}
	// Added should be untouched
	if len(res.Added) != 3 {
		t.Errorf("expected added to be unchanged, got %d", len(res.Added))
	}
}

func TestApplyUnique_SameCategory(t *testing.T) {
	res := ApplyUnique(makeUniqueResult(), []UniqueMode{UniqueModeSame})

	if _, ok := res.Same["S1"]; ok {
		t.Error("S1 should be removed: 'shared' appears once — wait, it is unique")
	}
	// Both S1 and S2 have unique values, so both should be retained
	if len(res.Same) != 2 {
		t.Errorf("expected 2 same entries, got %d", len(res.Same))
	}
}

func TestApplyUnique_ChangedByValue(t *testing.T) {
	res := ApplyUnique(makeUniqueResult(), nil)

	// K1 old="old1" appears in K2 as well, so K1 and K2 should be removed
	if _, ok := res.Changed["K1"]; ok {
		t.Error("expected K1 removed due to duplicate old value")
	}
	if _, ok := res.Changed["K2"]; ok {
		t.Error("expected K2 removed due to duplicate old value")
	}
	if _, ok := res.Changed["K3"]; !ok {
		t.Error("expected K3 retained (unique values)")
	}
}

func TestParseUniqueModes_Valid(t *testing.T) {
	modes, err := ParseUniqueModes("added,removed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(modes) != 2 {
		t.Errorf("expected 2 modes, got %d", len(modes))
	}
}

func TestParseUniqueModes_Empty(t *testing.T) {
	modes, err := ParseUniqueModes("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if modes != nil {
		t.Error("expected nil modes for empty spec")
	}
}

func TestParseUniqueModes_Invalid(t *testing.T) {
	_, err := ParseUniqueModes("added,bogus")
	if err == nil {
		t.Error("expected error for unknown mode")
	}
}
