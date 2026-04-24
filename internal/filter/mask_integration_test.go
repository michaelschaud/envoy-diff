package filter

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/diff"
)

// TestApplyMask_WithApplyMulti verifies that masking composes correctly with
// ApplyMulti so that secrets are redacted even when multiple env files are
// compared in one pass.
func TestApplyMask_WithApplyMulti(t *testing.T) {
	sets := []map[string]string{
		{"DB_PASSWORD": "prod-pass", "APP_ENV": "production"},
		{"DB_PASSWORD": "stage-pass", "APP_ENV": "staging"},
		{"DB_PASSWORD": "dev-pass", "APP_ENV": "development"},
	}

	results := ApplyMulti(sets)

	masked := make([]diff.Result, len(results))
	for i, r := range results {
		masked[i] = ApplyMask(r, []string{"password"})
	}

	for i, r := range masked {
		for k, v := range r.Same {
			if k == "DB_PASSWORD" && v != maskMarker {
				t.Errorf("result[%d]: expected DB_PASSWORD redacted in Same, got %q", i, v)
			}
		}
		for k, v := range r.Added {
			if k == "DB_PASSWORD" && v != maskMarker {
				t.Errorf("result[%d]: expected DB_PASSWORD redacted in Added, got %q", i, v)
			}
		}
	}
}

// TestApplyMask_MultipleSubstrings checks that any matching substring triggers
// redaction.
func TestApplyMask_MultipleSubstrings(t *testing.T) {
	r := diff.Result{
		Added: map[string]string{
			"DB_PASSWORD": "pass",
			"API_KEY":     "key123",
			"LOG_LEVEL":   "info",
		},
		Removed: map[string]string{},
		Same:    map[string]string{},
		Changed: map[string][2]string{},
	}

	out := ApplyMask(r, []string{"password", "key"})

	if out.Added["DB_PASSWORD"] != maskMarker {
		t.Errorf("expected DB_PASSWORD redacted, got %q", out.Added["DB_PASSWORD"])
	}
	if out.Added["API_KEY"] != maskMarker {
		t.Errorf("expected API_KEY redacted, got %q", out.Added["API_KEY"])
	}
	if out.Added["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL unchanged, got %q", out.Added["LOG_LEVEL"])
	}
}
