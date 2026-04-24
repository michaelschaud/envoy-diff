package filter

import (
	"fmt"
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// AnnotateRule describes a single annotation: a key substring to match and a
// label to attach to matching entries in the diff result.
type AnnotateRule struct {
	KeySubstring string
	Label        string
}

// ParseAnnotateRules parses a slice of "substring=label" strings into
// AnnotateRule values. Entries that are malformed or have an empty substring
// are silently skipped so that invalid CLI input does not crash the tool.
func ParseAnnotateRules(raw []string) []AnnotateRule {
	rules := make([]AnnotateRule, 0, len(raw))
	for _, r := range raw {
		idx := strings.Index(r, "=")
		if idx <= 0 {
			continue
		}
		substr := strings.TrimSpace(r[:idx])
		label := strings.TrimSpace(r[idx+1:])
		if substr == "" {
			continue
		}
		rules = append(rules, AnnotateRule{KeySubstring: substr, Label: label})
	}
	return rules
}

// ApplyAnnotate rewrites the values of matching keys so that the annotation
// label is appended in square brackets, e.g. "somevalue [deprecated]".
// Matching is case-insensitive on the key name. Keys that match multiple rules
// receive all matching labels in the order the rules were defined.
//
// The original diff.Result is not mutated; a new Result is returned.
func ApplyAnnotate(result diff.Result, rules []AnnotateRule) diff.Result {
	if len(rules) == 0 {
		return result
	}

	return diff.Result{
		Added:   annotateMap(result.Added, rules),
		Removed: annotateMap(result.Removed, rules),
		Changed: annotateChanged(result.Changed, rules),
		Same:    annotateMap(result.Same, rules),
	}
}

// annotateMap returns a copy of m with values annotated for any key that
// matches at least one rule.
func annotateMap(m map[string]string, rules []AnnotateRule) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = annotateValue(k, v, rules)
	}
	return out
}

// annotateChanged returns a copy of the changed map with values annotated.
func annotateChanged(m map[string][2]string, rules []AnnotateRule) map[string][2]string {
	out := make(map[string][2]string, len(m))
	for k, pair := range m {
		old := annotateValue(k, pair[0], rules)
		new_ := annotateValue(k, pair[1], rules)
		out[k] = [2]string{old, new_}
	}
	return out
}

// annotateValue appends labels for each rule whose KeySubstring is found
// (case-insensitively) in key. If no rules match, the original value is
// returned unchanged.
func annotateValue(key, value string, rules []AnnotateRule) string {
	lowerKey := strings.ToLower(key)
	labels := []string{}
	for _, r := range rules {
		if strings.Contains(lowerKey, strings.ToLower(r.KeySubstring)) {
			if r.Label != "" {
				labels = append(labels, r.Label)
			}
		}
	}
	if len(labels) == 0 {
		return value
	}
	return fmt.Sprintf("%s [%s]", value, strings.Join(labels, ", "))
}
