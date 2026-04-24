package filter

import (
	"sort"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// SortOrder defines the ordering direction for diff results.
type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// ApplySort returns a new diff.Result with all key slices sorted
// alphabetically in the specified direction.
func ApplySort(r diff.Result, order SortOrder) diff.Result {
	sorted := diff.Result{
		Added:   sortedCopy(r.Added, order),
		Removed: sortedCopy(r.Removed, order),
		Same:    sortedCopy(r.Same, order),
		Changed: sortedChangedCopy(r.Changed, order),
	}
	return sorted
}

func sortedCopy(m map[string]string, order SortOrder) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func sortedChangedCopy(m map[string][2]string, order SortOrder) map[string][2]string {
	if m == nil {
		return nil
	}
	out := make(map[string][2]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// SortedKeys returns the keys of a map[string]string in the given order.
func SortedKeys(m map[string]string, order SortOrder) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if order == SortDesc {
		reverse(keys)
	}
	return keys
}

// SortedChangedKeys returns the keys of a map[string][2]string in the given order.
func SortedChangedKeys(m map[string][2]string, order SortOrder) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if order == SortDesc {
		reverse(keys)
	}
	return keys
}

func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
