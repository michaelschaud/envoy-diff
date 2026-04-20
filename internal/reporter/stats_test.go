package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/envoy-diff/internal/diff"
)

func makeStats(added, removed, changed, total int) Stats {
	return Stats{
		Added:   added,
		Removed: removed,
		Changed: changed,
		Total:   total,
	}
}

func TestStats_HasDiff_True(t *testing.T) {
	s := makeStats(1, 0, 0, 5)
	if !s.HasDiff() {
		t.Error("expected HasDiff to be true")
	}
}

func TestStats_HasDiff_False(t *testing.T) {
	s := makeStats(0, 0, 0, 3)
	if s.HasDiff() {
		t.Error("expected HasDiff to be false")
	}
}

func TestStats_String(t *testing.T) {
	s := makeStats(2, 1, 3, 10)
	out := s.String()
	for _, sub := range []string{"total=10", "added=2", "removed=1", "changed=3"} {
		if !strings.Contains(out, sub) {
			t.Errorf("expected %q in output %q", sub, out)
		}
	}
}

func TestWriteStats_ContainsFields(t *testing.T) {
	var buf bytes.Buffer
	s := makeStats(3, 2, 1, 15)
	if err := WriteStats(&buf, s); err != nil {
		t.Fatalf("WriteStats error: %v", err)
	}
	out := buf.String()
	for _, sub := range []string{"Total keys", "Added", "Removed", "Changed", "15", "3", "2", "1"} {
		if !strings.Contains(out, sub) {
			t.Errorf("expected %q in output:\n%s", sub, out)
		}
	}
}

func TestStatsFromReport(t *testing.T) {
	r := New(
		diff.Result{
			Added:   map[string]string{"A": "1"},
			Removed: map[string]string{"B": "2", "C": "3"},
			Changed: map[string]diff.ChangedValue{},
			Same:    map[string]string{"D": "4"},
		},
		"staging.env",
		"production.env",
	)
	s := StatsFromReport(r)
	if s.Added != 1 {
		t.Errorf("expected Added=1, got %d", s.Added)
	}
	if s.Removed != 2 {
		t.Errorf("expected Removed=2, got %d", s.Removed)
	}
	if s.Changed != 0 {
		t.Errorf("expected Changed=0, got %d", s.Changed)
	}
}
