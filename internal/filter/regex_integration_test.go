package filter_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
	"github.com/yourorg/envoy-diff/internal/filter"
)

// TestApplyRegexExclude_WithApplyMulti verifies that regex exclusion composes
// correctly with ApplyMulti prefix filtering.
func TestApplyRegexExclude_WithApplyMulti(t *testing.T) {
	base := diff.Result{
		Added: map[string]string{
			"DB_PASSWORD": "newpass",
			"DB_HOST":     "db.prod",
			"APP_SECRET":  "topsecret",
			"APP_PORT":    "9090",
		},
		Removed: map[string]string{},
		Changed: map[string][2]string{},
		Same:    map[string]string{},
	}

	// First apply prefix filter to keep only DB_ and APP_ keys (all in this case)
	prefixed := filter.ApplyMulti([]diff.Result{base}, []string{"DB_", "APP_"}, false)

	// Then exclude sensitive patterns
	filtered := filter.ApplyRegexExclude(prefixed, []string{".*password.*", ".*secret.*"})

	if _, ok := filtered.Added["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been excluded by regex")
	}
	if _, ok := filtered.Added["APP_SECRET"]; ok {
		t.Error("APP_SECRET should have been excluded by regex")
	}
	if _, ok := filtered.Added["DB_HOST"]; !ok {
		t.Error("DB_HOST should be retained")
	}
	if _, ok := filtered.Added["APP_PORT"]; !ok {
		t.Error("APP_PORT should be retained")
	}
}

// TestApplyRegexExclude_SameKeysUnaffectedByDefault checks that Same keys are
// also filtered when they match sensitive patterns.
func TestApplyRegexExclude_SameKeysFiltered(t *testing.T) {
	result := diff.Result{
		Added:   map[string]string{},
		Removed: map[string]string{},
		Changed: map[string][2]string{},
		Same: map[string]string{
			"DB_PASSWORD": "same",
			"APP_ENV":     "production",
		},
	}
	out := filter.ApplyRegexExclude(result, []string{".*password.*"})
	if _, ok := out.Same["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD in Same should be excluded")
	}
	if _, ok := out.Same["APP_ENV"]; !ok {
		t.Error("APP_ENV should remain in Same")
	}
}
