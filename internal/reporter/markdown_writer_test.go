package reporter

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/envoy-diff/internal/diff"
)

func TestMarkdownReportWriter_ContainsMetadata(t *testing.T) {
	r := &Report{
		Metadata: Metadata{
			StagingFile:    "staging.env",
			ProductionFile: "production.env",
			GeneratedAt:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Result: diff.Result{},
	}
	var buf bytes.Buffer
	w := NewMarkdownReportWriter()
	if err := w.Write(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "staging.env") {
		t.Error("expected staging file in output")
	}
	if !strings.Contains(out, "production.env") {
		t.Error("expected production file in output")
	}
	if !strings.Contains(out, "2024-01-01") {
		t.Error("expected generated date in output")
	}
}

func TestMarkdownReportWriter_ContainsSections(t *testing.T) {
	r := makeResult()
	report := New(r, "s.env", "p.env")
	var buf bytes.Buffer
	w := NewMarkdownReportWriter()
	if err := w.Write(&buf, report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, section := range []string{"## Added", "## Removed", "## Changed", "## Stats"} {
		if !strings.Contains(out, section) {
			t.Errorf("expected section %q in output", section)
		}
	}
}

func TestMarkdownReportWriter_EmptyResult(t *testing.T) {
	r := &Report{
		Metadata: Metadata{StagingFile: "a.env", ProductionFile: "b.env"},
		Result:   diff.Result{},
	}
	var buf bytes.Buffer
	w := NewMarkdownReportWriter()
	if err := w.Write(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "## Added") {
		t.Error("did not expect Added section for empty result")
	}
}
