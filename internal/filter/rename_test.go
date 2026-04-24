package filter_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
	"github.com/yourorg/envoy-diff/internal/filter"
)

func makeRenameResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"APP_HOST": "localhost", "DB_PORT": "5432"},
		Removed: map[string]string{"OLD_SECRET": "abc"},
		Same:    map[string]string{"LOG_LEVEL": "info"},
		Changed: map[string]diff.ChangedValue{
			"APP_TIMEOUT": {Old: "30", New: "60"},
		},
	}
}

func TestApplyRename_NoRules(t *testing.T) {
	r := makeRenameResult()
	out := filter.ApplyRename(r, nil)
	if _, ok := out.Added["APP_HOST"]; !ok {
		t.Error("expected APP_HOST to remain unchanged")
	}
}

func TestApplyRename_ExactKey(t *testing.T) {
	r := makeRenameResult()
	rules := []filter.RenameRule{{From: "LOG_LEVEL", To: "LOG_SEVERITY"}}
	out := filter.ApplyRename(r, rules)
	if _, ok := out.Same["LOG_SEVERITY"]; !ok {
		t.Error("expected LOG_LEVEL to be renamed to LOG_SEVERITY")
	}
	if _, ok := out.Same["LOG_LEVEL"]; ok {
		t.Error("expected old key LOG_LEVEL to be absent")
	}
}

func TestApplyRename_PrefixKey(t *testing.T) {
	r := makeRenameResult()
	rules := []filter.RenameRule{{From: "APP", To: "SVC"}}
	out := filter.ApplyRename(r, rules)
	if _, ok := out.Added["SVC_HOST"]; !ok {
		t.Error("expected APP_HOST to be renamed to SVC_HOST")
	}
	if _, ok := out.Changed["SVC_TIMEOUT"]; !ok {
		t.Error("expected APP_TIMEOUT to be renamed to SVC_TIMEOUT")
	}
}

func TestApplyRename_CaseInsensitive(t *testing.T) {
	r := makeRenameResult()
	rules := []filter.RenameRule{{From: "db", To: "DATABASE"}}
	out := filter.ApplyRename(r, rules)
	if _, ok := out.Added["DATABASE_PORT"]; !ok {
		t.Error("expected DB_PORT to be renamed to DATABASE_PORT (case-insensitive)")
	}
}

func TestApplyRename_FirstRuleWins(t *testing.T) {
	r := diff.Result{
		Added:   map[string]string{"APP_KEY": "val"},
		Removed: map[string]string{},
		Same:    map[string]string{},
		Changed: map[string]diff.ChangedValue{},
	}
	rules := []filter.RenameRule{
		{From: "APP", To: "FIRST"},
		{From: "APP", To: "SECOND"},
	}
	out := filter.ApplyRename(r, rules)
	if _, ok := out.Added["FIRST_KEY"]; !ok {
		t.Error("expected first matching rule to win")
	}
}
