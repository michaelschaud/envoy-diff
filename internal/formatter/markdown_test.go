package formatter

import (
	"strings"
	"testing"

	"github.com/user/envoy-diff/internal/diff"
)

func TestMarkdownWriter_ContainsSections(t *testing.T) {
	result := makeResult()
	var buf strings.Builder
	w := NewMarkdownWriter(&buf)
	if err := w.Write(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	expected := []string{
		"# Environment Diff Report",
		"## Added",
		"## Removed",
		"## Changed",
		"| Key | Staging | Production |",
	}
	for _, s := range expected {
		if !strings.Contains(out, s) {
			t.Errorf("expected output to contain %q, got:\n%s", s, out)
		}
	}
}

func TestMarkdownWriter_AddedRow(t *testing.T) {
	result := diff.Result{
		Added:   map[string]string{"NEW_KEY": "newval"},
		Removed: map[string]string{},
		Changed: map[string]diff.ValuePair{},
	}
	var buf strings.Builder
	w := NewMarkdownWriter(&buf)
	if err := w.Write(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "`NEW_KEY`") {
		t.Errorf("expected NEW_KEY in output, got:\n%s", out)
	}
	if !strings.Contains(out, "`newval`") {
		t.Errorf("expected newval in output, got:\n%s", out)
	}
}

func TestMarkdownWriter_EmptyResult(t *testing.T) {
	result := diff.Result{
		Added:   map[string]string{},
		Removed: map[string]string{},
		Changed: map[string]diff.ValuePair{},
	}
	var buf strings.Builder
	w := NewMarkdownWriter(&buf)
	if err := w.Write(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "##") {
		t.Errorf("expected no sections for empty result, got:\n%s", out)
	}
}
