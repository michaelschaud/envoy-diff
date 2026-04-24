package filter

import "testing"

func TestParseDedupeMode_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  DedupeMode
	}{
		{"first", DedupeKeepFirst},
		{"last", DedupeKeepLast},
		{"all", DedupeRemoveAll},
		{"", ""},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseDedupeMode(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseDedupeMode_Invalid(t *testing.T) {
	_, err := ParseDedupeMode("unknown")
	if err == nil {
		t.Error("expected error for unknown mode")
	}
}

func TestParseDedupeMode_CaseSensitive(t *testing.T) {
	_, err := ParseDedupeMode("First")
	if err == nil {
		t.Error("expected error: mode matching should be case-sensitive")
	}
}
