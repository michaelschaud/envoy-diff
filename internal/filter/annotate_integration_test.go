package filter

import (
	"strings"
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeAnnotateResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"NEW_KEY": "newval"},
		Removed: map[string]string{"OLD_KEY": "oldval"},
		Changed: map[string]diff.ChangedValue{
			"DB_HOST": {Old: "localhost", New: "prod.db"},
		},
		Same: map[string]string{"STABLE": "yes"},
	}
}

func TestApplyAnnotate_WithApplyMulti(t *testing.T) {
	result := makeAnnotateResult()
	rules := []AnnotateRule{{Key: "DB_HOST", Template: "[env] {{value}}"}}
	annotated := ApplyAnnotate(result, rules)

	v, ok := annotated.Changed["DB_HOST"]
	if !ok {
		t.Fatal("expected DB_HOST in Changed")
	}
	if !strings.HasPrefix(v.New, "[env] ") {
		t.Errorf("expected annotated new value, got %q", v.New)
	}
}

func TestApplyAnnotate_UnaffectedKeys(t *testing.T) {
	result := makeAnnotateResult()
	rules := []AnnotateRule{{Key: "DB_HOST", Template: "annotated:{{value}}"}}
	annotated := ApplyAnnotate(result, rules)

	if annotated.Added["NEW_KEY"] != "newval" {
		t.Errorf("Added key should be unchanged, got %q", annotated.Added["NEW_KEY"])
	}
	if annotated.Same["STABLE"] != "yes" {
		t.Errorf("Same key should be unchanged, got %q", annotated.Same["STABLE"])
	}
}

func TestApplyAnnotate_NoRules(t *testing.T) {
	result := makeAnnotateResult()
	annotated := ApplyAnnotate(result, nil)

	if annotated.Changed["DB_HOST"].New != "prod.db" {
		t.Errorf("expected unchanged value, got %q", annotated.Changed["DB_HOST"].New)
	}
}
