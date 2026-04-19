package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envoy-diff/internal/reporter"
)

func TestTextReportWriter_ContainsMetadata(t *testing.T) {
	r := reporter.New("staging.env", "production.env", makeResult())
	var buf bytes.Buffer
	w := reporter.TextReportWriter{}
	if err := w.Write(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "staging.env") {
		t.Error("expected source file in output")
	}
	if !strings.Contains(out, "production.env") {
		t.Error("expected target file in output")
	}
}

func TestTextReportWriter_ContainsSections(t *testing.T) {
	r := reporter.New("a", "b", makeResult())
	var buf bytes.Buffer
	w := reporter.TextReportWriter{}
	_ = w.Write(&buf, r)
	out := buf.String()
	for _, section := range []string{"[ADDED]", "[REMOVED]", "[CHANGED]"} {
		if !strings.Contains(out, section) {
			t.Errorf("expected section %q in output", section)
		}
	}
}

func TestTextReportWriter_NoDiffOutput(t *testing.T) {
	import_diff "github.com/envoy-diff/internal/diff"
	result := import_diff.Result{Same: map[string]string{"PORT": "8080"}}
	r := reporter.New("a", "b", result)
	var buf bytes.Buffer
	w := reporter.TextReportWriter{}
	_ = w.Write(&buf, r)
	out := buf.String()
	if strings.Contains(out, "[ADDED]") || strings.Contains(out, "[REMOVED]") || strings.Contains(out, "[CHANGED]") {
		t.Error("expected no diff sections for identical envs")
	}
}
