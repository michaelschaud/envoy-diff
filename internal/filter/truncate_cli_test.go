package filter

import (
	"testing"
)

func TestParseTruncateConfig_ValidNoSuffix(t *testing.T) {
	cfg, err := ParseTruncateConfig("20")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxLen != 20 {
		t.Errorf("expected MaxLen=20, got %d", cfg.MaxLen)
	}
	if cfg.Suffix != "..." {
		t.Errorf("expected default suffix '...', got %q", cfg.Suffix)
	}
}

func TestParseTruncateConfig_ValidWithSuffix(t *testing.T) {
	cfg, err := ParseTruncateConfig("15:***")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxLen != 15 {
		t.Errorf("expected MaxLen=15, got %d", cfg.MaxLen)
	}
	if cfg.Suffix != "***" {
		t.Errorf("expected suffix '***', got %q", cfg.Suffix)
	}
}

func TestParseTruncateConfig_EmptySpec(t *testing.T) {
	_, err := ParseTruncateConfig("")
	if err == nil {
		t.Error("expected error for empty spec")
	}
}

func TestParseTruncateConfig_InvalidMaxLen(t *testing.T) {
	_, err := ParseTruncateConfig("abc")
	if err == nil {
		t.Error("expected error for non-numeric maxLen")
	}
}

func TestParseTruncateConfig_ZeroMaxLen(t *testing.T) {
	_, err := ParseTruncateConfig("0")
	if err == nil {
		t.Error("expected error for maxLen=0")
	}
}

func TestParseTruncateConfig_NegativeMaxLen(t *testing.T) {
	_, err := ParseTruncateConfig("-5")
	if err == nil {
		t.Error("expected error for negative maxLen")
	}
}

func TestParseTruncateConfig_EmptySuffixAllowed(t *testing.T) {
	cfg, err := ParseTruncateConfig("10:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Suffix != "" {
		t.Errorf("expected empty suffix, got %q", cfg.Suffix)
	}
}
