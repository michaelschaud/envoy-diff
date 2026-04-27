package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeThresholdResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"PORT": "8080", "TIMEOUT": "hello", "WORKERS": "2"},
		Removed: map[string]string{"OLD_PORT": "3000", "OLD_LABEL": "prod"},
		Changed: map[string][2]string{
			"MAX_CONN": {"50", "200"},
			"HOST":     {"localhost", "example.com"},
		},
		Same: map[string]string{"REGION": "us-east-1", "RETRY": "3"},
	}
}

func TestApplyThreshold_NoConfig(t *testing.T) {
	r := makeThresholdResult()
	out := ApplyThreshold(r, ThresholdConfig{})
	if len(out.Added) != len(r.Added) {
		t.Errorf("expected %d added, got %d", len(r.Added), len(out.Added))
	}
}

func TestApplyThreshold_MinOnly(t *testing.T) {
	r := makeThresholdResult()
	min := 5000.0
	out := ApplyThreshold(r, ThresholdConfig{Min: &min})
	// PORT=8080 kept, WORKERS=2 removed, TIMEOUT=hello kept (non-numeric)
	if _, ok := out.Added["PORT"]; !ok {
		t.Error("expected PORT to be kept")
	}
	if _, ok := out.Added["WORKERS"]; ok {
		t.Error("expected WORKERS to be removed (2 < 5000)")
	}
	if _, ok := out.Added["TIMEOUT"]; !ok {
		t.Error("expected non-numeric TIMEOUT to be kept")
	}
}

func TestApplyThreshold_MaxOnly(t *testing.T) {
	r := makeThresholdResult()
	max := 100.0
	out := ApplyThreshold(r, ThresholdConfig{Max: &max})
	if _, ok := out.Added["PORT"]; ok {
		t.Error("expected PORT to be removed (8080 > 100)")
	}
	if _, ok := out.Added["WORKERS"]; !ok {
		t.Error("expected WORKERS to be kept (2 <= 100)")
	}
}

func TestApplyThreshold_ChangedKeepIfEitherInRange(t *testing.T) {
	r := makeThresholdResult()
	min := 100.0
	out := ApplyThreshold(r, ThresholdConfig{Min: &min})
	// MAX_CONN: old=50 (out), new=200 (in) → kept
	if _, ok := out.Changed["MAX_CONN"]; !ok {
		t.Error("expected MAX_CONN to be kept because new value 200 >= 100")
	}
	// HOST is non-numeric → kept
	if _, ok := out.Changed["HOST"]; !ok {
		t.Error("expected non-numeric HOST to be kept")
	}
}

func TestParseThresholdConfig_MinMax(t *testing.T) {
	cfg, err := ParseThresholdConfig("min:10,max:500")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Min == nil || *cfg.Min != 10 {
		t.Errorf("expected Min=10, got %v", cfg.Min)
	}
	if cfg.Max == nil || *cfg.Max != 500 {
		t.Errorf("expected Max=500, got %v", cfg.Max)
	}
}

func TestParseThresholdConfig_Empty(t *testing.T) {
	cfg, err := ParseThresholdConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Min != nil || cfg.Max != nil {
		t.Error("expected nil Min and Max for empty spec")
	}
}

func TestParseThresholdConfig_InvalidNumber(t *testing.T) {
	_, err := ParseThresholdConfig("min:abc")
	if err == nil {
		t.Error("expected error for non-numeric value")
	}
}

func TestParseThresholdConfig_UnknownKey(t *testing.T) {
	_, err := ParseThresholdConfig("avg:50")
	if err == nil {
		t.Error("expected error for unknown key")
	}
}
