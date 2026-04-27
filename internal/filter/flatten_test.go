package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeFlattenResult() diff.Result {
	return diff.Result{
		Added: map[string]string{
			"APP__HOST":    "localhost",
			"APP__DB__PORT": "5432",
		},
		Removed: map[string]string{
			"APP__SECRET": "old",
		},
		Same: map[string]string{
			"PLAIN_KEY": "value",
		},
		Changed: map[string][2]string{
			"APP__LOG__LEVEL": {"info", "debug"},
		},
	}
}

func TestApplyFlatten_EmptyDelimiter(t *testing.T) {
	r := makeFlattenResult()
	out := ApplyFlatten(r, FlattenConfig{Delimiter: ""})
	if _, ok := out.Added["APP__HOST"]; !ok {
		t.Error("expected original key APP__HOST to be preserved when delimiter is empty")
	}
}

func TestApplyFlatten_UnlimitedDepth(t *testing.T) {
	r := makeFlattenResult()
	out := ApplyFlatten(r, FlattenConfig{Delimiter: "__", Depth: 0})

	if _, ok := out.Added["HOST"]; !ok {
		t.Error("expected APP__HOST -> HOST with unlimited depth")
	}
	if _, ok := out.Added["PORT"]; !ok {
		t.Error("expected APP__DB__PORT -> PORT with unlimited depth")
	}
	if _, ok := out.Removed["SECRET"]; !ok {
		t.Error("expected APP__SECRET -> SECRET with unlimited depth")
	}
	if _, ok := out.Same["PLAIN_KEY"]; !ok {
		t.Error("expected PLAIN_KEY to be preserved (no delimiter present)")
	}
	if _, ok := out.Changed["LEVEL"]; !ok {
		t.Error("expected APP__LOG__LEVEL -> LEVEL with unlimited depth")
	}
}

func TestApplyFlatten_DepthOne(t *testing.T) {
	r := makeFlattenResult()
	out := ApplyFlatten(r, FlattenConfig{Delimiter: "__", Depth: 1})

	if _, ok := out.Added["HOST"]; !ok {
		t.Error("expected APP__HOST -> HOST with depth 1")
	}
	// APP__DB__PORT with depth 1 should become DB__PORT
	if _, ok := out.Added["DB__PORT"]; !ok {
		t.Error("expected APP__DB__PORT -> DB__PORT with depth 1")
	}
	if _, ok := out.Changed["LOG__LEVEL"]; !ok {
		t.Error("expected APP__LOG__LEVEL -> LOG__LEVEL with depth 1")
	}
}

func TestApplyFlatten_PlainKeysUnchanged(t *testing.T) {
	r := diff.Result{
		Added: map[string]string{"NODELIM": "val"},
	}
	out := ApplyFlatten(r, FlattenConfig{Delimiter: "__", Depth: 0})
	if v, ok := out.Added["NODELIM"]; !ok || v != "val" {
		t.Error("expected key without delimiter to pass through unchanged")
	}
}

func TestApplyFlatten_EmptyResult(t *testing.T) {
	r := diff.Result{}
	out := ApplyFlatten(r, FlattenConfig{Delimiter: "__", Depth: 0})
	if len(out.Added) != 0 || len(out.Removed) != 0 || len(out.Changed) != 0 {
		t.Error("expected empty result to remain empty after flatten")
	}
}
