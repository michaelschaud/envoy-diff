package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// TestApplyValueFilter_WithApplyMulti verifies that value filtering composes
// correctly with ApplyMulti prefix filtering.
func TestApplyValueFilter_WithApplyMulti(t *testing.T) {
	results := map[string]diff.Result{
		"staging": {
			Added: map[string]string{
				"APP_SECRET": "REDACTED",
				"APP_PORT":   "8080",
			},
			Removed: map[string]string{},
			Changed: map[string][2]string{},
			Same:    map[string]string{},
		},
	}

	// First apply multi-prefix filter to keep only APP_ keys
	filtered := ApplyMulti(results, []string{"APP_"}, false)

	// Then apply value filter to remove REDACTED values
	for env, r := range filtered {
		filtered[env] = ApplyValueFilter(r, []string{"REDACTED"})
	}

	r := filtered["staging"]
	if _, ok := r.Added["APP_SECRET"]; ok {
		t.Error("expected APP_SECRET to be removed by value filter")
	}
	if _, ok := r.Added["APP_PORT"]; !ok {
		t.Error("expected APP_PORT to remain after value filter")
	}
}

// TestApplyValueFilter_EmptyResult ensures no panic on empty maps.
func TestApplyValueFilter_EmptyResult(t *testing.T) {
	r := diff.Result{
		Added:   map[string]string{},
		Removed: map[string]string{},
		Changed: map[string][2]string{},
		Same:    map[string]string{},
	}
	out := ApplyValueFilter(r, []string{"anything"})
	if len(out.Added) != 0 || len(out.Removed) != 0 || len(out.Changed) != 0 || len(out.Same) != 0 {
		t.Error("expected all empty maps on empty input")
	}
}
