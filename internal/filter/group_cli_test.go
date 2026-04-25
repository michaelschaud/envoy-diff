package filter

import (
	"strings"
	"testing"
)

func TestParseGroupConfig_EmptyDelimiter(t *testing.T) {
	cfg, err := ParseGroupConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Delimiter != "" {
		t.Errorf("expected empty delimiter, got %q", cfg.Delimiter)
	}
}

func TestParseGroupConfig_Underscore(t *testing.T) {
	cfg, err := ParseGroupConfig("_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Delimiter != "_" {
		t.Errorf("expected '_', got %q", cfg.Delimiter)
	}
}

func TestParseGroupConfig_Dot(t *testing.T) {
	cfg, err := ParseGroupConfig(".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Delimiter != "." {
		t.Errorf("expected '.', got %q", cfg.Delimiter)
	}
}

func TestParseGroupConfig_TooLong(t *testing.T) {
	_, err := ParseGroupConfig(strings.Repeat("x", 9))
	if err == nil {
		t.Fatal("expected error for delimiter longer than 8 chars")
	}
}

func TestParseGroupConfig_MaxLength(t *testing.T) {
	delim := strings.Repeat("x", 8)
	cfg, err := ParseGroupConfig(delim)
	if err != nil {
		t.Fatalf("unexpected error for 8-char delimiter: %v", err)
	}
	if cfg.Delimiter != delim {
		t.Errorf("expected %q, got %q", delim, cfg.Delimiter)
	}
}
