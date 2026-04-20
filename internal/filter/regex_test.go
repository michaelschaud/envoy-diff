package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func TestMatchesRegex_NoPatterns(t *testing.T) {
	if MatchesRegex("MY_KEY", nil) {
		t.Error("expected false for nil patterns")
	}
}

func TestMatchesRegex_ExactPattern(t *testing.T) {
	if !MatchesRegex("MY_KEY", []string{"^my_key$"}) {
		t.Error("expected case-insensitive exact match")
	}
}

func TestMatchesRegex_WildcardPattern(t *testing.T) {
	if !MatchesRegex("DB_PASSWORD", []string{".*password.*"}) {
		t.Error("expected wildcard pattern to match")
	}
}

func TestMatchesRegex_NoMatch(t *testing.T) {
	if MatchesRegex("APP_PORT", []string{"^db_.*"}) {
		t.Error("expected no match")
	}
}

func TestMatchesRegex_InvalidPatternSkipped(t *testing.T) {
	// invalid regex should be skipped, not panic
	if MatchesRegex("KEY", []string{"[invalid"}) {
		t.Error("expected false for invalid pattern")
	}
}

func TestApplyRegexExclude_NoPatterns(t *testing.T) {
	result := diff.Result{
		Added:   map[string]string{"NEW_KEY": "val"},
		Removed: map[string]string{"OLD_KEY": "val"},
		Changed: map[string][2]string{"DB_HOST": {"a", "b"}},
		Same:    map[string]string{"APP_ENV": "prod"},
	}
	out := ApplyRegexExclude(result, nil)
	if len(out.Added) != 1 || len(out.Removed) != 1 || len(out.Changed) != 1 {
		t.Error("expected result unchanged when no patterns provided")
	}
}

func TestApplyRegexExclude_RemovesMatchingKeys(t *testing.T) {
	result := diff.Result{
		Added:   map[string]string{"DB_PASSWORD": "secret", "APP_PORT": "8080"},
		Removed: map[string]string{"DB_SECRET": "old"},
		Changed: map[string][2]string{"API_TOKEN": {"x", "y"}, "LOG_LEVEL": {"info", "debug"}},
		Same:    map[string]string{"DB_HOST": "localhost"},
	}
	out := ApplyRegexExclude(result, []string{".*password.*", ".*secret.*", ".*token.*"})

	if _, ok := out.Added["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be excluded")
	}
	if _, ok := out.Added["APP_PORT"]; !ok {
		t.Error("APP_PORT should be retained")
	}
	if _, ok := out.Removed["DB_SECRET"]; ok {
		t.Error("DB_SECRET should be excluded")
	}
	if _, ok := out.Changed["API_TOKEN"]; ok {
		t.Error("API_TOKEN should be excluded")
	}
	if _, ok := out.Changed["LOG_LEVEL"]; !ok {
		t.Error("LOG_LEVEL should be retained")
	}
}
