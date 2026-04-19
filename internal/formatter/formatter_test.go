package formatter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envoy-diff/internal/diff"
	"github.com/envoy-diff/internal/formatter"
)

func makeResult() diff.Result {
	return diff.Result{
		OnlyInA: map[string]string{"STAGING_ONLY": "foo"},
		OnlyInB: map[string]string{"PROD_ONLY": "bar"},
		Changed: map[string][2]string{
			"DB_HOST": {"localhost", "db.prod.internal"},
		},
		Common: map[string]string{"APP_NAME": "envoy"},
	}
}

func TestTextWriter_ContainsSections(t *testing.T) {
	var buf bytes.Buffer
	formatter.TextWriter(&buf, makeResult())
	out := buf.String()

	expected := []string{
		"only in staging",
		"STAGING_ONLY=foo",
		"only in production",
		"PROD_ONLY=bar",
		"changed values",
		"DB_HOST",
		"localhost",
		"db.prod.internal",
	}
	for _, s := range expected {
		if !strings.Contains(out, s) {
			t.Errorf("expected output to contain %q\ngot:\n%s", s, out)
		}
	}
}

func TestTextWriter_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	formatter.TextWriter(&buf, diff.Result{
		Common: map[string]string{"X": "1"},
	})
	if !strings.Contains(buf.String(), "No differences found") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestSummary(t *testing.T) {
	s := formatter.Summary(makeResult())
	if !strings.Contains(s, "1 only-in-staging") {
		t.Errorf("unexpected summary: %s", s)
	}
	if !strings.Contains(s, "1 only-in-production") {
		t.Errorf("unexpected summary: %s", s)
	}
	if !strings.Contains(s, "1 changed") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_Identical(t *testing.T) {
	s := formatter.Summary(diff.Result{})
	if s != "identical" {
		t.Errorf("expected 'identical', got %q", s)
	}
}
