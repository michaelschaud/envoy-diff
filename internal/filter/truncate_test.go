package filter

import (
	"strings"
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeTruncateResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"NEW_KEY": "short", "LONG_ADD": "this-value-is-quite-long-indeed"},
		Removed: map[string]string{"OLD_KEY": "tiny", "LONG_REM": "another-very-long-value-here"},
		Same:    map[string]string{"SAME_KEY": "same"},
		Changed: map[string]diff.ChangedValue{
			"MOD_KEY": {Old: "old-but-long-value-string", New: "new-but-also-long-value"},
		},
	}
}

func TestApplyTruncate_NoMaxLen(t *testing.T) {
	r := makeTruncateResult()
	out := ApplyTruncate(r, 0, "...")
	if out.Added["LONG_ADD"] != r.Added["LONG_ADD"] {
		t.Errorf("expected value unchanged when maxLen=0")
	}
}

func TestApplyTruncate_ShortValuesUnchanged(t *testing.T) {
	r := makeTruncateResult()
	out := ApplyTruncate(r, 20, "...")
	if out.Added["NEW_KEY"] != "short" {
		t.Errorf("expected short value to be unchanged, got %q", out.Added["NEW_KEY"])
	}
	if out.Removed["OLD_KEY"] != "tiny" {
		t.Errorf("expected short value to be unchanged, got %q", out.Removed["OLD_KEY"])
	}
}

func TestApplyTruncate_LongValuesAreTruncated(t *testing.T) {
	r := makeTruncateResult()
	out := ApplyTruncate(r, 10, "...")
	v := out.Added["LONG_ADD"]
	if len(v) > 10 {
		t.Errorf("expected value length <= 10, got %d: %q", len(v), v)
	}
	if !strings.HasSuffix(v, "...") {
		t.Errorf("expected truncated value to end with '...', got %q", v)
	}
}

func TestApplyTruncate_ChangedValues(t *testing.T) {
	r := makeTruncateResult()
	out := ApplyTruncate(r, 10, "...")
	cv := out.Changed["MOD_KEY"]
	if len(cv.Old) > 10 {
		t.Errorf("Old value not truncated: %q", cv.Old)
	}
	if len(cv.New) > 10 {
		t.Errorf("New value not truncated: %q", cv.New)
	}
}

func TestApplyTruncate_DefaultSuffix(t *testing.T) {
	r := diff.Result{
		Added: map[string]string{"K": "hello-world-long"},
	}
	out := ApplyTruncate(r, 8, "")
	v := out.Added["K"]
	if !strings.HasSuffix(v, "...") {
		t.Errorf("expected default suffix '...', got %q", v)
	}
}

func TestApplyTruncate_EmptyResult(t *testing.T) {
	r := diff.Result{}
	out := ApplyTruncate(r, 10, "...")
	if len(out.Added)+len(out.Removed)+len(out.Same)+len(out.Changed) != 0 {
		t.Errorf("expected empty result, got non-empty")
	}
}
