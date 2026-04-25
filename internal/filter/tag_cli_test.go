package filter

import (
	"testing"
)

func TestParseTagRules_Valid(t *testing.T) {
	rules, err := ParseTagRules([]string{"database=db", "secret=pass"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Tag != "database" || rules[0].Substr != "db" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
	if rules[1].Tag != "secret" || rules[1].Substr != "pass" {
		t.Errorf("unexpected rule[1]: %+v", rules[1])
	}
}

func TestParseTagRules_Empty(t *testing.T) {
	rules, err := ParseTagRules(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected empty rules, got %d", len(rules))
	}
}

func TestParseTagRules_MissingEquals(t *testing.T) {
	_, err := ParseTagRules([]string{"nodivider"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseTagRules_EmptyTag(t *testing.T) {
	_, err := ParseTagRules([]string{"=somesubstr"})
	if err == nil {
		t.Fatal("expected error for empty tag")
	}
}

func TestParseTagRules_EmptySubstr(t *testing.T) {
	_, err := ParseTagRules([]string{"mytag="})
	if err == nil {
		t.Fatal("expected error for empty substr")
	}
}
