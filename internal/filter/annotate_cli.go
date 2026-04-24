package filter

import (
	"fmt"
	"strings"
)

// ParseAnnotateRules parses a slice of "key=template" strings into AnnotateRule pairs.
// The template may reference {{value}} as a placeholder for the original value.
func ParseAnnotateRules(rules []string) ([]AnnotateRule, error) {
	var parsed []AnnotateRule
	for _, r := range rules {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid annotate rule %q: expected key=template", r)
		}
		key := strings.TrimSpace(parts[0])
		template := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, fmt.Errorf("invalid annotate rule %q: key must not be empty", r)
		}
		if template == "" {
			return nil, fmt.Errorf("invalid annotate rule %q: template must not be empty", r)
		}
		parsed = append(parsed, AnnotateRule{Key: key, Template: template})
	}
	return parsed, nil
}
