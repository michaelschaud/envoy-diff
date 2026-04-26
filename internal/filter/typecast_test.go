package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeTypecastResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"PORT": "8080", "DEBUG": "true", "RATIO": "0.5", "NAME": "app"},
		Removed: map[string]string{"TIMEOUT": "30"},
		Changed: map[string][2]string{"WORKERS": {"2", "4"}},
		Same:    map[string]string{"HOST": "localhost"},
	}
}

func TestApplyTypecast_NoModes(t *testing.T) {
	r := makeTypecastResult()
	out := ApplyTypecast(r, nil)
	if out.Added["PORT"] != "8080" {
		t.Errorf("expected unchanged PORT, got %q", out.Added["PORT"])
	}
}

func TestApplyTypecast_IntMode(t *testing.T) {
	r := makeTypecastResult()
	out := ApplyTypecast(r, []TypecastMode{TypecastInt})
	if out.Added["PORT"] != "[int] 8080" {
		t.Errorf("expected [int] prefix for PORT, got %q", out.Added["PORT"])
	}
	if out.Added["NAME"] != "app" {
		t.Errorf("expected NAME unchanged, got %q", out.Added["NAME"])
	}
}

func TestApplyTypecast_BoolMode(t *testing.T) {
	r := makeTypecastResult()
	out := ApplyTypecast(r, []TypecastMode{TypecastBool})
	if out.Added["DEBUG"] != "[bool] true" {
		t.Errorf("expected [bool] prefix for DEBUG, got %q", out.Added["DEBUG"])
	}
	if out.Added["PORT"] != "8080" {
		t.Errorf("expected PORT unchanged, got %q", out.Added["PORT"])
	}
}

func TestApplyTypecast_FloatMode(t *testing.T) {
	r := makeTypecastResult()
	out := ApplyTypecast(r, []TypecastMode{TypecastFloat})
	if out.Added["RATIO"] != "[float] 0.5" {
		t.Errorf("expected [float] prefix for RATIO, got %q", out.Added["RATIO"])
	}
}

func TestApplyTypecast_ChangedValues(t *testing.T) {
	r := makeTypecastResult()
	out := ApplyTypecast(r, []TypecastMode{TypecastInt})
	pair := out.Changed["WORKERS"]
	if pair[0] != "[int] 2" || pair[1] != "[int] 4" {
		t.Errorf("expected [int] prefix on both changed values, got %v", pair)
	}
}

func TestApplyTypecast_MultiModes(t *testing.T) {
	r := makeTypecastResult()
	out := ApplyTypecast(r, []TypecastMode{TypecastInt, TypecastBool})
	if out.Added["PORT"] != "[int] 8080" {
		t.Errorf("expected [int] prefix for PORT, got %q", out.Added["PORT"])
	}
	if out.Added["DEBUG"] != "[bool] true" {
		t.Errorf("expected [bool] prefix for DEBUG, got %q", out.Added["DEBUG"])
	}
	if out.Added["RATIO"] != "0.5" {
		t.Errorf("expected RATIO unchanged, got %q", out.Added["RATIO"])
	}
}
