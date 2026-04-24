package filter

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeTransformResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"NEW_KEY": "HelloWorld"},
		Removed: map[string]string{"OLD_KEY": "ByeBye"},
		Changed: map[string][2]string{"HOST": {"staging.example.com", "prod.example.com"}},
		Same:    map[string]string{"REGION": "us-east-1"},
	}
}

func TestApplyTransform_NoMode(t *testing.T) {
	r := makeTransformResult()
	out := ApplyTransform(r, TransformNone)
	if out.Added["NEW_KEY"] != "HelloWorld" {
		t.Errorf("expected unchanged value, got %q", out.Added["NEW_KEY"])
	}
}

func TestApplyTransform_Upper(t *testing.T) {
	r := makeTransformResult()
	out := ApplyTransform(r, TransformUpper)

	if out.Added["NEW_KEY"] != "HELLOWORLD" {
		t.Errorf("Added: expected HELLOWORLD, got %q", out.Added["NEW_KEY"])
	}
	if out.Removed["OLD_KEY"] != "BYEBYE" {
		t.Errorf("Removed: expected BYEBYE, got %q", out.Removed["OLD_KEY"])
	}
	if out.Changed["HOST"][0] != "STAGING.EXAMPLE.COM" {
		t.Errorf("Changed old: expected STAGING.EXAMPLE.COM, got %q", out.Changed["HOST"][0])
	}
	if out.Changed["HOST"][1] != "PROD.EXAMPLE.COM" {
		t.Errorf("Changed new: expected PROD.EXAMPLE.COM, got %q", out.Changed["HOST"][1])
	}
	if out.Same["REGION"] != "US-EAST-1" {
		t.Errorf("Same: expected US-EAST-1, got %q", out.Same["REGION"])
	}
}

func TestApplyTransform_Lower(t *testing.T) {
	r := makeTransformResult()
	out := ApplyTransform(r, TransformLower)

	if out.Added["NEW_KEY"] != "helloworld" {
		t.Errorf("Added: expected helloworld, got %q", out.Added["NEW_KEY"])
	}
	if out.Changed["HOST"][0] != "staging.example.com" {
		t.Errorf("Changed old: expected staging.example.com, got %q", out.Changed["HOST"][0])
	}
}

func TestApplyTransform_UnknownMode(t *testing.T) {
	r := makeTransformResult()
	out := ApplyTransform(r, TransformMode("base64"))
	if out.Added["NEW_KEY"] != "HelloWorld" {
		t.Errorf("expected unchanged value for unknown mode, got %q", out.Added["NEW_KEY"])
	}
}

func TestApplyTransform_EmptyResult(t *testing.T) {
	r := diff.Result{}
	out := ApplyTransform(r, TransformUpper)
	if len(out.Added) != 0 || len(out.Removed) != 0 || len(out.Changed) != 0 || len(out.Same) != 0 {
		t.Error("expected all empty maps for empty result")
	}
}
