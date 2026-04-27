package filter

import "testing"

func TestParseMergeStrategy_Valid(t *testing.T) {
	cases := []struct {
		input    string
		want     MergeStrategy
	}{
		{"left", MergeStrategyLeft},
		{"right", MergeStrategyRight},
		{"union", MergeStrategyUnion},
		{"LEFT", MergeStrategyLeft},
		{" right ", MergeStrategyRight},
	}
	for _, tc := range cases {
		got, err := ParseMergeStrategy(tc.input)
		if err != nil {
			t.Errorf("ParseMergeStrategy(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseMergeStrategy(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseMergeStrategy_Invalid(t *testing.T) {
	_, err := ParseMergeStrategy("unknown")
	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestParseMergeOverlay_Valid(t *testing.T) {
	pairs := []string{"FOO=bar", "BAZ=qux"}
	out := ParseMergeOverlay(pairs)
	if out["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", out["FOO"])
	}
	if out["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", out["BAZ"])
	}
}

func TestParseMergeOverlay_SkipsMalformed(t *testing.T) {
	pairs := []string{"NOEQUALS", "=EMPTYKEY", "VALID=yes"}
	out := ParseMergeOverlay(pairs)
	if len(out) != 1 {
		t.Errorf("expected 1 valid entry, got %d: %v", len(out), out)
	}
	if out["VALID"] != "yes" {
		t.Errorf("expected VALID=yes, got %q", out["VALID"])
	}
}

func TestParseMergeOverlay_Empty(t *testing.T) {
	out := ParseMergeOverlay(nil)
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
