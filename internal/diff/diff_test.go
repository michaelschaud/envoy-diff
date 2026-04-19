package diff

import (
	"testing"
)

func TestCompare_AllCategories(t *testing.T) {
	left := EnvMap{
		"SHARED_SAME":    "value1",
		"SHARED_CHANGED": "old",
		"ONLY_LEFT":      "leftval",
	}
	right := EnvMap{
		"SHARED_SAME":    "value1",
		"SHARED_CHANGED": "new",
		"ONLY_RIGHT":     "rightval",
	}

	res := Compare(left, right)

	if v, ok := res.Unchanged["SHARED_SAME"]; !ok || v != "value1" {
		t.Errorf("expected SHARED_SAME in Unchanged, got %v", v)
	}

	if pair, ok := res.Modified["SHARED_CHANGED"]; !ok || pair[0] != "old" || pair[1] != "new" {
		t.Errorf("expected SHARED_CHANGED in Modified with [old new], got %v", pair)
	}

	if v, ok := res.OnlyInLeft["ONLY_LEFT"]; !ok || v != "leftval" {
		t.Errorf("expected ONLY_LEFT in OnlyInLeft, got %v", v)
	}

	if v, ok := res.OnlyInRight["ONLY_RIGHT"]; !ok || v != "rightval" {
		t.Errorf("expected ONLY_RIGHT in OnlyInRight, got %v", v)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	res := Compare(EnvMap{}, EnvMap{})
	if len(res.Unchanged)+len(res.Modified)+len(res.OnlyInLeft)+len(res.OnlyInRight) != 0 {
		t.Error("expected all empty result for empty inputs")
	}
}

func TestCompare_IdenticalMaps(t *testing.T) {
	m := EnvMap{"FOO": "bar", "BAZ": "qux"}
	res := Compare(m, m)
	if len(res.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged, got %d", len(res.Unchanged))
	}
	if len(res.Modified)+len(res.OnlyInLeft)+len(res.OnlyInRight) != 0 {
		t.Error("expected no diffs for identical maps")
	}
}
