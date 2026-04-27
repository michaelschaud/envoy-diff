package filter

import (
	"testing"
)

func TestParseCollapseConfig_EmptySpec(t *testing.T) {
	_, err := ParseCollapseConfig("")
	if err == nil {
		t.Fatal("expected error for empty spec")
	}
}

func TestParseCollapseConfig_DelimiterOnly(t *testing.T) {
	cfg, err := ParseCollapseConfig(",")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Delimiter != "," {
		t.Errorf("expected delimiter ',', got %q", cfg.Delimiter)
	}
	if len(cfg.Keys) != 0 {
		t.Errorf("expected no keys, got %v", cfg.Keys)
	}
}

func TestParseCollapseConfig_WithKeys(t *testing.T) {
	cfg, err := ParseCollapseConfig(":PATH,CLASSPATH")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Delimiter != ":" {
		t.Errorf("expected delimiter ':', got %q", cfg.Delimiter)
	}
	if len(cfg.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(cfg.Keys))
	}
	if cfg.Keys[0] != "PATH" || cfg.Keys[1] != "CLASSPATH" {
		t.Errorf("unexpected keys: %v", cfg.Keys)
	}
}

func TestParseCollapseConfig_EmptyDelimiterInSpec(t *testing.T) {
	_, err := ParseCollapseConfig(":")
	if err == nil {
		t.Fatal("expected error for empty delimiter part")
	}
}

func TestParseCollapseConfig_KeysWithSpaces(t *testing.T) {
	cfg, err := ParseCollapseConfig(", : KEY_A , KEY_B ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// delimiter is everything before the first colon
	if cfg.Delimiter != ", " {
		t.Errorf("unexpected delimiter: %q", cfg.Delimiter)
	}
	if len(cfg.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(cfg.Keys))
	}
	if cfg.Keys[0] != "KEY_A" || cfg.Keys[1] != "KEY_B" {
		t.Errorf("unexpected keys: %v", cfg.Keys)
	}
}
