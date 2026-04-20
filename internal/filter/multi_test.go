package filter_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
	"github.com/yourorg/envoy-diff/internal/filter"
)

func makeMultiResult() diff.Result {
	return diff.Result{
		Added: map[string]string{
			"APP_NEW":    "val1",
			"DB_HOST":    "localhost",
			"SECRET_KEY": "abc",
		},
		Removed: map[string]string{
			"APP_OLD": "val2",
			"DB_PORT": "5432",
		},
		Changed: map[string][2]string{
			"APP_ENV": {"staging", "production"},
			"DB_NAME": {"mydb_stg", "mydb_prod"},
		},
		Unchanged: map[string]string{
			"APP_VERSION": "1.0",
		},
	}
}

func TestApplyMulti_NoPrefixes(t *testing.T) {
	r := makeMultiResult()
	out := filter.ApplyMulti(r, []string{}, false)
	if len(out.Added) != 3 || len(out.Removed) != 2 || len(out.Changed) != 2 {
		t.Errorf("expected full result, got added=%d removed=%d changed=%d", len(out.Added), len(out.Removed), len(out.Changed))
	}
}

func TestApplyMulti_MultiplePrefixes(t *testing.T) {
	r := makeMultiResult()
	out := filter.ApplyMulti(r, []string{"APP_", "DB_"}, false)
	if _, ok := out.Added["SECRET_KEY"]; ok {
		t.Error("SECRET_KEY should be filtered out")
	}
	if _, ok := out.Added["APP_NEW"]; !ok {
		t.Error("APP_NEW should be present")
	}
	if _, ok := out.Added["DB_HOST"]; !ok {
		t.Error("DB_HOST should be present")
	}
}

func TestApplyMulti_OnlyChanged(t *testing.T) {
	r := makeMultiResult()
	out := filter.ApplyMulti(r, []string{}, true)
	if len(out.Added) != 0 || len(out.Removed) != 0 {
		t.Errorf("expected no added/removed with onlyChanged, got added=%d removed=%d", len(out.Added), len(out.Removed))
	}
	if len(out.Changed) != 2 {
		t.Errorf("expected 2 changed, got %d", len(out.Changed))
	}
}

func TestApplyMulti_PrefixAndOnlyChanged(t *testing.T) {
	r := makeMultiResult()
	out := filter.ApplyMulti(r, []string{"APP_"}, true)
	if len(out.Changed) != 1 {
		t.Errorf("expected 1 changed APP_ key, got %d", len(out.Changed))
	}
	if _, ok := out.Changed["APP_ENV"]; !ok {
		t.Error("APP_ENV should be in changed")
	}
}
