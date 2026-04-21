package filter_test

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/diff"
	"github.com/your-org/envoy-diff/internal/filter"
)

func makeLimitResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"A1": "v1", "A2": "v2", "A3": "v3"},
		Removed: map[string]string{"R1": "v1", "R2": "v2"},
		Changed: map[string]diff.ChangedValue{
			"C1": {Old: "old1", New: "new1"},
			"C2": {Old: "old2", New: "new2"},
			"C3": {Old: "old3", New: "new3"},
		},
		Same: map[string]string{"S1": "v1", "S2": "v2"},
	}
}

func TestApplyLimit_NoLimits(t *testing.T) {
	result := makeLimitResult()
	out := filter.ApplyLimit(result, filter.LimitOptions{})
	if len(out.Added) != 3 {
		t.Errorf("expected 3 added, got %d", len(out.Added))
	}
	if len(out.Changed) != 3 {
		t.Errorf("expected 3 changed, got %d", len(out.Changed))
	}
}

func TestApplyLimit_LimitsAdded(t *testing.T) {
	result := makeLimitResult()
	out := filter.ApplyLimit(result, filter.LimitOptions{MaxAdded: 2})
	if len(out.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(out.Added))
	}
	if len(out.Removed) != 2 {
		t.Errorf("removed should be unchanged, got %d", len(out.Removed))
	}
}

func TestApplyLimit_LimitsChanged(t *testing.T) {
	result := makeLimitResult()
	out := filter.ApplyLimit(result, filter.LimitOptions{MaxChanged: 1})
	if len(out.Changed) != 1 {
		t.Errorf("expected 1 changed, got %d", len(out.Changed))
	}
}

func TestApplyLimit_LimitLargerThanSet(t *testing.T) {
	result := makeLimitResult()
	out := filter.ApplyLimit(result, filter.LimitOptions{MaxAdded: 100, MaxRemoved: 100})
	if len(out.Added) != 3 {
		t.Errorf("expected all 3 added, got %d", len(out.Added))
	}
	if len(out.Removed) != 2 {
		t.Errorf("expected all 2 removed, got %d", len(out.Removed))
	}
}

func TestApplyLimit_EmptyResult(t *testing.T) {
	out := filter.ApplyLimit(diff.Result{}, filter.LimitOptions{MaxAdded: 5})
	if len(out.Added) != 0 {
		t.Errorf("expected empty added, got %d", len(out.Added))
	}
}
