package filter

import (
	"strings"
	"testing"

	"github.com/user/envoy-diff/internal/diff"
)

func makeHighlightResult() diff.Result {
	return diff.Result{
		Added:   map[string]string{"SECRET_KEY": "abc", "HOST": "prod"},
		Removed: map[string]string{"OLD_TOKEN": "xyz"},
		Same:    map[string]string{"PORT": "8080"},
		Changed: map[string][2]string{"DB_PASSWORD": {"old", "new"}},
	}
}

func TestApplyHighlight_NoSubstrings(t *testing.T) {
	r := makeHighlightResult()
	out := ApplyHighlight(r, HighlightConfig{Marker: " [!]"})
	if out.Added["SECRET_KEY"] != "abc" {
		t.Errorf("expected unchanged value, got %q", out.Added["SECRET_KEY"])
	}
}

func TestApplyHighlight_NoMarker(t *testing.T) {
	r := makeHighlightResult()
	out := ApplyHighlight(r, HighlightConfig{Substrings: []string{"SECRET"}})
	if out.Added["SECRET_KEY"] != "abc" {
		t.Errorf("expected unchanged value, got %q", out.Added["SECRET_KEY"])
	}
}

func TestApplyHighlight_MarksAdded(t *testing.T) {
	r := makeHighlightResult()
	out := ApplyHighlight(r, HighlightConfig{Marker: " [!]", Substrings: []string{"SECRET"}})
	if !strings.HasSuffix(out.Added["SECRET_KEY"], " [!]") {
		t.Errorf("expected marker on SECRET_KEY, got %q", out.Added["SECRET_KEY"])
	}
	if strings.HasSuffix(out.Added["HOST"], " [!]") {
		t.Errorf("HOST should not be marked, got %q", out.Added["HOST"])
	}
}

func TestApplyHighlight_MarksRemoved(t *testing.T) {
	r := makeHighlightResult()
	out := ApplyHighlight(r, HighlightConfig{Marker: " ***", Substrings: []string{"TOKEN"}})
	if !strings.HasSuffix(out.Removed["OLD_TOKEN"], " ***") {
		t.Errorf("expected marker on OLD_TOKEN, got %q", out.Removed["OLD_TOKEN"])
	}
}

func TestApplyHighlight_MarksChanged(t *testing.T) {
	r := makeHighlightResult()
	out := ApplyHighlight(r, HighlightConfig{Marker: " [!]", Substrings: []string{"PASSWORD"}})
	pair := out.Changed["DB_PASSWORD"]
	if !strings.HasSuffix(pair[0], " [!]") || !strings.HasSuffix(pair[1], " [!]") {
		t.Errorf("expected both sides of changed pair to be marked, got %v", pair)
	}
}

func TestApplyHighlight_CaseInsensitive(t *testing.T) {
	r := makeHighlightResult()
	out := ApplyHighlight(r, HighlightConfig{Marker: "!", Substrings: []string{"secret"}})
	if !strings.HasSuffix(out.Added["SECRET_KEY"], "!") {
		t.Errorf("expected case-insensitive match on SECRET_KEY, got %q", out.Added["SECRET_KEY"])
	}
}

func TestParseHighlightConfig_Valid(t *testing.T) {
	cfg, err := ParseHighlightConfig("[!]:SECRET,TOKEN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Marker != "[!]" {
		t.Errorf("expected marker [!], got %q", cfg.Marker)
	}
	if len(cfg.Substrings) != 2 {
		t.Errorf("expected 2 substrings, got %d", len(cfg.Substrings))
	}
}

func TestParseHighlightConfig_EmptySpec(t *testing.T) {
	_, err := ParseHighlightConfig("")
	if err == nil {
		t.Error("expected error for empty spec")
	}
}

func TestParseHighlightConfig_MissingColon(t *testing.T) {
	_, err := ParseHighlightConfig("[!]SECRET")
	if err == nil {
		t.Error("expected error for missing colon")
	}
}

func TestParseHighlightConfig_EmptyMarker(t *testing.T) {
	_, err := ParseHighlightConfig(":SECRET")
	if err == nil {
		t.Error("expected error for empty marker")
	}
}
