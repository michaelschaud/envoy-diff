package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeMergeResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"NEW_KEY": "new"},
		Removed: map[string]string{"OLD_KEY": "old"},
		Changed: map[string][2]string{"CHANGED": {"before", "after"}},
		Same:    map[string]string{"STABLE": "value"},
	}
}

func TestApplyMerge_NoOverlay(t *testing.T) {
	result := makeMergeResult()
	out := ApplyMerge(result, MergeConfig{})
	if len(out.Added) != 1 || out.Added["NEW_KEY"] != "new" {
		t.Errorf("expected Added unchanged, got %v", out.Added)
	}
}

func TestApplyMerge_RightStrategy(t *testing.T) {
	result := makeMergeResult()
	out := ApplyMerge(result, MergeConfig{
		Overlay:  map[string]string{"STABLE": "overridden"},
		Strategy: MergeStrategyRight,
	})
	if out.Same["STABLE"] != "overridden" {
		t.Errorf("expected STABLE overridden, got %q", out.Same["STABLE"])
	}
}

func TestApplyMerge_LeftStrategy(t *testing.T) {
	result := makeMergeResult()
	out := ApplyMerge(result, MergeConfig{
		Overlay:  map[string]string{"STABLE": "overridden"},
		Strategy: MergeStrategyLeft,
	})
	if out.Same["STABLE"] != "value" {
		t.Errorf("expected STABLE preserved, got %q", out.Same["STABLE"])
	}
}

func TestApplyMerge_UnionStrategy(t *testing.T) {
	result := makeMergeResult()
	out := ApplyMerge(result, MergeConfig{
		Overlay:  map[string]string{"STABLE": "overridden", "EXTRA": "added"},
		Strategy: MergeStrategyUnion,
	})
	if out.Same["STABLE"] != "overridden" {
		t.Errorf("expected STABLE overridden, got %q", out.Same["STABLE"])
	}
	if out.Same["EXTRA"] != "added" {
		t.Errorf("expected EXTRA injected, got %q", out.Same["EXTRA"])
	}
}

func TestApplyMerge_ChangedKeyOverlay(t *testing.T) {
	result := makeMergeResult()
	out := ApplyMerge(result, MergeConfig{
		Overlay:  map[string]string{"CHANGED": "patched"},
		Strategy: MergeStrategyRight,
	})
	if out.Changed["CHANGED"][1] != "patched" {
		t.Errorf("expected CHANGED right value patched, got %q", out.Changed["CHANGED"][1])
	}
}
