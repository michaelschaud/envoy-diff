package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeScoreResult() diff.Result {
	return diff.Result{
		Added: map[string]string{
			"DB_HOST":     "localhost",
			"DB_PASSWORD": "secret",
			"APP_NAME":    "myapp",
		},
		Removed: map[string]string{
			"OLD_DB_URL": "postgres://old",
			"CACHE_TTL":  "300",
		},
		Changed: map[string][2]string{
			"DB_PORT":    {"5432", "5433"},
			"LOG_LEVEL":  {"info", "debug"},
			"DB_REPLICA": {"host1", "host2"},
		},
		Same: map[string]string{
			"REGION": "us-east-1",
			"DB_SSL":  "true",
		},
	}
}

func TestApplyScore_NoSubstrings(t *testing.T) {
	r := makeScoreResult()
	out := ApplyScore(r, nil)
	if len(out.Added) != len(r.Added) {
		t.Errorf("expected %d added, got %d", len(r.Added), len(out.Added))
	}
	if len(out.Changed) != len(r.Changed) {
		t.Errorf("expected %d changed, got %d", len(r.Changed), len(out.Changed))
	}
}

func TestApplyScore_PreservesAllKeys(t *testing.T) {
	r := makeScoreResult()
	out := ApplyScore(r, []string{"DB"})
	if len(out.Added) != 3 {
		t.Errorf("expected 3 added keys, got %d", len(out.Added))
	}
	if len(out.Changed) != 3 {
		t.Errorf("expected 3 changed keys, got %d", len(out.Changed))
	}
	if len(out.Removed) != 2 {
		t.Errorf("expected 2 removed keys, got %d", len(out.Removed))
	}
}

func TestComputeScore_NoMatch(t *testing.T) {
	score := ComputeScore("APP_NAME", []string{"DB", "CACHE"})
	if score != 0 {
		t.Errorf("expected score 0, got %d", score)
	}
}

func TestComputeScore_SingleMatch(t *testing.T) {
	score := ComputeScore("DB_HOST", []string{"DB", "CACHE"})
	if score != 1 {
		t.Errorf("expected score 1, got %d", score)
	}
}

func TestComputeScore_MultipleMatches(t *testing.T) {
	score := ComputeScore("DB_CACHE_KEY", []string{"DB", "CACHE"})
	if score != 2 {
		t.Errorf("expected score 2, got %d", score)
	}
}

func TestComputeScore_CaseInsensitive(t *testing.T) {
	score := ComputeScore("db_host", []string{"DB"})
	if score != 1 {
		t.Errorf("expected case-insensitive score 1, got %d", score)
	}
}

func TestApplyScore_ChangedKeysPreserved(t *testing.T) {
	r := makeScoreResult()
	out := ApplyScore(r, []string{"DB", "LOG"})
	for k, v := range r.Changed {
		got, ok := out.Changed[k]
		if !ok {
			t.Errorf("changed key %q missing from output", k)
		}
		if got != v {
			t.Errorf("changed key %q: expected %v, got %v", k, v, got)
		}
	}
}
