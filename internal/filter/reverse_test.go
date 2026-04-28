package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeReverseResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"NEW_KEY": "new_val"},
		Removed: map[string]string{"OLD_KEY": "old_val"},
		Changed: map[string]diff.Change{
			"CHANGED_KEY": {Old: "before", New: "after"},
		},
		Same: map[string]string{"STABLE": "stable_val"},
	}
}

func TestApplyReverse_SwapsAddedAndRemoved(t *testing.T) {
	result := makeReverseResult()
	reversed := ApplyReverse(result)

	if _, ok := reversed.Added["OLD_KEY"]; !ok {
		t.Error("expected OLD_KEY to appear in Added after reverse")
	}
	if _, ok := reversed.Removed["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY to appear in Removed after reverse")
	}
}

func TestApplyReverse_SwapsChangedOldNew(t *testing.T) {
	result := makeReverseResult()
	reversed := ApplyReverse(result)

	c, ok := reversed.Changed["CHANGED_KEY"]
	if !ok {
		t.Fatal("expected CHANGED_KEY in Changed")
	}
	if c.Old != "after" {
		t.Errorf("expected Old=\"after\", got %q", c.Old)
	}
	if c.New != "before" {
		t.Errorf("expected New=\"before\", got %q", c.New)
	}
}

func TestApplyReverse_PreservesSame(t *testing.T) {
	result := makeReverseResult()
	reversed := ApplyReverse(result)

	if v, ok := reversed.Same["STABLE"]; !ok || v != "stable_val" {
		t.Errorf("expected Same to be preserved, got %v", reversed.Same)
	}
}

func TestApplyReverse_EmptyResult(t *testing.T) {
	empty := diff.Result{}
	reversed := ApplyReverse(empty)

	if len(reversed.Added) != 0 || len(reversed.Removed) != 0 || len(reversed.Changed) != 0 {
		t.Error("expected empty result to remain empty after reverse")
	}
}

func TestApplyReverse_DoubleReverseIsIdentity(t *testing.T) {
	result := makeReverseResult()
	double := ApplyReverse(ApplyReverse(result))

	if double.Added["NEW_KEY"] != result.Added["NEW_KEY"] {
		t.Error("double reverse should restore Added")
	}
	if double.Removed["OLD_KEY"] != result.Removed["OLD_KEY"] {
		t.Error("double reverse should restore Removed")
	}
	c := double.Changed["CHANGED_KEY"]
	if c.Old != "before" || c.New != "after" {
		t.Errorf("double reverse should restore Changed values, got Old=%q New=%q", c.Old, c.New)
	}
}
