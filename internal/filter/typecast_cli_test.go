package filter

import (
	"testing"
)

func TestParseTypecastModes_Valid(t *testing.T) {
	modes, err := ParseTypecastModes("int,bool,float,string")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(modes) != 4 {
		t.Fatalf("expected 4 modes, got %d", len(modes))
	}
}

func TestParseTypecastModes_Empty(t *testing.T) {
	modes, err := ParseTypecastModes("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(modes) != 0 {
		t.Fatalf("expected 0 modes, got %d", len(modes))
	}
}

func TestParseTypecastModes_Unknown(t *testing.T) {
	_, err := ParseTypecastModes("int,banana")
	if err == nil {
		t.Fatal("expected error for unknown mode, got nil")
	}
}

func TestParseTypecastModes_CaseInsensitive(t *testing.T) {
	modes, err := ParseTypecastModes("INT,Bool")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(modes) != 2 {
		t.Fatalf("expected 2 modes, got %d", len(modes))
	}
	if modes[0] != TypecastInt {
		t.Errorf("expected TypecastInt, got %q", modes[0])
	}
	if modes[1] != TypecastBool {
		t.Errorf("expected TypecastBool, got %q", modes[1])
	}
}

func TestParseTypecastModes_WhitespaceHandled(t *testing.T) {
	modes, err := ParseTypecastModes(" int , float ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(modes) != 2 {
		t.Fatalf("expected 2 modes, got %d", len(modes))
	}
}
