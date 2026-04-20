package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeValueResult() diff.Result {
	return diff.Result{
		Added: map[string]string{
			"NEW_KEY":      "somevalue",
			"SECRET_KEY":   "REDACTED",
		},
		Removed: map[string]string{
			"OLD_KEY":      "oldvalue",
			"HIDDEN_KEY":   "todo_placeholder",
		},
		Changed: map[string][2]string{
			"DB_URL":       {"postgres://old", "postgres://new"},
			"API_TOKEN":    {"REDACTED", "newtoken"},
		},
		Same: map[string]string{
			"STABLE_KEY":  "stablevalue",
			"MASKED":      "redacted_value",
		},
	}
}

func TestApplyValueFilter_NoSubstrings(t *testing.T) {
	r := makeValueResult()
	out := ApplyValueFilter(r, nil)
	if len(out.Added) != len(r.Added) {
		t.Errorf("expected Added unchanged, got %d", len(out.Added))
	}
}

func TestApplyValueFilter_FiltersAdded(t *testing.T) {
	r := makeValueResult()
	out := ApplyValueFilter(r, []string{"REDACTED"})
	if _, ok := out.Added["SECRET_KEY"]; ok {
		t.Error("expected SECRET_KEY to be filtered from Added")
	}
	if _, ok := out.Added["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY to remain in Added")
	}
}

func TestApplyValueFilter_FiltersRemoved(t *testing.T) {
	r := makeValueResult()
	out := ApplyValueFilter(r, []string{"todo"})
	if _, ok := out.Removed["HIDDEN_KEY"]; ok {
		t.Error("expected HIDDEN_KEY to be filtered from Removed")
	}
	if _, ok := out.Removed["OLD_KEY"]; !ok {
		t.Error("expected OLD_KEY to remain in Removed")
	}
}

func TestApplyValueFilter_FiltersChangedOnOldValue(t *testing.T) {
	r := makeValueResult()
	out := ApplyValueFilter(r, []string{"REDACTED"})
	if _, ok := out.Changed["API_TOKEN"]; ok {
		t.Error("expected API_TOKEN to be filtered from Changed (old value matches)")
	}
	if _, ok := out.Changed["DB_URL"]; !ok {
		t.Error("expected DB_URL to remain in Changed")
	}
}

func TestApplyValueFilter_FiltersSame(t *testing.T) {
	r := makeValueResult()
	out := ApplyValueFilter(r, []string{"redacted"})
	if _, ok := out.Same["MASKED"]; ok {
		t.Error("expected MASKED to be filtered from Same")
	}
	if _, ok := out.Same["STABLE_KEY"]; !ok {
		t.Error("expected STABLE_KEY to remain in Same")
	}
}

func TestApplyValueFilter_CaseInsensitive(t *testing.T) {
	r := makeValueResult()
	out := ApplyValueFilter(r, []string{"REDACTED"})
	if _, ok := out.Same["MASKED"]; ok {
		t.Error("expected case-insensitive match to filter MASKED")
	}
}
