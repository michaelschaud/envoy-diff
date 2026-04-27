package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func ptr(f float64) *float64 { return &f }

func makeClampResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"PORT": "9000", "TIMEOUT": "abc"},
		Removed: map[string]string{"OLD_PORT": "1"},
		Same:    map[string]string{"WORKERS": "50"},
		Changed: map[string][2]string{
			"MAX_CONN": {"200", "5000"},
		},
	}
}

func TestApplyClamp_NoBounds(t *testing.T) {
	r := makeClampResult()
	out := ApplyClamp(r, ClampConfig{})
	if out.Added["PORT"] != "9000" {
		t.Errorf("expected PORT=9000, got %s", out.Added["PORT"])
	}
}

func TestApplyClamp_MinOnly(t *testing.T) {
	r := makeClampResult()
	out := ApplyClamp(r, ClampConfig{Min: ptr(10)})
	if out.Removed["OLD_PORT"] != "10" {
		t.Errorf("expected OLD_PORT clamped to 10, got %s", out.Removed["OLD_PORT"])
	}
	if out.Added["PORT"] != "9000" {
		t.Errorf("expected PORT unchanged at 9000, got %s", out.Added["PORT"])
	}
}

func TestApplyClamp_MaxOnly(t *testing.T) {
	r := makeClampResult()
	out := ApplyClamp(r, ClampConfig{Max: ptr(100)})
	if out.Added["PORT"] != "100" {
		t.Errorf("expected PORT clamped to 100, got %s", out.Added["PORT"])
	}
	if out.Same["WORKERS"] != "50" {
		t.Errorf("expected WORKERS unchanged at 50, got %s", out.Same["WORKERS"])
	}
}

func TestApplyClamp_MinAndMax(t *testing.T) {
	r := makeClampResult()
	out := ApplyClamp(r, ClampConfig{Min: ptr(10), Max: ptr(100)})
	if out.Added["PORT"] != "100" {
		t.Errorf("expected PORT=100, got %s", out.Added["PORT"])
	}
	if out.Removed["OLD_PORT"] != "10" {
		t.Errorf("expected OLD_PORT=10, got %s", out.Removed["OLD_PORT"])
	}
	pair := out.Changed["MAX_CONN"]
	if pair[0] != "100" || pair[1] != "100" {
		t.Errorf("expected MAX_CONN both clamped to 100, got %v", pair)
	}
}

func TestApplyClamp_NonNumericUnchanged(t *testing.T) {
	r := makeClampResult()
	out := ApplyClamp(r, ClampConfig{Min: ptr(0), Max: ptr(100)})
	if out.Added["TIMEOUT"] != "abc" {
		t.Errorf("expected non-numeric TIMEOUT unchanged, got %s", out.Added["TIMEOUT"])
	}
}
