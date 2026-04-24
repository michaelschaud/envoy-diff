package filter

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/diff"
)

func makeMaskResult() diff.Result {
	return diff.Result{
		Added: map[string]string{
			"DB_PASSWORD": "s3cr3t",
			"APP_PORT":    "8080",
		},
		Removed: map[string]string{
			"API_SECRET": "old-secret",
			"LOG_LEVEL":  "debug",
		},
		Same: map[string]string{
			"REGION": "us-east-1",
		},
		Changed: map[string][2]string{
			"AUTH_TOKEN": {"token-old", "token-new"},
			"REPLICAS":   {"2", "4"},
		},
	}
}

func TestApplyMask_NoSubstrings(t *testing.T) {
	r := makeMaskResult()
	out := ApplyMask(r, nil)
	if out.Added["DB_PASSWORD"] != "s3cr3t" {
		t.Errorf("expected original value, got %q", out.Added["DB_PASSWORD"])
	}
}

func TestApplyMask_RedactsAdded(t *testing.T) {
	r := makeMaskResult()
	out := ApplyMask(r, []string{"password"})
	if out.Added["DB_PASSWORD"] != maskMarker {
		t.Errorf("expected redacted value, got %q", out.Added["DB_PASSWORD"])
	}
	if out.Added["APP_PORT"] != "8080" {
		t.Errorf("expected unchanged value, got %q", out.Added["APP_PORT"])
	}
}

func TestApplyMask_RedactsRemoved(t *testing.T) {
	r := makeMaskResult()
	out := ApplyMask(r, []string{"secret"})
	if out.Removed["API_SECRET"] != maskMarker {
		t.Errorf("expected redacted value, got %q", out.Removed["API_SECRET"])
	}
	if out.Removed["LOG_LEVEL"] != "debug" {
		t.Errorf("expected unchanged value, got %q", out.Removed["LOG_LEVEL"])
	}
}

func TestApplyMask_RedactsChanged(t *testing.T) {
	r := makeMaskResult()
	out := ApplyMask(r, []string{"token"})
	pair := out.Changed["AUTH_TOKEN"]
	if pair[0] != maskMarker || pair[1] != maskMarker {
		t.Errorf("expected both values redacted, got %v", pair)
	}
	if out.Changed["REPLICAS"][0] != "2" {
		t.Errorf("expected unchanged value, got %q", out.Changed["REPLICAS"][0])
	}
}

func TestApplyMask_CaseInsensitive(t *testing.T) {
	r := makeMaskResult()
	out := ApplyMask(r, []string{"PASSWORD"})
	if out.Added["DB_PASSWORD"] != maskMarker {
		t.Errorf("expected redacted value for case-insensitive match, got %q", out.Added["DB_PASSWORD"])
	}
}

func TestApplyMask_DoesNotMutateOriginal(t *testing.T) {
	r := makeMaskResult()
	ApplyMask(r, []string{"password", "secret", "token"})
	if r.Added["DB_PASSWORD"] != "s3cr3t" {
		t.Error("original result was mutated")
	}
}
