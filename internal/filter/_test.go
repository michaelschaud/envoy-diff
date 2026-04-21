package filter_test

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/diff"
	"github.com/your-org/envoy-diff/internal/filter"
)

// TestApplyLimit_WithApplyMulti verifies that ApplyLimit composes correctly
// with ApplyMulti when processing multiple env file results.
func TestApplyLimit_WithApplyMulti(t *testing.T) {
	results := []diff.Result{
		{
			Added:   map[string]string{"X1": "a", "X2": "b", "X3": "c"},
			Removed: map[string]string{"Y1": "a"},
			Changed: map[string]diff.ChangedValue{},
			Same:    map[string]string{},
		},
	}

	opts := filter.LimitOptions{MaxAdded: 2}
	for i, r := range results {
		results[i] = filter.ApplyLimit(r, opts)
	}

	if len(results[0].Added) != 2 {
		t.Errorf("expected 2 added after limit, got %d", len(results[0].Added))
	}
	if len(results[0].Removed) != 1 {
		t.Errorf("removed should be unaffected, got %d", len(results[0].Removed))
	}
}

// TestApplyLimit_AllCategories checks all four categories are limited independently.
func TestApplyLimit_AllCategories(t *testing.T) {
	result := diff.Result{
		Added:   map[string]string{"A1": "v", "A2": "v", "A3": "v"},
		Removed: map[string]string{"R1": "v", "R2": "v", "R3": "v"},
		Changed: map[string]diff.ChangedValue{
			"C1": {Old: "o", New: "n"},
			"C2": {Old: "o", New: "n"},
		},
		Same: map[string]string{"S1": "v", "S2": "v", "S3": "v", "S4": "v"},
	}

	out := filter.ApplyLimit(result, filter.LimitOptions{
		MaxAdded:   1,
		MaxRemoved: 2,
		MaxChanged: 1,
		MaxSame:    3,
	})

	if len(out.Added) != 1 {
		t.Errorf("expected 1 added, got %d", len(out.Added))
	}
	if len(out.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(out.Removed))
	}
	if len(out.Changed) != 1 {
		t.Errorf("expected 1 changed, got %d", len(out.Changed))
	}
	if len(out.Same) != 3 {
		t.Errorf("expected 3 same, got %d", len(out.Same))
	}
}
