package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeCapResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"A": "short", "B": "this-is-a-very-long-value"},
		Removed: map[string]string{"C": "x"},
		Same:    map[string]string{"D": "same-value-here-that-is-long"},
		Changed: map[string][2]string{
			"E": {"old-long-value-exceeds", "new"},
			"F": {"ok", "also-a-very-long-new-value"},
		},
	}
}

func TestApplyCap_NoMaxLen(t *testing.T) {
	result := makeCapResult()
	out := ApplyCap(result, CapConfig{})
	if out.Added["B"] != "this-is-a-very-long-value" {
		t.Errorf("expected unchanged value, got %q", out.Added["B"])
	}
}

func TestApplyCap_CapsLongAdded(t *testing.T) {
	out := ApplyCap(makeCapResult(), CapConfig{MaxLen: 8, Replacement: "[CAPPED]"})
	if out.Added["A"] != "short" {
		t.Errorf("short value should be unchanged, got %q", out.Added["A"])
	}
	if out.Added["B"] != "[CAPPED]" {
		t.Errorf("long value should be capped, got %q", out.Added["B"])
	}
}

func TestApplyCap_DefaultReplacement(t *testing.T) {
	out := ApplyCap(makeCapResult(), CapConfig{MaxLen: 3})
	if out.Added["B"] != "[CAPPED]" {
		t.Errorf("expected default replacement, got %q", out.Added["B"])
	}
}

func TestApplyCap_CapsChangedValues(t *testing.T) {
	out := ApplyCap(makeCapResult(), CapConfig{MaxLen: 5, Replacement: "***"})
	pair := out.Changed["E"]
	if pair[0] != "***" {
		t.Errorf("expected old value capped, got %q", pair[0])
	}
	if pair[1] != "new" {
		t.Errorf("expected new value unchanged, got %q", pair[1])
	}
	pairF := out.Changed["F"]
	if pairF[0] != "ok" {
		t.Errorf("expected old value unchanged, got %q", pairF[0])
	}
	if pairF[1] != "***" {
		t.Errorf("expected new value capped, got %q", pairF[1])
	}
}

func TestParseCapConfig_ValidNoReplacement(t *testing.T) {
	cfg, err := ParseCapConfig("20")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxLen != 20 {
		t.Errorf("expected MaxLen=20, got %d", cfg.MaxLen)
	}
}

func TestParseCapConfig_ValidWithReplacement(t *testing.T) {
	cfg, err := ParseCapConfig("10:REDACTED")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxLen != 10 || cfg.Replacement != "REDACTED" {
		t.Errorf("unexpected config: %+v", cfg)
	}
}

func TestParseCapConfig_InvalidMaxLen(t *testing.T) {
	_, err := ParseCapConfig("abc")
	if err == nil {
		t.Error("expected error for non-numeric maxlen")
	}
}

func TestParseCapConfig_ZeroMaxLen(t *testing.T) {
	_, err := ParseCapConfig("0")
	if err == nil {
		t.Error("expected error for zero maxlen")
	}
}

func TestParseCapConfig_Empty(t *testing.T) {
	cfg, err := ParseCapConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxLen != 0 {
		t.Errorf("expected zero config, got %+v", cfg)
	}
}
