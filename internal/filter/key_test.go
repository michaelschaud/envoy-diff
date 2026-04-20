package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeKeyResult() diff.Result {
	return diff.Result{
		Added: map[string]string{
			"AWS_ACCESS_KEY":  "abc123",
			"APP_PORT":        "8080",
			"DATABASE_SECRET": "hunter2",
		},
		Removed: map[string]string{
			"AWS_SECRET_KEY": "old",
			"LOG_LEVEL":      "debug",
		},
		Changed: map[string][2]string{
			"AWS_REGION": {"us-east-1", "eu-west-1"},
			"APP_ENV":    {"staging", "production"},
		},
		Same: map[string]string{
			"TZ": "UTC",
		},
	}
}

func TestApplyKeyFilter_NoSubstrings(t *testing.T) {
	r := makeKeyResult()
	out := ApplyKeyFilter(r, nil)
	if len(out.Added) != 3 || len(out.Removed) != 2 || len(out.Changed) != 2 {
		t.Error("expected result to be unchanged when no substrings provided")
	}
}

func TestApplyKeyFilter_FiltersAdded(t *testing.T) {
	r := makeKeyResult()
	out := ApplyKeyFilter(r, []string{"AWS"})
	if _, ok := out.Added["AWS_ACCESS_KEY"]; ok {
		t.Error("expected AWS_ACCESS_KEY to be excluded from Added")
	}
	if _, ok := out.Added["APP_PORT"]; !ok {
		t.Error("expected APP_PORT to remain in Added")
	}
}

func TestApplyKeyFilter_FiltersRemoved(t *testing.T) {
	r := makeKeyResult()
	out := ApplyKeyFilter(r, []string{"secret"})
	if _, ok := out.Removed["AWS_SECRET_KEY"]; ok {
		t.Error("expected AWS_SECRET_KEY to be excluded from Removed (case-insensitive)")
	}
	if _, ok := out.Added["DATABASE_SECRET"]; ok {
		t.Error("expected DATABASE_SECRET to be excluded from Added")
	}
	if _, ok := out.Removed["LOG_LEVEL"]; !ok {
		t.Error("expected LOG_LEVEL to remain in Removed")
	}
}

func TestApplyKeyFilter_FiltersChanged(t *testing.T) {
	r := makeKeyResult()
	out := ApplyKeyFilter(r, []string{"AWS"})
	if _, ok := out.Changed["AWS_REGION"]; ok {
		t.Error("expected AWS_REGION to be excluded from Changed")
	}
	if _, ok := out.Changed["APP_ENV"]; !ok {
		t.Error("expected APP_ENV to remain in Changed")
	}
}

func TestApplyKeyFilter_MultipleSubstrings(t *testing.T) {
	r := makeKeyResult()
	out := ApplyKeyFilter(r, []string{"AWS", "APP"})
	if len(out.Added) != 1 {
		t.Errorf("expected 1 added entry, got %d", len(out.Added))
	}
	if _, ok := out.Added["DATABASE_SECRET"]; !ok {
		t.Error("expected DATABASE_SECRET to remain")
	}
}

func TestKeyMatchesAny_CaseInsensitive(t *testing.T) {
	if !KeyMatchesAny("AWS_SECRET_KEY", []string{"secret"}) {
		t.Error("expected case-insensitive match for 'secret' in 'AWS_SECRET_KEY'")
	}
	if KeyMatchesAny("LOG_LEVEL", []string{"secret"}) {
		t.Error("expected no match for 'secret' in 'LOG_LEVEL'")
	}
}
