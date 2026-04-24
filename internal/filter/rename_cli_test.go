package filter_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/filter"
)

func TestParseRenameRules_Valid(t *testing.T) {
	rules, err := filter.ParseRenameRules([]string{"APP=SVC", "DB=DATABASE"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].From != "APP" || rules[0].To != "SVC" {
		t.Errorf("unexpected first rule: %+v", rules[0])
	}
	if rules[1].From != "DB" || rules[1].To != "DATABASE" {
		t.Errorf("unexpected second rule: %+v", rules[1])
	}
}

func TestParseRenameRules_Empty(t *testing.T) {
	rules, err := filter.ParseRenameRules(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected empty rules, got %d", len(rules))
	}
}

func TestParseRenameRules_MissingEquals(t *testing.T) {
	_, err := filter.ParseRenameRules([]string{"APPONLY"})
	if err == nil {
		t.Error("expected error for missing '=' separator")
	}
}

func TestParseRenameRules_EmptyFrom(t *testing.T) {
	_, err := filter.ParseRenameRules([]string{"=SVC"})
	if err == nil {
		t.Error("expected error for empty FROM")
	}
}

func TestParseRenameRules_EmptyTo(t *testing.T) {
	_, err := filter.ParseRenameRules([]string{"APP="})
	if err == nil {
		t.Error("expected error for empty TO")
	}
}
