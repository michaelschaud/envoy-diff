package filter

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/diff"
)

func TestMatchesGlob_NoPatterns(t *testing.T) {
	if MatchesGlob("SECRET_KEY", nil) {
		t.Error("expected false for empty pattern list")
	}
}

func TestMatchesGlob_ExactPattern(t *testing.T) {
	if !MatchesGlob("SECRET_KEY", []string{"SECRET_KEY"}) {
		t.Error("expected exact pattern to match")
	}
}

func TestMatchesGlob_WildcardSuffix(t *testing.T) {
	patterns := []string{"AWS_*"}
	keys := []string{"AWS_ACCESS_KEY", "AWS_SECRET", "AWS_REGION"}
	for _, k := range keys {
		if !MatchesGlob(k, patterns) {
			t.Errorf("expected %q to match pattern AWS_*", k)
		}
	}
}

func TestMatchesGlob_NoMatch(t *testing.T) {
	if MatchesGlob("DATABASE_URL", []string{"AWS_*"}) {
		t.Error("expected DATABASE_URL not to match AWS_*")
	}
}

func TestMatchesGlob_CaseInsensitive(t *testing.T) {
	if !MatchesGlob("aws_secret", []string{"AWS_*"}) {
		t.Error("expected case-insensitive match")
	}
}

func TestApplyGlobExclude_NoPatterns(t *testing.T) {
	r := DiffResult{
		Added:   map[string]string{"FOO": "bar"},
		Removed: map[string]string{"BAZ": "qux"},
	}
	out := ApplyGlobExclude(r, nil)
	if len(out.Added) != 1 || len(out.Removed) != 1 {
		t.Error("expected result unchanged when no patterns provided")
	}
}

func TestApplyGlobExclude_RemovesMatchingKeys(t *testing.T) {
	r := DiffResult{
		Added:     map[string]string{"AWS_KEY": "abc", "PORT": "8080"},
		Removed:   map[string]string{"AWS_SECRET": "old"},
		Changed:   map[string][2]string{"AWS_REGION": {"us-east-1", "eu-west-1"}, "HOST": {"a", "b"}},
		Unchanged: map[string]string{"AWS_ACCOUNT": "123", "DEBUG": "true"},
	}
	out := ApplyGlobExclude(r, []string{"AWS_*"})

	if _, ok := out.Added["AWS_KEY"]; ok {
		t.Error("AWS_KEY should have been excluded from Added")
	}
	if _, ok := out.Added["PORT"]; !ok {
		t.Error("PORT should remain in Added")
	}
	if _, ok := out.Removed["AWS_SECRET"]; ok {
		t.Error("AWS_SECRET should have been excluded from Removed")
	}
	if _, ok := out.Changed["AWS_REGION"]; ok {
		t.Error("AWS_REGION should have been excluded from Changed")
	}
	if _, ok := out.Changed["HOST"]; !ok {
		t.Error("HOST should remain in Changed")
	}
	if _, ok := out.Unchanged["AWS_ACCOUNT"]; ok {
		t.Error("AWS_ACCOUNT should have been excluded from Unchanged")
	}
}

// DiffResult alias for convenience in tests — mirrors the type used across the filter package.
type DiffResult = diff.Result
