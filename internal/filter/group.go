package filter

import (
	"sort"
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// GroupResult holds keys grouped by a common prefix delimiter.
type GroupResult struct {
	Groups map[string]*diff.Result
}

// ApplyGroup partitions a Result into sub-results keyed by the portion of each
// env-var name before the first occurrence of delimiter. Keys that do not
// contain the delimiter are placed under the empty-string group "".
func ApplyGroup(r diff.Result, delimiter string) GroupResult {
	if delimiter == "" {
		return GroupResult{Groups: map[string]*diff.Result{"": &r}}
	}

	groups := map[string]*diff.Result{}

	ensure := func(g string) {
		if _, ok := groups[g]; !ok {
			groups[g] = &diff.Result{
				Added:   map[string]string{},
				Removed: map[string]string{},
				Same:    map[string]string{},
				Changed: map[string]diff.ChangedValue{},
			}
		}
	}

	groupKey := func(k string) string {
		if idx := strings.Index(k, delimiter); idx >= 0 {
			return k[:idx]
		}
		return ""
	}

	for k, v := range r.Added {
		g := groupKey(k)
		ensure(g)
		groups[g].Added[k] = v
	}
	for k, v := range r.Removed {
		g := groupKey(k)
		ensure(g)
		groups[g].Removed[k] = v
	}
	for k, v := range r.Same {
		g := groupKey(k)
		ensure(g)
		groups[g].Same[k] = v
	}
	for k, v := range r.Changed {
		g := groupKey(k)
		ensure(g)
		groups[g].Changed[k] = v
	}

	return GroupResult{Groups: groups}
}

// SortedGroupKeys returns the group names from a GroupResult in sorted order.
func SortedGroupKeys(gr GroupResult) []string {
	keys := make([]string, 0, len(gr.Groups))
	for k := range gr.Groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// IsEmpty reports whether a GroupResult contains no groups or all groups are
// empty (i.e. every sub-result has no added, removed, changed, or same keys).
func (gr GroupResult) IsEmpty() bool {
	for _, r := range gr.Groups {
		if len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Changed) > 0 || len(r.Same) > 0 {
			return false
		}
	}
	return true
}
