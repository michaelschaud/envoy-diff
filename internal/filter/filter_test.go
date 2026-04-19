package filter_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
	"github.com/yourorg/envoy-diff/internal/filter"
)

func makeResult() diff.Result {
	return diff.Result{
		Added:     map[string]string{"DB_HOST": "localhost", "APP_PORT": "8080"},
		Removed:   map[string]string{"OLD_KEY": "value"},
		Changed:   map[string][2]string{"DB_PASS": {"old", "new"}, "APP_ENV": {"dev", "prod"}},
		Unchanged: map[string]string{"LOG_LEVEL": "info", "APP_NAME": "myapp"},
	}
}

func TestApply_NoPrefixNoOnlyChanged(t *testing.T) {
	r := filter.Apply(makeResult(), filter.Options{})
	if len(r.Added) != 2 || len(r.Removed) != 1 || len(r.Changed) != 2 || len(r.Unchanged) != 2 {
		t.Errorf("expected no filtering, got added=%d removed=%d changed=%d unchanged=%d",
			len(r.Added), len(r.Removed), len(r.Changed), len(r.Unchanged))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	r := filter.Apply(makeResult(), filter.Options{Prefixes: []string{"DB_"}})
	if _, ok := r.Added["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in added")
	}
	if _, ok := r.Added["APP_PORT"]; ok {
		t.Error("expected APP_PORT to be filtered out")
	}
	if _, ok := r.Changed["DB_PASS"]; !ok {
		t.Error("expected DB_PASS in changed")
	}
	if _, ok := r.Changed["APP_ENV"]; ok {
		t.Error("expected APP_ENV to be filtered out")
	}
}

func TestApply_OnlyChanged(t *testing.T) {
	r := filter.Apply(makeResult(), filter.Options{OnlyChanged: true})
	if len(r.Unchanged) != 0 {
		t.Errorf("expected unchanged to be empty, got %d", len(r.Unchanged))
	}
	if len(r.Added) != 2 {
		t.Errorf("expected added to remain, got %d", len(r.Added))
	}
}

func TestApply_PrefixAndOnlyChanged(t *testing.T) {
	r := filter.Apply(makeResult(), filter.Options{Prefixes: []string{"APP_"}, OnlyChanged: true})
	if _, ok := r.Added["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in added")
	}
	if len(r.Unchanged) != 0 {
		t.Error("expected unchanged empty")
	}
	if _, ok := r.Changed["DB_PASS"]; ok {
		t.Error("expected DB_PASS filtered out")
	}
}
