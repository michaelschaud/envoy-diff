package filter_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
	"github.com/yourorg/envoy-diff/internal/filter"
)

func makeSortResult() diff.Result {
	return diff.Result{
		Added: map[string]string{
			"ZEBRA": "1",
			"APPLE": "2",
			"MANGO": "3",
		},
		Removed: map[string]string{
			"BETA":  "old",
			"ALPHA": "old2",
		},
		Changed: map[string][2]string{
			"PORT": {"8080", "9090"},
			"HOST": {"localhost", "prod.example.com"},
		},
		Same: map[string]string{
			"REGION": "us-east-1",
		},
	}
}

func TestApplySort_AscPreservesAllKeys(t *testing.T) {
	r := makeSortResult()
	out := filter.ApplySort(r, filter.SortAsc)

	if len(out.Added) != len(r.Added) {
		t.Errorf("expected %d added keys, got %d", len(r.Added), len(out.Added))
	}
	if len(out.Removed) != len(r.Removed) {
		t.Errorf("expected %d removed keys, got %d", len(r.Removed), len(out.Removed))
	}
	if len(out.Changed) != len(r.Changed) {
		t.Errorf("expected %d changed keys, got %d", len(r.Changed), len(out.Changed))
	}
}

func TestSortedKeys_Ascending(t *testing.T) {
	m := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	keys := filter.SortedKeys(m, filter.SortAsc)

	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %s, got %s", i, expected[i], k)
		}
	}
}

func TestSortedKeys_Descending(t *testing.T) {
	m := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	keys := filter.SortedKeys(m, filter.SortDesc)

	expected := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %s, got %s", i, expected[i], k)
		}
	}
}

func TestSortedChangedKeys_Ascending(t *testing.T) {
	m := map[string][2]string{
		"PORT": {"8080", "9090"},
		"HOST": {"localhost", "prod"},
	}
	keys := filter.SortedChangedKeys(m, filter.SortAsc)

	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "HOST" || keys[1] != "PORT" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestSortedChangedKeys_Descending(t *testing.T) {
	m := map[string][2]string{
		"PORT": {"8080", "9090"},
		"HOST": {"localhost", "prod"},
	}
	keys := filter.SortedChangedKeys(m, filter.SortDesc)

	if keys[0] != "PORT" || keys[1] != "HOST" {
		t.Errorf("unexpected order: %v", keys)
	}
}
