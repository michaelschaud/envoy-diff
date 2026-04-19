package filter

import (
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// ApplyMulti applies multiple prefix filters using OR logic: a key is kept if
// it matches any of the provided prefixes. If prefixes is empty, all keys are
// kept. onlyChanged behaviour is identical to Apply.
func ApplyMulti(result diff.Result, prefixes []string, onlyChanged bool) diff.Result {
	if len(prefixes) == 0 {
		if onlyChanged {
			return diff.Result{
				Added:   result.Added,
				Removed: result.Removed,
				Changed: result.Changed,
				Same:    map[string]string{},
			}
		}
		return result
	}

	matches := func(key string) bool {
		for _, p := range prefixes {
			if strings.HasPrefix(strings.ToUpper(key), strings.ToUpper(p)) {
				return true
			}
		}
		return false
	}

	filterMap := func(m map[string]string) map[string]string {
		out := map[string]string{}
		for k, v := range m {
			if matches(k) {
				out[k] = v
			}
		}
		return out
	}

	filterChanged := func(m map[string][2]string) map[string][2]string {
		out := map[string][2]string{}
		for k, v := range m {
			if matches(k) {
				out[k] = v
			}
		}
		return out
	}

	out := diff.Result{
		Added:   filterMap(result.Added),
		Removed: filterMap(result.Removed),
		Changed: filterChanged(result.Changed),
	}
	if onlyChanged {
		out.Same = map[string]string{}
	} else {
		out.Same = filterMap(result.Same)
	}
	return out
}
