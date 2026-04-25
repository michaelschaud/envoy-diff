package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeNormalizeResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"A": "  hello  ", "B": `"world"`},
		Removed: map[string]string{"C": "  'quoted'  "},
		Same:    map[string]string{"D": "  plain  "},
		Changed: map[string][2]string{
			"E": {"  old  ", `"new"`},
		},
	}
}

func TestApplyNormalize_NoMode(t *testing.T) {
	r := makeNormalizeResult()
	out := ApplyNormalize(r, "")
	if out.Added["A"] != "  hello  " {
		t.Errorf("expected unchanged value, got %q", out.Added["A"])
	}
}

func TestApplyNormalize_TrimMode(t *testing.T) {
	r := makeNormalizeResult()
	out := ApplyNormalize(r, NormalizeTrimSpace)
	if out.Added["A"] != "hello" {
		t.Errorf("expected trimmed value, got %q", out.Added["A"])
	}
	if out.Same["D"] != "plain" {
		t.Errorf("expected trimmed same value, got %q", out.Same["D"])
	}
	if out.Changed["E"][0] != "old" {
		t.Errorf("expected trimmed old changed value, got %q", out.Changed["E"][0])
	}
}

func TestApplyNormalize_UnquoteMode(t *testing.T) {
	r := makeNormalizeResult()
	out := ApplyNormalize(r, NormalizeTrimQuotes)
	if out.Added["B"] != "world" {
		t.Errorf("expected unquoted value, got %q", out.Added["B"])
	}
	if out.Changed["E"][1] != "new" {
		t.Errorf("expected unquoted new changed value, got %q", out.Changed["E"][1])
	}
}

func TestApplyNormalize_BothMode(t *testing.T) {
	r := makeNormalizeResult()
	out := ApplyNormalize(r, NormalizeBoth)
	if out.Removed["C"] != "quoted" {
		t.Errorf("expected trim+unquote value, got %q", out.Removed["C"])
	}
	if out.Added["A"] != "hello" {
		t.Errorf("expected trimmed value, got %q", out.Added["A"])
	}
}

func TestApplyNormalize_UnknownMode(t *testing.T) {
	r := makeNormalizeResult()
	out := ApplyNormalize(r, NormalizeMode("bogus"))
	if out.Added["A"] != "  hello  " {
		t.Errorf("expected unchanged value for unknown mode, got %q", out.Added["A"])
	}
}
