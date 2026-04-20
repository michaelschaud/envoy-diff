package filter

import "testing"

func TestMatchesPrefix_EmptyPrefix(t *testing.T) {
	if !matchesPrefix("ANY_KEY", "") {
		t.Error("empty prefix should match any key")
	}
}

func TestMatchesPrefix_ExactMatch(t *testing.T) {
	if !matchesPrefix("DB_HOST", "DB_") {
		t.Error("expected DB_HOST to match prefix DB_")
	}
}

func TestMatchesPrefix_NoMatch(t *testing.T) {
	if matchesPrefix("APP_HOST", "DB_") {
		t.Error("expected APP_HOST not to match prefix DB_")
	}
}

func TestMatchesPrefix_CaseInsensitive(t *testing.T) {
	if !matchesPrefix("db_host", "DB_") {
		t.Error("expected case-insensitive match for db_host with DB_")
	}
}
