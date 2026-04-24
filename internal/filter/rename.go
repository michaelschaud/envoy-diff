package filter

import (
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// RenameRule maps an old key prefix (or exact key) to a new one.
type RenameRule struct {
	From string
	To   string
}

// ApplyRename rewrites keys in the result according to the provided rename rules.
// Rules are applied in order; the first matching rule wins.
// Matching is case-insensitive prefix replacement.
func ApplyRename(result diff.Result, rules []RenameRule) diff.Result {
	if len(rules) == 0 {
		return result
	}
	return diff.Result{
		Added:   renameMap(result.Added, rules),
		Removed: renameMap(result.Removed, rules),
		Same:    renameMap(result.Same, rules),
		Changed: renameChanged(result.Changed, rules),
	}
}

func renameMap(m map[string]string, rules []RenameRule) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[applyRules(k, rules)] = v
	}
	return out
}

func renameChanged(m map[string]diff.ChangedValue, rules []RenameRule) map[string]diff.ChangedValue {
	out := make(map[string]diff.ChangedValue, len(m))
	for k, v := range m {
		out[applyRules(k, rules)] = v
	}
	return out
}

func applyRules(key string, rules []RenameRule) string {
	lower := strings.ToLower(key)
	for _, r := range rules {
		from := strings.ToLower(r.From)
		if lower == from {
			return r.To
		}
		if strings.HasPrefix(lower, from+"_") {
			return r.To + "_" + key[len(r.From)+1:]
		}
	}
	return key
}
