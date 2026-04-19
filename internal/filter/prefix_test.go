package filter

import "testing"

func TestMatchesPrefix_EmptyPrefix(t *testing.T) {
	if !matchesPrefix("", "ANY_KEY") {
		t.Error("empty prefix should match any key")
	}
}

func TestMatchesPrefix_ExactMatch(t *testing.T) {
	if !matchesPrefix("DB_", "DB_HOST") {
		t.Error("expected match")
	}
}

func TestMatchesPrefix_NoMatch(t *testing.T) {
	if matchesPrefix("DB_", "APP_PORT") {
		t.Error("expected no match")
	}
}

func TestMatchesPrefix_CaseInsensitive(t *testing.T) {
	if !matchesPrefix("db_", "DB_HOST") {
		t.Error("expected case-insensitive match")
	}
}
