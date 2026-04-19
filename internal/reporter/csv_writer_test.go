package reporter

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func TestCSVReportWriter_Headers(t *testing.T) {
	r := makeResult()
	report := New("staging.env", "production.env", r)

	var buf bytes.Buffer
	w := NewCSVReportWriter()
	if err := w.Write(&buf, report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reader := csv.NewReader(strings.NewReader(buf.String()))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("invalid csv: %v", err)
	}

	if len(records) == 0 {
		t.Fatal("expected at least a header row")
	}

	header := records[0]
	expected := []string{"category", "key", "staging_value", "production_value"}
	for i, col := range expected {
		if header[i] != col {
			t.Errorf("header[%d] = %q, want %q", i, header[i], col)
		}
	}
}

func TestCSVReportWriter_ContainsRows(t *testing.T) {
	res := diff.Result{
		Added:   map[string]string{"NEW_KEY": "newval"},
		Removed: map[string]string{"OLD_KEY": "oldval"},
		Changed: map[string]diff.ValuePair{"MOD_KEY": {From: "v1", To: "v2"}},
		Same:    map[string]string{},
	}
	report := New("a.env", "b.env", res)

	var buf bytes.Buffer
	w := NewCSVReportWriter()
	if err := w.Write(&buf, report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := buf.String()
	for _, want := range []string{"added", "NEW_KEY", "removed", "OLD_KEY", "changed", "MOD_KEY"} {
		if !strings.Contains(body, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestCSVReportWriter_EmptyResult(t *testing.T) {
	res := diff.Result{
		Added:   map[string]string{},
		Removed: map[string]string{},
		Changed: map[string]diff.ValuePair{},
		Same:    map[string]string{},
	}
	report := New("a.env", "b.env", res)

	var buf bytes.Buffer
	w := NewCSVReportWriter()
	if err := w.Write(&buf, report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reader := csv.NewReader(strings.NewReader(buf.String()))
	records, _ := reader.ReadAll()
	if len(records) != 1 {
		t.Errorf("expected only header row, got %d rows", len(records))
	}
}
