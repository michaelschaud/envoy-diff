package formatter

import (
	"bytes"
	"testing"
)

func TestNewWriter_Text(t *testing.T) {
	w, err := NewWriter("text", &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := w.(*TextWriter); !ok {
		t.Errorf("expected *TextWriter, got %T", w)
	}
}

func TestNewWriter_JSON(t *testing.T) {
	w, err := NewWriter("json", &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := w.(*JSONWriter); !ok {
		t.Errorf("expected *JSONWriter, got %T", w)
	}
}

func TestNewWriter_CSV(t *testing.T) {
	w, err := NewWriter("csv", &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := w.(*CSVWriter); !ok {
		t.Errorf("expected *CSVWriter, got %T", w)
	}
}

func TestNewWriter_Unknown(t *testing.T) {
	_, err := NewWriter("yaml", &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestNewWriter_WritesWithoutError(t *testing.T) {
	formats := []string{"text", "json", "csv"}
	result := makeResult()
	for _, f := range formats {
		var buf bytes.Buffer
		w, err := NewWriter(f, &buf)
		if err != nil {
			t.Fatalf("format %q: unexpected error creating writer: %v", f, err)
		}
		if err := w.Write(result); err != nil {
			t.Errorf("format %q: Write returned error: %v", f, err)
		}
	}
}
