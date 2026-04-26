package reporter_test

import (
	"testing"

	"github.com/envoy-diff/internal/diff"
	"github.com/envoy-diff/internal/reporter"
)

func makeResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"NEW_KEY": "value"},
		Removed: map[string]string{"OLD_KEY": "value"},
		Changed: map[string][2]string{"HOST": {"localhost", "prod.example.com"}},
		Same:    map[string]string{"PORT": "8080"},
	}
}

func TestNew_PopulatesStats(t *testing.T) {
	r := reporter.New("staging.env", "production.env", makeResult())
	if r.Stats.Added != 1 {
		t.Errorf("expected Added=1, got %d", r.Stats.Added)
	}
	if r.Stats.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", r.Stats.Removed)
	}
	if r.Stats.Changed != 1 {
		t.Errorf("expected Changed=1, got %d", r.Stats.Changed)
	}
	if r.Stats.Same != 1 {
		t.Errorf("expected Same=1, got %d", r.Stats.Same)
	}
}

func TestNew_PopulatesMetadata(t *testing.T) {
	r := reporter.New("staging.env", "production.env", makeResult())
	if r.SourceFile != "staging.env" {
		t.Errorf("unexpected SourceFile: %s", r.SourceFile)
	}
	if r.TargetFile != "production.env" {
		t.Errorf("unexpected TargetFile: %s", r.TargetFile)
	}
	if r.GeneratedAt.IsZero() {
		t.Error("GeneratedAt should not be zero")
	}
}

func TestHasDiff_True(t *testing.T) {
	r := reporter.New("a", "b", makeResult())
	if !r.HasDiff() {
		t.Error("expected HasDiff to be true")
	}
}

func TestHasDiff_False(t *testing.T) {
	result := diff.Result{Same: map[string]string{"PORT": "8080"}}
	r := reporter.New("a", "b", result)
	if r.HasDiff() {
		t.Error("expected HasDiff to be false")
	}
}

func TestHasDiff_EmptyResult(t *testing.T) {
	r := reporter.New("a", "b", diff.Result{})
	if r.HasDiff() {
		t.Error("expected HasDiff to be false for empty result")
	}
	if r.Stats.Added != 0 || r.Stats.Removed != 0 || r.Stats.Changed != 0 || r.Stats.Same != 0 {
		t.Errorf("expected all stats to be 0 for empty result, got %+v", r.Stats)
	}
}
