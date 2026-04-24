package filter

import (
	"testing"
)

func TestParseAnnotateRules_Valid(t *testing.T) {
	rules, err := ParseAnnotateRules([]string{"DB_HOST=host:{{value}}", "API_KEY=secret:{{value}}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", rules[0].Key)
	}
	if rules[0].Template != "host:{{value}}" {
		t.Errorf("unexpected template: %q", rules[0].Template)
	}
}

func TestParseAnnotateRules_Empty(t *testing.T) {
	rules, err := ParseAnnotateRules([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected 0 rules, got %d", len(rules))
	}
}

func TestParseAnnotateRules_MissingEquals(t *testing.T) {
	_, err := ParseAnnotateRules([]string{"DB_HOSThost"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseAnnotateRules_EmptyKey(t *testing.T) {
	_, err := ParseAnnotateRules([]string{"=template"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestParseAnnotateRules_EmptyTemplate(t *testing.T) {
	_, err := ParseAnnotateRules([]string{"DB_HOST="})
	if err == nil {
		t.Fatal("expected error for empty template")
	}
}
