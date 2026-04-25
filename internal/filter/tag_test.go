package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeTagResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"DB_HOST": "localhost", "APP_PORT": "8080"},
		Removed: map[string]string{"OLD_KEY": "value"},
		Same:    map[string]string{"LOG_LEVEL": "info"},
		Changed: map[string][2]string{"DB_PASS": {"old", "new"}},
	}
}

func TestApplyTag_NoRules(t *testing.T) {
	r := makeTagResult()
	out := ApplyTag(r, nil)
	if out.Added["DB_HOST"] != "localhost" {
		t.Errorf("expected unchanged value, got %q", out.Added["DB_HOST"])
	}
}

func TestApplyTag_TagsAdded(t *testing.T) {
	r := makeTagResult()
	rules := []TagRule{{Tag: "database", Substr: "db"}}
	out := ApplyTag(r, rules)
	if out.Added["DB_HOST"] != "localhost [database]" {
		t.Errorf("expected tagged value, got %q", out.Added["DB_HOST"])
	}
	if out.Added["APP_PORT"] != "8080" {
		t.Errorf("expected unchanged value, got %q", out.Added["APP_PORT"])
	}
}

func TestApplyTag_TagsRemoved(t *testing.T) {
	r := makeTagResult()
	rules := []TagRule{{Tag: "legacy", Substr: "old"}}
	out := ApplyTag(r, rules)
	if out.Removed["OLD_KEY"] != "value [legacy]" {
		t.Errorf("expected tagged removed value, got %q", out.Removed["OLD_KEY"])
	}
}

func TestApplyTag_TagsChanged(t *testing.T) {
	r := makeTagResult()
	rules := []TagRule{{Tag: "secret", Substr: "pass"}}
	out := ApplyTag(r, rules)
	pair := out.Changed["DB_PASS"]
	if pair[0] != "old [secret]" || pair[1] != "new [secret]" {
		t.Errorf("expected tagged changed pair, got %v", pair)
	}
}

func TestApplyTag_CaseInsensitive(t *testing.T) {
	r := diff.Result{
		Added:   map[string]string{"api_token": "abc"},
		Removed: map[string]string{},
		Same:    map[string]string{},
		Changed: map[string][2]string{},
	}
	rules := []TagRule{{Tag: "auth", Substr: "TOKEN"}}
	out := ApplyTag(r, rules)
	if out.Added["api_token"] != "abc [auth]" {
		t.Errorf("expected case-insensitive match, got %q", out.Added["api_token"])
	}
}
