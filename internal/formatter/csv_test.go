package formatter

import (
	"bytes"
	"strings"
	"testing"
)

func TestCSVWriter_Headers(t *testing.T) {
	result := makeResult()
	var buf bytes.Buffer
	w := NewCSVWriter(&buf)
	if err := w.Write(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.HasPrefix(output, "category,key,staging_value,production_value") {
		t.Errorf("expected CSV header as first line, got: %q", output)
	}
}

func TestCSVWriter_ContainsRows(t *testing.T) {
	result := makeResult()
	var buf bytes.Buffer
	w := NewCSVWriter(&buf)
	if err := w.Write(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()

	expected := []string{"added", "removed", "changed", "identical"}
	for _, cat := range expected {
		if !strings.Contains(output, cat) {
			t.Errorf("expected category %q in CSV output", cat)
		}
	}
}

func TestCSVWriter_EmptyResult(t *testing.T) {
	result := makeResult()
	result.Added = nil
	result.Removed = nil
	result.Changed = nil
	result.Identical = nil

	var buf bytes.Buffer
	w := NewCSVWriter(&buf)
	if err := w.Write(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line, got %d lines", len(lines))
	}
}
