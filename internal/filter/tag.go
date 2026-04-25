package filter

import (
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// TagRule associates a tag label with a key substring match.
type TagRule struct {
	Tag    string
	Substr string
}

// ApplyTag adds a tag suffix to values whose keys match any of the given rules.
// The tag is appended as " [tag]" to the value string.
func ApplyTag(result diff.Result, rules []TagRule) diff.Result {
	if len(rules) == 0 {
		return result
	}
	return diff.Result{
		Added:   tagMap(result.Added, rules),
		Removed: tagMap(result.Removed, rules),
		Same:    tagMap(result.Same, rules),
		Changed: tagChanged(result.Changed, rules),
	}
}

func tagMap(m map[string]string, rules []TagRule) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if tag, ok := matchTagRule(k, rules); ok {
			out[k] = v + " [" + tag + "]"
		} else {
			out[k] = v
		}
	}
	return out
}

func tagChanged(m map[string][2]string, rules []TagRule) map[string][2]string {
	out := make(map[string][2]string, len(m))
	for k, pair := range m {
		if tag, ok := matchTagRule(k, rules); ok {
			out[k] = [2]string{
				pair[0] + " [" + tag + "]",
				pair[1] + " [" + tag + "]",
			}
		} else {
			out[k] = pair
		}
	}
	return out
}

func matchTagRule(key string, rules []TagRule) (string, bool) {
	lower := strings.ToLower(key)
	for _, r := range rules {
		if strings.Contains(lower, strings.ToLower(r.Substr)) {
			return r.Tag, true
		}
	}
	return "", false
}
